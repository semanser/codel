package executor

import (
	"fmt"
	"log"
	"strings"

	gorillaWs "github.com/gorilla/websocket"
	"github.com/semanser/ai-coder/agent"
	"github.com/semanser/ai-coder/models"
	"github.com/semanser/ai-coder/websocket"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var queue = make(chan *agent.Command, 1000)

func AddCommand(cmd *agent.Command) {
	queue <- cmd
	log.Printf("Command %d added to the queue", cmd.ID)
}

func ProcessQueue(db *gorm.DB) {
	log.Println("Starting tasks processor")

	go func() {
		for {
			log.Println("Waiting for a task")
			cmd := <-queue

			log.Printf("Processing command %d", cmd.ID)
			flowId := fmt.Sprint(cmd.FlowID)

			if cmd.Type == agent.Ask {
				args := cmd.Args.(agent.AskArgs)

				task := models.Task{
					ID:      cmd.ID,
					Type:    models.Action,
					Message: cmd.Description,
					Status:  models.Finished,
					Args:    datatypes.JSON(args.Input),
				}

				tx := db.Save(&task)

				if tx.Error != nil {
					fmt.Errorf("failed to find task with id %d: %w", cmd.ID, tx.Error)
				}

				return
			}

			if cmd.Type == agent.Input {
				flow := models.Flow{}
				tx := db.First(&models.Flow{}, cmd.FlowID).Preload("Tasks").Find(&flow)

				if tx.Error != nil {
					fmt.Errorf("failed to find flow with id %d: %w", cmd.ID, tx.Error)
					return
				}

				var commands []agent.Command

				for _, task := range flow.Tasks {
					commands = append(commands, agent.Command{
						ID:          task.ID,
						FlowID:      task.FlowID,
						Type:        agent.CommandType(task.Type),
						Args:        task.Args,
						Results:     task.Results,
						Description: task.Message,
					})
				}

				c, err := agent.NextCommand(agent.AgentPrompt{
					Commands: commands,
				})

				if err != nil {
					fmt.Errorf("failed to get next command: %w", err)
					return
				}

				task := models.Task{
					ID:      c.ID,
					Message: c.Description,
					Type:    models.Action,
					Status:  models.InProgress,
					FlowID:  c.FlowID,
				}

				tx = db.Save(&task)

				if tx.Error != nil {
					fmt.Errorf("failed to save command: %w", tx.Error)
				}

				log.Printf("Next command: %d", task.ID)
				c.Args = task.ID

				AddCommand(c)
				return
			}

			if cmd.Type != agent.Terminal {
				log.Printf("Command %d is not supported command", cmd.ID)
				return
			}

			args := cmd.Args.(agent.TerminalArgs)

			// Send the input to the websocket channel
			err := websocket.SendToChannel(flowId, websocket.FormatTerminalInput(args.Input))

			if err != nil {
				fmt.Errorf("failed to send message to channel: %w", err)
			}

			conn, err := websocket.GetConnection(flowId)
			if err != nil {
				fmt.Errorf("failed to get connection: %w", err)
				return
			}
			w, err := conn.NextWriter(gorillaWs.BinaryMessage)

			if err != nil {
				fmt.Errorf("failed to get writer: %w", err)
				return
			}

			splitArgs := strings.Split(args.Input, " ")
			err = execCommand("", splitArgs, w)

			if err != nil {
				fmt.Errorf("failed to execute command %d: %w", cmd.ID, err)
				return
			} else {
				log.Printf("Command %d executed successfully", cmd.ID)
			}

			err = w.Close()

			if err != nil {
				fmt.Errorf("failed to send message to channel: %w", err)
				return
			}

			if err != nil {
				fmt.Errorf("failed to execute command %d: %w", cmd.ID, err)
			} else {
				log.Printf("Command %d executed successfully", cmd.ID)
			}
		}
	}()
}
