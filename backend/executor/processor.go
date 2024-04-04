package executor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/semanser/ai-coder/database"
	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/graph/subscriptions"
	"github.com/semanser/ai-coder/providers"
	"github.com/semanser/ai-coder/websocket"
)

func processBrowserTask(db *database.Queries, task database.Task) error {
	var args = providers.BrowserArgs{}
	err := json.Unmarshal([]byte(task.Args.String), &args)
	if err != nil {
		return fmt.Errorf("failed to unmarshal args: %v", err)
	}

	var url = args.Url
	var screenshotName string

	if args.Action == providers.Read {
		content, screenshot, err := Content(url)

		if err != nil {
			return fmt.Errorf("failed to get content: %w", err)
		}

		log.Println("Screenshot taken")
		screenshotName = screenshot

		_, err = db.UpdateTaskResults(context.Background(), database.UpdateTaskResultsParams{
			ID:      task.ID,
			Results: database.StringToNullString(content),
		})

		if err != nil {
			return fmt.Errorf("failed to update task results: %w", err)
		}
	}

	if args.Action == providers.Url {
		content, screenshot, err := URLs(url)

		if err != nil {
			return fmt.Errorf("failed to get content: %w", err)
		}

		screenshotName = screenshot

		_, err = db.UpdateTaskResults(context.Background(), database.UpdateTaskResultsParams{
			ID:      task.ID,
			Results: database.StringToNullString(content),
		})

		if err != nil {
			return fmt.Errorf("failed to update task results: %w", err)
		}
	}

	subscriptions.BroadcastBrowserUpdated(task.FlowID.Int64, &gmodel.Browser{
		URL: url,
		// TODO Use a dynamic URL
		ScreenshotURL: "http://localhost:8080/browser/" + screenshotName,
	})

	return nil
}

func processDoneTask(db *database.Queries, task database.Task) error {
	flow, err := db.UpdateFlowStatus(context.Background(), database.UpdateFlowStatusParams{
		ID:     task.FlowID.Int64,
		Status: database.StringToNullString("finished"),
	})

	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	subscriptions.BroadcastFlowUpdated(task.FlowID.Int64, &gmodel.Flow{
		ID:       uint(flow.ID),
		Status:   gmodel.FlowStatus("finished"),
		Terminal: &gmodel.Terminal{},
	})

	return nil
}

