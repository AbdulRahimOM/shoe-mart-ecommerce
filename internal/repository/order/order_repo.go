package orderrepo

import (
	"context"

	e "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/customErrors"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/domain/entities"
	response "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/models/responseModels"
	repoInterface "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/interface"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"gorm.io/gorm"
)

var (
	errOrderDontExist_400 = e.Error{StatusCode: 400, Status: "Failed", Msg: "No such order exists", Err: nil}
)

type OrderRepo struct {
	DB  *gorm.DB
	Cld *cloudinary.Cloudinary
}

func NewOrderRepository(db *gorm.DB, cloudinary *cloudinary.Cloudinary) repoInterface.IOrderRepo {
	return &OrderRepo{
		DB:  db,
		Cld: cloudinary,
	}
}

// MakeOrder implements repo.IOrderRepo.
func (repo *OrderRepo) MakeOrder(order *entities.Order, orderItems *[]entities.OrderItem) (uint, *e.Error) {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			tx.Rollback()
		}
	}()

	//add order
	result = tx.Create(&order)
	if result.Error != nil {
		tx.Rollback()
		return 0, e.DBQueryError_500(&result.Error)
	}

	//create order items
	for _, item := range *orderItems {
		//update orderItems with orderID
		item.OrderID = order.ID

		//add order item to db
		result = tx.Create(&item)
		if result.Error != nil {
			tx.Rollback()
			return 0, e.DBQueryError_500(&result.Error)
		}
	}

	//clear cart
	result = tx.Where("user_id = ?", order.UserID).Delete(&entities.Cart{})
	if result.Error != nil {
		tx.Rollback()
		return 0, e.DBQueryError_500(&result.Error)
	}

	//commit transaction
	tx.Commit()

	return order.ID, nil
}

// GetOrdersOfUser
func (repo *OrderRepo) GetOrdersOfUser(userID uint, resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, *e.Error) {
	var orders []entities.Order
	query := repo.DB.
		Preload("FkAddress").
		Where("user_id = ?", userID).
		Order("id desc").
		Limit(resultLimit).
		Offset(resultOffset).
		Find(&orders)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	var orderInfos []entities.DetailedOrderInfo
	for _, order := range orders {
		//get order items
		var orderItems []entities.OrderItem
		query := repo.DB.
			Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
			Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
			Where("order_id = ?", order.ID).
			Find(&orderItems)

		if query.Error != nil {
			return nil, e.DBQueryError_500(&query.Error)
		}

		var orderInfo entities.DetailedOrderInfo
		orderInfo.OrderDetails = order
		orderInfo.OrderItems = orderItems

		orderInfos = append(orderInfos, orderInfo)
	}
	return &orderInfos, nil
}

// Get All orders (for admin)
func (repo *OrderRepo) GetOrders(resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, *e.Error) {
	var orders []entities.Order
	query := repo.DB.
		Preload("FkAddress").
		Limit(resultLimit).
		Offset(resultOffset).
		Order("id desc").
		Find(&orders)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	var orderInfos []entities.DetailedOrderInfo
	for _, order := range orders {
		//get order items
		var orderItems []entities.OrderItem
		query := repo.DB.
			Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkBrand").
			Preload("FkProduct.FkDimensionalVariation.FkColourVariant.FkModel.FkCategory").
			Where("order_id = ?", order.ID).
			Find(&orderItems)

		if query.Error != nil {
			return nil, e.DBQueryError_500(&query.Error)
		}

		var orderInfo entities.DetailedOrderInfo
		orderInfo.OrderDetails = order
		orderInfo.OrderItems = orderItems

		orderInfos = append(orderInfos, orderInfo)
	}
	return &orderInfos, nil
}

// CancelOrder implements repo.IOrderRepo.
func (repo *OrderRepo) CancelOrder(orderID uint) *e.Error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			tx.Rollback()
		}
	}()

	//update order status
	result = tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "cancelled")
	if result.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&result.Error)
	}

	//get order items
	var orderItems []entities.OrderItem
	query := tx.
		Where("order_id = ?", orderID).
		Find(&orderItems)

	if query.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&query.Error)
	}

	//update stock
	for _, item := range orderItems {
		result := tx.Model(&entities.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", item.Quantity))
		if result.Error != nil {
			tx.Rollback()
			return e.DBQueryError_500(&result.Error)
		}
	}

	//get order's final amount (to update wallet)
	var order entities.Order
	query = tx.
		Select("final_amount,payment_status").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&query.Error)
	}
	if order.PaymentStatus == "paid" {
		//update wallet , update payment status to refunded
		result = tx.Model(&entities.User{}).Where("id = ?", order.UserID).Update("wallet_balance", gorm.Expr("wallet_balance + ?", order.FinalAmount))
		if result.Error != nil {
			tx.Rollback()
			return e.DBQueryError_500(&result.Error)
		}

		//update payment status to refunded
		result = tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("payment_status", "refunded")
		if result.Error != nil {
			tx.Rollback()
			return e.DBQueryError_500(&result.Error)
		}
	}

	//commit transaction
	tx.Commit()

	return nil
}

