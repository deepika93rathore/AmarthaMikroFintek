# AmarthaMikroFintek

Project Overview
This project is a billing engine designed to manage loan schedules, track outstanding balances, determine delinquency status, and process payments. It’s built using Go, Gin (for HTTP routing), and GORM (for database interactions with PostgreSQL). The system supports creating loans, generating a repayment schedule, accepting weekly payments, and calculating the outstanding amount.

#### Create .env file for database connection details ############

Key Features
Loan Schedule Generation:
When a new loan is created, the system calculates the total repayable amount (including a flat annual interest rate), determines a fixed weekly payment, and sets the initial outstanding balance.

Outstanding Balance Tracking:
The outstanding balance is updated as payments are made. It represents the remaining amount that the borrower needs to repay.

Delinquency Detection:
If the borrower misses two consecutive payments, they are marked as delinquent. The system checks the recent payment history to identify such cases.

Payment Processing:
Borrowers can make payments that are exactly equal to the weekly installment amount. The system records each payment and updates the loan's outstanding balance accordingly.

API Responses:
All API responses follow a consistent structure with a status, a message, and data payload. This ensures clarity for the clients consuming the API.

{
    "status": "success",
    "message": "fetched outstanding amount",
    "data": {
        "loan_id": 1,
        "outstanding": 5261538.461538462
    }
}

{
    "status": "success",
    "message": "fetched delinquency status",
    "data": {
        "delinquent": false,
        "loan_id": 1
    }
}

{
    "status": "success",
    "message": "payment successful",
    "data": {
        "amount": 110000,
        "loan_id": 2
    }
}

Seed Data:
The application includes seeding logic to populate initial loan and payment data for testing and development purposes.

Project Structure

billing-engine/
├── config/
│   └── database.go       // Handles DB connection, migrations, and seeding
├── handlers/
│   └── loan_handler.go   // Contains the HTTP handlers for API endpoints
├── models/
│   ├── model.go           // Loan data model and schedule calculations,Payment data model,API response structures for consistent responses
├── routes/
│   └── routes.go         // Defines API routes and groups endpoints (e.g., /api)
├── services/
│   └── loan_service.go   // Business logic for loan operations (outstanding, payment, delinquency)
├── main.go               // Application entry point; sets up Gin and routes
├── go.mod                // Module declaration and dependency management
└── go.sum                // Dependency versions and integrity checks

How It Works

Initialization and Database Connection:

The application starts in main.go, where the Gin router is initialized.
The config/database.go file establishes a connection to PostgreSQL, auto-migrates the schema based on your models, and seeds initial data if needed.

API Endpoints:
The project exposes several endpoints. For example:

GET /loan/:id/outstanding: Returns the current outstanding balance for a loan.
GET /loan/:id/delinquent: Checks whether a borrower is delinquent (based on missed payments).
POST /loan/:id/payment: Accepts a payment amount, updates the loan balance, and records the payment.

Business Logic:
The services/loan_service.go file contains the core logic:

It retrieves and calculates the outstanding balance.
It determines delinquency by checking recent payment history.
It processes payments, ensuring that each payment matches the scheduled installment and updating the outstanding amount.

Response Models:
To maintain consistency across API responses, a response model (models/APIResponse) is used:

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
This ensures that every API response has a clear status (success or error), a descriptive message, and any necessary data payload.

Testing the API:
Use tools like Postman to hit the endpoints. For instance, you can test the outstanding balance endpoint with:

GET http://localhost:8080/loan/1/outstanding
And for processing a payment:

Edit
POST http://localhost:8080/loan/1/payment
Content-Type: application/json

{
    "amount": 110000
}


Summary

Technologies:
Built using Go, Gin, GORM, and PostgreSQL.

Design:
The project is structured with clear separation between configuration, models, services, handlers, and routes. This modularity ensures maintainability and scalability.

User Benefit:
The consistent API response model and clear endpoint definitions make it easy for clients (such as web or mobile apps) to integrate with and consume the loan billing service.
