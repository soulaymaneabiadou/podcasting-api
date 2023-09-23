package routes

import (
	"podcast/controllers"
	"podcast/middleware"

	"github.com/gin-gonic/gin"
)

func authRoutes(r *gin.RouterGroup) {
	g := r.Group("/auth")
	// TODO: DI
	ac := controllers.NewAuthController()

	g.POST("/signup", ac.SignUp)
	g.POST("/signin", ac.SignIn)
	g.GET("/signout", middleware.Authenticate(), ac.SignOut)
	g.GET("/refreshtoken", middleware.AuthenticateRefreshToken(), ac.RefreshToken)
	g.GET("/me", middleware.Authenticate(), ac.CurrentUser)
	g.PATCH("/updatedetails", middleware.Authenticate(), ac.UpdateDetails)
	g.PATCH("/updatepassword", middleware.Authenticate(), ac.UpdatePassword)
	g.POST("/forgotpassword", ac.ForgotPassword)
	g.PUT("/resetpassword/:resettoken", ac.ResetPassword)
}
