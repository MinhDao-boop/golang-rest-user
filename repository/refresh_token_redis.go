package repository

import (
	"errors"
	"strconv"
	"time"

	"golang-rest-user/config"
)

type RefreshTokenRedis interface {
	Create(string, uint, time.Duration) error
	FindValidByHash(string) error
	Revoke(string) error
	RevokeAllByUser(uint) error
}

type RefreshTokenRedisRepo struct{}

func NewRefreshTokenRedisRepo() RefreshTokenRedis {
	return &RefreshTokenRedisRepo{}
}

func (r *RefreshTokenRedisRepo) Create(tokenHash string, userID uint, ttl time.Duration) error {
	key := "refresh:" + tokenHash
	userKey := "user_refresh:" + strconv.Itoa(int(userID))

	pipe := config.RDB.TxPipeline()
	pipe.HSet(config.Ctx, key, "user_id", userID)
	pipe.Expire(config.Ctx, key, ttl)
	pipe.SAdd(config.Ctx, userKey, tokenHash)
	pipe.Expire(config.Ctx, userKey, ttl)

	_, err := pipe.Exec(config.Ctx)
	return err
}

func (r *RefreshTokenRedisRepo) FindValidByHash(tokenHash string) error {
	key := "refresh:" + tokenHash

	exists, err := config.RDB.Exists(config.Ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("refresh token revoked")
	}
	return nil
}

func (r *RefreshTokenRedisRepo) Revoke(tokenHash string) error {
	return config.RDB.Del(config.Ctx, "refresh:"+tokenHash).Err()
}

func (r *RefreshTokenRedisRepo) RevokeAllByUser(userID uint) error {
	userKey := "user_refresh:" + strconv.Itoa(int(userID))

	tokens, err := config.RDB.SMembers(config.Ctx, userKey).Result()
	if err != nil {
		return err
	}

	pipe := config.RDB.TxPipeline()

	for _, t := range tokens {
		pipe.Del(config.Ctx, "refresh:"+t)
	}
	pipe.Del(config.Ctx, userKey)

	_, err = pipe.Exec(config.Ctx)
	return err
}
