package handler

import (
	"wind-surf-go/internal/model"
	"wind-surf-go/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	// Check if username already exists
	var existingUser model.User
	if result := h.db.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		utils.ErrorResponse(c, "username already exists")
		return
	}

	// Create new user
	user := model.User{
		Username: req.Username,
		Password: req.Password,
	}

	// Hash password before saving
	if err := user.HashPassword(); err != nil {
		utils.ServerErrorResponse(c, "failed to hash password")
		return
	}

	// Save user to database
	if result := h.db.Create(&user); result.Error != nil {
		utils.ServerErrorResponse(c, "failed to create user")
		return
	}

	// Generate token for the new user
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ServerErrorResponse(c, "failed to generate token")
		return
	}

	utils.CreatedResponse(c, LoginResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	// Find user by username
	var user model.User
	if result := h.db.Where("username = ?", req.Username).First(&user); result.Error != nil {
		utils.ErrorResponse(c, "invalid credentials")
		return
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		utils.ErrorResponse(c, "invalid credentials")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.ServerErrorResponse(c, "failed to generate token")
		return
	}

	utils.SuccessResponse(c, LoginResponse{
		Token:    token,
		Username: user.Username,
		UserID:   user.ID,
	})
}
