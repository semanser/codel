package executor

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/semanser/ai-coder/database"
)

var (
	dockerClient *client.Client
)

func InitDockerClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	dockerClient = cli
	info, err := dockerClient.Info(context.Background())

	if err != nil {
		return err
	}

	log.Printf("Docker client initialized: %s", info.Name)

	return nil
}

func SpawnContainer(ctx context.Context, name string, dockerImage string, db *database.Queries) (containerID string, err error) {
	log.Printf("Spawning container %s \"%s\"\n", dockerImage, name)

	dbContainer, err := db.CreateContainer(ctx, database.CreateContainerParams{
		Name:   database.StringToPgText(name),
		Image:  database.StringToPgText(dockerImage),
		Status: database.StringToPgText("starting"),
	})

	if err != nil {
		return "", fmt.Errorf("Error creating container in database: %w", err)
	}

	defer func() {
		status := "failed"

		if err != nil {
			err := StopContainer(containerID)

			if err != nil {
				log.Printf("Error stopping failed container %s: %s\n", containerID, err)
			}
		} else {
			status = "running"
		}

		db.UpdateContainerStatus(ctx, database.UpdateContainerStatusParams{
			ID:     dbContainer.ID,
			Status: database.StringToPgText(status),
		})
	}()

	filters := filters.NewArgs()
	filters.Add("reference", dockerImage)
	images, err := dockerClient.ImageList(ctx, types.ImageListOptions{
		Filters: filters,
	})

	if err != nil {
		return "", fmt.Errorf("Error listing images: %w", err)
	}

	imageFound := len(images) > 0

	log.Printf("Image %s found: %t\n", dockerImage, imageFound)

	if !imageFound {
		log.Printf("Pulling image %s...\n", dockerImage)
		readCloser, err := dockerClient.ImagePull(ctx, dockerImage, types.ImagePullOptions{})

		if err != nil {
			return "", fmt.Errorf("Error pulling image: %w", err)
		}

		// Wait for the pull to finish
		_, err = io.Copy(io.Discard, readCloser)

		if err != nil {
			return "", fmt.Errorf("Error waiting for image pull: %w", err)
		}
	}

	log.Printf("Creating container %s...\n", name)
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: dockerImage,
		Cmd:   []string{"tail", "-f", "/dev/null"},
	}, nil, nil, nil, name)

	if err != nil {
		return "", fmt.Errorf("Error creating container: %w", err)
	}

	log.Printf("Container %s created\n", name)

	containerID = resp.ID
	err = dockerClient.ContainerStart(ctx, containerID, container.StartOptions{})

	if err != nil {
		return "", fmt.Errorf("Error starting container: %w", err)
	}
	log.Printf("Container %s started\n", name)

	return containerID, nil
}

func StopContainer(containerID string) error {
	if err := dockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		return err
	}
	log.Printf("Container %s stopped\n", containerID)
	return nil
}

func DeleteContainer(containerID string) error {
	log.Printf("Deleting container %s...\n", containerID)

	if err := StopContainer(containerID); err != nil {
		return err
	}

	if err := dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		return err
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
			id := container.LocalID.String
			if err := DeleteContainer(id); err != nil {
				log.Printf("Error deleting container %s: %s\n", id, err)
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

func ExecCommand(container string, command string, dst io.Writer) (err error) {
	// Create options for starting the exec process
	cmd := []string{
		"sh",
		"-c",
		command,
	}

	// Check if container is running
	isRunning, err := IsContainerRunning(container)

	if err != nil {
		return fmt.Errorf("Error inspecting container: %w", err)
	}

	if !isRunning {
		return fmt.Errorf("Container is not running")
	}

	createResp, err := dockerClient.ContainerExecCreate(context.Background(), container, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		return fmt.Errorf("Error creating exec process: %w", err)
	}

	// Attach to the exec process
	resp, err := dockerClient.ContainerExecAttach(context.Background(), createResp.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return fmt.Errorf("Error attaching to exec process: %w", err)
	}
	defer resp.Close()

	_, err = io.Copy(dst, resp.Reader)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Error copying output: %w", err)
	}

	// Wait for the exec process to finish
	_, err = dockerClient.ContainerExecInspect(context.Background(), createResp.ID)
	if err != nil {
		return fmt.Errorf("Error inspecting exec process: %w", err)
	}

	return nil
}

func GenerateContainerName(flowID int64) string {
	return fmt.Sprintf("flow-%d", flowID)
}
