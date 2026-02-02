package response

type Response struct {
	Err         error
	Data        interface{}
	MessageCode string
}

type EmptyData struct{}

type ListResponse struct {
	Page     int         `form:"page,default=1" json:"page"`
	PageSize int         `form:"page_size,default=50" json:"page_size"`
	Total    int64       `json:"total"`
	Contents interface{} `form:"contents,default=50" json:"contents"`
}

type DeleteResponse struct {
	Deleted int64 `json:"deleted"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	AccessExpiresIn  int    `json:"access_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

type NumberResponse struct {
	NumberAffected int64 `json:"number_affected"`
}

func NewResponse() *Response {
	return &Response{
		Err:         nil,
		Data:        nil,
		MessageCode: ERR0001,
	}
}
