package docker

import (
	"sync"

	"github.com/docker/docker/client"
)

var (
	dockerCli  *client.Client
	dockerOnce sync.Once
)

func GetDockerClient() (*client.Client, error) {
	var err error
	dockerOnce.Do(func() {
		dockerCli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	})
	return dockerCli, err
}
