package response

type BaseResponse struct {
	Code       string      `json:"code"`
	DebugStack interface{} `json:"debug_stack"`
	Message    string      `json:"message"`
	RequestID  string      `json:"request_id"`
	Response   interface{} `json:"response"`
	Version    string      `json:"version"`
}
