package enums

type UserPerm string

const (
	PermOwner  UserPerm = "owner"
	PermEditor UserPerm = "editor"
	PermViewer UserPerm = "viewer"
)
