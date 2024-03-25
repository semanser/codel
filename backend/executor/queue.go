package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	gorillaWs "github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/semanser/ai-coder/agent"
	"github.com/semanser/ai-coder/database"
	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/graph/subscriptions"
	"github.com/semanser/ai-coder/services"
	"github.com/semanser/ai-coder/websocket"
)

var queue = make(chan database.Task, 1000)

func AddCommand(task database.Task) {
	queue <- task
	log.Printf("Command %d added to the queue", task.ID)
}

func ProcessQueue(db *database.Queries) {
	log.Println("Starting tasks processor")

	go func() {
		for {
			log.Println("Waiting for a task")
			task := <-queue

			log.Printf("Processing command %d of type %s", task.ID, task.Type.String)

			// Input tasks are added by the user optimistically on the client
			// so they should not be broadcasted back to the client
			if task.Type.String != "input" {
				subscriptions.BroadcastTaskAdded(task.FlowID.Int64, &gmodel.Task{
					ID:        uint(task.ID),
					Message:   task.Message.String,
					Type:      gmodel.TaskType(task.Type.String),
					CreatedAt: task.CreatedAt.Time,
					Status:    gmodel.TaskStatus(task.Status.String),
					Args:      string(task.Args),
					Results:   task.Results.String,
				})
			}

			if task.Type.String == "input" {
				err := processInputTask(db, task)

				if err != nil {
					log.Printf("failed to process input: %w", err)
					continue
				}

				nextTask, err := getNextTask(db, task.FlowID.Int64)

				if err != nil {
					log.Printf("failed to get next task: %w", err)
					continue
				}

				AddCommand(*nextTask)
			}

			if task.Type.String == "ask" {
				err := processAskTask(db, task)

				if err != nil {
					log.Printf("failed to process ask: %w", err)
					continue
				}
			}

			if task.Type.String == "terminal" {
				err := processTerminalTask(db, task)

				if err != nil {
					log.Printf("failed to process terminal: %w", err)
					continue
				}
				nextTask, err := getNextTask(db, task.FlowID.Int64)

				if err != nil {
					log.Printf("failed to get next task: %w", err)
					continue
				}

				AddCommand(*nextTask)
			}

			if task.Type.String == "code" {
				err := processCodeTask(db, task)

				if err != nil {
					log.Printf("failed to process code: %w", err)
					continue
				}

				nextTask, err := getNextTask(db, task.FlowID.Int64)

				if err != nil {
					log.Printf("failed to get next task: %w", err)
					continue
				}

				AddCommand(*nextTask)
			}

			if task.Type.String == "done" {
				err := processDoneTask(db, task)

				if err != nil {
					log.Printf("failed to process done: %w", err)
					continue
				}
			}
		}
	}()
}

func processDoneTask(db *database.Queries, task database.Task) error {
	flow, err := db.UpdateFlowStatus(context.Background(), database.UpdateFlowStatusParams{
		ID:     task.FlowID.Int64,
		Status: database.StringToPgText("finished"),
	})

	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	subscriptions.BroadcastFlowUpdated(task.FlowID.Int64, &gmodel.Flow{
		ID:     uint(flow.ID),
		Status: gmodel.FlowStatus("finished"),
	})

	return nil
}

func processInputTask(db *database.Queries, task database.Task) error {
	tasks, err := db.ReadTasksByFlowId(context.Background(), pgtype.Int8{
		Int64: task.FlowID.Int64,
		Valid: true,
	})

	if err != nil {
		return fmt.Errorf("failed to get tasks by flow id: %w", err)
	}

	// This is the first task in the flow.
	// We need to get the basic flow data as well as spin up the container
	if len(tasks) == 1 {
		summary, err := services.GetMessageSummary(task.Message.String, 10)

		if err != nil {
			return fmt.Errorf("failed to get message summary: %w", err)
		}

		dockerImage, err := services.GetDockerImageName(task.Message.String)

		if err != nil {
			return fmt.Errorf("failed to get docker image name: %w", err)
		}

		flow, err := db.UpdateFlowName(context.Background(), database.UpdateFlowNameParams{
			ID:   task.FlowID.Int64,
			Name: database.StringToPgText(summary),
		})

		if err != nil {
			return fmt.Errorf("failed to update flow: %w", err)
		}

		subscriptions.BroadcastFlowUpdated(flow.ID, &gmodel.Flow{
			ID:            uint(flow.ID),
			Name:          summary,
			ContainerName: dockerImage,
		})

		flowId := fmt.Sprint(task.FlowID)
		msg := fmt.Sprintf("Initializing the docker image %s...", dockerImage)
		err = websocket.SendToChannel(flowId, websocket.FormatTerminalSystemOutput(msg))
		if err != nil {
			log.Printf("failed to send message to channel: %w", err)
		}

		containerName := GenerateContainerName(flow.ID)

		containerID, err := SpawnContainer(context.Background(), containerName, dockerImage, db)

		if err != nil {
			return fmt.Errorf("failed to spawn container: %w", err)
		}

		_, err = db.UpdateFlowContainer(context.Background(), database.UpdateFlowContainerParams{
			ID:          flow.ID,
			ContainerID: pgtype.Int8{Int64: containerID, Valid: true},
		})

		if err != nil {
			return fmt.Errorf("failed to update flow container: %w", err)
		}

		err = websocket.SendToChannel(flowId, websocket.FormatTerminalSystemOutput("Container initialized. Ready to execute commands."))
		if err != nil {
			log.Printf("failed to send message to channel: %w", err)
		}
	}

	return nil
}

