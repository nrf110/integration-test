package gcs

import (
	"context"
	"errors"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"net/http"
	"strings"
)

const exposedPort = "4443/tcp"

type Container struct {
	testcontainers.Container
}

func (c *Container) setExternalUrl(ctx context.Context) error {
	url, err := c.Url(ctx)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/_internal/config", url),
		strings.NewReader(`{"externalUrl":"`+url+`"}`),
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

func (c *Container) Url(ctx context.Context) (string, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", err
	}

	mappedPort, err := c.MappedPort(ctx, exposedPort)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%d", host, mappedPort.Int()), nil
}

func Run(ctx context.Context, image string, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image:      image,
		Entrypoint: []string{"/bin/fake-gcs-server", "-scheme", "http"},
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
