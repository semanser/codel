package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	gorillaWs "github.com/gorilla/websocket"
	"github.com/semanser/ai-coder/agent"
	"github.com/semanser/ai-coder/models"
	"github.com/semanser/ai-coder/websocket"
	"gorm.io/gorm"
)

var queue = make(chan models.Task, 1000)

func AddCommand(cmd models.Task) {
	queue <- cmd
	log.Printf("Command %d added to the queue", cmd.ID)
}

func ProcessQueue(db *gorm.DB) {
	log.Println("Starting tasks processor")

	go func() {
		for {
			log.Println("Waiting for a task")
			cmd := <-queue

			log.Printf("Processing command %d of type %s", cmd.ID, cmd.Type)
			flowId := fmt.Sprint(cmd.FlowID)

			if cmd.Type == models.Ask {
				// TODO Send the subscription with the ask to the client

				task := models.Task{
					ID:     cmd.ID,
					Status: models.Finished,
				}

				tx := db.Updates(task)

				if tx.Error != nil {
					log.Printf("failed to find task with id %d: %w", cmd.ID, tx.Error)
				}

				continue
			}

			if cmd.Type == models.Input {
				flow := models.Flow{}
				tx := db.First(&models.Flow{}, cmd.FlowID).Preload("Tasks").Find(&flow)

				if tx.Error != nil {
					log.Printf("failed to find flow with id %d: %w", cmd.ID, tx.Error)
					continue
				}

				c, err := agent.NextTask(agent.AgentPrompt{
					Tasks: flow.Tasks,
				})

				if err != nil {
					log.Printf("failed to get next command: %w", err)
					continue
				}

				nextTask := models.Task{
					Args:    c.Args,
					Message: c.Message,
					Type:    c.Type,
					Status:  models.Finished,
					FlowID:  cmd.FlowID,
				}

				tx = db.Save(&nextTask)

				if tx.Error != nil {
					log.Printf("failed to save command: %w", tx.Error)
					continue
				}

				log.Printf("The next command is %d", nextTask.ID)
				AddCommand(nextTask)
				continue
			}

			if cmd.Type == models.Terminal {
				var args = agent.TerminalArgs{}
				err := json.Unmarshal([]byte(cmd.Args), &args)
				if err != nil {
					log.Printf("failed to unmarshal args: %v", err)
					continue
				}

				// Send the input to the websocket channel
				log.Printf("Sending input to the websocket channel")
				log.Printf("The input is %s", args.Input)
				err = websocket.SendToChannel(flowId, websocket.FormatTerminalInput(args.Input))

				if err != nil {
					log.Printf("failed to send message to channel: %w", err)
				}

				conn, err := websocket.GetConnection(flowId)
				if err != nil {
					log.Printf("failed to get connection: %w", err)
					continue
				}
				w, err := conn.NextWriter(gorillaWs.BinaryMessage)

				if err != nil {
					log.Printf("failed to get writer: %w", err)
					continue
				}

				splitArgs := strings.Split(args.Input, " ")
				err = execCommand("", splitArgs, w)

				if err != nil {
					log.Printf("failed to execute command %d: %w", cmd.ID, err)
					continue
				} else {
					log.Printf("Command %d executed successfully", cmd.ID)
				}

				err = w.Close()

				if err != nil {
					log.Printf("failed to send message to channel: %w", err)
					continue
				}

				if err != nil {
					log.Printf("failed to execute command %d: %w", cmd.ID, err)
				} else {
					log.Printf("Command %d executed successfully", cmd.ID)
				}
			}
		}
	}()
}
