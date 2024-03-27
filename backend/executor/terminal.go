package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/semanser/ai-coder/database"
	gmodel "github.com/semanser/ai-coder/graph/model"
	"github.com/semanser/ai-coder/graph/subscriptions"
	"github.com/semanser/ai-coder/websocket"
)

func ExecCommand(flowID int64, command string, db *database.Queries) (result string, err error) {
	container := TerminalName(flowID)

	// Create options for starting the exec process
	cmd := []string{
		"sh",
		"-c",
		command,
	}

	// Check if container is running
	isRunning, err := IsContainerRunning(container)

	if err != nil {
		return "", fmt.Errorf("Error inspecting container: %w", err)
	}

	if !isRunning {
		return "", fmt.Errorf("Container is not running")
	}

	// TODO avoid duplicating here and in the flows table
	log, err := db.CreateLog(context.Background(), database.CreateLogParams{
		FlowID:  pgtype.Int8{Int64: flowID, Valid: true},
		Message: command,
		Type:    "input",
	})

	if err != nil {
		return "", fmt.Errorf("Error creating log: %w", err)
	}

	subscriptions.BroadcastTerminalLogsAdded(flowID, &gmodel.Log{
		ID:   uint(log.ID),
		Text: websocket.FormatTerminalInput(command),
	})

	createResp, err := dockerClient.ContainerExecCreate(context.Background(), container, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		return "", fmt.Errorf("Error creating exec process: %w", err)
	}

	// Attach to the exec process
	resp, err := dockerClient.ContainerExecAttach(context.Background(), createResp.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return "", fmt.Errorf("Error attaching to exec process: %w", err)
	}
	defer resp.Close()

	dst := bytes.Buffer{}
	_, err = io.Copy(&dst, resp.Reader)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("Error copying output: %w", err)
	}

	// Wait for the exec process to finish
	_, err = dockerClient.ContainerExecInspect(context.Background(), createResp.ID)
	if err != nil {
		return "", fmt.Errorf("Error inspecting exec process: %w", err)
	}

	results := dst.String()

	// TODO avoid duplicating here and in the flows table
	log, err = db.CreateLog(context.Background(), database.CreateLogParams{
		FlowID:  pgtype.Int8{Int64: flowID, Valid: true},
		Message: results,
		Type:    "output",
	})

	if err != nil {
		return "", fmt.Errorf("Error creating log: %w", err)
	}

	subscriptions.BroadcastTerminalLogsAdded(flowID, &gmodel.Log{
		ID:   uint(log.ID),
		Text: results,
	})

	return dst.String(), nil
}

func TerminalName(flowID int64) string {
	return fmt.Sprintf("codel-terminal-%d", flowID)
}
