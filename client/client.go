package client

import (
	"context"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Client struct {
	*hcloud.Client
	context.Context
}

func (c Client) Await(action *hcloud.Action) error {
	_, errors := c.Client.Action.WatchProgress(c.Context, action)
	return <-errors
}
