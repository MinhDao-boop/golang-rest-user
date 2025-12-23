package database

import (
	"sync"

	"gorm.io/gorm"
)

var (
	TenantDBs = map[string]*gorm.DB{}
	mu        sync.RWMutex
)

func GetTenantDB(tenantCode string) (*gorm.DB, bool) {
	mu.RLock()
	defer mu.RUnlock()
	db, exists := TenantDBs[tenantCode]
	return db, exists
}

func SetTenantDB(tenantCode string, db *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()
	TenantDBs[tenantCode] = db
}

func SwapTenantDB(tenantCode string, newDB *gorm.DB) *gorm.DB {
	mu.Lock()
	defer mu.Unlock()
	old := TenantDBs[tenantCode]
	TenantDBs[tenantCode] = newDB
	return old
}

func RemoveTenantDB(tenantCode string) *gorm.DB {
	mu.Lock()
	defer mu.Unlock()
	old := TenantDBs[tenantCode]
	delete(TenantDBs, tenantCode)
	return old
}
