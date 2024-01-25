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
	engine.Use(middleware.ClearCache)
	engine.GET("/login", middleware.NotLoggedOutCheck, admin.GetAdminLogin)

	engine.GET("/", middleware.AdminAuth, admin.GetAdminHome)
	engine.GET("/home", middleware.AdminAuth, admin.GetAdminHome)
	engine.POST("/login", admin.PostLogIn)

	engine.GET("/userslist", middleware.AdminAuth, admin.GetUsersList)
	engine.POST("/blockuser", middleware.AdminAuth, admin.BlockUser)
	engine.POST("/unblockuser", middleware.AdminAuth, admin.UnblockUser)

	engine.GET("/sellerslist", middleware.AdminAuth, admin.GetSellersList)
	engine.POST("/blockseller", middleware.AdminAuth, admin.BlockSeller)
	engine.POST("/unblockseller", middleware.AdminAuth, admin.UnblockSeller)

	//viewing
	engine.GET("/categories", middleware.AdminAuth, category.GetCategories)
	engine.GET("/brands", middleware.AdminAuth, brand.GetBrands)
	engine.GET("/models", middleware.AdminAuth, model.GetModelsByBrandsAndCategories)
	engine.GET("/products", middleware.AdminAuth, product.GetProducts)

	//get all orders
	engine.GET("/orders", middleware.AdminAuth, order.GetOrders)
	//cancel order
	engine.PATCH("/cancelorder", middleware.AdminAuth, order.CancelOrderByAdmin)

	// dashBoardData
	engine.GET("/dashboarddata/:range", middleware.AdminAuth, reports.GetDashBoardData)

	//salesreport
	// engine.GET("/exportsalesreport/:range", middleware.AdminAuth, reports.ExportSalesReport)
	engine.GET("/exportsalesreport/:range", reports.ExportSalesReport)

	//mark order as delivered
	engine.PATCH("/markdelivery", middleware.AdminAuth, order.MarkOrderAsDelivered)
	//mark order as returned
	engine.PATCH("/markorderasreturned", middleware.AdminAuth, order.MarkOrderAsReturned)
}
