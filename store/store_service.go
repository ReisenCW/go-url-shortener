package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// Define the struct wrapper around raw Redis client
type StorageService struct {
	redisClient *redis.Client
}

// Top level declarations for the storeService and Redis context
var (
	storeService = &StorageService{} // 单例模式
    ctx = context.Background()	
)

// 完成 redis 客户端的创建、连接
func InitializeStore() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 无密码
		DB:       0,  // 默认DB
	})

	// 发送 Ping 命令，若失败则panic, 若成功则打印pong消息
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)
	storeService.redisClient = redisClient
	return storeService
}

// 保存原url与生成的短url的映射关系
func SaveUrlMapping(shortUrl string, originalUrl string, userId string) { 
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key url | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
}

// 根据短url获取原始url
func RetrieveInitialUrl(shortUrl string) string {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | Error: %v - shortUrl: %s\n", err, shortUrl))
	}
	return result
}

// 实际应用中，缓存不应该设置过期时间，应该设置LRU策略
// 让那些不常用的值被自动清理掉，并在缓存满时存回RDBMS
const CacheDuration = 6 * time.Hour