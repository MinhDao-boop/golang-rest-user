package repository

import (
	"errors"
	"strings"

	"golang-rest-user/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetList(page, pageSize int, search string) (users []models.User, total int64, err error)
	Update(user *models.User) error
	DeleteByID(id uint) error
	DeleteByIDs(ids []uint) (deleted int64, err error)
	GetByUsername(username string) (*models.User, error)
	GetByUUID(string) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetList(page, pageSize int, search string) (users []models.User, total int64, err error) {
	offset := (page - 1) * pageSize
	query := r.db.Model(&models.User{})
	if strings.TrimSpace(search) != "" {
		query = query.Where("username LIKE ?", "%"+search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("id asc").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) DeleteByID(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepo) DeleteByIDs(ids []uint) (int64, error) {
	res := r.db.Delete(&models.User{}, ids)
	return res.RowsAffected, res.Error
}

func (r *userRepo) GetByUUID(uuid string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &u, nil
}
