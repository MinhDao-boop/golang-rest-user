package repository

import (
	"context"
	"golang-rest-user/models"
	"golang-rest-user/utils"
)

type UserRepo interface {
	Create(context.Context, *models.User) error
	GetByID(context.Context, uint) (*models.User, error)
	GetList(ctx context.Context, page, pageSize int, search string) (users []models.User, total int64, err error)
	Update(context.Context, *models.User) error
	DeleteByIDs(context.Context, []uint) (deleted int64, err error)
	GetByUsername(context.Context, string) (*models.User, error)
	GetByUUID(context.Context, string) (*models.User, error)
}

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	db, err := utils.GetTenantDBFromContext(ctx)
	if err != nil {
		return err
	}
	return db.Create(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	db, err := utils.GetTenantDBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var u models.User
	if err := db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	db, err := utils.GetTenantDBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var u models.User
	if err := db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetList(ctx context.Context, page, pageSize int, search string) (users []models.User, total int64, err error) {
	db, err := utils.GetTenantDBFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	query := db.Model(&models.User{})
	query = query.Where("username LIKE ?", "%"+search+"%")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("id asc").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	db, err := utils.GetTenantDBFromContext(ctx)
	if err != nil {
		return err
	}
	return db.Save(user).Error
}

func (r *userRepo) DeleteByIDs(ctx context.Context, ids []uint) (int64, error) {
	db, err := utils.GetTenantDBFromContext(ctx)
	if err != nil {
		return 0, err
	}
	res := db.Delete(&models.User{}, ids)
	return res.RowsAffected, res.Error
}

func (r *userRepo) GetByUUID(ctx context.Context, uuid string) (*models.User, error) {
	db, err := utils.GetTenantDBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var u models.User
	if err := db.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
