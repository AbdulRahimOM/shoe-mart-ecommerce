package routes

import (
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	ordermanagementHandlers "MyShoo/internal/handlers/orderManagementHandlers"
	productHandlers "MyShoo/internal/handlers/productManagementHandlers"
	"MyShoo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SellerRoutes(engine *gin.RouterGroup,
	seller *accHandlers.SellerHandler,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *ordermanagementHandlers.CartHandler,
	wishList *ordermanagementHandlers.WishListHandler,
) {
	engine.Use(middleware.ClearCache)
	loggedOutGroup := engine.Group("/")
	loggedOutGroup.Use(middleware.NotLoggedOutCheck)
	{
		loggedOutGroup.GET("/login", middleware.NotLoggedOutCheck, seller.GetLogin)
		loggedOutGroup.POST("/signup", middleware.NotLoggedOutCheck, seller.PostSignUp)
		loggedOutGroup.POST("/login", middleware.NotLoggedOutCheck, seller.PostLogIn)
	}

	authSeller := engine.Group("/")
	authSeller.Use(middleware.SellerAuth,middleware.VerifySellerStatus)
	//product management//////////////=============================================================//////
	//categories
	{
		authSeller.GET("/categories", category.GetCategories)
		authSeller.POST("/addcategory", category.AddCategory)
		authSeller.PATCH("/editcategory", category.EditCategory)

		//brands
		authSeller.GET("/brands", brand.GetBrands)
		authSeller.POST("/addbrand", brand.AddBrand)
		authSeller.PATCH("/editbrand", brand.EditBrand)

		//models
		authSeller.GET("/models", model.GetModelsByBrandsAndCategories)
		authSeller.POST("/addmodel", model.AddModel)
		authSeller.PATCH("/editmodel", model.EditModel)

		//Products
		//colour variants
		authSeller.POST("/addcolourvariant", product.AddColourVariant)
		authSeller.PATCH("/editcolourvariant", product.EditColourVariant)

		//dimensional variants
		authSeller.POST("/adddimensionalvariant", product.AddDimensionalVariant)

		//stock management
		//get stocks => get products is enough
		authSeller.GET("/products", product.GetProducts)
		authSeller.POST("/addstock", product.AddStock)    //add to stock (need not know existing stock)
		authSeller.PATCH("/editstock", product.EditStock) //need update: Add handler
	}
}
