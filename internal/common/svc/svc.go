package svc

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// DB 全局数据库实例
var svc = new(ServiceContext)

type ServiceContext struct {
	Cfg   *Config
	DB    *gorm.DB
	Redis *redis.Client
	Minio *minio.Client
	ES    *elasticsearch.Client
}

func Init(configFile string) (*ServiceContext, error) {
	var err error

	svc.Cfg, err = LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	svc.DB, err = InitDB(svc.Cfg)
	if err != nil {
		return nil, err
	}

	// 初始化Redis客户端
	svc.Redis, err = InitRedis(svc.Cfg)
	if err != nil {
		return nil, err
	}

	svc.Minio, err = InitMinio(svc.Cfg)
	if err != nil {
		return nil, err
	}
	// 初始化 Elasticsearch
	svc.ES, err = InitES(svc.Cfg)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func Get() *ServiceContext {
	return svc
}

func DB() *gorm.DB {
	return svc.DB
}

func Cfg() *Config {
	return svc.Cfg
}

func Redis() *redis.Client {
	return svc.Redis
}

// Client 获取MinIO客户端
func Minio() *minio.Client {
	return svc.Minio
}

// ES 返回 Elasticsearch 客户端
func ES() *elasticsearch.Client {
	return svc.ES
}
