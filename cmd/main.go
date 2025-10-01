package main

import (
	"gofsen/internal/middlewares"
	"gofsen/internal/router"
	"log"
	"net/http"
)

func main() {

	router := router.NewRouter()

	// Add global middlewares
	router.Use(middlewares.RecoveryMiddleware) // Add recovery first
	router.Use(middlewares.LoggerMiddleware)

	log.Println("🚀 Gofsen API running at http://localhost:8081")
	log.Println("📋 Test all functionalities at: http://localhost:8081/test/all")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}
