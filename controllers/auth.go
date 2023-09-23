package controllers

import (
	"errors"

	"podcast/database"
	"podcast/repositories"
	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	as *services.AuthService
}

// TODO: DI
func NewAuthController() *AuthController {
	repo := repositories.NewUsersRepository(database.DB)
	srv := services.NewAuthService(repo)

	return &AuthController{as: srv}
}

func (ac *AuthController) SignUp(c *gin.Context) {
	var data types.SignupInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	user, err := ac.as.Signup(data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SendTokensResponse(c, user)
}

func (ac *AuthController) SignIn(c *gin.Context) {
	var data types.SigninInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	user, err := ac.as.Signin(data)
	if err != nil {
		utils.ErrorsResponse(c, errors.New("invalid credentials"))
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
