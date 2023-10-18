package controllers

import (
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

func (sc *StripeController) CreateAccount(c *gin.Context) {
	var acct *stripe.Account

	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorResponse(c, err, "The user does not exist")
		return
	}

	if user.DetailsSubmitted && user.ChargesEnabled && user.PayoutsEnabled {
		utils.ErrorResponse(c, err, "The user already has a completed account")
		return
	}

	if user.StripeAccountId == "" {
		acct, err = sc.ss.CreateAccount(user)
		if err != nil {
			utils.ErrorResponse(c, err, "Unable to create a stripe account, please try again later")
			return
		}

		user, err = sc.us.SetStripeAccountId(user, fmt.Sprint(acct.ID))
		if err != nil {
			utils.ErrorResponse(c, err, "Unable to set the user's stripe account id, please try again later")
			return
		}
	} else {
		acct, err = sc.ss.GetAccount(user.StripeAccountId)
		if err != nil {
			utils.ErrorResponse(c, err, "Unable to find a stripe account for the provided user")
			return
		}
	}

	link, err := sc.ss.CreateAccountOnboardingLink(acct)
	if err != nil {
		utils.ErrorResponse(c, err, "unable to create a stripe onboarding link, please check in later")
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, link.URL)
	c.JSON(http.StatusOK, link.URL)
}

func (sc *StripeController) OnboardAccount(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorResponse(c, err, "User not found")
		return
	}

	if user.DetailsSubmitted && user.ChargesEnabled && user.PayoutsEnabled {
		utils.ErrorResponse(c, err, "The user's account is already complete")
		return
	}

	acct, err := sc.ss.GetAccount(user.StripeAccountId)
	if err != nil {
		utils.ErrorResponse(c, err, "Unable to find a stripe account, please try to connect a new one")
		return
	}

	link, err := sc.ss.CreateAccountOnboardingLink(acct)
	if err != nil {
		utils.ErrorResponse(c, err, "Unable to create an onboarding link, please try again later")
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, link.URL)
	c.JSON(http.StatusOK, link.URL)
}

func (sc *StripeController) CreateAccountLogin(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorResponse(c, err, "User not found")
		return
	}

	if user.StripeAccountId == "" {
		utils.ErrorResponse(c, err, "No connect account was found, please start by creating one through the connect flow")
		return
	}

	link, err := sc.ss.CreateAccountLoginLink(user.StripeAccountId)
	if err != nil {
		utils.ErrorResponse(c, err, "Unable to create a login link for the connect account, please try again later")
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, link.URL)
	c.JSON(http.StatusOK, link.URL)
}

func (sc *StripeController) CreateCustomerPortal(c *gin.Context) {
	id, _ := utils.GetCtxUser(c)

	user, err := sc.us.GetUserById(id)
	if err != nil {
		utils.ErrorResponse(c, err, "User not found")
		return
	}

	if user.StripeCustomerId == "" {
		utils.ErrorResponse(c, err, "No customer account was found, please create a customer account first")
		return
	}

	session, err := sc.ss.CreateCustomerPortalSession(user.StripeCustomerId)
	if err != nil {
		utils.ErrorResponse(c, err, "Unable to create a customer portal, please retry later")
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, session.URL)
	c.JSON(http.StatusOK, session.URL)
}
