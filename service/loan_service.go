package service

import (
	"errors"
	"billing-engine/model"
	"gorm.io/gorm"
)

// Get Outstanding Loan Amount
func GetOutstanding(db *gorm.DB, loanID uint) (float64, error) {
	var loan model.Loan
	if err := db.Preload("Payments").First(&loan, loanID).Error; err != nil {
		return 0, err
	}
	return loan.Outstanding, nil
}

// Check if Borrower is Delinquent
func IsDelinquent(db *gorm.DB, loanID uint) (bool, error) {
	var payments []model.Payment
	if err := db.Where("loan_id = ?", loanID).Order("week desc").Limit(2).Find(&payments).Error; err != nil {
		return false, err
	}

	if len(payments) < 2 {
		return true, nil
	}

	return false, nil
}

// Make a Payment
func MakePayment(db *gorm.DB, loanID uint, amount float64) error {
	var loan model.Loan
	if err := db.First(&loan, loanID).Error; err != nil {
		return err
	}

	if loan.Outstanding < amount {
		return errors.New("payment exceeds outstanding balance")
	}

	payment := model.Payment{
		LoanID: loanID,
		Amount: amount,
		Week:   loan.Weeks - int(loan.Outstanding/loan.WeeklyPayment),
	}

	if err := db.Create(&payment).Error; err != nil {
		return err
	}

	loan.Outstanding -= amount
	return db.Save(&loan).Error
}
