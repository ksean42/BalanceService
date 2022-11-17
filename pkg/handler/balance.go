package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ksean42/BalanceService/pkg/entities"
	"log"
	"net/http"
	"time"
)

// @Summary balance user balance
// @Tags user balance
// @Accept json
// @Produce json
// @Param input body entities.UserBalanceRequest true "user id for its balance"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/get [get]
func (h *Handler) GetBalance(c *gin.Context) {
	var req entities.UserBalanceRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: "Bad request"})
		return
	}
	balance, err := h.Balance.GetBalance(req.Id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: balance})
}

// @Summary Add amount to user account
// @Tags user balance
// @Accept json
// @Produce json
// @Param input body entities.AddRequest true "user id and amount"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/add [post]
func (h *Handler) Add(c *gin.Context) {
	var req entities.AddRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: "Bad request"})
		return
	}
	err = h.Balance.Add(req.Id, req.Amount)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: "OK"})
}

// @Summary Reserve funds on user account
// @Tags user balance
// @Accept json
// @Produce json
// @Param input body entities.Request true "user id,order_id, service_id to reserve funds for service"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/reserve [post]
func (h *Handler) Reserve(c *gin.Context) {
	req := &entities.Request{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	err = h.Balance.Reserve(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: "OK"})
}

// @Summary Approve payment transaction
// @Tags user balance
// @Accept json
// @Produce json
// @Param input body entities.Request true "user id,order_id, service_id to reserve funds for approve payment service"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/approve [post]
func (h *Handler) Approve(c *gin.Context) {
	req := &entities.Request{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	err = h.Balance.Approve(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: "OK"})
}

// @Summary Create and get path to csv file with report
// @Tags report
// @Accept json
// @Produce json
// @Param input body entities.ReportRequest true "Month for report. Format: 2022-11"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/report [get]
func (h *Handler) GetReport(c *gin.Context) {
	req := &entities.ReportRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	t, err := time.Parse("2006-01", req.Date)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	fileLink, err := h.Balance.GetReport(t)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: fileLink})
}

// @Summary Create transfer between users
// @Tags user balance
// @Accept json
// @Produce json
// @Param input body entities.TransferRequest true "src user id,dest user id, amount to transfer money"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/transfer [post]
func (h *Handler) Transfer(c *gin.Context) {
	req := &entities.TransferRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	err = h.Balance.Transfer(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: "OK"})
}

// @Summary Reject reserving and refund money
// @Tags user balance
// @Accept json
// @Produce json
// @Param input body entities.ReserveReject true "user id and order id for reject the reservation"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/reject [post]
func (h *Handler) Reject(c *gin.Context) {
	req := &entities.ReserveReject{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	err = h.Balance.Reject(req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: "OK"})
}

// @Summary Get user transactions/transfers/reserver report
// @Tags report
// @Accept json
// @Produce json
// @Param input body entities.ReportRequest true "User id for report"
// @Success 200 {object} entities.ResultResponse
// @Failure 400,404 {object} entities.ErrorResponse
// @Router /api/userReport [get]
func (h *Handler) GetUserReport(c *gin.Context) {
	req := &entities.UserReportRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.Balance.GetUserReport(req.Id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, entities.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, entities.ResultResponse{Result: result})
}
