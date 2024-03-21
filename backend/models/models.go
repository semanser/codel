package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Flow struct {
	gorm.Model
	ID    uint
	Name  string
	Tasks []Task
}

type TaskType string

const (
	Input    TaskType = "input"
	Terminal TaskType = "terminal"
	Browser  TaskType = "browser"
	Code     TaskType = "code"
	Ask      TaskType = "ask"
	Done     TaskType = "done"
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
	Message string
	Type    TaskType
	Status  TaskStatus
	Args    datatypes.JSON
	Results string
	FlowID  uint
	Flow    Flow
}
