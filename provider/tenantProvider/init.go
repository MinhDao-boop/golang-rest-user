package tenantProvider

import (
	"fmt"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/provider/serviceProvider"
)

var instance map[string]*TenantInfo

func Init() {
	service := serviceProvider.GetInstance()
	instance = make(map[string]*TenantInfo)

	service.TenantService.SetCallBackFunction(HandleTenant)
	data, err := service.TenantService.ListAllTenantConnect()
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
	temp := instance[tenantCode]
	temp.Destruction()
	instance[tenantCode] = nil
	delete(instance, tenantCode)
}

func EditInstance(tenant *models.Tenant) {
	DeleteInstance(tenant.Code)
	AddInstance(tenant)
}

func HandleTenant(mode enums.HandleTenant, tenantCode string, tenant *models.Tenant) {
	switch mode {
	case enums.AddTenantConnect:
		AddInstance(tenant)
		break
	case enums.EditTenantConnect:
		EditInstance(tenant)
		break
	case enums.DeleteTenantConnect:
		DeleteInstance(tenantCode)
		break
	default:
		fmt.Println("Cannot handle tenant mode", mode)
	}
}
