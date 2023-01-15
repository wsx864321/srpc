package client

import "time"

type FailMode int

const (
	FailOver   FailMode = iota // 尝试另外一个节点
	FailFast                   // 立即失败
	FailTry                    // 对当前节点进行重试，直到超过重试次数
	FailBackup                 // 对冲请求（当前对冲请求需要设置latency，之后会仿照prometheus的Histograms指标类型进行分位的收集，业务只需要传递类似P95 P99这样的配置即可）
)

type callOptions struct {
	target        string        // 不采用服务发现的方式直接进行请求
	failMode      FailMode      // 失败模式
	retries       int           // retry次数
	backupLatency time.Duration // 使用备选节点时候的延时
}

type CallOption func(opt *callOptions)

// WithTarget 设置target
func WithTarget(addr string) CallOption {
	return func(opt *callOptions) {
		opt.target = addr
	}
}

// WithFailMode 设置失败模式
func WithFailMode(mode FailMode) CallOption {
	return func(opt *callOptions) {
		opt.failMode = mode
	}
}

// WithRetries 设置FailTry模式下的retry次数
func WithRetries(retries int) CallOption {
	return func(opt *callOptions) {
		opt.retries = retries
	}
}

// WithBackupLatency 设置对冲请求的延时时间限制
func WithBackupLatency(backupLatency time.Duration) CallOption {
	return func(opt *callOptions) {
		opt.backupLatency = backupLatency
	}
}

// NewCallOptions ...
func NewCallOptions(opts ...CallOption) *callOptions {
	var o callOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &o
}
