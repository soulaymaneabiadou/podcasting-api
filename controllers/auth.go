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
		utils.ErrorsResponse(c, err)
		return
	}

	_, err := ac.as.Signup(data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "thanks for signing up, please check your inbox for a verification email")
}

func (ac *AuthController) Join(c *gin.Context) {
	var data types.SignupInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	_, err := ac.as.Join(data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "thanks for joining, please check your inbox for a verification email")
}

func (ac *AuthController) SignIn(c *gin.Context) {
	var data types.SigninInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	user, err := ac.as.Signin(data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SendTokensResponse(c, user)
}

func (ac *AuthController) CurrentUser(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := ac.as.GetUser(id)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

func (ac *AuthController) UpdateDetails(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	var data types.UpdateDetailsInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	user, err := ac.as.UpdateDetails(id, data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

func (ac *AuthController) UpdatePassword(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	var data types.UpdatePasswordInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	_, err := ac.as.UpdatePassword(id, data)
	if err != nil {
		utils.ErrorsResponse(c, errors.New("invalid credentials"))
		return
	}

	utils.ClearAuthCookies(c)
	utils.SuccessResponse(c, "the password has been updated")
}

func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var data types.ForgotPasswordInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	token, err := ac.as.ForgotPassword(data)
	if err != nil {
		utils.ErrorsResponse(c, errors.New("the user does not exist"))
		return
	}

	utils.SuccessResponse(c, token)
}

func (ac *AuthController) ResetPassword(c *gin.Context) {
	var data types.ResetPasswordInput

	token := c.Param("resettoken")

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	_, err := ac.as.ResetPassword(token, data.Password)
	if err != nil {
		utils.ErrorsResponse(c, errors.New("the reset token is invalid or has expired"))
		return
	}

	utils.SuccessResponse(c, "the password has been reseted")
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
		utils.ErrorsResponse(c, errors.New("the verification token is invalid or has expired"))
		return
	}

	utils.SuccessResponse(c, "the account has been verified, you may login now.")
}
