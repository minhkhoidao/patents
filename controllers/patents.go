package controllers

import (
	"go-nginx/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PatentController struct {
	service *services.PatentService
}

func NewPatentController(service *services.PatentService) *PatentController {
	return &PatentController{service: service}
}

func (c *PatentController) GetDataGraph(ctx *gin.Context) {
	start_date := ctx.Query("start_date")
	end_date := ctx.Query("end_date")

	patents, err := c.service.GetDataGraph(ctx, start_date, end_date)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, patents)
}
