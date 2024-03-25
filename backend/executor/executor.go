package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/semanser/ai-coder/database"
)

var (
	dockerClient *client.Client
)

func InitDockerClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("Error initializing docker client: %w", err)
	}
	dockerClient = cli
	info, err := dockerClient.Info(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting docker info: %w", err)
	}

	log.Printf("Docker client initialized: %s", info.Name)

	return nil
}

func SpawnContainer(ctx context.Context, name string, dockerImage string, db *database.Queries) (dbContainerID int64, err error) {
	log.Printf("Spawning container %s \"%s\"\n", dockerImage, name)

	dbContainer, err := db.CreateContainer(ctx, database.CreateContainerParams{
		Name:   database.StringToPgText(name),
		Image:  database.StringToPgText(dockerImage),
		Status: database.StringToPgText("starting"),
	})

	if err != nil {
		return dbContainer.ID, fmt.Errorf("Error creating container in database: %w", err)
	}

	localContainerID := ""

	defer func() {
		status := "failed"

		if err != nil {
			err := StopContainer(localContainerID, dbContainerID, db)

			if err != nil {
				log.Printf("Error stopping failed container %s: %s\n", dbContainerID, err)
			}
		} else {
			status = "running"
		}

		_, err := db.UpdateContainerStatus(ctx, database.UpdateContainerStatusParams{
			ID:     dbContainer.ID,
			Status: database.StringToPgText(status),
		})

		if err != nil {
			log.Printf("Error updating container status: %s\n", err)
		}

		_, err = db.UpdateContainerLocalId(ctx, database.UpdateContainerLocalIdParams{
			ID:      dbContainer.ID,
			LocalID: database.StringToPgText(localContainerID),
		})
	}()

	filters := filters.NewArgs()
	filters.Add("reference", dockerImage)
	images, err := dockerClient.ImageList(ctx, types.ImageListOptions{
		Filters: filters,
	})

	if err != nil {
		return dbContainer.ID, fmt.Errorf("Error listing images: %w", err)
	}

	imageFound := len(images) > 0

	log.Printf("Image %s found: %t\n", dockerImage, imageFound)

	if !imageFound {
		log.Printf("Pulling image %s...\n", dockerImage)
		readCloser, err := dockerClient.ImagePull(ctx, dockerImage, types.ImagePullOptions{})

		if err != nil {
			return dbContainer.ID, fmt.Errorf("Error pulling image: %w", err)
		}

		// Wait for the pull to finish
		_, err = io.Copy(io.Discard, readCloser)

		if err != nil {
			return dbContainer.ID, fmt.Errorf("Error waiting for image pull: %w", err)
		}
	}

	log.Printf("Creating container %s...\n", name)
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: dockerImage,
		Cmd:   []string{"tail", "-f", "/dev/null"},
	}, nil, nil, nil, name)

	if err != nil {
		return dbContainer.ID, fmt.Errorf("Error creating container: %w", err)
	}

	log.Printf("Container %s created\n", name)

	localContainerID = resp.ID
	err = dockerClient.ContainerStart(ctx, localContainerID, container.StartOptions{})

	if err != nil {
		return dbContainer.ID, fmt.Errorf("Error starting container: %w", err)
	}
	log.Printf("Container %s started\n", name)

	return dbContainer.ID, nil
}

func StopContainer(containerID string, dbID int64, db *database.Queries) error {
	if err := dockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		if client.IsErrNotFound(err) {
			log.Printf("Container %s not found. Marking it as stopped.\n", containerID)
			db.UpdateContainerStatus(context.Background(), database.UpdateContainerStatusParams{
				Status: database.StringToPgText("stopped"),
				ID:     dbID,
			})

			return nil
		} else {
			return fmt.Errorf("Error stopping container: %w", err)
		}
	}

	_, err := db.UpdateContainerStatus(context.Background(), database.UpdateContainerStatusParams{
		Status: database.StringToPgText("stopped"),
		ID:     dbID,
	})

	if err != nil {
		return fmt.Errorf("Error updating container status to stopped: %w", err)
	}

	log.Printf("Container %s stopped\n", containerID)
	return nil
}

func DeleteContainer(containerID string, dbID int64, db *database.Queries) error {
	log.Printf("Deleting container %s...\n", containerID)

	if err := StopContainer(containerID, dbID, db); err != nil {
		return fmt.Errorf("Error stopping container: %w", err)
	}

	if err := dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("Error removing container: %w", err)
	}
	log.Printf("Container %s removed\n", containerID)
	return nil
}

func Cleanup(db *database.Queries) error {
	log.Println("Cleaning up containers")

	var wg sync.WaitGroup

	containers, err := db.GetAllRunningContainers(context.Background())

	if err != nil {
		return fmt.Errorf("Error getting running containers: %w", err)
	}

	for _, container := range containers {
		wg.Add(1)
		go func() {
			localId := container.LocalID.String
			if err := DeleteContainer(localId, container.ID, db); err != nil {
				log.Printf("Error deleting container %s: %s\n", localId, err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return nil
}

func IsContainerRunning(containerID string) (bool, error) {
	containerInfo, err := dockerClient.ContainerInspect(context.Background(), containerID)
	return containerInfo.State.Running, err
}

func ExecCommand(id int64, command string, db *database.Queries) (result string, err error) {
	container := GenerateContainerName(id)

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
	db.CreateLog(context.Background(), database.CreateLogParams{
		FlowID:  pgtype.Int8{Int64: id, Valid: true},
		Message: command,
		Type:    "input",
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
	db.CreateLog(context.Background(), database.CreateLogParams{
		FlowID:  pgtype.Int8{Int64: id, Valid: true},
		Message: results,
		Type:    "output",
	})

	return dst.String(), nil
}

func GenerateContainerName(flowID int64) string {
	return fmt.Sprintf("flow-%d", flowID)
}
