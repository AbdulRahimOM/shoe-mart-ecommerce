package routes

import (
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	orderhandler "MyShoo/internal/handlers/orderHandlers"
	productHandlers "MyShoo/internal/handlers/productHandlers"
	"MyShoo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SellerRoutes(engine *gin.RouterGroup,
	seller *accHandlers.SellerHandler,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *orderhandler.CartHandler,
	wishList *orderhandler.WishListHandler,
) {
	engine.Use(middleware.ClearCache)

	{
		// viewing whole content from public perspective
		engine.GET("/categories", category.GetCategories)
		engine.GET("/brands", brand.GetBrands)
		engine.GET("/models", model.GetModelsByBrandsAndCategories)
		engine.GET("/products", product.GetProducts)
	}

	loggedOutGroup := engine.Group("/")
	loggedOutGroup.Use(middleware.NotLoggedOutCheck)
	{
		loggedOutGroup.GET("/login", middleware.NotLoggedOutCheck, seller.GetLogin)
		loggedOutGroup.POST("/signup", middleware.NotLoggedOutCheck, seller.PostSignUp)
		loggedOutGroup.POST("/login", middleware.NotLoggedOutCheck, seller.PostLogIn)
	}

	authSeller := engine.Group("/")
	authSeller.Use(middleware.SellerAuth, middleware.VerifySellerStatus)
	{
		//product management (Role: Adding)(Edit access is only for admin)
		authSeller.POST("/addbrand", brand.AddBrand)
		authSeller.POST("/addmodel", model.AddModel)
		authSeller.POST("/addcolourvariant", product.AddColourVariant)
		authSeller.POST("/adddimensionalvariant", product.AddDimensionalVariant)

		//stock management
		authSeller.POST("/addstock", product.AddStock)    //add to stock (need not know existing stock)
		authSeller.PATCH("/editstock", product.EditStock) //need update: Add handler
		//for get stocks => get products is enough
	}
}
