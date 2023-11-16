package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khhini/riset-gcp-api-gateway-auth/internal/modules/user/entity"
)

type Handler struct {
	UserService entity.Service
}

func NewHandler(r *gin.RouterGroup, userService entity.Service) {
	handler := Handler{
		UserService: userService,
	}
	r.POST("/users", handler.CreateUser)
	r.GET("/users/:id", handler.GetUserByID)
	r.PUT("/users/:id", handler.UpdateUser)
	r.DELETE("/users/:id", handler.DeleteUser)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser entity.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	createdUserID, err := h.UserService.InsertUser(c.Request.Context(), &newUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": createdUserID})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.UserService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updatedUser entity.User

	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	if err := h.UserService.UpdateUser(c.Request.Context(), userID, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})

}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := h.UserService.DeleteUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
