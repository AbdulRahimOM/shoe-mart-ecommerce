package routes

import (
	orderhandler "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/orderHandlers"
	productHandlers "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/productHandlers"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(engine *gin.RouterGroup,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *orderhandler.CartHandler,
	wishList *orderhandler.WishListHandler,
) {
	engine.GET("/categories", category.GetCategories)
	engine.GET("/brands", brand.GetBrands)
	engine.GET("/models", model.GetModelsByBrandsAndCategories)
	engine.GET("/products", product.GetProducts)
	engine.GET("/colourvariants/:modelID", product.GetColourVariantsUnderModel)
}
