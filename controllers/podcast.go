package controllers

import (
	"errors"
	"log"
	"strconv"

	"podcast/database"
	"podcast/repositories"
	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type PodcastsController struct {
	as *services.PodcastsService
}

// TODO: DI
func NewPodcastsController() *PodcastsController {
	repo := repositories.NewPodcastsRepository(database.DB)
	srv := services.NewPodcastsService(repo)

	return &PodcastsController{as: srv}
}

func (pc *PodcastsController) GetPodcasts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // offset
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // per page

	count, podcasts, err := pc.as.GetPodcasts(types.Paginator{Limit: limit, Page: page})
	pagination := utils.PaginationInput{Page: page, Limit: limit, Count: count}

	if err != nil {
		log.Println(err.Error())
		utils.ErrorsResponse(c, errors.New("an error occured while getting podcasts, please try again later"))
	}

	utils.PaginatedResponse(c, podcasts, pagination)
}

func (pc *PodcastsController) GetPodcast(c *gin.Context) {
	id := c.Param("id")

	podcast, err := pc.as.GetPodcastById(id)
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

	podcast, err := pc.as.CreatePodcast(data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (pc *PodcastsController) UpdatePodcast(c *gin.Context) {
	id := c.Param("id")
	uid, _ := utils.GetCtxUser(c)

	var data types.UpdatePodcastInput
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	podcast, err := pc.as.UpdatePodcast(uid, id, data)
	if err != nil {
		utils.ErrorsResponse(c, err)
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (pc *PodcastsController) DeletePodcast(c *gin.Context) {
	id := c.Param("id")
	uid, _ := utils.GetCtxUser(c)

	res, err := pc.as.DeletePodcast(uid, id)
	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	utils.SuccessResponse(c, res)
}
