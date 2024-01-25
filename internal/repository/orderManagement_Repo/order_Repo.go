package ordermanagementrepo

import (
	"MyShoo/internal/domain/entities"
	repoInterface "MyShoo/internal/repository/interface"
	"fmt"

	"gorm.io/gorm"
)

type OrderRepo struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repoInterface.IOrderRepo {
	return &OrderRepo{DB: db}
}

// MakeOrder implements repository_interface.IOrderRepo.
func (repo *OrderRepo) MakeOrder(order *entities.Order, orderItems *[]entities.OrderItem) (*entities.Order, error) {
	//add order
	result := repo.DB.Create(&order)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add order. query.Error= ", result.Error, "\n----")
		return order, result.Error
	}

	//preload order
	result = repo.DB.Preload("FkAddress").First(&order)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't preload order. query.Error= ", result.Error, "\n----")
		return order, result.Error
	}

	//create order items
	for _, item := range *orderItems {
		//update orderItems with orderID
		item.OrderID = order.ID

		//add order item to db
		result := repo.DB.Create(&item)
		if result.Error != nil {
			fmt.Println("-------\nquery error happened. couldn't add order item. query.Error= ", result.Error, "\n----")
			return order, result.Error
		}
	}

	return order, nil
}

// GetOrdersOfUser
func (repo *OrderRepo) GetOrdersOfUser(userID uint, resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, error) {
	var orders []entities.Order
	query := repo.DB.
		Preload("FkAddress").
		Where("user_id = ?", userID).
		Order("id desc").
		Limit(resultLimit).
		Offset(resultOffset).
		Find(&orders)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
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
			fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
			return nil, query.Error
		}

		var orderInfo entities.DetailedOrderInfo
		orderInfo.OrderDetails = order
		orderInfo.OrderItems = orderItems

		orderInfos = append(orderInfos, orderInfo)
	}
	return &orderInfos, nil
}

// Get All orders (for admin)
func (repo *OrderRepo) GetOrders(resultOffset int, resultLimit int) (*[]entities.DetailedOrderInfo, error) {
	var orders []entities.Order
	query := repo.DB.
		Preload("FkAddress").
		Limit(resultLimit).
		Offset(resultOffset).
		Order("id desc").
		Find(&orders)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
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
			fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
			return nil, query.Error
		}

		var orderInfo entities.DetailedOrderInfo
		orderInfo.OrderDetails = order
		orderInfo.OrderItems = orderItems

		orderInfos = append(orderInfos, orderInfo)
	}
	return &orderInfos, nil
}

// CancelOrder implements repository_interface.IOrderRepo.
func (repo *OrderRepo) CancelOrder(orderID uint) error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			fmt.Println("-------\npanic happened. couldn't cancel order. r= ", r, "query.Error= ", result.Error, "\n----")
			tx.Rollback()
		}
	}()

	//update order status
	result = tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "cancelled")
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't cancel order. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return result.Error
	}

	//get order items
	var orderItems []entities.OrderItem
	query := tx.
		Where("order_id = ?", orderID).
		Find(&orderItems)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		tx.Rollback()
		return query.Error
	}

	//update stock
	for _, item := range orderItems {
		result := tx.Model(&entities.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", item.Quantity))
		if result.Error != nil {
			fmt.Println("-------\nquery error happened. couldn't update stock. query.Error= ", result.Error, "\n----")
			tx.Rollback()
			return result.Error
		}
	}

	//get order's final amount (to update wallet)
	var order entities.Order
	query = tx.
		Select("final_amount").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		tx.Rollback()
		return query.Error
	}

	//update wallet
	result = tx.Model(&entities.User{}).Where("id = ?", order.UserID).Update("wallet_balance", gorm.Expr("wallet_balance + ?", order.FinalAmount))
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't update wallet. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return result.Error
	}

	//commit transaction
	tx.Commit()

	return nil
}

func (repo *OrderRepo) DoOrderExistByID(orderID uint) (bool, error) {
	var temp entities.Order
	query := repo.DB.Raw(`
		SELECT *
		FROM orders
		WHERE id = ?`,
		orderID).Scan(&temp)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't check if-order is existing or not. query.Error= ", query.Error, "\n----")
		return false, query.Error
	}

	if query.RowsAffected == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo *OrderRepo) GetOrderStatusByID(orderID uint) (string, error) {
	var order entities.Order
	query := repo.DB.
		Select("status").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return "", query.Error
	}

	return order.Status, nil
}

func (repo *OrderRepo) GetUserIDByOrderID(orderID uint) (uint, error) {
	var order entities.Order
	query := repo.DB.
		Select("user_id").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return 0, query.Error
	}

	return order.UserID, nil
}

