package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"wind-surf-go/internal/handler"
	"wind-surf-go/internal/middleware"
)

// SetupRouter initializes all routes and returns the router engine
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize handlers
	userHandler := handler.NewUserHandler(db)

	// Base route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Wind Surf API Server")
	})

	// API v1 routes
	v1 := r.Group("/v1")
	{
		api := v1.Group("/api")
		{
			// Public user routes
			users := api.Group("/users")
			{
				users.POST("/register", userHandler.Register)
				users.POST("/login", userHandler.Login)
			}

			// Protected user routes
			authUsers := api.Group("/users", middleware.AuthMiddleware())
			{
				authUsers.GET("", userHandler.QueryUsers)
			}

			// Add more API groups here as needed
			// posts := api.Group("/posts")
			// {
			//     posts.GET("/", postHandler.List)
			//     posts.POST("/", postHandler.Create)
			// }
		}
	}

	return r
}
