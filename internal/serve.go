package myhttp

import (
	"MyShoo/internal/config"
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	orderHandlers "MyShoo/internal/handlers/orderHandlers"
	"MyShoo/internal/handlers/paymentHandlers"
	productHandlers "MyShoo/internal/handlers/productHandlers"
	reporthandlers "MyShoo/internal/handlers/reportHandlers"
	"MyShoo/internal/routes"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHttp struct {
	eng *gin.Engine
}

func NewServerHTTP(
	user *accHandlers.UserHandler,
	admin *accHandlers.AdminHandler,
	seller *accHandlers.SellerHandler,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *orderHandlers.CartHandler,
	wishList *orderHandlers.WishListHandler,
	order *orderHandlers.OrderHandler,
	reports *reporthandlers.ReportsHandler,
	payment *paymentHandlers.PaymentHandler,
) *ServerHttp {
	engine := gin.Default()
	// engine.LoadHTMLGlob("./internal/view/*.html")
	// engine.LoadHTMLGlob("./internal/templates/*.html")	//server showing error with this while running binary executable (don't know exact reason, busy) 
	engine.LoadHTMLFiles("./internal/templates/payment.html")

	routes.UserRoutes(engine.Group("/"), user, category, brand, model, product, cart, wishList, order, payment)
	routes.AdminRoutes(engine.Group("/admin"), admin, category, brand, model, product, cart, wishList, order, reports)

	routes.SellerRoutes(engine.Group("/seller"), seller, category, brand, model, product, cart, wishList)
	routes.PublicRoutes(engine.Group("/"), category, brand, model, product, cart, wishList)

	//add swagger
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &ServerHttp{eng: engine}

}
func (serveHttp *ServerHttp) Start() {
	serveHttp.eng.Run(":"+config.Port)
}
