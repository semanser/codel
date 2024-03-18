package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Flow struct {
	gorm.Model
	ID    uint
	Tasks []Task
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
	ID          uint
	Description string
	Type        TaskType
	Status      TaskStatus
	Args        datatypes.JSON
	Results     datatypes.JSON
	FlowID      uint
	Flow        Flow
}