func (repo *OrderRepo) GetOrderStatusByID(orderID uint) (string, *e.Error) {
	var status string
	query := repo.DB.
		Table("orders").
		Select("status").
		Where("id = ?", orderID).
		Find(&status)

	if query.Error != nil {
		return "", e.DBQueryError_500(&query.Error)
	}
	if query.RowsAffected == 0 {
		return "", &errOrderDontExist_400
	}

	return status, nil
}

func (repo *OrderRepo) GetUserIDByOrderID(orderID uint) (uint, *e.Error) {
	var order entities.Order
	query := repo.DB.
		Table("orders").
		Select("user_id").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		return 0, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return 0, &errOrderDontExist_400
	}

	return order.UserID, nil
}

func (repo *OrderRepo) MakeOrder_UpdateStock_ClearCart(order *entities.Order, orderItems *[]entities.OrderItem) (uint, *e.Error) {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			tx.Rollback()
		}
	}()

	//add order
	result = tx.Create(&order)
	if result.Error != nil {
		tx.Rollback()
		return 0, e.DBQueryError_500(&result.Error)
	}

	//create order items
	for _, item := range *orderItems {
		//update orderItems with orderID
		item.OrderID = order.ID

		//add order item to db
		result := tx.Create(&item)
		if result.Error != nil {
			tx.Rollback()
			return 0, e.DBQueryError_500(&result.Error)
		}

		//update stock
		result = tx.Model(&entities.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", item.Quantity))
		if result.Error != nil {
			tx.Rollback()
			return 0, e.DBQueryError_500(&result.Error)
		}
	}

	//clear cart
	result = tx.Where("user_id = ?", order.UserID).Delete(&entities.Cart{})
	if result.Error != nil {
		tx.Rollback()
		return 0, e.DBQueryError_500(&result.Error)
	}

	//commit transaction
	tx.Commit()

	return order.ID, nil
}

// ReturnOrder
func (repo *OrderRepo) ReturnOrderRequest(orderID uint) *e.Error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			tx.Rollback()
		}
	}()

	//update order status
	result = tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "return requested")
	if result.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&result.Error)
	}

	//commit transaction
	tx.Commit()

	return nil
}

// MarkOrderAsReturned
func (repo *OrderRepo) MarkOrderAsReturned(orderID uint) *e.Error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			tx.Rollback()
		}
	}()

	//update order status
	result = tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "returned")
	if result.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&result.Error)
	}

	//get order items
	var orderItems []entities.OrderItem
	query := tx.
		Where("order_id = ?", orderID).
		Find(&orderItems)

	if query.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&result.Error)
	}

	//update stock
	for _, item := range orderItems {
		result := tx.Model(&entities.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", item.Quantity))
		if result.Error != nil {
			tx.Rollback()
			return e.DBQueryError_500(&result.Error)
		}
	}

	//get order's final amount and userID (to update wallet)
	var order entities.Order
	query = tx.
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&result.Error)
	}

	//update wallet
	result = tx.Model(&entities.User{}).Where("id = ?", order.UserID).Update("wallet_balance", gorm.Expr("wallet_balance + ?", order.FinalAmount))
	if result.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&result.Error)
	}

	//commit transaction
	tx.Commit()

	return nil
}

// MarkOrderAsDelivered
func (repo *OrderRepo) MarkOrderAsDelivered(orderID uint) *e.Error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			tx.Rollback()
		}
	}()

	//update order status and delivered_date
	result = tx.Model(&entities.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{"status": "delivered", "delivered_date": gorm.Expr("CURRENT_TIMESTAMP"), "payment_status": "paid"})
	if result.Error != nil {
		tx.Rollback()
		return e.DBQueryError_500(&result.Error)
	}

	//commit transaction
	tx.Commit()

	return nil
}

// GetAllOrders
func (repo *OrderRepo) GetAllOrders() (*[]entities.Order, *e.Error) {
	var orders []entities.Order
	query := repo.DB.
		Find(&orders)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &orders, nil
}

