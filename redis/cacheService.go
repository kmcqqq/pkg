package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gitlab.bobbylive.cn/kongmengcheng/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type ICacheService interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
}

func NewCacheService(client *redis.Client) ICacheService {
	return &cacheService{
		client: client,
		ctx:    context.Background(),
	}
}

func GetCache[T any](cache ICacheService, redisKey string) ([]T, error) {
	var result []T

	// 尝试从缓存获取数据
	val, err := cache.Get(redisKey)
	if err != nil {
		return nil, err
	}
	if val == "" {
		return nil, nil
	}

	err = utils.Json2Struct(val, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetOrSetCache[T any](cache ICacheService, db *gorm.DB, redisKey string, queryFunc func(db *gorm.DB) ([]T, error), expiration time.Duration) ([]T, error) {
	var result []T

	// 尝试从缓存获取数据
	val, err := cache.Get(redisKey)
	if err == nil && val != "" {
		// 反序列化数据
		err = utils.Json2Struct(val, &result)
		if err == nil {
			return result, nil
		}

	}

	// 缓存未命中或反序列化失败，查询数据库
	result, err = queryFunc(db)
	if err != nil {
		return nil, err
	}

	// 将数据序列化并存入缓存
	data, err := utils.Struct2Json(result)
	if err == nil {
		err = cache.Set(redisKey, data, expiration)
	}

	return result, err
}

type cacheService struct {
	client *redis.Client
	ctx    context.Context
}

func (c cacheService) Get(key string) (string, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func (c cacheService) Set(key string, value string, expiration time.Duration) error {
	return c.client.Set(c.ctx, key, value, expiration).Err()
}

func (c cacheService) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c cacheService) Exists(key string) (bool, error) {
	count, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
