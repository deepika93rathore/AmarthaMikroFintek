package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"billing-engine/config"
	"billing-engine/service"
	"billing-engine/model"
)

// GetOutstanding handles GET requests to fetch the outstanding amount for a loan.
func GetOutstanding(c *gin.Context) {
	loanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: "invalid loan id",
		})
		return
	}

	outstanding, err := service.GetOutstanding(config.DB, uint(loanID))
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Status:  "error",
			Message: "loan not found",
		})
		return
	}

	response := model.OutstandingResponse{
		LoanID:      uint(loanID),
		Outstanding: outstanding,
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Status:  "success",
		Message: "fetched outstanding amount",
		Data:    response,
	})
}

// IsDelinquent handles GET requests to check if a borrower is delinquent.
func IsDelinquent(c *gin.Context) {
	loanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: "invalid loan id",
		})
		return
	}

	isDelinquent, err := service.IsDelinquent(config.DB, uint(loanID))
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Status:  "error",
			Message: "loan not found",
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Status:  "success",
		Message: "fetched delinquency status",
		Data: gin.H{
			"loan_id":    loanID,
			"delinquent": isDelinquent,
		},
	})
}

// MakePayment handles POST requests to process a payment.
func MakePayment(c *gin.Context) {
	loanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: "invalid loan id",
		})
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	if err := service.MakePayment(config.DB, uint(loanID), req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Status:  "success",
		Message: "payment successful",
		Data: gin.H{
			"loan_id": loanID,
			"amount":  req.Amount,
		},
	})
}
