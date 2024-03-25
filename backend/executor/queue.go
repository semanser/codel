package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	gorillaWs "github.com/gorilla/websocket"
	"github.com/semanser/ai-coder/agent"
	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/graph/subscriptions"
	"github.com/semanser/ai-coder/models"
	"github.com/semanser/ai-coder/services"
	"github.com/semanser/ai-coder/websocket"
	"gorm.io/gorm"
)

var queue = make(chan models.Task, 1000)

func AddCommand(task models.Task) {
	queue <- task
	log.Printf("Command %d added to the queue", task.ID)
}

func ProcessQueue(db *gorm.DB) {
	log.Println("Starting tasks processor")

	go func() {
		for {
			log.Println("Waiting for a task")
			task := <-queue

			log.Printf("Processing command %d of type %s", task.ID, task.Type)

			// Input tasks are added by the user optimistically on the client
			// so they should not be broadcasted back to the client
			if task.Type != models.Input {
				subscriptions.BroadcastTaskAdded(task.FlowID, &gmodel.Task{
					ID:        task.ID,
					Message:   task.Message,
					Type:      gmodel.TaskType(task.Type),
					CreatedAt: task.CreatedAt,
					Status:    gmodel.TaskStatus(task.Status),
					Args:      task.Args.String(),
					Results:   task.Results,
				})
			}

			if task.Type == models.Input {
				err := processInputTask(db, task)

				if err != nil {
					log.Printf("failed to process input: %w", err)
					continue
				}

				nextTask, err := getNextTask(db, task.FlowID)

				if err != nil {
					log.Printf("failed to get next task: %w", err)
					continue
				}

				AddCommand(*nextTask)
			}

			if task.Type == models.Ask {
				err := processAskTask(db, task)

				if err != nil {
					log.Printf("failed to process ask: %w", err)
					continue
				}
			}

			if task.Type == models.Terminal {
				err := processTerminalTask(db, task)

				if err != nil {
					log.Printf("failed to process terminal: %w", err)
					continue
				}
				nextTask, err := getNextTask(db, task.FlowID)

				if err != nil {
					log.Printf("failed to get next task: %w", err)
					continue
				}

				AddCommand(*nextTask)
			}

			if task.Type == models.Code {
				err := processCodeTask(db, task)

				if err != nil {
					log.Printf("failed to process code: %w", err)
					continue
				}

				nextTask, err := getNextTask(db, task.FlowID)

				if err != nil {
					log.Printf("failed to get next task: %w", err)
					continue
				}

				AddCommand(*nextTask)
			}
		}
	}()
}

func processInputTask(db *gorm.DB, task models.Task) error {
	flow := &models.Flow{
		ID: task.FlowID,
	}
	tx := db.Preload("Tasks").First(flow)

	if tx.Error != nil {
		return fmt.Errorf("failed to fetch flow: %w", tx.Error)
	}

	// This is the first task in the flow.
	// We need to get the basic flow data as well as spin up the container
	if len(flow.Tasks) == 1 {
		summary, err := services.GetMessageSummary(task.Message, 10)

		if err != nil {
			return fmt.Errorf("failed to get message summary: %w", err)
		}

		dockerImage, err := services.GetDockerImageName(task.Message)

		if err != nil {
			return fmt.Errorf("failed to get docker image name: %w", err)
		}

		tx := db.Updates(models.Flow{
			ID:          flow.ID,
			Name:        summary,
			DockerImage: dockerImage,
		})

		if tx.Error != nil {
			return fmt.Errorf("failed to update flow: %w", tx.Error)
		}

		subscriptions.BroadcastFlowUpdated(flow.ID, &gmodel.Flow{
			ID:            flow.ID,
			Name:          summary,
			ContainerName: dockerImage,
		})

		flowId := fmt.Sprint(task.FlowID)
		msg := fmt.Sprintf("Initializing the docker image %s...", dockerImage)
		err = websocket.SendToChannel(flowId, websocket.FormatTerminalSystemOutput(msg))
		if err != nil {
			log.Printf("failed to send message to channel: %w", err)
		}
		_, err = SpawnContainer(context.Background(), GenerateContainerName(flow.ID), dockerImage)

		err = websocket.SendToChannel(flowId, websocket.FormatTerminalSystemOutput("Container initialized. Ready to execute commands."))
		if err != nil {
			log.Printf("failed to send message to channel: %w", err)
		}

		if err != nil {
			return fmt.Errorf("failed to spawn container: %w", err)
		}
	}

	return nil
}

func processAskTask(db *gorm.DB, task models.Task) error {
	tx := db.Updates(models.Task{
		ID:     task.ID,
		Status: models.Finished,
	})

	if tx.Error != nil {
		return fmt.Errorf("failed to find task with id %d: %w", task.ID, tx.Error)
	}

	return nil
}

func processTerminalTask(db *gorm.DB, task models.Task) error {
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

	err = ExecCommand(GenerateContainerName(task.FlowID), args.Input, multi)

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	db.Updates(models.Task{
		ID:      task.ID,
		Results: result.String(),
		Status:  models.Finished,
	})

	err = w.Close()

	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

func processCodeTask(db *gorm.DB, task models.Task) error {
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

	err = ExecCommand(GenerateContainerName(task.FlowID), cmd, &r)

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	db.Updates(models.Task{
		ID:      task.ID,
		Results: r.String(),
	})

	return nil
}

func getNextTask(db *gorm.DB, flowId uint) (*models.Task, error) {
	flow := models.Flow{}
	tx := db.First(&models.Flow{}, flowId).Preload("Tasks").Find(&flow)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to find flow with id %d: %w", flowId, tx.Error)
	}

	const maxResultsLength = 4000
	for i, task := range flow.Tasks {
		// Limit the number of result characters since some output commands can have a lot of output
		if len(task.Results) > maxResultsLength {
			// Get the last N symbols from the output
			flow.Tasks[i].Results = flow.Tasks[i].Results[len(flow.Tasks[i].Results)-maxResultsLength:]
		}
	}

	c, err := agent.NextTask(agent.AgentPrompt{
		Tasks:       flow.Tasks,
		DockerImage: flow.DockerImage,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get next command: %w", err)
	}

	nextTask := &models.Task{
		Args:    c.Args,
		Message: c.Message,
		Type:    c.Type,
		Status:  models.InProgress,
		FlowID:  flowId,
	}

	tx = db.Save(nextTask)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to save command: %w", tx.Error)
	}

	return nextTask, nil
}
