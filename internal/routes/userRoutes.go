package routes

import (
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	ordermanagementHandlers "MyShoo/internal/handlers/orderManagementHandlers"
	"MyShoo/internal/handlers/paymentHandlers"
	productHandlers "MyShoo/internal/handlers/productManagementHandlers"
	"MyShoo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup,
	user *accHandlers.UserHandler,
	category *productHandlers.CategoryHandler,
	brand *productHandlers.BrandsHandler,
	model *productHandlers.ModelHandler,
	product *productHandlers.ProductHandler,
	cart *ordermanagementHandlers.CartHandler,
	wishList *ordermanagementHandlers.WishListHandler,
	order *ordermanagementHandlers.OrderHandler,
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
			authUser.GET("/", user.GetHome)
			authUser.GET("/home", user.GetHome)

			//cart related________________________________________________
			authUser.GET("/cart", cart.GetCart)
			authUser.PUT("/cart", cart.AddToCart)
			authUser.DELETE("/cart", cart.DeleteFromCart)
			//clear entire cart
			authUser.DELETE("/clearcart", cart.ClearCart)

			//order_____________________________________________________
			authUser.GET("/myorders", order.GetOrdersOfUser)
			authUser.POST("/makeorder", order.MakeOrder)
			authUser.PATCH("/cancelorder", order.CancelMyOrder)
			authUser.PATCH("/returnorder", order.ReturnMyOrder)

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
		engine.POST("/payment", payment.ProceedToPayViaRazorPay)
		engine.POST("/payment/verify", payment.VerifyPayment)

	}
}
