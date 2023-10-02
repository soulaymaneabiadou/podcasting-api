package utils

import (
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"podcast/types"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	PrevPage   int     `json:"prev_page"`
	NextPage   int     `json:"next_page"`
	TotalPages float64 `json:"total_pages"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
}

type PaginationInput struct {
	Page  int
	Limit int
	Count int64
}

type response struct {
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Errors     interface{} `json:"errors"`
	Pagination Pagination  `json:"pagination"`
}

func createPagination(i PaginationInput) Pagination {
	total := float64(i.Count) / float64(i.Page*i.Limit)
	return Pagination{
		Page:       i.Page,
		Limit:      i.Limit,
		PrevPage:   i.Page - 1,
		NextPage:   i.Page + 1,
		TotalPages: math.Ceil(total),
	}
}

func PaginatedResponse(c *gin.Context, data interface{}, i PaginationInput) {
	pagination := createPagination(i)
	res := response{
		Message:    "",
		Success:    true,
		Data:       data,
		Errors:     nil,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, res)
}

func SetAccessCookie(c *gin.Context, access string) {
	accessCookieExp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_COOKIE_EXPIRE"))
	accessCookieExp = int(60 * accessCookieExp) // minutes
	c.SetCookie("access_token", access, accessCookieExp, "/", os.Getenv("DOMAIN"), false, true)
}

func SetRefreshCookie(c *gin.Context, refresh string) {
	refreshCookieExp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_COOKIE_EXPIRE"))
	refreshCookieExp = int(60 * 60 * 24 * refreshCookieExp) // days
	c.SetCookie("refresh_token", refresh, refreshCookieExp, "/", os.Getenv("DOMAIN"), false, true)
}

func SendAccessTokenResponse(c *gin.Context, token string) {
	SetAccessCookie(c, token)

	res := gin.H{"access_token": token}
	c.JSON(http.StatusOK, res)
}

func SendTokensResponse(c *gin.Context, user types.User) {
	accessToken, refreshToken, err := CreateTokens("", JwtPayload{ID: user.ID})
	if err != nil {
		log.Println(err)
		ErrorsResponse(c, err)
		return
	}

	SetAccessCookie(c, accessToken)
	SetRefreshCookie(c, refreshToken)

	res := gin.H{"access_token": accessToken, "refresh_token": refreshToken, "user": user}
	c.JSON(http.StatusOK, res)
}

func ClearAuthCookies(c *gin.Context) {
	keys := []string{"access_token", "refresh_token"}
	for _, v := range keys {
		c.SetCookie(v, "", 10, "/", os.Getenv("DOMAIN"), false, true)
	}
	// c.SetCookie("access_token", "", 0, "/", os.Getenv("DOMAIN"), false, true)
	// c.SetCookie("refresh_token", "", 0, "/", os.Getenv("DOMAIN"), false, true)
}

func SuccessResponse(c *gin.Context, data interface{}) {
	res := response{
		Message: "",
		Success: true,
		Data:    data,
		Errors:  nil,
	}

	c.JSON(http.StatusOK, res)
}

func NotFoundResponse(c *gin.Context) {
	res := response{
		Message: "Resource not found",
		Success: false,
		Data:    nil,
		Errors:  map[string]interface{}{},
	}

	c.JSON(http.StatusNotFound, res)
}

func UnauthorizedResponse(c *gin.Context) {
	res := response{
		Message: "Unauthorized",
		Success: false,
		Data:    nil,
		Errors:  map[string]interface{}{},
	}

	c.JSON(http.StatusUnauthorized, res)
	c.Abort()
}

func ForbiddenResponse(c *gin.Context, msg string) {
	res := response{
		Message: msg,
		Success: false,
		Data:    nil,
		Errors:  map[string]interface{}{},
	}

	c.JSON(http.StatusForbidden, res)
	c.Abort()
}

func ErrorsResponse(c *gin.Context, err error) {
	res := response{
		Message: "Invalid data",
		Success: false,
		Data:    nil,
		Errors:  parseError(err),
	}

	c.JSON(http.StatusBadRequest, res)
}

func MessageResponse(c *gin.Context, msg string) {
	res := response{
		Message: msg,
		Success: true,
		Data:    nil,
		Errors:  nil,
	}

	c.JSON(http.StatusOK, res)
}

func ErrorResponse(c *gin.Context, err error, msg string) {
	res := response{
		Message: msg,
		Success: false,
		Data:    nil,
		Errors:  parseError(err),
	}

	c.JSON(http.StatusBadRequest, res)
}
