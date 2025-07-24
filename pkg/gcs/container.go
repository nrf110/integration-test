package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const exposedPort = "4443/tcp"

type Container struct {
	testcontainers.Container
}

func (c *Container) setExternalUrl(ctx context.Context) error {
	hostPort, err := c.HostAndPort(ctx)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/", hostPort)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("http://%s/_internal/config", hostPort),
		strings.NewReader(`{"externalUrl":"`+url+`","publicHost":"`+hostPort+`"}`),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, jsonErr := io.ReadAll(res.Body)
		if jsonErr != nil {
			return jsonErr
		}
		return errors.New("failed to update fake-gcs-server with new external url: " + string(body))
	}

	return nil
}

func (c *Container) HostAndPort(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	mappedPort, err := c.MappedPort(ctx, exposedPort)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%d", host, mappedPort.Int()), nil
}

func Run(ctx context.Context, image string, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image:      image,
		Entrypoint: []string{"/bin/fake-gcs-server", "-scheme", "http"},
		WaitingFor: wait.ForLog("server started"),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		err := opt.Customize(&genericContainerReq)
		if err != nil {
			return nil, err
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	var c *Container
	if container == nil {
		return nil, errors.New("container not found")
	}

	err = container.Start(ctx)
	if err != nil {
		return nil, err
	}
	c = &Container{container}

	err = c.setExternalUrl(ctx)
	if err != nil {
		return c, err
	}

	return c, nil
}
