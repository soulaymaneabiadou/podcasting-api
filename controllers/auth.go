package controllers

import (
	"errors"

	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	as *services.AuthService
}

func NewAuthController(as *services.AuthService) *AuthController {
	return &AuthController{as: as}
}

func (ac *AuthController) SignUp(c *gin.Context) {
	var data types.SignupInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide valid information to sign up")
		return
	}

	_, err := ac.as.Signup(data)
	if err != nil {
		utils.ErrorResponse(c, err, err.Error())
		return
	}

	utils.MessageResponse(c, "Thanks for signing up, please check your inbox for a verification email")
}

func (ac *AuthController) Join(c *gin.Context) {
	var data types.SignupInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide valid information to create your account")
		return
	}

	_, err := ac.as.Join(data)
	if err != nil {
		utils.ErrorResponse(c, err, err.Error())
		return
	}

	utils.MessageResponse(c, "Thanks for joining, please check your inbox for a verification email")
}

func (ac *AuthController) SignIn(c *gin.Context) {
	var data types.SigninInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide a valid email and a password to sign into your account")
		return
	}

	data.IP = c.ClientIP()

	user, err := ac.as.Signin(data)
	if err != nil {
		utils.ErrorResponse(c, err, err.Error())
		return
	}

	utils.SendTokensResponse(c, user)
}

func (ac *AuthController) CurrentUser(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := ac.as.GetUser(id)
	if err != nil {
		utils.ErrorResponse(c, err, "Unauthorized access")
		return
	}

	utils.SuccessResponse(c, user)
}

func (ac *AuthController) UpdateDetails(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	var data types.UpdateDetailsInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide valid data to update your profile details")
		return
	}

	user, err := ac.as.UpdateDetails(id, data)
	if err != nil {
		utils.ErrorResponse(c, err, "Unable to update your profile details, please check the data you provided and retry later")
		return
	}

	utils.SuccessResponse(c, user)
}

func (ac *AuthController) UpdatePassword(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	var data types.UpdatePasswordInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide your current password and a new valid password")
		return
	}

	_, err := ac.as.UpdatePassword(id, data)
	if err != nil {
		utils.ErrorResponse(c, errors.New("invalid credentials"), "The provided current password is invalid")
		return
	}

	// utils.ClearAuthCookies(c)
	utils.MessageResponse(c, "your password has been updated successfully")
}

func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var data types.ForgotPasswordInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide a valid email to request a password reset")
		return
	}

	_, err := ac.as.ForgotPassword(data)
	if err != nil {
		utils.ErrorResponse(c, err, err.Error())
		return
	}

	utils.MessageResponse(c, "An email with instructions has been sent to your inbox.")
}

func (ac *AuthController) ResetPassword(c *gin.Context) {
	var data types.ResetPasswordInput

	token := c.Param("resettoken")

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide a valid password to reset access to your account")
		return
	}

	_, err := ac.as.ResetPassword(token, data.Password)
	if err != nil {
		utils.ErrorResponse(c, errors.New("the reset token is invalid or has expired"), "The link has expired, please request a new reset email")
		return
	}

	utils.MessageResponse(c, "Your password has been resetted, you might login to your account now")
}

func (ac *AuthController) SignOut(c *gin.Context) {
	// delete token if saved anywhere
	utils.ClearAuthCookies(c)
	utils.SuccessResponse(c, gin.H{})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	_, user := utils.GetCtxUser(c)

	// delete token if saved anywhere
	token, err := utils.CreateAccessToken("", utils.JwtPayload{
		ID: user.ID,
	})
	if err != nil {
		utils.UnauthorizedResponse(c)
		return
	}

	utils.SendAccessTokenResponse(c, token)
}

func (ac *AuthController) Verify(c *gin.Context) {
	token := c.Param("verificationtoken")

	_, err := ac.as.Verify(token)
	if err != nil {
		utils.ErrorResponse(c, err, "The provided verification token is invalid or has expired, please check your inbox for the correct one.")
		return
	}

	utils.MessageResponse(c, "Your account has been verified, you may login now.")
}
