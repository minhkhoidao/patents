package routes

import (
	"go-nginx/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(patentController *controllers.PatentController) *gin.Engine {
	router := gin.Default()
	router.GET("/patents", patentController.GetDataGraph)
	return router
}
