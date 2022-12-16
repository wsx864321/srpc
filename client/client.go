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
	// 1.获取服务地址
	// 2.进行负载均衡策略，选择node
	// 3.从连接池获取链接
	// 4.数据进行序列化
	// 5.组装invoke(发送数据+接手数据)
	// 6.执行中间件以及handler
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
	defer c.opts.pool.Put(endpoint.Network, endpoint.GetAddr(), conn)

	conn.Write(req.([]byte))
	conn.Read(resp.([]byte))

	return nil
}
