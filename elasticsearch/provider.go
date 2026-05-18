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

func New(cfg configx.ElasticSearchConfig) (ESXService, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}

	client, err := elasticsearch.NewTypedClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	_, err = client.Info().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connect to elasticsearch: %w", err)
	}

	return &ESX{client: client}, nil
}
