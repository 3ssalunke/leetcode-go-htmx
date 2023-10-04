package docker

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/api/types/container"
// 	"github.com/docker/docker/client"
// 	"github.com/docker/docker/pkg/archive"
// )

// type DockerClient struct {
// 	Client *client.Client
// }

// func NewDockerClient() (*DockerClient, error) {
// 	os.Setenv("DOCKER_HOST", "tcp://localhost:2375")
// 	client, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &DockerClient{client}, nil
// }

// func (client *DockerClient) CreateDockerImage(ctx context.Context, lang string, imageTag string) (*types.ImageBuildResponse, error) {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return nil, err
// 	}

// 	dockerfile := "Dockerfile"

// 	buildContextPath := fmt.Sprintf("%s\\docker\\runtimes\\%s", wd, lang)
// 	buildCtx, _ := archive.TarWithOptions(buildContextPath, &archive.TarOptions{})

// 	image, err := client.Client.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
// 		Context:    buildCtx,
// 		Dockerfile: dockerfile,
// 		Tags:       []string{imageTag},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = io.Copy(os.Stdout, image.Body)
// 	return &image, err
// }

// func (client *DockerClient) RemoveDockerImage(ctx context.Context, imageTag string) error {
// 	_, err := client.Client.ImageRemove(ctx, imageTag, types.ImageRemoveOptions{
// 		Force:         true,
// 		PruneChildren: true,
// 	})

// 	return err
// }

// func (client *DockerClient) RunDockerContainer(ctx context.Context, imageTag string) (*container.CreateResponse, error) {
// 	containerConfig := &container.Config{
// 		Image: imageTag,
// 	}

// 	container, err := client.Client.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = client.Client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
// 		return nil, err
// 	}

// 	return &container, nil
// }

// func (client *DockerClient) RemoveDockerContainer(ctx context.Context, containerID string) error {
// 	if err := client.Client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
// 		Force:         true,
// 		RemoveVolumes: true,
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (client *DockerClient) GetContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
// 	logOptions := types.ContainerLogsOptions{
// 		ShowStderr: true,
// 		ShowStdout: true,
// 		Timestamps: false,
// 		Follow:     false,
// 		Tail:       "40",
// 	}

// 	reader, err := client.Client.ContainerLogs(ctx, containerID, logOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return reader, err
// }
