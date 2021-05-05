package otel

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTraceProvider(isDev bool, serviceName string, jobID, endpointHost, endpointPort string) (func(), error) {
	var endpoint jaeger.EndpointOption
	if isDev {
		endpoint = jaeger.WithAgentEndpoint(fmt.Sprintf("%s:%s", endpointHost, endpointPort))
	} else {
		endpoint = jaeger.WithCollectorEndpoint(fmt.Sprintf("http://%s:%s/api/traces", endpointHost, endpointPort))
	}

	flush, err := jaeger.InstallNewPipeline(
		endpoint,
		jaeger.WithProcess(jaeger.Process{
			ServiceName: serviceName,
			Tags: []attribute.KeyValue{
				attribute.String("jobID", jobID),
			},
		}),
		jaeger.WithSDK(&sdkTrace.Config{
			DefaultSampler: sdkTrace.AlwaysSample(),
		}),
	)

	if err != nil {
		return nil, err
	}

	return flush, nil
}
