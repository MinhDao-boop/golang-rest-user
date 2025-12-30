package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang-rest-user/config"
)

type RefreshTokenRedis interface {
	Create(tokenHash string, userID uint, tenantCode string, ttl time.Duration) error
	FindValidByHash(tokenHash string, tenantCode string, userID uint) error
	Revoke(tokenHash string, tenantCode string, userID uint) error
	RevokeAllByUser(tenantCode string, userID uint) error
}

type RefreshTokenRedisRepo struct{}

func NewRefreshTokenRedisRepo() RefreshTokenRedis {
	return &RefreshTokenRedisRepo{}
}

/* -------------------- key helpers -------------------- */

func refreshKey(tenant string, userID uint, tokenHash string) string {
	return fmt.Sprintf(
		"auth:{%s}:user:%d:refresh:%s",
		tenant,
		userID,
		tokenHash,
	)
}

func userRefreshSetKey(tenant string, userID uint) string {
	return fmt.Sprintf(
		"auth:{%s}:user:%d:refresh_tokens",
		tenant,
		userID,
	)
}

/* -------------------- methods -------------------- */

func (r *RefreshTokenRedisRepo) Create(tokenHash string, userID uint, tenantCode string, ttl time.Duration) error {

	ctx := context.Background()

	refreshKey := refreshKey(tenantCode, userID, tokenHash)
	userSetKey := userRefreshSetKey(tenantCode, userID)

	pipe := config.RDB.TxPipeline()

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

func (r *RefreshTokenRedisRepo) FindValidByHash(tokenHash string, tenantCode string, userID uint) error {

	ctx := context.Background()
	key := refreshKey(tenantCode, userID, tokenHash)

	exists, err := config.RDB.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("refresh token revoked or expired")
	}
	return nil
}

func (r *RefreshTokenRedisRepo) Revoke(tokenHash string, tenantCode string, userID uint) error {

	ctx := context.Background()

	refreshKey := refreshKey(tenantCode, userID, tokenHash)
	userSetKey := userRefreshSetKey(tenantCode, userID)

	pipe := config.RDB.TxPipeline()
	pipe.Del(ctx, refreshKey)
	pipe.SRem(ctx, userSetKey, tokenHash)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *RefreshTokenRedisRepo) RevokeAllByUser(tenantCode string, userID uint) error {

	ctx := context.Background()
	userSetKey := userRefreshSetKey(tenantCode, userID)

	tokens, err := config.RDB.SMembers(ctx, userSetKey).Result()
	if err != nil {
		return err
	}

	pipe := config.RDB.TxPipeline()
	for _, token := range tokens {
		pipe.Del(ctx, refreshKey(tenantCode, userID, token))
	}
	pipe.Del(ctx, userSetKey)

	_, err = pipe.Exec(ctx)
	return err
}
