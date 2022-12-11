package client

import "context"

type Client struct {
	opts *Options
}

// NewClient 生成client对象
func NewClient(opts ...Option) *Client {
	return &Client{
		opts: NewOptions(opts...),
	}
}

func (c *Client) Call(ctx context.Context, methodName string, req, resp interface{}) error {
	return nil
}
