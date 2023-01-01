package trace

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	tp             *tracesdk.TracerProvider
	once           sync.Once
	defaultOptions = Options{
		Enable:      true,
		Url:         "http://127.0.0.1:14268/api/traces",
		Sampler:     1.0,
		ServiceName: "srpc",
	}
)

type Options struct {
	Enable      bool
	Url         string
	Sampler     float64
	ServiceName string
}

type Option func(opt *Options)

func newOptions(opts ...Option) *Options {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &o
}

func WithEnable(enable bool) Option {
	return func(opt *Options) {
		opt.Enable = enable
	}
}

func WithUrl(url string) Option {
	return func(opt *Options) {
		opt.Url = url
	}
}

func WithSampler(sampler float64) Option {
	return func(opt *Options) {
		opt.Sampler = sampler
	}
}

func WithServiceName(name string) Option {
	return func(opt *Options) {
		opt.ServiceName = name
	}
}

// StartAgent 开启trace collector
func StartAgent(opts ...Option) {
	o := newOptions(opts...)
	if !o.Enable {
		return
	}

	once.Do(func() {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Url)))
		if err != nil {
			return
		}

		tp = tracesdk.NewTracerProvider(
			tracesdk.WithSampler(tracesdk.TraceIDRatioBased(o.Sampler)),
			tracesdk.WithBatcher(exp),
			tracesdk.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(o.ServiceName),
			)),
		)

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{}))
	})
}

// StopAgent 关闭trace collector,在服务停止时调用StopAgent，不然可能造成trace数据的丢失
func StopAgent() {
	_ = tp.Shutdown(context.TODO())
}
