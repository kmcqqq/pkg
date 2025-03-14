package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"momo-server-go/pkg/config"
	"momo-server-go/pkg/utils"
	"time"
)

var (
	client *redis.Client
)

// Initialize 初始化 Redis 客户端
func Initialize(cfg *config.ServerInfo) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     cfg.User,
		Password:     cfg.Pwd,
		DB:           utils.StringToInt(cfg.Db),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		//PoolSize:     10,
		PoolTimeout: 4 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

// GetClient 获取 Redis 客户端
func GetClient() *redis.Client {
	return client
}

// Close 关闭 Redis 连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
