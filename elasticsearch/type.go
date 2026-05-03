package esx

import (
	"context"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Response[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
}

type BulkConfig struct {
	FlushInterval time.Duration
	FlushBytes    int
	NumWorkers    int
}

type ElasticsearchService interface {
	BootstrapIndex(ctx context.Context, alias string, mappingJson string) error
	Sync(ctx context.Context, alias, id string, data any) error
	BulkWorker(alias string, cfg BulkConfig) (esutil.BulkIndexer, error)
	GetClient() *elasticsearch.TypedClient
	Close(ctx context.Context) error
}
