package routes

import (
	"go-nginx/controllers"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func SetupRouter(patentController *controllers.PatentController) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"}, // Specify your frontend URLs
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Set to true if you need to send cookies
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/patents/", patentController.GetDataGraph)
	return router
}
