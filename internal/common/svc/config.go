package svc

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env        string           `yaml:"Env"`
	Server     ServerConfig     `yaml:"Server"`
	DB         DBConfig         `yaml:"DB"`
	Log        LogConfig        `yaml:"Log"`
	JWT        JWTConfig        `yaml:"JWT"`
	MinIO      MinIOConfig      `yaml:"MinIO"`
	Redis      RedisConfig      `yaml:"Redis"`
	ES         ESConfig         `yaml:"ES"`
	VolcEngine VolcEngineConfig `yaml:"VolcEngine"`
}

type ServerConfig struct {
	Port int64 `yaml:"Port"`
}
type DBConfig struct {
	Username     string `yaml:"Username"`
	Password     string `yaml:"Password"`
	Host         string `yaml:"Host"`
	Port         int64  `yaml:"Port"`
	Database     string `yaml:"Database"`
	MaxOpenConns int64  `yaml:"MaxOpenConns"`
	MaxIdleConns int64  `yaml:"MaxIdleConns"`
}

type LogConfig struct {
	Output string `yaml:"Output"`
	Format string `yaml:"Format"`
}

type JWTConfig struct {
	SecretKey string `yaml:"SecretKey"`
	TokenExp  string `yaml:"TokenExp"` // 例如: "24h", "1h30m"
}

type MinIOConfig struct {
	Endpoint        string `yaml:"Endpoint"`
	AccessKeyID     string `yaml:"AccessKeyID"`
	SecretAccessKey string `yaml:"SecretAccessKey"`
	UseSSL          bool   `yaml:"UseSSL"`
	Bucket          string `yaml:"Bucket"`
	Region          string `yaml:"Region"`
}

type RedisConfig struct {
	Addr     string `yaml:"Addr"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

type ESConfig struct {
	Addresses []string `yaml:"Addresses"`
	Username  string   `yaml:"Username"`
	Password  string   `yaml:"Password"`
}

type VolcEngineConfig struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	SecretAccessKey string `yaml:"SecretAccessKey"`
	Region          string `yaml:"Region"`
	Endpoint        string `yaml:"Endpoint"`
}

func LoadConfig(configFile string) (*Config, error) {
	// 优先级：显式参数 > 环境变量 CONFIG_FILE > 默认 config/local.yml
	if configFile == "" {
		if env := os.Getenv("CONFIG_FILE"); env != "" {
			configFile = env
		} else {
			configFile = "config/local.yml"
		}
	}

	cfg := new(Config)

	// // 1) 先尝试加载同目录下的 base.yml 作为基础配置
	// basePath := filepath.Join(filepath.Dir(configFile), "base.yml")
	// if info, err := os.Stat(basePath); err == nil && !info.IsDir() {
	// 	if bs, err := os.ReadFile(basePath); err == nil {
	// 		_ = yaml.Unmarshal(bs, cfg)
	// 	} else {
	// 		return nil, fmt.Errorf("read base config %s failed: %w", basePath, err)
	// 	}
	// }

	// 2) 加载指定配置文件，覆盖基础配置
	bs, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config file %s failed: %w", configFile, err)
	}
	if err := yaml.Unmarshal(bs, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config %s failed: %w", configFile, err)
	}

	return cfg, nil
}
