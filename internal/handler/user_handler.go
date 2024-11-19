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

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// QueryUsersRequest represents the query parameters for listing users
type QueryUsersRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Username string `form:"username"`
}

// QueryUsersResponse represents the response for listing users
type QueryUsersResponse struct {
	Total int64         `json:"total"`
	Items []UserResponse `json:"items"`
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

// QueryUsers handles the request to list all users with pagination and filtering
func (h *UserHandler) QueryUsers(c *gin.Context) {
	var req QueryUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, err.Error())
		return
	}

	// Initialize the database query
	query := h.db.Model(&model.User{})

	// Apply username filter if provided
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.ServerErrorResponse(c, "failed to count users")
		return
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Get users with pagination
	var users []model.User
	if err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		utils.ServerErrorResponse(c, "failed to query users")
		return
	}

	// Convert to response format
	items := make([]UserResponse, len(users))
	for i, user := range users {
		items[i] = UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	utils.SuccessResponse(c, QueryUsersResponse{
		Total: total,
		Items: items,
	})
}
