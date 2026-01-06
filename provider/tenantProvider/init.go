package tenantProvider

import (
	"golang-rest-user/models"
	"golang-rest-user/service/tenantSvc"
)

var instance map[string]*TenantInfo

func Init(tenantConnectSvc tenantSvc.TenantConnect) {
	instance = make(map[string]*TenantInfo)

	data, err := tenantConnectSvc.ListAllTenantConnect()
	if err == nil {
		for _, item := range data {
			temp := &TenantInfo{
				Info: &item,
			}
			_ = temp.Init()
			instance[item.Code] = temp
		}
	}
}

func GetTenantInfo(tenantCode string) *TenantInfo {
	return instance[tenantCode]
}

func AddInstance(tenant *models.Tenant) {
	temp := &TenantInfo{
		Info: tenant,
	}
	_ = temp.Init()
	instance[tenant.Code] = temp
}

func DeleteInstance(tenantCode string) {
	instance[tenantCode] = nil
	delete(instance, tenantCode)
}

func EditInstance(tenant *models.Tenant) {
	DeleteInstance(tenant.Code)
	AddInstance(tenant)
}
