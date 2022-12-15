package client

import (
	"context"
	"github.com/wsx864321/srpc/lb"
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
	service, err := c.opts.dis.GetService(ctx, c.opts.serviceName)
	if err != nil {
		return err
	}

	endpoint, err := lb.GetLB(c.opts.lbName).Pick(service)
	if err != nil {
		return err
	}

	conn, err := c.opts.pool.Get(endpoint.Network, endpoint.GetAddr())
	if err != nil {
		return err
	}

	conn.Write(req.([]byte))
	conn.Read(resp.([]byte))

	return nil
}
