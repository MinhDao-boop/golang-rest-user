package repository

import (
	"golang-rest-user/models"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepo interface {
	Create(rToken *models.RefreshToken) error
	FindValidByHash(hash string) (*models.RefreshToken, error)
	Revoke(id uint) error
	RevokeAllByUser(userID uint) error
}

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) RefreshTokenRepo {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(rToken *models.RefreshToken) error {
	return r.db.Create(rToken).Error
}

func (r *refreshTokenRepo) FindValidByHash(hash string) (*models.RefreshToken, error) {
	var rToken models.RefreshToken
	err := r.db.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ? ", hash, time.Now()).
		First(&rToken).Error
	if err != nil {
		return nil, err
	}
	return &rToken, nil
}

func (r *refreshTokenRepo) Revoke(id uint) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at = ?", now).Error
}

func (r *refreshTokenRepo) RevokeAllByUser(userID uint) error {
	now := time.Now()
	return r.db.Model(models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at = ?", now).Error
}
