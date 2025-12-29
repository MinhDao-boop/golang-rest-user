package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"golang-rest-user/config"
)

type RefreshTokenRedis interface {
	Create(string, uint, string, time.Duration) error
	FindValidByHash(string) error
	Revoke(string) error
	RevokeAllByUser(uint) error
}

type RefreshTokenRedisRepo struct{}

func NewRefreshTokenRedisRepo() RefreshTokenRedis {
	return &RefreshTokenRedisRepo{}
}

func (r *RefreshTokenRedisRepo) Create(tokenHash string, userID uint, tenantCode string, ttl time.Duration) error {
	key := "refresh:" + tokenHash
	userKey := "user_refresh:" + strconv.Itoa(int(userID))

	ctx := context.Background()
	pipe := config.RDB.TxPipeline()
	pipe.HSet(ctx, key, map[string]interface{}{
		"user_id":     userID,
		"tenant_code": tenantCode,
	})
	pipe.Expire(ctx, key, ttl)
	pipe.SAdd(ctx, userKey, tokenHash)
	pipe.Expire(ctx, userKey, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRedisRepo) FindValidByHash(tokenHash string) error {
	key := "refresh:" + tokenHash

	ctx := context.Background()
	exists, err := config.RDB.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("refresh token revoked")
	}
	return nil
}

func (r *RefreshTokenRedisRepo) Revoke(tokenHash string) error {
	ctx := context.Background()
	return config.RDB.Del(ctx, "refresh:"+tokenHash).Err()
}

func (r *RefreshTokenRedisRepo) RevokeAllByUser(userID uint) error {
	userKey := "user_refresh:" + strconv.Itoa(int(userID))

	ctx := context.Background()
	tokens, err := config.RDB.SMembers(ctx, userKey).Result()
	if err != nil {
		return err
	}

	pipe := config.RDB.TxPipeline()

	for _, t := range tokens {
		pipe.Del(ctx, "refresh:"+t)
	}
	pipe.Del(ctx, userKey)

	_, err = pipe.Exec(ctx)
	return err
}
