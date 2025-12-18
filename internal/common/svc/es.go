package svc

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func InitES(cfg *Config) (*elasticsearch.Client, error) {
	if len(cfg.ES.Addresses) == 0 {
		return new(elasticsearch.Client), nil // 未配置则不启用
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.ES.Addresses,
		Username:  cfg.ES.Username,
		Password:  cfg.ES.Password,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
