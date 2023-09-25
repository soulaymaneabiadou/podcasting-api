package routes

import (
	"podcast/middleware"

	"github.com/gin-gonic/gin"
)

func authRoutes(r *gin.RouterGroup) {
	g := r.Group("/auth")
	ac := CreateAuthController()

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
