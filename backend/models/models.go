package models

type FlowStatus string

const (
	FlowInProgress FlowStatus = "in_progress"
	FlowFinished   FlowStatus = "finished"
)

type Flow struct {
	ID          uint
	Name        string
	Tasks       []Task
	Status      FlowStatus
	ContainerID uint
	Container   Container
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
	ID      uint
	Message string
	Type    TaskType
	Status  TaskStatus
	// Args    datatypes.JSON
	Results string
	FlowID  uint
	Flow    Flow
}

type ContainerStatus = string

const (
	ContainerStarting ContainerStatus = "starting"
	ContainerRunning  ContainerStatus = "running"
	ContainerStopped  ContainerStatus = "stopped"
	ContainerFailed   ContainerStatus = "failed"
)

type Container struct {
	ID     uint
	Name   string
	Image  string
	Status ContainerStatus
}
