package pullbuddy

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type dockerImagePuller struct {
	dockerClient *client.Client
}

func newDockerImagePuller(dockerClient *client.Client) *dockerImagePuller {
	return &dockerImagePuller{
		dockerClient: dockerClient,
	}
}

func (puller *dockerImagePuller) pull(id string) error {
	readCloser, err := puller.dockerClient.ImagePull(
		context.Background(),
		id,
		types.ImagePullOptions{},
	)
	if readCloser != nil {
		defer readCloser.Close()
	}
	return err
}
