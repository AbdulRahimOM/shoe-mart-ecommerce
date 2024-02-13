package routes

import (
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	orderhandler "MyShoo/internal/handlers/orderHandlers"
	"MyShoo/internal/handlers/paymentHandlers"
	productHandlers "MyShoo/internal/handlers/productHandlers"
	"MyShoo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup,
	user *accHandlers.UserHandler,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *orderhandler.CartHandler,
	wishList *orderhandler.WishListHandler,
	order *orderhandler.OrderHandler,
	payment *paymentHandlers.PaymentHandler,
) {
	engine.Use(middleware.ClearCache)
	{
		loggedOutGroup := engine.Group("/")
		loggedOutGroup.Use(middleware.NotLoggedOutCheck)
		{
			loggedOutGroup.GET("/login", user.GetLogin)

			loggedOutGroup.POST("/signup", user.PostSignUp)
			loggedOutGroup.POST("/login", user.PostLogIn)

			loggedOutGroup.POST("/resetpasswordsendotp", user.SendOtpForPWChange)
			loggedOutGroup.POST("/resetpasswordverifyotp", user.VerifyOtpForPWChange)
			loggedOutGroup.POST("/resetpassword", user.ResetPassword)
		}

		signinUpGroup := engine.Group("/")
		signinUpGroup.Use(middleware.UserAuth, middleware.UserAwaitingVerification)
		{
			signinUpGroup.GET("/sendotp", user.SendOtp)
			signinUpGroup.POST("/verifyotp", user.VerifyOtp)
		}

		authUser := engine.Group("/")
		authUser.Use(middleware.UserAuth, middleware.VerifyUserStatus)
		{

			//cart related________________________________________________
			authUser.GET("/cart", cart.GetCart)
			authUser.PUT("/cart", cart.AddToCart)
			authUser.DELETE("/cart", cart.DeleteFromCart)
			//clear entire cart
			authUser.DELETE("/clearcart", cart.ClearCart)

			// checkOut
			authUser.GET("/checkout/selectaddress", order.GetAddressForCheckout)
			authUser.POST("/checkout/setaddr-selectcoupon", order.SetAddressGetCoupons)
			authUser.POST("/checkout/setcoupon-getpaymentmethods", order.SetCouponGetPaymentMethods)

			//order_____________________________________________________
			authUser.POST("/makeorder", order.MakeOrder)
			authUser.GET("/myorders", order.GetOrdersOfUser)
			authUser.PATCH("/cancelorder", order.CancelMyOrder)
			authUser.PATCH("/returnorder", order.ReturnMyOrder)
			authUser.GET("/order-invoice", order.GetInvoiceOfOrder)

			//user address related_______________________________________
			authUser.GET("/addresses", user.GetUserAddresses)
			authUser.POST("/addaddress", user.AddUserAddress)
			authUser.PATCH("/editaddress", user.EditUserAddress)
			authUser.DELETE("/deleteaddress", user.DeleteUserAddress)

			//user profile related
			authUser.GET("/profile", user.GetProfile)
			authUser.PATCH("/editprofile", user.EditProfile)

			//wishlist related____________________________________________
			authUser.GET("/mywishlists", wishList.GetAllWishLists)
			authUser.GET("/wishlist", wishList.GetWishListByID)
			authUser.POST("/createwishlist", wishList.CreateWishList)
			authUser.POST("/addtowishlist", wishList.AddToWishList)
			authUser.DELETE("/removefromwishlist", wishList.RemoveFromWishList)
		}

		//payment related_____________________________________________
		engine.POST("/payment", middleware.UserAuth, middleware.VerifyUserStatus, payment.ProceedToPayViaRazorPay)
		engine.POST("/payment/verify", payment.VerifyPayment)
		engine.POST("/retrypayment", middleware.UserAuth, middleware.VerifyUserStatus, payment.RetryPayment)
		// engine.GET("/order-invoice", order.GetInvoiceOfOrder)

	}
}
