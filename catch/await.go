package catch

import (
	"context"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

func Await(client *hcloud.Client, ctx context.Context, action *hcloud.Action) error {
	_, errors := client.Action.WatchProgress(ctx, action)
	return <-errors
}
