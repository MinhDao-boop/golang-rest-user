package utils

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func GetTenantDBFromContext(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value("TENANT_DB").(*gorm.DB)
	if !ok || db == nil {
		return nil, errors.New("tenant db connection not found in context")
	}
	return db.WithContext(ctx), nil
}
