package routes

import (
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	ordermanagementHandlers "MyShoo/internal/handlers/orderManagementHandlers"
	productHandlers "MyShoo/internal/handlers/productManagementHandlers"
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
	cart *ordermanagementHandlers.CartHandler,
	wishList *ordermanagementHandlers.WishListHandler,
	order *ordermanagementHandlers.OrderHandler,
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

		//get all orders
		authAdmin.GET("/orders", order.GetOrders)
		//mark order as delivered
		authAdmin.PATCH("/markdelivery", order.MarkOrderAsDelivered)
		//mark order as returned
		authAdmin.PATCH("/markorderasreturned", order.MarkOrderAsReturned)
		//cancel order
		authAdmin.PATCH("/cancelorder", order.CancelOrderByAdmin)

		// dashBoardData
		authAdmin.GET("/dashboarddata/:range", reports.GetDashBoardData)

		//salesreport
		authAdmin.GET("/salesreport/:range", reports.ExportSalesReport)

		//coupon
		authAdmin.GET("/coupons", order.GetCoupons)
		authAdmin.POST("/newcoupon", order.NewCouponHandler)
		authAdmin.PATCH("/blockcoupon", order.BlockCouponHandler)
		authAdmin.PATCH("/unblockcoupon", order.UnblockCouponHandler)

	}
}