func processInputTask(provider providers.Provider, db *database.Queries, task database.Task) error {
	tasks, err := db.ReadTasksByFlowId(context.Background(), sql.NullInt64{
		Int64: task.FlowID.Int64,
		Valid: true,
	})

	if err != nil {
		return fmt.Errorf("failed to get tasks by flow id: %w", err)
	}

	// This is the first task in the flow.
	// We need to get the basic flow data as well as spin up the container
	if len(tasks) == 1 {
		summary, err := provider.Summary(task.Message.String, 10)

		if err != nil {
			return fmt.Errorf("failed to get message summary: %w", err)
		}

		dockerImage, err := provider.DockerImageName(task.Message.String)

		if err != nil {
			return fmt.Errorf("failed to get docker image name: %w", err)
		}

		flow, err := db.UpdateFlowName(context.Background(), database.UpdateFlowNameParams{
			ID:   task.FlowID.Int64,
			Name: database.StringToNullString(summary),
		})

		if err != nil {
			return fmt.Errorf("failed to update flow: %w", err)
		}

		subscriptions.BroadcastFlowUpdated(flow.ID, &gmodel.Flow{
			ID:   uint(flow.ID),
			Name: summary,
			Terminal: &gmodel.Terminal{
				ContainerName: dockerImage,
				Connected:     false,
			},
		})

		msg := websocket.FormatTerminalSystemOutput(fmt.Sprintf("Initializing the docker image %s...", dockerImage))
		l, err := db.CreateLog(context.Background(), database.CreateLogParams{
			FlowID:  task.FlowID,
			Message: msg,
			Type:    "system",
		})

		if err != nil {
			return fmt.Errorf("error creating log: %w", err)
		}

		subscriptions.BroadcastTerminalLogsAdded(flow.ID, &gmodel.Log{
			ID:   uint(l.ID),
			Text: msg,
		})

		terminalContainerName := TerminalName(flow.ID)
		terminalContainerID, err := SpawnContainer(context.Background(),
			terminalContainerName,
			&container.Config{
				Image: dockerImage,
				Cmd:   []string{"tail", "-f", "/dev/null"},
			},
			&container.HostConfig{},
			db,
		)

		if err != nil {
			return fmt.Errorf("failed to spawn container: %w", err)
		}

		subscriptions.BroadcastFlowUpdated(flow.ID, &gmodel.Flow{
			ID:   uint(flow.ID),
			Name: summary,
			Terminal: &gmodel.Terminal{
				Connected:     true,
				ContainerName: dockerImage,
			},
		})

		_, err = db.UpdateFlowContainer(context.Background(), database.UpdateFlowContainerParams{
			ID:          flow.ID,
			ContainerID: sql.NullInt64{Int64: terminalContainerID, Valid: true},
		})

		if err != nil {
			return fmt.Errorf("failed to update flow container: %w", err)
		}

		msg = websocket.FormatTerminalSystemOutput("Container initialized. Ready to execute commands.")
		l, err = db.CreateLog(context.Background(), database.CreateLogParams{
			FlowID:  task.FlowID,
			Message: msg,
			Type:    "system",
		})

		if err != nil {
			return fmt.Errorf("error creating log: %w", err)
		}
		subscriptions.BroadcastTerminalLogsAdded(flow.ID, &gmodel.Log{
			ID:   uint(l.ID),
			Text: msg,
		})
	}

	return nil
}

func processAskTask(db *database.Queries, task database.Task) error {
	task, err := db.UpdateTaskStatus(context.Background(), database.UpdateTaskStatusParams{
		Status: database.StringToNullString("finished"),
		ID:     task.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to find task with id %d: %w", task.ID, err)
	}

	return nil
}

func processTerminalTask(db *database.Queries, task database.Task) error {
	var args = providers.TerminalArgs{}
	err := json.Unmarshal([]byte(task.Args.String), &args)
	if err != nil {
		return fmt.Errorf("failed to unmarshal args: %v", err)
	}

	results, err := ExecCommand(task.FlowID.Int64, args.Input, db)

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	_, err = db.UpdateTaskResults(context.Background(), database.UpdateTaskResultsParams{
		ID:      task.ID,
		Results: database.StringToNullString(results),
	})

	if err != nil {
		return fmt.Errorf("failed to update task results: %w", err)
	}

	return nil
}

func processCodeTask(db *database.Queries, task database.Task) error {
	var args = providers.CodeArgs{}
	err := json.Unmarshal([]byte(task.Args.String), &args)
	if err != nil {
		return fmt.Errorf("failed to unmarshal args: %v", err)
	}

	var cmd = ""
	var results = ""

	if args.Action == providers.ReadFile {
		// TODO consider using dockerClient.CopyFromContainer command instead
		cmd = fmt.Sprintf("cat %s", args.Path)
		results, err = ExecCommand(task.FlowID.Int64, cmd, db)

		if err != nil {
			return fmt.Errorf("error executing cat command: %w", err)
		}
	}

	if args.Action == providers.UpdateFile {
		err = WriteFile(task.FlowID.Int64, args.Content, args.Path, db)

		if err != nil {
			return fmt.Errorf("error writing a file: %w", err)
		}

		results = "File updated"
	}

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	_, err = db.UpdateTaskResults(context.Background(), database.UpdateTaskResultsParams{
		ID:      task.ID,
		Results: database.StringToNullString(results),
	})

	if err != nil {
		return fmt.Errorf("failed to update task results: %w", err)
	}

	return nil
}
