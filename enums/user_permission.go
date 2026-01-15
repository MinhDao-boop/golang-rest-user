package enums

type UserPermission string

const (
	UserOwner  UserPermission = "owner"
	UserEditor UserPermission = "editor"
	UserViewer UserPermission = "viewer"
)

func IsValidUserPermission(permission string) bool {
	switch permission {
	case "owner":
		return true
	case "editor":
		return true
	case "viewer":
		return true
	default:
		return false
	}
}
