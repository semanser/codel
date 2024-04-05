package executor

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
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

const defaultImage = "debian:latest"

func InitClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error initializing docker client: %w", err)
	}
	cli.NegotiateAPIVersion(context.Background())

	dockerClient = cli
	info, err := dockerClient.Info(context.Background())

	if err != nil {
		return fmt.Errorf("error getting docker info: %w", err)
	}

	log.Printf("Docker client initialized: %s, %s", info.Name, info.Architecture)
	log.Printf("Docker server API version: %s", info.ServerVersion)
	log.Printf("Docker client API version: %s", dockerClient.ClientVersion())

	return nil
}

func SpawnContainer(ctx context.Context, name string, config *container.Config, hostConfig *container.HostConfig, db *database.Queries) (dbContainerID int64, err error) {
	if config == nil {
		return 0, fmt.Errorf("no config found for container %s", name)
	}

	log.Printf("Spawning container %s \"%s\"\n", config.Image, name)

	dbContainer, err := db.CreateContainer(ctx, database.CreateContainerParams{
		Name:   database.StringToNullString(name),
		Image:  database.StringToNullString(config.Image),
		Status: database.StringToNullString("starting"),
	})

	if err != nil {
		return dbContainer.ID, fmt.Errorf("error creating container in database: %w", err)
	}

	localContainerID := ""

	defer func() {
		status := "failed"

		if err != nil {
			err := StopContainer(localContainerID, dbContainerID, db)

			if err != nil {
				log.Printf("error stopping failed container %d: %v\n", dbContainerID, err)
			}
		} else {
			status = "running"
		}

		_, err := db.UpdateContainerStatus(ctx, database.UpdateContainerStatusParams{
			ID:     dbContainer.ID,
			Status: database.StringToNullString(status),
		})

		if err != nil {
			log.Printf("error updating container status: %s\n", err)
		}

		_, err = db.UpdateContainerLocalId(ctx, database.UpdateContainerLocalIdParams{
			ID:      dbContainer.ID,
			LocalID: database.StringToNullString(localContainerID),
		})

		if err != nil {
			log.Printf("error updating container local id: %s\n", err)
		}
	}()

	filters := filters.NewArgs()
	filters.Add("reference", config.Image)
	images, err := dockerClient.ImageList(ctx, types.ImageListOptions{
		Filters: filters,
	})

	if err != nil {
		return dbContainer.ID, fmt.Errorf("error listing images: %w", err)
	}

	imageExistsLocally := len(images) > 0

	log.Printf("Image %s found locally: %t\n", config.Image, imageExistsLocally)

	if !imageExistsLocally {
		log.Printf("Pulling image %s...\n", config.Image)
		readCloser, err := dockerClient.ImagePull(ctx, config.Image, types.ImagePullOptions{})

		if err != nil {
			config.Image = defaultImage
			log.Printf("Error pulling image: %s. Using default image %s\n", err, defaultImage)
		}

		if err == nil {
			// Wait for the pull to finish
			_, err = io.Copy(io.Discard, readCloser)

			if err != nil {
				log.Printf("Error waiting for image pull: %s\n", err)
			}
		}
	}

	log.Printf("Creating container %s...\n", name)
	resp, err := dockerClient.ContainerCreate(ctx, config, hostConfig, nil, nil, name)

	if err != nil {
		return dbContainer.ID, fmt.Errorf("error creating container: %w", err)
	}

	log.Printf("Container %s created\n", name)

	localContainerID = resp.ID
	err = dockerClient.ContainerStart(ctx, localContainerID, container.StartOptions{})

	if err != nil {
		return dbContainer.ID, fmt.Errorf("error starting container: %w", err)
	}
	log.Printf("Container %s started\n", name)

	return dbContainer.ID, nil
}

func StopContainer(containerID string, dbID int64, db *database.Queries) error {
	if err := dockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{}); err != nil {
		if client.IsErrNotFound(err) {
			log.Printf("Container %s not found. Marking it as stopped.\n", containerID)
			db.UpdateContainerStatus(context.Background(), database.UpdateContainerStatusParams{
				Status: database.StringToNullString("stopped"),
				ID:     dbID,
			})

			return nil
		} else {
			return fmt.Errorf("error stopping container: %w", err)
		}
	}

	_, err := db.UpdateContainerStatus(context.Background(), database.UpdateContainerStatusParams{
		Status: database.StringToNullString("stopped"),
		ID:     dbID,
	})

	if err != nil {
		return fmt.Errorf("error updating container status to stopped: %w", err)
	}

	log.Printf("Container %s stopped\n", containerID)
	return nil
}

func DeleteContainer(containerID string, dbID int64, db *database.Queries) error {
	log.Printf("Deleting container %s...\n", containerID)

	if err := StopContainer(containerID, dbID, db); err != nil {
		return fmt.Errorf("error stopping container: %w", err)
	}

	if err := dockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("error removing container: %w", err)
	}
	log.Printf("Container %s removed\n", containerID)
	return nil
}

func Cleanup(db *database.Queries) error {
	// Remove tmp files
	log.Println("Removing tmp files...")
	err := os.RemoveAll("./tmp/")
	if err != nil {
		return fmt.Errorf("error removing tmp files: %w", err)
	}

	log.Println("Cleaning up containers and making all flows finished...")

	var wg sync.WaitGroup

	containers, err := db.GetAllRunningContainers(context.Background())

	if err != nil {
		return fmt.Errorf("error getting running containers: %w", err)
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

	flows, err := db.ReadAllFlows(context.Background())

	if err != nil {
		return fmt.Errorf("error getting all flows: %w", err)
	}

	for _, flow := range flows {
		if flow.Status.String == "in_progress" {
			_, err := db.UpdateFlowStatus(context.Background(), database.UpdateFlowStatusParams{
				Status: database.StringToNullString("finished"),
				ID:     flow.ID,
			})

			if err != nil {
				log.Printf("Error updating flow status: %s\n", err)
			}
		}
	}

	return nil
}

func IsContainerRunning(containerID string) (bool, error) {
	containerInfo, err := dockerClient.ContainerInspect(context.Background(), containerID)

	if err != nil {
		return false, fmt.Errorf("error inspecting container: %w", err)
	}

	return containerInfo.State.Running, err
}
