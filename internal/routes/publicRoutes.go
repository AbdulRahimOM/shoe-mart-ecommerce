package routes

import (
	ordermanagementHandlers "MyShoo/internal/handlers/orderManagementHandlers"
	productHandlers "MyShoo/internal/handlers/productManagementHandlers"
	"MyShoo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(engine *gin.RouterGroup,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *ordermanagementHandlers.CartHandler,
	wishList *ordermanagementHandlers.WishListHandler,
) {
	engine.Use(middleware.ClearCache)

	engine.GET("/categories", category.GetCategories)
	engine.GET("/brands", brand.GetBrands)
	engine.GET("/models", model.GetModelsByBrandsAndCategories)

	engine.GET("/products", product.GetProducts)

}
