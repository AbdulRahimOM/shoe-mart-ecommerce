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
	engine.GET("/login", middleware.NotLoggedOutCheck, seller.GetLogin)

	engine.GET("/", middleware.SellerAuth, middleware.VerifySellerStatus, seller.GetHome)
	engine.GET("/home", middleware.SellerAuth, middleware.VerifySellerStatus, seller.GetHome)

	engine.POST("/signup", middleware.NotLoggedOutCheck, seller.PostSignUp)
	engine.POST("/login", seller.PostLogIn)

	//product management//////////////=============================================================//////
	//categories
	engine.GET("/categories", middleware.SellerAuth, category.GetCategories)
	engine.POST("/addcategory", middleware.SellerAuth, category.AddCategory)
	engine.PATCH("/editcategory", middleware.SellerAuth, category.EditCategory)

	//brands
	engine.GET("/brands",  middleware.SellerAuth, brand.GetBrands)
	engine.POST("/addbrand", middleware.SellerAuth, brand.AddBrand)
	engine.PATCH("/editbrand", middleware.SellerAuth, brand.EditBrand)

	//models
	engine.GET("/models", middleware.SellerAuth,  model.GetModelsByBrandsAndCategories)
	engine.POST("/addmodel", middleware.SellerAuth, model.AddModel)
	engine.PATCH("/editmodel", middleware.SellerAuth, model.EditModel)

	//Products
	//colour variants
	engine.POST("/addcolourvariant", middleware.SellerAuth, product.AddColourVariant)
	engine.PATCH("/editcolourvariant", middleware.SellerAuth, product.EditColourVariant)

	//dimensional variants
	engine.POST("/adddimensionalvariant", middleware.SellerAuth, product.AddDimensionalVariant)

	//stock management
		//get stocks => get products is enough
	engine.GET("/products", middleware.SellerAuth,product.GetProducts)
	engine.POST("/addstock", middleware.SellerAuth, product.AddStock)	//add to stock (need not know existing stock)
	engine.PATCH("/editstock", middleware.SellerAuth, product.EditStock) //need update: Add handler
}
