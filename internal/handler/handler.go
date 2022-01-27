package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/EMus88/go-musthave-shortener-tpl/internal/app/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
	//
}
type Request struct {
	LongURL string `json:"url"`
}
type Result struct {
	Result string `json:"result"`
}

func NewHandler(service *service.Service) *Handler {
	s
	return &Handler{service: service}
}

//=================================================================
func (h *Handler) HandlerGet(c *gin.Context) {
	id := c.Param("id")
	longURL, err := h.service.GetURL(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusTemporaryRedirect)
	c.Header("Location", longURL)
}

//==================================================================
func (h *Handler) HandlerPostText(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.String(http.StatusBadRequest, "Not allowed request")
		return
	}
	id, err := h.service.SaveURL(string(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	}
	c.String(http.StatusCreated, h.service.Config.BaseURL+"/"+id)
}

//===================================================================
func (h *Handler) HandlerPostJSON(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	if request.LongURL == "" || c.GetHeader("content-type") != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}

	id, err := h.service.SaveURL(request.LongURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	}
	shortURL := fmt.Sprint(h.service.Config.BaseURL, "/", id)
	var result Result
	result.Result = shortURL
	c.JSON(http.StatusCreated, result)
}
