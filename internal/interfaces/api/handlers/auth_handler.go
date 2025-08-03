package handlers

import (
	"net/http"

	"github.com/anggitaputra13/mcs-auth/internal/domain/entities"
	"github.com/anggitaputra13/mcs-auth/internal/domain/services"
	"github.com/anggitaputra13/mcs-auth/pkg/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
	jwt         *auth.JWT
}

func NewAuthHandler(authService services.AuthService, jwt *auth.JWT) *AuthHandler {
	return &AuthHandler{authService: authService, jwt: jwt}
}

// ErrorResponse represents an error response
// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
// swagger:model
type SuccessResponse struct {
	Message string `json:"message"`
}

// RegisterRequest represents user registration request
// swagger:model
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents user login request
// swagger:model
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse represents token response
// swagger:model
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// ValidateResponse represents token validation response
// swagger:model
type ValidateResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "User registration data"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.authService.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "User login data"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"access_token": token,
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetString("token")
	if err := h.authService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// Validate godoc
// @Summary Validate token
// @Description Validate user token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ValidateResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/validate [get]
func (h *AuthHandler) Validate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "message": "Token is valid"})
}
