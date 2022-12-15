package docker

import (
	"net/http"
	"os"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/client"
	"github.com/luno/jettison/errors"
)

func NewDockerClient() (*client.Client, error) {
	var clientOpts []client.Opt

	clientOpts = append(clientOpts,
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	host := os.Getenv("DOCKER_HOST")
	if host != "" {
		helper, err := connhelper.GetConnectionHelper(host)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				DialContext: helper.Dialer,
			},
		}

		clientOpts = append(clientOpts,
			client.WithHTTPClient(httpClient),
			client.WithHost(helper.Host),
			client.WithDialContext(helper.Dialer),
		)
	}

	cli, err := client.NewClientWithOpts(clientOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return cli, nil
}
