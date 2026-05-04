package esx

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func (x *ESX) BootstrapIndex(ctx context.Context, alias string, mappingJson string) error {
	exists, err := x.client.Indices.ExistsAlias(alias).Do(ctx)
	if err != nil {
		return fmt.Errorf("check alias %s: %w", alias, err)
	}

	if !exists {
		realIndex := fmt.Sprintf("%s_%s", alias, time.Now().Format("20060102"))

		_, err := x.client.Indices.Create(realIndex).
			Raw(strings.NewReader(mappingJson)).Do(ctx)
		if err != nil {
			return fmt.Errorf("create index %s: %w", realIndex, err)
		}

		_, err = x.client.Indices.PutAlias(realIndex, alias).Do(ctx)
		if err != nil {
			return fmt.Errorf("put alias %s to index %s: %w", alias, realIndex, err)
		}

		return nil
	}
	var body map[string]interface{}
	if err := json.Unmarshal([]byte(mappingJson), &body); err != nil {
		return fmt.Errorf("failed to parse mapping json: %w", err)
	}

	mappings, ok := body["mappings"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("mapping file missing 'mappings' root element")
	}

	mappingsOnlyJson, _ := json.Marshal(mappings)

	_, err = x.client.Indices.PutMapping(alias).
		Raw(strings.NewReader(string(mappingsOnlyJson))).Do(ctx)
	if err != nil {
		return fmt.Errorf("update mapping for alias %s: %w", alias, err)
	}

	return nil
}

func (x *ESX) Sync(ctx context.Context, alias, id string, data any) error {
	_, err := x.client.Update(alias, id).Doc(data).DocAsUpsert(true).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync data to ES: %w", err)
	}

	return nil
}

func (x *ESX) BulkWorker(alias string, cfg BulkConfig) (esutil.BulkIndexer, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         alias,
		Client:        &x.client.BaseClient,
		FlushInterval: cfg.FlushInterval,
		FlushBytes:    cfg.FlushBytes,
		NumWorkers:    cfg.NumWorkers,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating bulk indexer: %w", err)
	}

	return bi, nil
}

func (x *ESX) GetClient() *elasticsearch.TypedClient {
	return x.client
}

func (x *ESX) Close(ctx context.Context) error {
	return x.client.Close(ctx)
}
