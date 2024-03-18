package main

import (
	"github.com/jackc/pgx/pgtype"
	"gorm.io/gorm"
)

type Flow struct {
	gorm.Model
	ID uint
}

type TaskType string

const (
	Input  TaskType = "input"
	Action TaskType = "action"
)

type TaskStatus = string

const (
	InProgress TaskStatus = "in_progress"
	Finished   TaskStatus = "finished"
	Stopped    TaskStatus = "stopped"
	Failed     TaskStatus = "failed"
)

type Task struct {
	gorm.Model
	ID      uint
	Type    TaskType
	Status  TaskStatus
	Args    pgtype.JSONB
	Results pgtype.JSONB
}
