package controllers

import (
	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type SubscriptionsController struct {
	ss *services.SubscriptionsService
}

func NewSubscriptionsController(ss *services.SubscriptionsService) *SubscriptionsController {
	return &SubscriptionsController{ss: ss}
}

func (pc *SubscriptionsController) SubscribeToPodcast(c *gin.Context) {
	uid, _ := utils.GetCtxUser(c)

	var data types.SubscribeInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide a valid podcast id")
		return
	}

	url, err := pc.ss.Subscribe(uid, data.PodcastId)
	if err != nil {
		utils.ErrorResponse(c, err, "Unable to subscribe to this podcast, please try again later with valid information")
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, url)
	utils.SuccessResponse(c, url)
}
