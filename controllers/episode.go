package controllers

import (
	"fmt"
	"log"
	"strconv"

	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type EpisodesController struct {
	es *services.EpisodesService
	us *services.UsersService
}

func NewEpisodesController(es *services.EpisodesService, us *services.UsersService) *EpisodesController {
	return &EpisodesController{es: es, us: us}
}

func (pc *EpisodesController) GetPodcastEpisodes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pid := c.Param("pid")

	count, podcasts, err := pc.es.GetPodcastEpisodes(pid, types.Paginator{Limit: limit, Page: page})
	pagination := utils.PaginationInput{Page: page, Limit: limit, Count: count}

	if err != nil {
		log.Println(err.Error())
		utils.ErrorResponse(c, err, "An error has occured while getting this podcast's episodes, please try again later")
		return
	}

	utils.PaginatedResponse(c, podcasts, pagination)
}

func (pc *EpisodesController) GetPodcastEpisodesBySlug(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pslug := c.Param("pslug")

	count, podcasts, err := pc.es.GetPodcastEpisodesBySlug(pslug, types.Paginator{Limit: limit, Page: page})
	pagination := utils.PaginationInput{Page: page, Limit: limit, Count: count}

	if err != nil {
		log.Println(err.Error())
		utils.ErrorResponse(c, err, "An error has occured while getting this podcast's episodes, please try again later")
		return
	}

	utils.PaginatedResponse(c, podcasts, pagination)
}

func (pc *EpisodesController) GetPodcastEpisode(c *gin.Context) {
	id := c.Param("eid")

	episode, err := pc.es.GetPodcastEpisodeById(id)
	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	// episode is public, no need to check for a subscription
	if episode.Visibility == "public" {
		utils.SuccessResponse(c, episode)
		return
	}

	uid, _ := utils.GetCtxUser(c)
	subscribed, err := pc.us.IsUserSubscribedToPodcast(uid, fmt.Sprint(episode.PodcastId))
	if err != nil {
		utils.ErrorResponse(c, err, "An error has occured while checking for eligibility, please try again later")
		return
	}

	if !subscribed {
		utils.ErrorResponse(c, err, "you haven't subscribed to this episode's podcast yet.")
		return
	}

	utils.SuccessResponse(c, episode)
}

func (pc *EpisodesController) GetPodcastEpisodeBySlug(c *gin.Context) {
	slug := c.Param("eslug")

	episode, err := pc.es.GetPodcastEpisodeBySlug(slug)
	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	utils.SuccessResponse(c, episode)
}

func (pc *EpisodesController) CreatePodcastEpisode(c *gin.Context) {
	var data types.CreateEpisodeInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide valid data to create this episode")
		return
	}

	user, _ := c.Get("user")
	data.CreatorId = user.(utils.JwtPayload).ID

	pid, _ := strconv.ParseUint(c.Param("pid"), 0, 64)
	data.PodcastId = uint(pid)

	episode, err := pc.es.CreatePodcastEpisode(data)
	if err != nil {
		utils.ErrorResponse(c, err, "Please provide valid information to create this episode")
		return
	}

	utils.SuccessResponse(c, episode)
}

func (pc *EpisodesController) UpdatePodcastEpisode(c *gin.Context) {
	id := c.Param("eid")
	uid, _ := utils.GetCtxUser(c)

	var data types.UpdateEpisodeInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Please provide valid data to update this episode")
		return
	}

	episode, err := pc.es.UpdatePodcastEpisode(uid, id, data)
	if err != nil {
		utils.ErrorResponse(c, err, "Please provide valid information to update this episode")
		return
	}

	utils.SuccessResponse(c, episode)
}

func (pc *EpisodesController) PublishPodcastEpisode(c *gin.Context) {
	id := c.Param("eid")
	uid, _ := utils.GetCtxUser(c)

	var data types.PublishEpisodeInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorResponse(c, err, "Unable to publish this episode, please check the provided data")
		return
	}

	episode, err := pc.es.PublishPodcastEpisode(uid, id, data)
	if err != nil {
		utils.ErrorResponse(c, err, "Please provide valid information to be able to publish this episode")
		return
	}

	utils.SuccessResponse(c, episode)
}

func (pc *EpisodesController) DeletePodcastEpisode(c *gin.Context) {
	id := c.Param("eid")
	uid, _ := utils.GetCtxUser(c)

	res, err := pc.es.DeletePodcastEpisode(uid, id)

	if err != nil && res == true {
		utils.ErrorResponse(c, err, "You cannot a non draft episode, please unlist it and try again")
		return
	}

	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	utils.SuccessResponse(c, res)
}
