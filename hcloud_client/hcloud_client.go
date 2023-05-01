package hcloud_client

import (
	"context"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Client struct {
	*hcloud.Client
	context.Context
}

func New() Client {
	return Client{
		Client:  hcloud.NewClient(hcloud.WithToken(env.HCloudToken())),
		Context: context.Background(),
	}
}
