package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/3ssalunke/leetcode-clone-exen/docker"
)

type ExecutionPayload struct {
	ProblemId string `json:"problem_id"`
	Lang      string `json:"lang"`
	TypedCode string `json:"typed_code"`
}

func ExecuteCode(payload *ExecutionPayload) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return "", fmt.Errorf("failed to create docker client - %w", err)
	}
	defer dockerClient.Client.Close()

	dockerImageTag := fmt.Sprintf("%s-docker-img", payload.Lang)

	image, err := dockerClient.CreateDockerImage(ctx, payload.Lang, dockerImageTag)
	if err != nil {
		return "", fmt.Errorf("failed to create docker image - %w", err)
	}
	defer image.Body.Close()
	log.Println("docker image created successfully.")

	container, err := dockerClient.RunDockerContainer(ctx, dockerImageTag)
	if err != nil {
		return "", fmt.Errorf("failed to run docker container - %w", err)
	}
	log.Println("docker container started successfully.")

	logs, err := dockerClient.GetContainerLogs(ctx, container.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get docker container logs - %w", err)
	}
	defer logs.Close()

	var logBuffer strings.Builder
	io.Copy(&logBuffer, logs)

	if err = dockerClient.RemoveDockerContainer(ctx, container.ID); err != nil {
		return "", fmt.Errorf("failed to remove docker container - %w", err)
	}
	log.Println("container removed successfully.")

	if err = dockerClient.RemoveDockerImage(ctx, dockerImageTag); err != nil {
		return "", fmt.Errorf("failed to remove docker image - %w", err)
	}

	log.Println("image removed successfully.")

	return logBuffer.String(), err
}