package repository

import (
	"golang-rest-user/models"
	"golang-rest-user/utils"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(*models.User) error
	GetByID(uint) (*models.User, error)
	GetList(page, pageSize int, search string) (users []models.User, total int64, err error)
	Update(*models.User) error
	DeleteByUUIDs([]string) (deleted int64, err error)
	GetByUsername(string) (*models.User, error)
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
	query := r.db.Model(&models.User{})
	query = query.Where("username LIKE ?", "%"+search+"%")
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query = query.Scopes(utils.PaginateGORM(page, pageSize))
	if err = query.Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) DeleteByUUIDs(uuids []string) (int64, error) {
	res := r.db.Where("uuid IN ?", uuids).Delete(&models.User{})
	return res.RowsAffected, res.Error
}

func (r *userRepo) GetByUUID(uuid string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
