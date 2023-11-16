package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler(r *gin.RouterGroup) {
	handler := Handler{}
	r.GET("/hello", handler.helloWorld)
}

func (h *Handler) helloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