func (repo *OrderRepo) MakeOrder_UpdateStock_ClearCart(order *entities.Order, orderItems *[]entities.OrderItem) (*entities.Order, error) {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			fmt.Println("-------\npanic happened. couldn't cancel order. r= ", r, "query.Error= ", result.Error, "\n----")
			tx.Rollback()
		}
	}()

	//add order
	result = tx.Create(&order)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't add order. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return order, result.Error
	}

	//preload order
	result = tx.Preload("FkAddress").First(&order)
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't preload order. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return order, result.Error
	}

	//create order items
	for _, item := range *orderItems {
		//update orderItems with orderID
		item.OrderID = order.ID

		//add order item to db
		result := tx.Create(&item)
		if result.Error != nil {
			fmt.Println("-------\nquery error happened. couldn't add order item. query.Error= ", result.Error, "\n----")
			tx.Rollback()
			return order, result.Error
		}

		//update stock
		result = tx.Model(&entities.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", item.Quantity))
		if result.Error != nil {
			fmt.Println("-------\nquery error happened. couldn't update stock. query.Error= ", result.Error, "\n----")
			tx.Rollback()
			return order, result.Error
		}
	}

	//clear cart
	result = tx.Where("user_id = ?", order.UserID).Delete(&entities.Cart{})
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't clear cart. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return order, result.Error
	}

	//commit transaction
	tx.Commit()

	return order, nil
}

// ReturnOrder
func (repo *OrderRepo) ReturnOrderRequest(orderID uint) error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			fmt.Println("-------\npanic happened. couldn't return order. r= ", r, "query.Error= ", result.Error, "\n----")
			tx.Rollback()
		}
	}()

	//update order status
	result = tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "return requested")
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't return order. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return result.Error
	}

	//commit transaction
	tx.Commit()

	return nil
}

//MarkOrderAsReturned
func (repo *OrderRepo) MarkOrderAsReturned(orderID uint) error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			fmt.Println("-------\npanic happened. couldn't return order. r= ", r, "query.Error= ", result.Error, "\n----")
			tx.Rollback()
		}
	}()

	//update order status
	result = tx.Model(&entities.Order{}).Where("id = ?", orderID).Update("status", "returned")
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't return order. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return result.Error
	}

	//get order items
	var orderItems []entities.OrderItem
	query := tx.
		Where("order_id = ?", orderID).
		Find(&orderItems)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		tx.Rollback()
		return query.Error
	}

	//update stock
	for _, item := range orderItems {
		result := tx.Model(&entities.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", item.Quantity))
		if result.Error != nil {
			fmt.Println("-------\nquery error happened. couldn't update stock. query.Error= ", result.Error, "\n----")
			tx.Rollback()
			return result.Error
		}
	}

	//get order's final amount and userID (to update wallet)
	var order entities.Order
	query = tx.
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		tx.Rollback()
		return query.Error
	}
	
	//update wallet
	result = tx.Model(&entities.User{}).Where("id = ?", order.UserID).Update("wallet_balance", gorm.Expr("wallet_balance + ?", order.FinalAmount))
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't update wallet. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return result.Error
	}

	//commit transaction
	tx.Commit()

	return nil
}

//MarkOrderAsDelivered
func (repo *OrderRepo) MarkOrderAsDelivered(orderID uint) error {
	//start transaction
	tx := repo.DB.Begin()
	var result *gorm.DB

	//defer rollback if error happened
	defer func() {
		if r := recover(); r != nil || result.Error != nil {
			fmt.Println("-------\npanic happened. couldn't return order. r= ", r, "query.Error= ", result.Error, "\n----")
			tx.Rollback()
		}
	}()

	//update order status and delivered_date
	result = tx.Model(&entities.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{"status": "delivered", "delivered_date": gorm.Expr("CURRENT_TIMESTAMP")})
	if result.Error != nil {
		fmt.Println("-------\nquery error happened. couldn't return order. query.Error= ", result.Error, "\n----")
		tx.Rollback()
		return result.Error
	}

	//commit transaction
	tx.Commit()

	return nil
}

//GetAllOrders
func (repo *OrderRepo) GetAllOrders() (*[]entities.Order, error) {
	var orders []entities.Order
	query := repo.DB.
		Find(&orders)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &orders, nil
}

func (repo *OrderRepo) 	GetOrderSummaryByID(orderID uint) (*entities.Order, error){
	var order entities.Order
	query := repo.DB.
		Preload("FkAddress").
		Where("id = ?", orderID).
		Find(&order)

	if query.Error != nil {
		fmt.Println("-------\nquery error happened. query.Error= ", query.Error, "\n----")
		return nil, query.Error
	}

	return &order, nil
}