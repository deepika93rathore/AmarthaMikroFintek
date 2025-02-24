package config

import (
	"log"
	"gorm.io/driver/postgres"
	"github.com/lpernett/godotenv"
	"gorm.io/gorm"
	"os"
	"fmt"
	"billing-engine/model"
)

var DB *gorm.DB

func ConnectDatabase() {

		// Load environment variables
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	
		// Format PostgreSQL DSN
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	
		// Open Database Connection
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to the database:", err)
		}
	
		// Assign db to global variable
		DB = db
	
		fmt.Println("Connected to PostgreSQL!")
		err = DB.AutoMigrate(&model.Loan{},&model.Payment{}) // Create 'domains' table first
		if err != nil {
			log.Fatalf("Failed to migrate Domain: %v", err)
		}

		SeedDatabase()
}


func SeedDatabase() {
    // Check if we already have Loans in the table
    var loanCount int64
    DB.Model(&model.Loan{}).Count(&loanCount)
    if loanCount > 0 {
        fmt.Println("Seed data already exists. Skipping seeding.")
        return
    }

    // ---------------------
    // SEED LOANS
    // ---------------------
    loans := []model.Loan{
        {
            BorrowerID:    101,
            Amount:        5000000,
            InterestRate:  0.10, // 10% annual
            Weeks:         50,
            // Will calculate TotalAmount, WeeklyPayment, Outstanding via CalculateSchedule()
        },
        {
            BorrowerID:    102,
            Amount:        3000000,
            InterestRate:  0.10,
            Weeks:         50,
        },
        {
            BorrowerID:    103,
            Amount:        8000000,
            InterestRate:  0.10,
            Weeks:         50,
        },
    }

    // Compute the schedule for each loan, then insert them
    for i := range loans {
        loans[i].CalculateSchedule()
        DB.Create(&loans[i])
    }

    fmt.Println("Loans seeded successfully.")

    // ---------------------
    // SEED PAYMENTS
    // ---------------------
    // Let's assume the first loan (BorrowerID=101) made a few payments
    payments := []model.Payment{
        {
            LoanID: loans[0].ID,
            Amount: loans[0].WeeklyPayment,
            Week:   1,
        },
        {
            LoanID: loans[0].ID,
            Amount: loans[0].WeeklyPayment,
            Week:   2,
        },
        {
            LoanID: loans[1].ID,
            Amount: loans[1].WeeklyPayment,
            Week:   1,
        },
    }

    for _, p := range payments {
        DB.Create(&p)
    }

    fmt.Println("Payments seeded successfully.")

    // Update the Outstanding for each loan based on seeded payments
    updateOutstanding(loans)
}

func updateOutstanding(loans []model.Loan) {
    for _, loan := range loans {
        var totalPaid float64
        var payments []model.Payment
        DB.Where("loan_id = ?", loan.ID).Find(&payments)
        for _, p := range payments {
            totalPaid += p.Amount
        }
        // Recalculate the outstanding
        updatedOutstanding := loan.TotalAmount - totalPaid

        DB.Model(&model.Loan{}).Where("id = ?", loan.ID).Update("outstanding", updatedOutstanding)
    }
}
