package redisProvider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang-rest-user/dto"
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

func sosContact(tenant string, zoneId string) string {
	return fmt.Sprintf(
		"{%s}:sos_contact:zone:%s:all",
		tenant,
		zoneId)
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

func SetContactKeys(tenantCode string, zoneId string, contacts []dto.SOSContactResponse, ttl time.Duration) error {
	key := sosContact(tenantCode, zoneId)

	bytes, err := json.Marshal(contacts)
	if err != nil {
		return err
	}

	return client.Set(ctx, key, bytes, ttl).Err()
}

func GetFullContactKeys(tenantCode string, zoneId string) ([]dto.SOSContactResponse, error) {
	key := sosContact(tenantCode, zoneId)

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, err
	}

	var contacts []dto.SOSContactResponse
	if err := json.Unmarshal([]byte(val), &contacts); err != nil {
		return nil, err
	}

	return contacts, nil
}

func RevokeAllContacts(tenantCode, zoneId string) error {
	key := sosContact(tenantCode, zoneId)
	return client.Del(ctx, key).Err()
}