func processAskTask(db *database.Queries, task database.Task) error {
	task, err := db.UpdateTaskStatus(context.Background(), database.UpdateTaskStatusParams{
		Status: database.StringToPgText("finished"),
		ID:     task.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to find task with id %d: %w", task.ID, err)
	}

	return nil
}

func processTerminalTask(db *database.Queries, task database.Task) error {
	flowId := fmt.Sprint(task.FlowID)
	var args = agent.TerminalArgs{}
	err := json.Unmarshal(task.Args, &args)
	if err != nil {
		return fmt.Errorf("failed to unmarshal args: %v", err)
	}

	// Send the input to the websocket channel
	err = websocket.SendToChannel(flowId, websocket.FormatTerminalInput(args.Input))

	if err != nil {
		log.Printf("failed to send message to channel: %w", err)
	}

	conn, err := websocket.GetConnection(flowId)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	w, err := conn.NextWriter(gorillaWs.BinaryMessage)

	// Write the terminal output to both to the websocket and to the database
	var result = &bytes.Buffer{}
	multi := io.MultiWriter(w, result)

	if err != nil {
		return fmt.Errorf("failed to get writer: %w", err)
	}

	err = ExecCommand(GenerateContainerName(task.FlowID.Int64), args.Input, multi)

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	_, err = db.UpdateTaskResults(context.Background(), database.UpdateTaskResultsParams{
		ID:      task.ID,
		Results: database.StringToPgText(result.String()),
	})

	if err != nil {
		return fmt.Errorf("failed to update task results: %w", err)
	}

	err = w.Close()

	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

func processCodeTask(db *database.Queries, task database.Task) error {
	var args = agent.CodeArgs{}
	err := json.Unmarshal(task.Args, &args)
	if err != nil {
		return fmt.Errorf("failed to unmarshal args: %v", err)
	}

	var cmd = ""
	var r = bytes.Buffer{}

	if args.Action == agent.ReadFile {
		cmd = fmt.Sprintf("cat %s", args.Path)
	}

	if args.Action == agent.UpdateFile {
		cmd = fmt.Sprintf("echo %s > %s", args.Content, args.Path)
	}

	err = ExecCommand(GenerateContainerName(task.FlowID.Int64), cmd, &r)

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	_, err = db.UpdateTaskResults(context.Background(), database.UpdateTaskResultsParams{
		ID:      task.ID,
		Results: database.StringToPgText(r.String()),
	})

	if err != nil {
		return fmt.Errorf("failed to update task results: %w", err)
	}

	return nil
}

func getNextTask(db *database.Queries, flowId int64) (*database.Task, error) {
	flow, err := db.ReadFlow(context.Background(), flowId)

	if err != nil {
		return nil, fmt.Errorf("failed to get flow: %w", err)
	}

	tasks, err := db.ReadTasksByFlowId(context.Background(), pgtype.Int8{
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
			tasks[i].Results = database.StringToPgText(results)
		}
	}

	c, err := agent.NextTask(agent.AgentPrompt{
		Tasks:       tasks,
		DockerImage: flow.ContainerImage.String,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get next command: %w", err)
	}

	nextTask, err := db.CreateTask(context.Background(), database.CreateTaskParams{
		Args:    c.Args,
		Message: c.Message,
		Type:    c.Type,
		Status:  database.StringToPgText("in_progress"),
		FlowID:  pgtype.Int8{Int64: flowId, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to save command: %w", err)
	}

	return &nextTask, nil
}
