package di

import (
	myhttp "MyShoo/internal"
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	orderManagementHandlers "MyShoo/internal/handlers/orderManagementHandlers"
	productHandlers "MyShoo/internal/handlers/productManagementHandlers"
	infra "MyShoo/internal/infrastructure"
	accRepository "MyShoo/internal/repository/accounts_Repo"
	ordermanagementrepo "MyShoo/internal/repository/orderManagement_Repo"
	productRepository "MyShoo/internal/repository/productManagement_Repo"
	accountsUsecase "MyShoo/internal/usecase/accountsUsecases"
	orderManageUseCase "MyShoo/internal/usecase/orderManageUseCase"
	prodManageUsecase "MyShoo/internal/usecase/productManagementUsecases"
	"fmt"
)

func InitializeAndStartAPI() {
	fmt.Println("Handler ::: InitializeAndStartAPI in di package")

	userRepository := accRepository.NewUserRepository(infra.DB)
	userUseCase := accountsUsecase.NewUserUseCase(userRepository)
	userHandler := accHandlers.NewUserHandler(userUseCase)

	adminRepository := accRepository.NewAdminRepository(infra.DB)
	adminUseCase := accountsUsecase.NewAdminUseCase(adminRepository)
	adminHandler := accHandlers.NewAdminHandler(adminUseCase)

	sellerRepository := accRepository.NewSellerRepository(infra.DB)
	sellerUseCase := accountsUsecase.NewSellerUseCase(sellerRepository)
	sellerHandler := accHandlers.NewSellerHandler(sellerUseCase)

	categoryRepository := productRepository.NewCategoryRepository(infra.DB)
	categoryUseCase := prodManageUsecase.NewCategoryUseCase(categoryRepository)
	categoryHandler := productHandlers.NewCategoryHandler(categoryUseCase)

	brandRepository := productRepository.NewBrandRepository(infra.DB)
	brandUseCase := prodManageUsecase.NewBrandUseCase(brandRepository)
	brandHandler := productHandlers.NewBrandHandler(brandUseCase)

	modelRepository := productRepository.NewModelRepository(infra.DB)
	modelUseCase := prodManageUsecase.NewModelUseCase(modelRepository)
	modelHandler := productHandlers.NewModelHandler(modelUseCase)

	productRepository := productRepository.NewProductRepository(infra.DB)
	productUseCase := prodManageUsecase.NewProductUseCase(productRepository)
	productHandler := productHandlers.NewProductHandler(productUseCase)

	//order management related_____________________________________
	//cart
	cartRepository := ordermanagementrepo.NewCartRepository(infra.DB)
	cartUseCase := orderManageUseCase.NewCartUseCase(cartRepository)
	cartHandler := orderManagementHandlers.NewCartHandler(cartUseCase)

	//wishList
	wishListRepository := ordermanagementrepo.NewWishListRepository(infra.DB)
	wishListUseCase := orderManageUseCase.NewWishListUseCase(wishListRepository, productRepository)
	wishListHandler := orderManagementHandlers.NewWishListHandler(wishListUseCase)

	//order
	orderRepository := ordermanagementrepo.NewOrderRepository(infra.DB)
	orderUseCase := orderManageUseCase.NewOrderUseCase(userRepository, orderRepository, cartRepository, productRepository)
	orderHandler := orderManagementHandlers.NewOrderHandler(orderUseCase)

	serveHttp := myhttp.NewServerHTTP(
		userHandler, adminHandler, sellerHandler,
		categoryHandler, brandHandler, modelHandler, productHandler,
		cartHandler, wishListHandler, orderHandler,
	)
	serveHttp.Start()
}
