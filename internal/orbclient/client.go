package orbclient

import (
	"context"

	"github.com/orbcorp/orb-go"
	"github.com/orbcorp/orb-go/option"
)

type Client struct {
	orb *orb.Client
}

func New(apiKey, baseURL string) *Client {
	opts := []option.RequestOption{}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return &Client{orb: orb.NewClient(opts...)}
}

func (c *Client) Ping(ctx context.Context) (interface{}, error) {
	return c.orb.TopLevel.Ping(ctx)
}
