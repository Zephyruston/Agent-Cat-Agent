package container

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerClient struct {
	cli *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli: cli}, nil
}

func (d *DockerClient) PullImage(ctx context.Context, imageName string) error {
	reader, err := d.cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)
	return nil
}

func (d *DockerClient) RunContainer(ctx context.Context, imageName string, cmd []string) (string, error) {
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   cmd,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", err
	}
	if err := d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}
	statusCh, errCh := d.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}
	out, err := d.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return resp.ID, nil
}

func (d *DockerClient) RemoveContainer(ctx context.Context, containerID string) error {
	return d.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

func (d *DockerClient) RunContainerWithMount(ctx context.Context, imageName string, cmd []string, hostDir, targetDir string) (string, error) {
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image:      imageName,
		Cmd:        cmd,
		Tty:        false,
		WorkingDir: targetDir,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: hostDir,
				Target: targetDir,
			},
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}
	defer d.RemoveContainer(ctx, resp.ID)
	if err := d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}
	statusCh, errCh := d.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}
	out, err := d.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return resp.ID, nil
}

func (d *DockerClient) Close() error {
	return d.cli.Close()
}
