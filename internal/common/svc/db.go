package svc

import (
	"errors"
	"fmt"
	"godir/internal/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/sharding"
)

// Init 初始化数据库连接
func InitDB(c *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DB.Username, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "source",
		NumberOfShards:      2,
		PrimaryKeyGenerator: sharding.PKSnowflake,
		ShardingAlgorithm:   ShardingAlgorithm,
		ShardingSuffixs:     ShardingSuffixs,
	}, "user"))

	db.AutoMigrate(
		&model.User{},
		&model.GodirUser{},
		&model.GodirMaterial{},
		&model.GodirPublishedMaterial{},
		&model.GodirPublishedLike{},
	)

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connected successfully")
	return db, nil
}

// Close 关闭数据库连接
func Close() error {
	if svc.DB != nil {
		sqlDB, err := svc.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(models ...interface{}) error {
	if svc.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	return svc.DB.AutoMigrate(models...)
}

func ShardingAlgorithm(columnValue any) (suffix string, err error) {
	source, ok := columnValue.(int64)
	if !ok {
		return "", errors.New("invalid source")
	}
	switch source {
	case 1:
		return "_page", nil
	case 2:
		return "_openapi", nil
	default:
		return "", errors.New("invalid source")
	}
}

func ShardingSuffixs() []string {
	return []string{"_page", "_openapi"}
}
