package enums

type SOSContactStatus int

const (
	ContactInactive = iota
	ContactActive
)

func IsValidStatus(status SOSContactStatus) bool {
	switch status {
	case ContactInactive, ContactActive:
		return true
	}
	return false
}
