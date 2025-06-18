package docker

import (
	"context"

	"github.com/docker/docker/client"
)

type ContainerManager struct {
	cli *client.Client
}

func NewContainerManager(cli *client.Client) *ContainerManager {
	return &ContainerManager{cli: cli}
}

func (m *ContainerManager) CreateContainer(ctx context.Context, name, image string, cmd []string) (string, error) {
	// TODO: 实现容器创建
	return "", nil
}

func (m *ContainerManager) StartContainer(ctx context.Context, id string) error {
	// TODO: 实现容器启动
	return nil
}

func (m *ContainerManager) StopContainer(ctx context.Context, id string) error {
	// TODO: 实现容器停止
	return nil
}

func (m *ContainerManager) RemoveContainer(ctx context.Context, id string) error {
	// TODO: 实现容器清理
	return nil
}
