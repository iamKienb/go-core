package esx

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	configx "github.com/iamKienb/go-core/config"
)

type ESX struct {
	client *elasticsearch.TypedClient
}

func New(ctx context.Context, cfg configx.ElasticSearchConfig) (ESXService, error) {
	esCfg := elasticsearch.Config{
		Addresses:       cfg.Addresses,
		Instrumentation: elasticsearch.NewOpenTelemetryInstrumentation(nil, false),
	}

	client, err := elasticsearch.NewTypedClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	_, err = client.Info().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to elasticsearch: %w", err)
	}

	return &ESX{client: client}, nil
}
