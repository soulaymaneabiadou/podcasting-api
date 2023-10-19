package controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"podcast/gateways/upload"
	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type EpisodesController struct {
	es *services.EpisodesService
	us *services.UsersService
	fh upload.FileHandler
}

func NewEpisodesController(
	es *services.EpisodesService,
	us *services.UsersService,
	fh upload.FileHandler,
) *EpisodesController {
	return &EpisodesController{es: es, us: us, fh: fh}
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
		utils.ErrorResponse(c, errors.New("missing subscription"), "You haven't subscribed to this episode's podcast yet.")
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

	media, err := utils.HandleFileValidation(c, "media", []string{"mp3"}, 10)
	if err != nil {
		utils.ErrorResponse(c, err, "Please include an mp3 media that does not exceed 10MB")
		return
	}

	mediaLink, err := pc.fh.Upload(media)
	if err != nil {
		utils.ErrorResponse(c, err, "An error occured while uploading the provided media file, please try again later")
		return
	}

	data = types.CreateEpisodeInput{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		MediaLink:   mediaLink,
		Visibility:  c.PostForm("visibility"),
		Tags:        strings.Split(c.PostForm("tags"), ", "),
	}

	if data.Visibility != "draft" {
		data.PublishedAt = time.Now()
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

	creatorId, err := strconv.ParseUint(id, 0, 32)
	if err != nil || !utils.IsOwner(c, uint(creatorId)) {
		utils.NotFoundResponse(c)
		return
	}

	var data types.UpdateEpisodeInput
	var mediaLink string

	media, err := utils.HandleFileValidation(c, "media", []string{"mp3"}, 10)
	if err != nil {
		log.Println("no media file was uploaded")
	}

	if media != nil {
		// TODO: delete old media if new one got uploaded
		mediaLink, err = pc.fh.Upload(media)
		if err != nil {
			utils.ErrorResponse(c, err, "An error occured while uploading the provided media file, please try again later")
			return
		}
	}

	data = types.UpdateEpisodeInput{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		MediaLink:   mediaLink,
		Visibility:  c.PostForm("visibility"),
		Tags:        strings.Split(c.PostForm("tags"), ", "),
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
