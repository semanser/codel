package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/semanser/ai-coder/assets"
	"github.com/semanser/ai-coder/config"
	"github.com/semanser/ai-coder/database"
	"github.com/semanser/ai-coder/executor"
	"github.com/semanser/ai-coder/router"
)

//go:embed templates/prompts/*.tmpl
var promptTemplates embed.FS

//go:embed templates/scripts/*.js
var scriptTemplates embed.FS

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	config.Init()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	db, err := sql.Open("sqlite3", config.Config.DatabaseURL)

	queries := database.New(db)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatalf("Unable to set dialect: %v\n", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Unable to run migrations: %v\n", err)
	}

	log.Println("Migrations ran successfully")

	port := strconv.Itoa(config.Config.Port)

	r := router.New(queries)

	assets.Init(promptTemplates, scriptTemplates)

	err = executor.InitClient()
	if err != nil {
		log.Fatalf("failed to initialize Docker client: %v", err)
	}

	err = executor.InitBrowser(queries)
	if err != nil {
		log.Fatalf("failed to initialize browser container: %v", err)
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
	if err := executor.Cleanup(queries); err != nil {
		log.Printf("Error during cleanup: %v", err)
	}

	log.Println("Shutdown complete")
}
