package docker

import "testing"

func TestGetDockerClient(t *testing.T) {
	cli, err := GetDockerClient()
	if err != nil {
		t.Fatalf("failed to get docker client: %v", err)
	}
	if cli == nil {
		t.Fatal("docker client should not be nil")
	}
}
