package handler

import (
	"avito_test_task/pkg/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services.Balance
}

func NewHandler(service *services.BalanceService) *Handler {
	return &Handler{service}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	api.GET("/get/:id", h.GetBalance)
	api.GET("/report", h.GetReport)
	api.POST("/add", h.Add)
	api.POST("/reserve", h.Reserve)
	api.POST("/approve", h.Approve)
	return router
}
