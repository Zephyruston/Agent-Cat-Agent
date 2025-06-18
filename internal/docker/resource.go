package docker

import "github.com/docker/docker/api/types/container"

type ResourceConfig struct {
	MemoryLimit int64 // bytes
	CPUShares   int64
}

func ToHostConfig(res ResourceConfig) *container.HostConfig {
	return &container.HostConfig{
		Resources: container.Resources{
			Memory:    res.MemoryLimit,
			CPUShares: res.CPUShares,
		},
	}
}
