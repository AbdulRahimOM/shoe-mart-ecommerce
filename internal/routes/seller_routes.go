package routes

import (
	accHandlers "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/accountHandlers"
	orderhandler "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/orderHandlers"
	productHandlers "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/productHandlers"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/middleware"

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
	{
		loggedOutGroup.GET("/login", seller.GetLogin)
		loggedOutGroup.POST("/signup", seller.PostSignUp)
		loggedOutGroup.POST("/login", seller.PostLogIn)
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
