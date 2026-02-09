package redisProvider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		panic("redis connection failed:" + err.Error())
	}
}

func GetClient() *redis.Client {
	return client
}

func refreshKey(tenant string, userID uint, tokenHash string) string {
	return fmt.Sprintf(
		"{%s}:auth:user:%d:refresh:%s",
		tenant,
		userID,
		tokenHash,
	)
}

func userRefreshSetKey(tenant string, userID uint) string {
	return fmt.Sprintf(
		"{%s}:auth:user:%d:refresh_tokens",
		tenant,
		userID,
	)
}

func userTokenVersion(tenant string, userID uint) string {
	return fmt.Sprintf(
		"{%s}:auth:user:%d:token_ver",
		tenant,
		userID,
	)
}

func Create(tokenHash string, userID uint, tenantCode string, ttl time.Duration) error {

	refreshKey := refreshKey(tenantCode, userID, tokenHash)
	userSetKey := userRefreshSetKey(tenantCode, userID)

	pipe := client.TxPipeline()

	pipe.HSet(ctx, refreshKey, map[string]interface{}{
		"user_id":     userID,
		"tenant_code": tenantCode,
	})

	pipe.Expire(ctx, refreshKey, ttl)
	pipe.SAdd(ctx, userSetKey, tokenHash)
	pipe.Expire(ctx, userSetKey, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func FindValidByHash(tokenHash string, tenantCode string, userID uint) error {

	key := refreshKey(tenantCode, userID, tokenHash)

	exists, err := client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("refresh token revoked or expired")
	}
	return nil
}

func Revoke(tokenHash string, tenantCode string, userID uint) error {

	refreshKey := refreshKey(tenantCode, userID, tokenHash)
	userSetKey := userRefreshSetKey(tenantCode, userID)

	pipe := client.TxPipeline()
	pipe.Del(ctx, refreshKey)
	pipe.SRem(ctx, userSetKey, tokenHash)

	_, err := pipe.Exec(ctx)
	return err
}

func RevokeAllByUser(tenantCode string, userID uint) error {

	userSetKey := userRefreshSetKey(tenantCode, userID)

	tokens, err := client.SMembers(ctx, userSetKey).Result()
	if err != nil {
		return err
	}

	pipe := client.TxPipeline()
	for _, token := range tokens {
		pipe.Del(ctx, refreshKey(tenantCode, userID, token))
	}
	pipe.Del(ctx, userSetKey)

	_, err = pipe.Exec(ctx)
	return err
}

func GetTokenVer(userID uint, tenantCode string) int {
	key := userTokenVersion(tenantCode, userID)
	val, err := client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		client.Set(ctx, key, 1, 0)
		return 1
	}
	return val
}

func IncreaseTokenVer(userID uint, tenantCode string) error {
	key := userTokenVersion(tenantCode, userID)
	return client.Incr(ctx, key).Err()
}

func GetCache(key string) ([]byte, error) {
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(value), nil
}

func SetCache(key string, value interface{}, duration int, isJson bool) error {
	if isJson {
		temp, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return client.Set(ctx, key, temp, time.Second*time.Duration(duration)).Err()
	}
	return client.Set(ctx, key, value, time.Second*time.Duration(duration)).Err()
}

func DeleteCache(key string) error {
	_, err := client.Del(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	return err
}

func DeleteCachePattern(key string) error {
	key = fmt.Sprintf("%s:*", key)
	iter := client.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		if err := DeleteCache(iter.Val()); err != nil {
			return err
		}
	}
	return iter.Err()
}
