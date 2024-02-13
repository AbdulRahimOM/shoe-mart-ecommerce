package di

import (
	myhttp "MyShoo/internal"
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	orderHandlers "MyShoo/internal/handlers/orderHandlers"
	"MyShoo/internal/handlers/paymentHandlers"
	productHandlers "MyShoo/internal/handlers/productHandlers"
	reporthandlers "MyShoo/internal/handlers/reportHandlers"
	infra "MyShoo/internal/infrastructure"
	accountrepo "MyShoo/internal/repository/accounts"
	orderrepo "MyShoo/internal/repository/order"

	// paymentrepo "MyShoo/internal/repository/payment_repo"
	productRepository "MyShoo/internal/repository/productManagement"
	reportsrepo "MyShoo/internal/repository/reports"
	accountsusecase "MyShoo/internal/usecase/accountsusecases"
	orderusecase "MyShoo/internal/usecase/orderUseCase"
	paymentusecase "MyShoo/internal/usecase/paymentUsecase"
	productusecase "MyShoo/internal/usecase/productManagementUsecases"
	reportsusecases "MyShoo/internal/usecase/reportsUsecases"
	"fmt"
)

func InitializeAndStartAPI() {
	fmt.Println("Handler ::: InitializeAndStartAPI in di package")

	userRepository := accountrepo.NewUserRepository(infra.DB)
	userUseCase := accountsusecase.NewUserUseCase(userRepository)
	userHandler := accHandlers.NewUserHandler(userUseCase)

	adminRepository := accountrepo.NewAdminRepository(infra.DB)
	adminUseCase := accountsusecase.NewAdminUseCase(adminRepository)
	adminHandler := accHandlers.NewAdminHandler(adminUseCase)

	sellerRepository := accountrepo.NewSellerRepository(infra.DB)
	sellerUseCase := accountsusecase.NewSellerUseCase(sellerRepository)
	sellerHandler := accHandlers.NewSellerHandler(sellerUseCase)

	categoryRepository := productRepository.NewCategoryRepository(infra.DB)
	categoryUseCase := productusecase.NewCategoryUseCase(categoryRepository)
	categoryHandler := productHandlers.NewCategoryHandler(categoryUseCase)

	brandRepository := productRepository.NewBrandRepository(infra.DB)
	brandUseCase := productusecase.NewBrandUseCase(brandRepository)
	brandHandler := productHandlers.NewBrandHandler(brandUseCase)

	modelRepository := productRepository.NewModelRepository(infra.DB)
	modelUseCase := productusecase.NewModelUseCase(modelRepository)
	modelHandler := productHandlers.NewModelHandler(modelUseCase)

	productRepository := productRepository.NewProductRepository(infra.DB, infra.CloudinaryClient)
	productUseCase := productusecase.NewProductUseCase(productRepository, modelRepository)
	productHandler := productHandlers.NewProductHandler(productUseCase)

	//order management related_____________________________________
	//cart
	cartRepository := orderrepo.NewCartRepository(infra.DB)
	cartUseCase := orderusecase.NewCartUseCase(cartRepository)
	cartHandler := orderHandlers.NewCartHandler(cartUseCase)

	//wishList
	wishListRepository := orderrepo.NewWishListRepository(infra.DB)
	wishListUseCase := orderusecase.NewWishListUseCase(wishListRepository, productRepository)
	wishListHandler := orderHandlers.NewWishListHandler(wishListUseCase)

	//order
	orderRepository := orderrepo.NewOrderRepository(infra.DB, infra.CloudinaryClient)
	orderUseCase := orderusecase.NewOrderUseCase(userRepository, orderRepository, cartRepository, productRepository)
	orderHandler := orderHandlers.NewOrderHandler(orderUseCase)

	//reports
	reportsRepository := reportsrepo.NewReportRepository(infra.DB, infra.CloudinaryClient)
	reportsUseCase := reportsusecases.NewReportsUseCase(reportsRepository, orderRepository)
	reportsHandler := reporthandlers.NewReportsHandler(reportsUseCase)

	//payment
	// paymentRepository := paymentrepo.NewPaymentRepository(infra.DB)	//not needed
	paymentUseCase := paymentusecase.NewPaymentUseCase(orderRepository)
	paymentHandler := paymentHandlers.NewPaymentHandler(paymentUseCase)

	serveHttp := myhttp.NewServerHTTP(
		userHandler, adminHandler, sellerHandler,
		categoryHandler, brandHandler, modelHandler, productHandler,
		cartHandler, wishListHandler, orderHandler,
		reportsHandler, paymentHandler,
	)
	serveHttp.Start()
}
