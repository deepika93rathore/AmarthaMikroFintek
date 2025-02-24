package model

import (
	"time"
)

type Loan struct {
	ID             uint           `gorm:"primaryKey"`
	BorrowerID     uint
	Amount         float64
	InterestRate   float64
	TotalAmount    float64
	WeeklyPayment  float64
	Weeks          int
	Outstanding    float64
	StartDate      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Payments       []Payment
}

type Payment struct {
	ID        uint      `gorm:"primaryKey"`
	LoanID    uint
	Amount    float64
	Week      int
	CreatedAt time.Time
}

// APIResponse is a generic structure for all responses.
type APIResponse struct {
	Status  string      `json:"status"`  // "success" or "error"
	Message string      `json:"message"` // a human-readable message
	Data    interface{} `json:"data"`    // payload data
}

// OutstandingResponse is used to return the outstanding amount for a loan.
type OutstandingResponse struct {
	LoanID      uint    `json:"loan_id"`
	Outstanding float64 `json:"outstanding"`
}

func (loan *Loan) CalculateSchedule() {
	loan.TotalAmount = loan.Amount + (loan.Amount * loan.InterestRate * float64(loan.Weeks) / 52)
	loan.WeeklyPayment = loan.TotalAmount / float64(loan.Weeks)
	loan.Outstanding = loan.TotalAmount
}