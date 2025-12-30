package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTenantDB(c *gin.Context) (*gorm.DB, error) {
	dbRaw, ok := c.Get("TENANT_DB")
	if !ok {
		return nil, errors.New("tenant code is required")
	}
	return dbRaw.(*gorm.DB), nil
}
