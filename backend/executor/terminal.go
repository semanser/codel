package executor

import (
	"archive/tar"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
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
		FlowID:  sql.NullInt64{Int64: flowID, Valid: true},
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
		FlowID:  sql.NullInt64{Int64: flowID, Valid: true},
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

	result = dst.String()

	if result == "" {
		result = "Command executed successfully"
	}

	return result, nil
}

func WriteFile(flowID int64, content string, path string, db *database.Queries) (err error) {
	container := TerminalName(flowID)

	// Check if container is running
	isRunning, err := IsContainerRunning(container)

	if err != nil {
		return fmt.Errorf("Error inspecting container: %w", err)
	}

	if !isRunning {
		return fmt.Errorf("Container is not running")
	}

	// TODO avoid duplicating here and in the flows table
	log, err := db.CreateLog(context.Background(), database.CreateLogParams{
		FlowID:  sql.NullInt64{Int64: flowID, Valid: true},
		Message: content,
		Type:    "input",
	})

	if err != nil {
		return fmt.Errorf("Error creating log: %w", err)
	}

	subscriptions.BroadcastTerminalLogsAdded(flowID, &gmodel.Log{
		ID:   uint(log.ID),
		Text: websocket.FormatTerminalInput(content),
	})

	// Put content into a tar archive
	archive := &bytes.Buffer{}
	tarWriter := tar.NewWriter(archive)
	filename := filepath.Base(path)
	tarHeader := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(content)),
	}
	err = tarWriter.WriteHeader(tarHeader)
	if err != nil {
		return fmt.Errorf("Error writing tar header: %w", err)
	}

	_, err = tarWriter.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("Error writing tar content: %w", err)
	}

	dir := filepath.Dir(path)
	err = dockerClient.CopyToContainer(context.Background(), container, dir, archive, types.CopyToContainerOptions{})

	if err != nil {
		return fmt.Errorf("Error writing file: %w", err)
	}

	message := fmt.Sprintf("Wrote to %s", path)

	// TODO avoid duplicating here and in the flows table
	log, err = db.CreateLog(context.Background(), database.CreateLogParams{
		FlowID:  sql.NullInt64{Int64: flowID, Valid: true},
		Message: message,
		Type:    "output",
	})

	if err != nil {
		return fmt.Errorf("Error creating log: %w", err)
	}

	subscriptions.BroadcastTerminalLogsAdded(flowID, &gmodel.Log{
		ID:   uint(log.ID),
		Text: message,
	})

	return nil
}

func TerminalName(flowID int64) string {
	return fmt.Sprintf("codel-terminal-%d", flowID)
}
