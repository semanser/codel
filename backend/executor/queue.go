package executor

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/semanser/ai-coder/database"
	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/graph/subscriptions"
	"github.com/semanser/ai-coder/providers"
)

var queue = make(map[int64]chan database.Task)
var stopChannels = make(map[int64]chan any)

func AddQueue(flowId int64, db *database.Queries) {
	if _, ok := queue[flowId]; !ok {
		queue[flowId] = make(chan database.Task, 1000)

		stop := make(chan any)
		stopChannels[flowId] = stop
		ProcessQueue(flowId, db)
	}
}

func AddCommand(flowId int64, task database.Task) {
	if queue[flowId] != nil {
		queue[flowId] <- task
	}
	log.Printf("Command %d added to the queue %d", task.ID, flowId)
}

func CleanQueue(flowId int64) {
	if _, ok := queue[flowId]; ok {
		queue[flowId] = nil
	}

	if _, ok := stopChannels[flowId]; ok {
		close(stopChannels[flowId])
		stopChannels[flowId] = nil
	}

	log.Printf("Queue %d cleaned", flowId)
}

func ProcessQueue(flowId int64, db *database.Queries) {
	log.Println("Starting tasks processor for queue", flowId)

	flow, err := db.ReadFlow(context.Background(), flowId)

	if err != nil {
		log.Printf("failed to get provider: %v", err)
		CleanQueue(flowId)
		return
	}

	provider, err := providers.ProviderFactory(providers.ProviderType(flow.ModelProvider.String))

	if err != nil {
		log.Printf("failed to get provider: %v", err)
		CleanQueue(flowId)
		return
	}

	log.Printf("Using provider: %s. Model: %s\n", provider.Name(), flow.Model.String)

	go func() {
		for {
			select {
			case <-stopChannels[flowId]:
				log.Printf("Stopping task processor for queue %d", flowId)
				return
			default:

				log.Println("Waiting for a task")
				task := <-queue[flowId]

				log.Printf("Processing command %d of type %s", task.ID, task.Type.String)

				// Input tasks are added by the user optimistically on the client
				// so they should not be broadcasted back to the client
				subscriptions.BroadcastTaskAdded(task.FlowID.Int64, &gmodel.Task{
					ID:        uint(task.ID),
					Message:   task.Message.String,
					Type:      gmodel.TaskType(task.Type.String),
					CreatedAt: task.CreatedAt.Time,
					Status:    gmodel.TaskStatus(task.Status.String),
					Args:      task.Args.String,
					Results:   task.Results.String,
				})

				if task.Type.String == "input" {
					err := processInputTask(provider, db, task)

					if err != nil {
						log.Printf("failed to process input: %v", err)
						continue
					}

					nextTask, err := getNextTask(provider, db, task.FlowID.Int64)

					if err != nil {
						log.Printf("failed to get next task: %v", err)
						continue
					}

					AddCommand(flowId, *nextTask)
				}

				if task.Type.String == "ask" {
					err := processAskTask(db, task)

					if err != nil {
						log.Printf("failed to process ask: %v", err)
						continue
					}
				}

				if task.Type.String == "terminal" {
					err := processTerminalTask(db, task)

					if err != nil {
						log.Printf("failed to process terminal: %v", err)
						continue
					}
					nextTask, err := getNextTask(provider, db, task.FlowID.Int64)

					if err != nil {
						log.Printf("failed to get next task: %v", err)
						continue
					}

					AddCommand(flowId, *nextTask)
				}

				if task.Type.String == "code" {
					err := processCodeTask(db, task)

					if err != nil {
						log.Printf("failed to process code: %v", err)
						continue
					}

					nextTask, err := getNextTask(provider, db, task.FlowID.Int64)

					if err != nil {
						log.Printf("failed to get next task: %v", err)
						continue
					}

					AddCommand(flowId, *nextTask)
				}

				if task.Type.String == "done" {
					err := processDoneTask(db, task)

					if err != nil {
						log.Printf("failed to process done: %v", err)
						continue
					}
				}

				if task.Type.String == "browser" {
					err := processBrowserTask(db, task)

					if err != nil {
						log.Printf("failed to process browser: %v", err)
						continue
					}

					nextTask, err := getNextTask(provider, db, task.FlowID.Int64)

					if err != nil {
						log.Printf("failed to get next task: %v", err)
						continue
					}

					AddCommand(flowId, *nextTask)
				}
			}
		}
	}()
}

func getNextTask(provider providers.Provider, db *database.Queries, flowId int64) (*database.Task, error) {
	flow, err := db.ReadFlow(context.Background(), flowId)

	if err != nil {
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}

	tasks, err := db.ReadTasksByFlowId(context.Background(), sql.NullInt64{
		Int64: flowId,
		Valid: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by flow id: %w", err)
	}

	const maxResultsLength = 4000
	for i, task := range tasks {
		// Limit the number of result characters since some output commands can have a lot of output
		if len(task.Results.String) > maxResultsLength {
			// Get the last N symbols from the output
			results := task.Results.String[len(task.Results.String)-maxResultsLength:]
			tasks[i].Results = database.StringToNullString(results)
		}
	}

	c := provider.NextTask(providers.NextTaskOptions{
		Tasks:       tasks,
		DockerImage: flow.ContainerImage.String,
	})

	lastTask := tasks[len(tasks)-1]

	_, err = db.UpdateTaskToolCallId(context.Background(), database.UpdateTaskToolCallIdParams{
		ToolCallID: c.ToolCallID,
		ID:         lastTask.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update task tool call id: %w", err)
	}

	nextTask, err := db.CreateTask(context.Background(), database.CreateTaskParams{
		Args:       c.Args,
		Message:    c.Message,
		Type:       c.Type,
		Status:     database.StringToNullString("in_progress"),
		FlowID:     sql.NullInt64{Int64: flowId, Valid: true},
		ToolCallID: c.ToolCallID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to save command: %w", err)
	}

	return &nextTask, nil
}
