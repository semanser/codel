package main

import (
  "github.com/jackc/pgx/pgtype"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Flow struct {
  gorm.Model
  ID uint
}

type TaskType string

const (
  Input TaskType = "input"
  Action TaskType = "action"
)

type TaskStatus = string

const (
  InProgress TaskStatus = "in_progress"
  Finished TaskStatus = "finished"
  Stoped TaskStatus = "stoped"
  Failed TaskStatus = "failed"
)

type Task struct {
  gorm.Model
  ID uint
  Type TaskType
  Status TaskStatus
  Args pgtype.JSONB
  Results pgtype.JSONB
}

func main() {
  dsn := "postgresql://postgres@localhost/ai-coder"
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // Migrate the schema
  db.AutoMigrate(&Flow{}, &Task{})

}
