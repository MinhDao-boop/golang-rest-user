package enums

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

func (t TokenType) IsValid() bool {
	switch t {
	case TokenTypeAccess, TokenTypeRefresh:
		return true
	default:
		return false
	}
}
