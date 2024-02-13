package routes

import (
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	orderhandler "MyShoo/internal/handlers/orderHandlers"
	productHandlers "MyShoo/internal/handlers/productHandlers"
	reporthandlers "MyShoo/internal/handlers/reportHandlers"
	"MyShoo/internal/middleware"

	"github.com/gin-gonic/gin"
)

// AdminRoutes
func AdminRoutes(engine *gin.RouterGroup,
	admin *accHandlers.AdminHandler,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *orderhandler.CartHandler,
	wishList *orderhandler.WishListHandler,
	order *orderhandler.OrderHandler,
	reports *reporthandlers.ReportsHandler,
) {

	//viewing
	{
		engine.GET("/categories", category.GetCategories)
		engine.GET("/brands", brand.GetBrands)
		engine.GET("/models", model.GetModelsByBrandsAndCategories)
		engine.GET("/products", product.GetProducts)
	}

	loggedOutGroup := engine.Group("/")
	loggedOutGroup.Use(middleware.ClearCache, middleware.NotLoggedOutCheck)
	{
		loggedOutGroup.GET("/login", admin.GetAdminLogin)
		loggedOutGroup.POST("/login", admin.PostLogIn)
	}

	authAdmin := engine.Group("/")
	authAdmin.Use(middleware.ClearCache, middleware.AdminAuth)
	{
		//system related
		authAdmin.GET("/system/restart-Configuration", admin.RestartConfig)

		//user related
		authAdmin.GET("/userslist", admin.GetUsersList)
		authAdmin.POST("/blockuser", admin.BlockUser)
		authAdmin.POST("/unblockuser", admin.UnblockUser)

		//seller related
		authAdmin.GET("/sellerslist", admin.GetSellersList)
		authAdmin.POST("/blockseller", admin.BlockSeller)
		authAdmin.POST("/unblockseller", admin.UnblockSeller)
		authAdmin.PATCH("/verify-seller", admin.VerifySeller)

		// category management
		authAdmin.POST("/addcategory", category.AddCategory) //(add category access only to admin)
		authAdmin.PATCH("/editcategory", category.EditCategory)

		//product edit (exclusive "edit" access, but no "add" access to admin)
		authAdmin.PATCH("/editbrand", brand.EditBrand)
		authAdmin.PATCH("/editmodel", model.EditModel)
		authAdmin.PATCH("/editcolourvariant", product.EditColourVariant)

		//order management
		authAdmin.GET("/orders", order.GetOrders)                          //get all orders
		authAdmin.PATCH("/markdelivery", order.MarkOrderAsDelivered)       //mark order as delivered
		authAdmin.PATCH("/markorderasreturned", order.MarkOrderAsReturned) //mark order as returned
		authAdmin.PATCH("/cancelorder", order.CancelOrderByAdmin)          //cancel order

		//reports
		authAdmin.GET("/dashboarddata/:range", reports.GetDashBoardData) // dashBoardData
		authAdmin.GET("/salesreport/:range", reports.ExportSalesReport)  //salesreport

		//coupon
		authAdmin.GET("/coupons", order.GetCoupons)
		authAdmin.POST("/newcoupon", order.NewCouponHandler)
		authAdmin.PATCH("/blockcoupon", order.BlockCouponHandler)
		authAdmin.PATCH("/unblockcoupon", order.UnblockCouponHandler)

		//top performers
		authAdmin.GET("/top-products", reports.TopProductsHandler)
		authAdmin.GET("/top-sellers/", reports.TopSellersHandler)
		authAdmin.GET("/top-brands/", reports.TopBrandsHandler)
		authAdmin.GET("/top-models/", reports.TopModelsHandler)

	}
}
