package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/ksean42/BalanceService/docs"
	"github.com/ksean42/BalanceService/pkg/services"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services.Balance
}

func NewHandler(service *services.BalanceService) *Handler {
	return &Handler{service}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		api.GET("/balance", h.GetBalance)
		api.GET("/report", h.GetReport)
		api.GET("/userReport", h.GetUserReport)
		api.POST("/add", h.Add)
		api.POST("/reserve", h.Reserve)
		api.POST("/approve", h.Approve)
		api.POST("/transfer", h.Transfer)
		api.POST("/reject", h.Reject)
	}
	return router
}
