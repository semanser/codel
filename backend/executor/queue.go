package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	gorillaWs "github.com/gorilla/websocket"
	"github.com/semanser/ai-coder/agent"
	"github.com/semanser/ai-coder/models"
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

			if task.Type == models.Ask {
				err := processAsk(db, task)

				if err != nil {
					log.Printf("failed to process ask: %w", err)
				}
			}

			if task.Type == models.Input {
				nextTask, err := processInput(db, task)

				if err != nil {
					log.Printf("failed to process input: %w", err)
				}

				AddCommand(*nextTask)
			}

			if task.Type == models.Terminal {
				nextTask, err := processTerminal(db, task)

				if err != nil {
					log.Printf("failed to process terminal: %w", err)
				}

				AddCommand(*nextTask)
			}
		}
	}()
}

func processAsk(db *gorm.DB, task models.Task) error {
	// TODO Send the subscription with the ask to the client
	tx := db.Updates(models.Task{
		ID:     task.ID,
		Status: models.Finished,
	})

	if tx.Error != nil {
		return fmt.Errorf("failed to find task with id %d: %w", task.ID, tx.Error)
	}

	return nil
}

func processInput(db *gorm.DB, task models.Task) (nextTask *models.Task, err error) {
	flow := models.Flow{}
	tx := db.First(&models.Flow{}, task.FlowID).Preload("Tasks").Find(&flow)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to find flow with id %d: %w", task.ID, tx.Error)
	}

	c, err := agent.NextTask(agent.AgentPrompt{
		Tasks: flow.Tasks,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get next command: %w", err)
	}

	nextTask = &models.Task{
		Args:    c.Args,
		Message: c.Message,
		Type:    c.Type,
		Status:  models.Finished,
		FlowID:  task.FlowID,
	}

	tx = db.Save(nextTask)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to save command: %w", tx.Error)
	}

	return nextTask, nil
}

func processTerminal(db *gorm.DB, task models.Task) (nextTask *models.Task, err error) {
	flowId := fmt.Sprint(task.FlowID)
	var args = agent.TerminalArgs{}
	err = json.Unmarshal(task.Args, &args)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal args: %v", err)
	}

	// Send the input to the websocket channel
	err = websocket.SendToChannel(flowId, websocket.FormatTerminalInput(args.Input))

	if err != nil {
		log.Printf("failed to send message to channel: %w", err)
	}

	conn, err := websocket.GetConnection(flowId)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	w, err := conn.NextWriter(gorillaWs.BinaryMessage)

	// Write the terminal output to both to the websocket and to the database
	result := new(bytes.Buffer)
	multi := io.MultiWriter(w, result)

	if err != nil {
		return nil, fmt.Errorf("failed to get writer: %w", err)
	}

	err = ExecCommand(GenerateContainerName(task.FlowID), args.Input, multi)

	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	// Mark the current task as finished and save the output
	db.Updates(models.Task{
		ID:      task.ID,
		Results: result.Bytes(),
		Status:  models.Finished,
	})

	err = w.Close()

	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	flow := models.Flow{}
	tx := db.First(&models.Flow{}, task.FlowID).Preload("Tasks").Find(&flow)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to find flow with id %d: %w", task.ID, tx.Error)
	}

	c, err := agent.NextTask(agent.AgentPrompt{
		Tasks: flow.Tasks,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get next command: %w", err)
	}

	nextTask = &models.Task{
		Args:    c.Args,
		Message: c.Message,
		Type:    c.Type,
		Status:  models.InProgress,
		FlowID:  task.FlowID,
	}

	tx = db.Save(nextTask)

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to save command: %w", tx.Error)
	}

	return nextTask, nil
}
