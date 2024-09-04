package di

import (
	myhttp "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal"
	accHandlers "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/accountHandlers"
	orderHandlers "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/orderHandlers"
	"github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/paymentHandlers"
	productHandlers "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/productHandlers"
	reporthandlers "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/handlers/reportHandlers"
	infra "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/infrastructure"
	accountrepo "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/accounts"
	orderrepo "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/order"

	// paymentrepo "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/payment_repo"
	"fmt"

	productRepository "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/productManagement"
	reportsrepo "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/repository/reports"
	accountsusecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/accountsUsecases"
	orderusecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/orderUseCase"
	paymentusecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/paymentUsecase"
	productusecase "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/productManagementUsecases"
	reportsusecases "github.com/AbdulRahimOM/shoe-mart-ecommerce/internal/usecase/reportsUsecases"
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
