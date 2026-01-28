package enums

type UserPermission uint

const (
	UserViewer UserPermission = iota
	UserEditor
	UserOwner
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

func HasPermission(userPermission, requiredPermission UserPermission) bool {
	return userPermission >= requiredPermission
}
