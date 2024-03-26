package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/config"
	"github.com/semanser/ai-coder/executor"
	"github.com/semanser/ai-coder/models"
	"github.com/semanser/ai-coder/router"
	"github.com/semanser/ai-coder/services"
)

//go:embed templates/prompts/*.tmpl
var promptTemplates embed.FS

func main() {
	config.Init()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	db, err := gorm.Open(postgres.Open(config.Config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.Flow{}, &models.Task{})

	port := strconv.Itoa(config.Config.Port)

	r := router.New(db)

	assets.Init(promptTemplates)
	services.Init()

	err = executor.InitDockerClient()
	if err != nil {
		log.Fatalf("failed to initialize Docker client: %v", err)
	}

	executor.ProcessQueue(db)

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
