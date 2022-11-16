package handler

import (
	"avito_test_task/pkg/entities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error string
}

type ResultResponse struct {
	Result interface{}
}

func (h *Handler) GetBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	balance, err := h.Balance.GetBalance(id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{balance})
}

func (h *Handler) Add(c *gin.Context) {
	var req entities.Request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	err = h.Balance.Add(req.Id, req.Amount)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{"OK"})
}
func (h *Handler) Reserve(c *gin.Context) {
	req := &entities.Request{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	err = h.Balance.Reserve(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}
	c.JSON(http.StatusOK, ResultResponse{"OK"})
}
func (h *Handler) Approve(c *gin.Context) {

}
