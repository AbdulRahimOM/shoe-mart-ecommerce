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
		authAdmin.GET("/userslist", admin.GetUsersList)
		authAdmin.POST("/blockuser", admin.BlockUser)
		authAdmin.POST("/unblockuser", admin.UnblockUser)

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
	}
}
