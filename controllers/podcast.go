package controllers

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"podcast/gateways/upload"
	"podcast/services"
	"podcast/types"
	"podcast/utils"

	"github.com/gin-gonic/gin"
)

type PodcastsController struct {
	ps *services.PodcastsService
	fh upload.FileHandler
}

func NewPodcastsController(ps *services.PodcastsService, fh upload.FileHandler) *PodcastsController {
	return &PodcastsController{ps: ps, fh: fh}
}

func (pc *PodcastsController) GetPodcasts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	count, podcasts, err := pc.ps.GetPodcasts(types.Paginator{Limit: limit, Page: page})
	pagination := utils.PaginationInput{Page: page, Limit: limit, Count: count}

	if err != nil {
		log.Println(err.Error())
		utils.ErrorResponse(c, err, "An error has occured while getting all podcasts, please try again later")
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

func (pc *PodcastsController) GetPodcastBySlug(c *gin.Context) {
	slug := c.Param("pslug")

	podcast, err := pc.ps.GetPodcastBySlug(slug)
	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (pc *PodcastsController) CreatePodcast(c *gin.Context) {
	c.Request.Body = utils.LimitRequestSize(c, 6)
	var data types.CreatePodcastInput

	picture, err := utils.HandleFileValidation(c, "picture", []string{"jpg", "png"}, 5)
	if err != nil {
		utils.ErrorResponse(c, err, "Please include a picture of type png or jpg, that does not exceed 5MB")
		return
	}

	picturePath, err := pc.fh.Upload(picture)
	if err != nil {
		utils.ErrorResponse(c, err, "An error occured while uploading the provided picture, please try again later")
		return
	}

	var socials types.SocialLinks
	json.Unmarshal([]byte(c.PostForm("social_links")), &socials)

	data = types.CreatePodcastInput{
		Name:        c.PostForm("name"),
		Headline:    c.PostForm("headline"),
		Description: c.PostForm("description"),
		Hosts:       strings.Split(c.PostForm("hosts"), ", "),
		Tags:        strings.Split(c.PostForm("tags"), ", "),
		Picture:     picturePath,
		SocialLinks: socials,
	}

	u, _ := c.Get("user")
	data.CreatorId = u.(utils.JwtPayload).ID

	podcast, err := pc.ps.CreatePodcast(data)
	if err != nil {
		utils.ErrorResponse(c, err, "Please provide valid information in order to create this podcast")
		return
	}

	utils.SuccessResponse(c, podcast)
}

func (pc *PodcastsController) UpdatePodcast(c *gin.Context) {
	id := c.Param("pid")
	uid, _ := utils.GetCtxUser(c)

	var data types.UpdatePodcastInput
	var picturePath string

	picture, err := c.FormFile("picture")
	if err != nil {
		log.Println("no picture was uploaded")
	}

	if picture != nil {
		// TODO: delete old picture if new one got uploaded
		picturePath, err = pc.fh.Upload(picture)
		if err != nil {
			utils.ErrorResponse(c, err, "An error occured while uploading the provided picture, please try again later")
			return
		}
	}

	var socials types.SocialLinks
	json.Unmarshal([]byte(c.PostForm("social_links")), &socials)

	data = types.UpdatePodcastInput{
		Headline:    c.PostForm("headline"),
		Description: c.PostForm("description"),
		Hosts:       strings.Split(c.PostForm("hosts"), ", "),
		Tags:        strings.Split(c.PostForm("tags"), ", "),
		Picture:     picturePath,
		SocialLinks: socials,
	}

	podcast, err := pc.ps.UpdatePodcast(uid, id, data)
	if err != nil {
		utils.ErrorResponse(c, err, "Please provide valid information to be able to update this podcast")
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

	_, err := pc.ps.GetPodcastById(pid)
	if err != nil {
		utils.ErrorResponse(c, err, "The podcast you are trying to subscribe to does not exist")
		return
	}

	url, err := pc.ps.Subscribe(uid, pid)
	if err != nil {
		utils.ErrorResponse(c, err, "Unable to subscribe to this podcast, please try again later with valid information")
		return
	}

	// TODO: enable
	// c.Redirect(http.StatusTemporaryRedirect, url)
	utils.SuccessResponse(c, url)
}

func (pc *PodcastsController) GetPodcastByCreator(c *gin.Context) {
	cid := c.Param("cid")

	podcast, err := pc.ps.GetPodcastByCreatorId(cid)
	if err != nil {
		utils.NotFoundResponse(c)
		return
	}

	utils.SuccessResponse(c, podcast)
}
