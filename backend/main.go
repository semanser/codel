package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/semanser/ai-coder/executor"
	"github.com/semanser/ai-coder/models"
	"github.com/semanser/ai-coder/router"
)

const defaultPort = "8080"

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	dsn := "postgresql://postgres@localhost/ai-coder"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.Flow{}, &models.Task{})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := router.New(db)

	err = executor.InitDockerClient()
	if err != nil {
		log.Fatalf("failed to initialize Docker client: %v", err)
	}

	// Run the server in a separate goroutine
	go func() {
		log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigChan
	log.Println("Shutting down...")

	// Cleanup resources
	if err := executor.Cleanup(); err != nil {
		log.Printf("Error during cleanup: %v", err)
	}

	log.Println("Shutdown complete")
}