func (repo *OrderRepo) GetOrderSummaryByID(orderID uint) (*entities.Order, *e.Error) {
	var order entities.Order
	query := repo.DB.
		Preload("FkAddress").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	if query.RowsAffected == 0 {
		return nil, &errOrderDontExist_400
	}

	return &order, nil
}

func (repo *OrderRepo) UpdateOrderToPaid_UpdateStock_ClearCart(orderID uint) (*entities.Order, *e.Error) {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			tx.Rollback()
		}
	}()

	//update order status to "placed" and payment status to "paid"
	result = tx.Model(&entities.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{"status": "placed", "payment_status": "paid"})
	if result.Error != nil {
		tx.Rollback()
		return nil, e.DBQueryError_500(&result.Error)
	}

	//get order
	var order entities.Order
	query := tx.
		Preload("FkAddress").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		tx.Rollback()
		return nil, e.DBQueryError_500(&query.Error)
	}

	//get order items
	var orderItems *[]entities.OrderItem
	query = tx.
		Where("order_id = ?", order.ID).
		Find(&orderItems)

	if query.Error != nil {
		tx.Rollback()
		return nil, e.DBQueryError_500(&query.Error)
	}

	//update stock
	for _, item := range *orderItems {
		result := tx.Model(&entities.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", item.Quantity))
		if result.Error != nil {
			tx.Rollback()
			return nil, e.DBQueryError_500(&result.Error)
		}
	}

	//clear cart
	result = tx.Where("user_id = ?", order.UserID).Delete(&entities.Cart{})
	if result.Error != nil {
		tx.Rollback()
		return nil, e.DBQueryError_500(&result.Error)
	}

	//commit transaction
	tx.Commit()

	return &order, nil
}

func (repo *OrderRepo) GetOrderByTransactionID(transactionID string) (uint, *e.Error) {
	var order entities.Order
	query := repo.DB.
		Where("transaction_id = ?", transactionID).
		Find(&order)

	if query.Error != nil {
		return 0, e.DBQueryError_500(&query.Error)
	}

	return order.ID, nil
}

func (repo *OrderRepo) UpdateOrderTransactionID(orderID uint, transactionID string) *e.Error {
	result := repo.DB.Model(&entities.Order{}).Where("id = ?", orderID).Update("transaction_id", transactionID)
	if result.Error != nil {
		return e.DBQueryError_500(&result.Error)
	} else if result.RowsAffected == 0 {
		return &errOrderDontExist_400
	}

	return nil
}

// GetPaymentStatusByID implements repo.IOrderRepo.
func (repo *OrderRepo) GetPaymentStatusByID(orderID uint) (string, *e.Error) {
	var status string
	query := repo.DB.
		Table("orders").
		Select("payment_status").
		Where("id = ?", orderID).
		Find(&status)

	if query.Error != nil {
		return "", e.DBQueryError_500(&query.Error)
	} else if query.RowsAffected == 0 {
		return "", &errOrderDontExist_400
	}

	return status, nil
}

// GetOrderItemsPQRByOrderID implements repo.IOrderRepo.
func (repo *OrderRepo) GetOrderItemsPQRByOrderID(orderID uint) (*[]response.PQMS, *e.Error) {
	var orderItems []response.PQMS
	query := repo.DB.Raw(`
		SELECT
			product.id AS "productID",
			product.name AS "productName",
			order_items.quantity AS "quantity",
			colour_variants.mrp AS "mrp",
			colour_variants."salePrice" AS "salePrice"
		FROM order_items
		JOIN product ON order_items."product_id" = product.id
		JOIN dimensional_variants ON product."dimensionalVariationID" = dimensional_variants.id
		JOIN colour_variants ON dimensional_variants."colourVariantId" = colour_variants.id
		WHERE order_items."order_id" = ?`,
		orderID).Scan(&orderItems)

	if query.Error != nil {
		return nil, e.DBQueryError_500(&query.Error)
	}

	return &orderItems, nil
}

func (repo *OrderRepo) UploadInvoice(file string, fileName string) (*string, *e.Error) {
	result, err := repo.Cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:    "github.com/AbdulRahimOM/shoe-mart-ecommerce/invoices",
		PublicID:  fileName,
		Overwrite: true,
	})
	if err != nil {
		return nil, e.SetError("error while uploading file to cloudinary. err: ", err, 500)
	}

	if result.Error.Message != "" {
		return nil, e.SetError("error while uploading file to cloudinary. result.Error: "+result.Error.Message, nil, 500)
	}

	return &result.SecureURL, nil
}
