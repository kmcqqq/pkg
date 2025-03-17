package mongo

import (
	"context"
	"fmt"
	"github.com/kmcqqq/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	client *mongo.Client
)

// Initialize 初始化 Mongo 客户端
func Initialize(cfg *config.ServerInfo) error {
	addr := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.User, cfg.Pwd, cfg.Host, cfg.Port, cfg.Db)
	clientOptions := options.Client().ApplyURI(addr)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

// GetClient 获取 Mongo 客户端
func GetClient() *mongo.Client {
	return client
}
