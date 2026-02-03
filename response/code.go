package response

import "errors"

const (
	// TEN0001: "Cần cung cấp thông tin định danh"
	TEN0001 = "TEN0001"

	//ERR0001 = "Có lỗi xảy ra, vui lòng thử lại"
	ERR0001 = "ERR0001"

	//FBD0000 = "Không có quyền truy cập"
	FBD0000 = "FBD0000"

	//AUT0003 = "Tên đăng nhập hoặc mật khẩu không đúng"
	AUT0003 = "AUT0003"

	//AUT0006 = "Không được phép truy cập"
	AUT0006 = "AUT0006"

	//REG0007 = "Tài khoản đã tồn tại"
	REG0007 = "REG0007"

	//SUS0000 = "Thành công"
	SUS0000 = "SUS0000"

	//VAD0000: "Đầu vào không hợp lệ"
	VAD0000 = "VAD0000"
)

const (
	GetError   = "Lấy dữ liệu lỗi"
	GetSuccess = "Lấy dữ liệu thành công"
)

var (
	ErrMissingIdentifier = errors.New("identifier is required")

	// ErrInvalidKey indicates that the encryption key is unreadable
	ErrInvalidKey = errors.New("invalid encryption key")

	// ErrShortCipher indicates that the ciphertext is too short
	ErrShortCipher = errors.New("ciphertext too short")

	// ErrForbidden indicates that the user does not have the required permission
	ErrForbidden = errors.New("permission denied")

	ErrUnauthorized = errors.New("unauthorized")

	// ErrInvalidPermission indicates that the client requests an invalid permission
	ErrInvalidPermission = errors.New("invalid permission")

	// ErrExistingUsername indicates that the client requests an existing username
	ErrExistingUsername = errors.New("username already exists")

	// ErrExistingDBName indicates that the client requests an existing database name
	ErrExistingDBName = errors.New("database name already exists")

	// ErrInvalidDBName indicates that the client requests an invalid database name
	ErrInvalidDBName = errors.New("invalid database name")

	// ErrExistingTenantCode indicates that the client requests an existing database name
	ErrExistingTenantCode = errors.New("tenant code already exists")

	// ErrInvalidStatus indicates that the client requests an invalid status
	ErrInvalidStatus = errors.New("invalid status")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect username or password")

	// ErrTokenType indicates an incorrect token type is requested
	ErrTokenType = errors.New("wrong token type")

	// ErrInvalidToken happens when an invalid token is requested
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenRevoked indicates that the client requests a revoked token
	ErrTokenRevoked = errors.New("token revoked")

	// ErrTokenExpired indicates that the client requests an expired token
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidSharing indicates that users share the zone to themself
	ErrInvalidSharing = errors.New("invalid sharing")

	// ErrInvalidJsonObj indicates that the client sends an invalid JSON object
	ErrInvalidJsonObj = errors.New("invalid json object")

	// ErrInvalidJsonArray indicates that the client sends an invalid JSON array
	ErrInvalidJsonArray = errors.New("invalid json array")
)
