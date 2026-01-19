package enums

type UserPermission string

const (
	UserOwner  UserPermission = "owner"
	UserEditor UserPermission = "editor"
	UserViewer UserPermission = "viewer"
)

func IsValidUserPermission(permission UserPermission) bool {
	switch permission {
	case UserOwner:
		return true
	case UserEditor:
		return true
	case UserViewer:
		return true
	default:
		return false
	}
}
