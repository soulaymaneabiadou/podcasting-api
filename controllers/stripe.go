package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"podcast/services"
	"podcast/utils"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
)

type StripeController struct {
	ss *services.StripeService
	us *services.UsersService
}

func NewStripeController(ss *services.StripeService, us *services.UsersService) *StripeController {
	return &StripeController{ss: ss, us: us}
}

func (sc *StripeController) Connect(c *gin.Context) {
	var acct *stripe.Account

	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	if user.DetailsSubmitted && user.ChargesEnabled && user.TransfersEnabled {
		utils.ErrorsResponse(c, errors.New("user account already setted up"))
		return
	}

	if user.StripeAccountId == "" {
		acct, err = sc.ss.CreateAccount(user)
		if err != nil {
			utils.ErrorsResponse(c, err)
			return
		}

		user, err = sc.us.SetStripeAccountId(user, fmt.Sprint(acct.ID))
		if err != nil {
			utils.ErrorsResponse(c, err)
			return
		}
	} else {
		acct, err = sc.ss.GetAccount(user.StripeAccountId)
		if err != nil {
			utils.ErrorsResponse(c, err)
			return
		}
	}

	link, err := sc.ss.CreateAccountLink(acct)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, link.URL)
	c.JSON(http.StatusOK, link.URL)
}

func (sc *StripeController) Onboard(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	if user.DetailsSubmitted && user.ChargesEnabled && user.TransfersEnabled {
		utils.ErrorsResponse(c, errors.New("user account already setted up"))
		return
	}

	acct, err := sc.ss.GetAccount(user.StripeAccountId)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	link, err := sc.ss.CreateAccountLink(acct)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, link.URL)
	c.JSON(http.StatusOK, link.URL)
}

func (sc *StripeController) CustomerPortal(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	if user.StripeCustomerId == "" {
		utils.ErrorsResponse(c, errors.New("user is not a customer yet"))
		return
	}

	session, err := sc.ss.CreateCustomerPortalSession(user.StripeCustomerId)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, session.URL)
	c.JSON(http.StatusOK, session.URL)
}

func (sc *StripeController) ConnectAccount(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	if user.StripeAccountId == "" {
		utils.ErrorsResponse(c, errors.New("user is not a customer yet"))
		return
	}

	link, err := sc.ss.CreateConnectAccountLink(user.StripeAccountId)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, link.URL)
	c.JSON(http.StatusOK, link.URL)
}
