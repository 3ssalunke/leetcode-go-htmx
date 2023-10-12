package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/3ssalunke/leetcode-clone-exen/docker"
	"github.com/3ssalunke/leetcode-clone-exen/util"
)

type ProblemDetails struct {
	ExecutionId  string   `json:"execution_id"`
	ProblemId    string   `json:"problem_id"`
	Lang         string   `json:"lang"`
	TypedCode    string   `json:"typed_code"`
	FunctionName string   `json:"function_name"`
	TestCases    []string `json:"test_cases"`
	TestAnswers  []string `json:"test_answers"`
}

func ExecuteCode(payload *ProblemDetails) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	err := util.WriteCodeInExecutionFile(payload.Lang, payload.TypedCode, payload.FunctionName, payload.TestCases)
	if err != nil {
		return false, err
	}

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return false, fmt.Errorf("failed to create docker client - %w", err)
	}
	defer dockerClient.Client.Close()

	dockerImageTag := fmt.Sprintf("%s-docker-img", payload.Lang)

	image, err := dockerClient.CreateDockerImage(ctx, payload.Lang, dockerImageTag)
	if err != nil {
		return false, fmt.Errorf("failed to create docker image - %w", err)
	}
	defer image.Body.Close()
	log.Println("docker image created successfully.")

	container, err := dockerClient.RunDockerContainer(ctx, dockerImageTag)
	if err != nil {
		return false, fmt.Errorf("failed to run docker container - %w", err)
	}
	log.Println("docker container started successfully.")

	logs, err := dockerClient.GetContainerLogs(ctx, container.ID)
	if err != nil {
		return false, fmt.Errorf("failed to get docker container logs - %w", err)
	}
	defer logs.Close()

	var logBuffer strings.Builder
	io.Copy(&logBuffer, logs)

	if err = dockerClient.RemoveDockerContainer(ctx, container.ID); err != nil {
		return false, fmt.Errorf("failed to remove docker container - %w", err)
	}
	log.Println("container removed successfully.")

	if err = dockerClient.RemoveDockerImage(ctx, dockerImageTag); err != nil {
		return false, fmt.Errorf("failed to remove docker image - %w", err)
	}
	log.Println("image removed successfully.")

	sequenceToBeRemoved := []byte{1, 0, 0, 0, 0, 0, 0, 9}
	sanitizedLogs := util.SanitizeDockerLogs([]byte(logBuffer.String()), sequenceToBeRemoved)

	result := util.DetermineResult(string(sanitizedLogs), payload.TestAnswers)

	return result, nil
}
