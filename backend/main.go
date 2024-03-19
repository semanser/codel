package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/semanser/ai-coder/models"
)

const defaultPort = "8080"

func main() {
	dsn := "postgresql://postgres@localhost/ai-coder"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Flow{}, &models.Task{})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := newRouter(db)

	// Run the server
	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
