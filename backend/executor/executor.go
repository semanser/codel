package executor

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var (
	dockerClient *client.Client
	containers   []string
)

const imageName = "alpine"

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

func SpawnContainer(name string) (containerID string, err error) {
	log.Printf("Spawning container %s\n", name)

	resp, err := dockerClient.ContainerCreate(context.Background(), &container.Config{
		Image: imageName,
		Cmd:   []string{"tail", "-f", "/dev/null"},
	}, nil, nil, nil, name)

	if err != nil {
		return "", err
	}
	log.Printf("Container %s created\n", name)

	containerID = resp.ID
	if err := dockerClient.ContainerStart(context.Background(), containerID, container.StartOptions{}); err != nil {
		return "", err
	}
	log.Printf("Container %s started\n", name)

	containers = append(containers, containerID)
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
	if err := StopContainer(containerID); err != nil {
		return err
	}

	if err := dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		return err
	}
	log.Printf("Container %s removed\n", containerID)
	return nil
}

func Cleanup() error {
	log.Println("Cleaning up containers")

	for _, containerID := range containers {
		if err := DeleteContainer(containerID); err != nil {
			return err
		}
	}
	return nil
}

func execCommand(container string, cmd []string, dst io.Writer) (err error) {
	// Create options for starting the exec process
	createResp, err := dockerClient.ContainerExecCreate(context.Background(), container, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return fmt.Errorf("Error creating exec process: %w", err)
	}

	// Attach to the exec process
	resp, err := dockerClient.ContainerExecAttach(context.Background(), createResp.ID, types.ExecStartCheck{})
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

func GenerateContainerName(flowID uint) string {
	return fmt.Sprintf("flow-%d", flowID)
}
