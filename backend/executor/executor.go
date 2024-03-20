package executor

import (
	"context"
	"log"

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
