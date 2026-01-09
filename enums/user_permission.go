package enums

type UserPerm string

const (
	PermOwner UserPerm = "owner"
	PermRead  UserPerm = "read"
	PermWrite UserPerm = "write"
)
