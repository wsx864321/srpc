package client

import (
	"context"
)

type Client struct {
	opts *Options
}

// NewClient 生成client对象
func NewClient(opts ...Option) *Client {
	client := &Client{
		opts: NewOptions(opts...),
	}

	return client
}

func (c *Client) Call(ctx context.Context, methodName string, req, resp interface{}) error {
	return nil
}
