package client

import (
	"context"
	"github.com/wsx864321/srpc/codec"
	"github.com/wsx864321/srpc/codec/serialize"
	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/lb"
	"github.com/wsx864321/srpc/util"
	"time"
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
	service, err := c.opts.dis.GetService(ctx, c.opts.serviceName)
	if err != nil {
		return err
	}

	// 2.进行负载均衡策略，选择node
	endpoint, err := lb.GetLB(c.opts.lbName).Pick(service)
	if err != nil {
		return err
	}

	// 3.从连接池获取链接
	conn, err := c.opts.pool.Get(endpoint.Network, endpoint.GetAddr())
	if err != nil {
		return err
	}
	defer c.opts.pool.Put(endpoint.Network, endpoint.GetAddr(), conn)

	// 4.组装invoke(数据序列化+数据编码+发送数据+接手数据)
	invoker := func(ctx context.Context, req, resp interface{}) error {
		// 4.1 数据序列化
		raw, err := serialize.GetSerialize(serialize.SerializeType(endpoint.Serialize)).Marshal(req)
		if err != nil {
			return err
		}
		// 4.2 获取metadata（trace、级联超时、用户其它自定义的数据等等）

		// 4.3 执行编码
		request, err := codec.NewCodec().Encode(
			codec.GeneralMsgType,
			codec.CompressTypeNot,
			uint64(time.Now().Unix()),
			[]byte(c.opts.serviceName),
			[]byte(methodName),
			[]byte(""),
			raw,
		)
		if err != nil {
			return err
		}

		// 4.4 发送请求数据
		if c.opts.writeTimeout > 0 {
			if err := conn.SetReadDeadline(time.Now().Add(c.opts.writeTimeout)); err != nil {
				return err
			}
		}

		if err := util.Write(conn, request); err != nil {
			return err
		}

		// 4.5 接受数据

		//conn.Read(resp.([]byte))

		return nil
	}

	// 5.执行中间件以及invoker函数
	return interceptor.ClientIntercept(ctx, req, resp, c.opts.interceptors, invoker)
}
