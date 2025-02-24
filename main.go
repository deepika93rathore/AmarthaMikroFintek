package main

import (
	"github.com/gin-gonic/gin"
	"billing-engine/config"
	"billing-engine/routes"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Connect to the database and migrate models
	config.ConnectDatabase()

	// Register routes
	routes.RegisterRoutes(r)

	// Run server on port 8080
	r.Run(":3000")
}
