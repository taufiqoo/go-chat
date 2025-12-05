package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taufiqoo/go-chat/internal/domain"
	"github.com/taufiqoo/go-chat/internal/service"
	"github.com/taufiqoo/go-chat/internal/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewsUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", user)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	response := domain.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", response)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	// Remove password from response
	var userResponses []domain.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, domain.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", userResponses)
}
