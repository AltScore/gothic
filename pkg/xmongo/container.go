package xmongo

import (
	"context"
	"io"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// Container is a wrapper around dockertest.Resource to allow run commands against an existent container
type Container struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	ctx      context.Context
	cancel   context.CancelFunc
}

// TailLogs tails the logs of the container.
func (c *Container) TailLogs(ctx context.Context, wr io.Writer, follow bool) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	opts := docker.LogsOptions{
		Context: ctx,

		Stderr:      true,
		Stdout:      true,
		Follow:      follow,
		Timestamps:  true,
		RawTerminal: true,

		Container: c.resource.Container.ID,

		OutputStream: wr,
	}

	return c.pool.Client.Logs(opts)
}

func (c *Container) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}
