package enums

type SosRoleKey string

const (
	FireCommand        SosRoleKey = "FIRE_COMMAND"        // Chỉ huy
	FireResponse       SosRoleKey = "FIRE_RESPONSE"       // Ứng cứu
	FireTech           SosRoleKey = "FIRE_TECH"           // Kỹ thuật
	FireControl        SosRoleKey = "FIRE_CONTROL"        // Trung tâm
	BuildingManagement SosRoleKey = "BUILDING_MANAGEMENT" // Quản lý
	Security           SosRoleKey = "SECURITY"            // An ninh
)

func IsValidRoleKey(key SosRoleKey) bool {
	switch key {
	case FireCommand, FireResponse, FireTech, FireControl, BuildingManagement, Security:
		return true
	}
	return false
}
