package handler

import (
	"family_budget/internal/dto"
	"family_budget/internal/entities/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService users.UserService
}

func NewUserHandler(userService users.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса", "details": err.Error()})
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать пользователя", "details": err.Error()})
		return
	}

	response := dto.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Login:   user.Login,
		Email:   user.Email,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accesstoken, refreshToken, err := h.userService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: accesstoken,
		RefreshToken: refreshToken,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccessToken, newRefreshToken, err := h.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
