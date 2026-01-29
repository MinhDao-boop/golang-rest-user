package response

type BaseResponse struct {
	Code       string      `json:"code"`
	DebugStack interface{} `json:"debug_stack"`
	Message    string      `json:"message"`
	RequestID  string      `json:"request_id"`
	Response   interface{} `json:"response"`
	Version    string      `json:"version"`
}

type Response struct {
	Err         error
	Data        interface{}
	MessageCode string
}

type EmptyData struct{}

type ListResponse struct {
	Page     int         `form:"page,default=0" json:"page"`
	PageSize int         `form:"page_size,default=50" json:"page_size"`
	Total    int64       `json:"total"`
	Contents interface{} `form:"contents,default=50" json:"contents"`
}

type NumberResponse struct {
	NumberAffected int64 `json:"number_affected"`
}

func NewResponse() *Response {
	return &Response{
		Err:         nil,
		Data:        nil,
		MessageCode: CodeBadRequest,
	}
}
