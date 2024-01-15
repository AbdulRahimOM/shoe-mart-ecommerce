package routes

import (
	accHandlers "MyShoo/internal/handlers/accountHandlers"
	ordermanagementHandlers "MyShoo/internal/handlers/orderManagementHandlers"
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
) {
	engine.Use(middleware.ClearCache)
	engine.GET("/login", middleware.NotLoggedOutCheck, user.GetLogin)

	engine.POST("/signup", middleware.NotLoggedOutCheck, user.PostSignUp)
	engine.POST("/login", middleware.NotLoggedOutCheck, user.PostLogIn)

	engine.GET("/sendotp", middleware.UserAuth, middleware.UserAwaitingVerification, user.SendOtp)
	engine.POST("/verifyotp", middleware.UserAuth, middleware.UserAwaitingVerification, user.VerifyOtp)

	engine.POST("/resetpasswordsendotp", middleware.NotLoggedOutCheck, user.SendOtpForPWChange)
	engine.POST("/resetpasswordverifyotp", middleware.NotLoggedOutCheck, user.VerifyOtpForPWChange)
	engine.POST("/resetpassword", middleware.NotLoggedOutCheck, user.ResetPassword)
	
	engine.GET("/", middleware.UserAuth, middleware.VerifyUserStatus, user.GetHome)
	engine.GET("/home", middleware.UserAuth, middleware.VerifyUserStatus, user.GetHome)

	//order management related_____________________________________
	
	engine.GET("/colourvariants/:modelID", middleware.UserAuth, product.GetColourVariantsUnderModel)

	//cart related________________________________________________
	engine.GET("/cart", middleware.UserAuth, middleware.VerifyUserStatus, cart.GetCart)
	engine.PUT("/cart", middleware.UserAuth, middleware.VerifyUserStatus, cart.AddToCart)
	engine.DELETE("/cart", middleware.UserAuth, middleware.VerifyUserStatus, cart.DeleteFromCart)
	//clear entire cart
	engine.DELETE("/clearcart", middleware.UserAuth, middleware.VerifyUserStatus, cart.ClearCart)

	//order_____________________________________________________
	//get orders of user
	engine.GET("/myorders", middleware.UserAuth, middleware.VerifyUserStatus, order.GetOrdersOfUser)
	//makeorder
	engine.POST("/makeorder", middleware.UserAuth, middleware.VerifyUserStatus, order.MakeOrder)
	//cancel order
	engine.PATCH("/cancelorder", middleware.UserAuth, middleware.VerifyUserStatus, order.CancelMyOrder)

	//get user addresses
	engine.GET("/addresses", middleware.UserAuth, middleware.VerifyUserStatus, user.GetUserAddresses)
	//add new address
	engine.POST("/addaddress", middleware.UserAuth, middleware.VerifyUserStatus, user.AddUserAddress)
	//edit address
	engine.PATCH("/editaddress", middleware.UserAuth, middleware.VerifyUserStatus, user.EditUserAddress)
	//delete address
	engine.DELETE("/deleteaddress", middleware.UserAuth, middleware.VerifyUserStatus, user.DeleteUserAddress)

	//user profile related
	engine.GET("/profile", middleware.UserAuth, middleware.VerifyUserStatus, user.GetProfile)
	engine.PATCH("/editprofile", middleware.UserAuth, middleware.VerifyUserStatus, user.EditProfile)

	//wishlist related____________________________________________
	engine.GET("/mywishlists", middleware.UserAuth, middleware.VerifyUserStatus, wishList.GetAllWishLists)
	engine.GET("/wishlist", middleware.UserAuth, middleware.VerifyUserStatus, wishList.GetWishListByID)
	engine.POST("/createwishlist", middleware.UserAuth, middleware.VerifyUserStatus, wishList.CreateWishList)
	engine.POST("/addtowishlist", middleware.UserAuth, middleware.VerifyUserStatus, wishList.AddToWishList)
	engine.DELETE("/removefromwishlist", middleware.UserAuth, middleware.VerifyUserStatus, wishList.RemoveFromWishList)

}
