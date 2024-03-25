package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FlowStatus string

const (
	FlowInProgress FlowStatus = "in_progress"
	FlowFinished   FlowStatus = "finished"
)

type Flow struct {
	gorm.Model
	ID          uint
	Name        string
	Tasks       []Task
	Status      FlowStatus
	DockerImage string
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
	TaskInProgress TaskStatus = "in_progress"
	TaskFinished   TaskStatus = "finished"
	TaskStopped    TaskStatus = "stopped"
	TaskFailed     TaskStatus = "failed"
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
