package tracing

import (
	"cobaApp/config"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
)

func GenerateTracing(config config.IConfig, log *logrus.Logger, serviceName string) (opentracing.Tracer, io.Closer) {
	jaegerCfg := jaegerConfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans: false,
			LocalAgentHostPort: fmt.Sprintf("%v:%v",
				config.GetConfig().Jaeger.Host, config.GetConfig().Jaeger.Port),
		},
	}

	tracer, closer, err := jaegerCfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("cant connect jaeger : %v", err)
	}

	return tracer, closer
}
