package controllers

import (
	"errors"
	"log"
	"strconv"

	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type PodcastsController struct {
	ps *services.PodcastsService
}

func NewPodcastsController(ps *services.PodcastsService) *PodcastsController {
	return &PodcastsController{ps: ps}
}

func (pc *PodcastsController) GetPodcasts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	count, podcasts, err := pc.ps.GetPodcasts(types.Paginator{Limit: limit, Page: page})
	pagination := utils.PaginationInput{Page: page, Limit: limit, Count: count}

	if err != nil {
		log.Println(err.Error())
		utils.ErrorsResponse(c, errors.New("an error occured while getting podcasts, please try again later"))
		return
	}

	utils.PaginatedResponse(c, podcasts, pagination)
}

func (pc *PodcastsController) GetPodcast(c *gin.Context) {
	id := c.Param("pid")

	podcast, err := pc.ps.GetPodcastById(id)
	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (pc *PodcastsController) CreatePodcast(c *gin.Context) {
	var data types.CreatePodcastInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	u, _ := c.Get("user")
	data.CreatorId = u.(utils.JwtPayload).ID

	podcast, err := pc.ps.CreatePodcast(data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (pc *PodcastsController) UpdatePodcast(c *gin.Context) {
	id := c.Param("pid")
	uid, _ := utils.GetCtxUser(c)

	var data types.UpdatePodcastInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	podcast, err := pc.ps.UpdatePodcast(uid, id, data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (pc *PodcastsController) DeletePodcast(c *gin.Context) {
	id := c.Param("pid")
	uid, _ := utils.GetCtxUser(c)

	res, err := pc.ps.DeletePodcast(uid, id)
	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	utils.SuccessResponse(c, res)
}

func (pc *PodcastsController) Subscribe(c *gin.Context) {
	pid := c.Param("pid")
	uid, _ := utils.GetCtxUser(c)

	url, err := pc.ps.Subscribe(uid, pid)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, url)
	utils.SuccessResponse(c, url)
}
