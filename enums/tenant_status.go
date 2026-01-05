package enums

type TenantStatus string

const (
	TenantStatusActive   TenantStatus = "active"
	TenantStatusInactive TenantStatus = "inactive"
)

func (t TenantStatus) IsValid() bool {
	switch t {
	case TenantStatusActive, TenantStatusInactive:
		return true
	default:
		return false
	}
}
