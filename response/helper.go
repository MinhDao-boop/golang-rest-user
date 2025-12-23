package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, BaseResponse{
		Code:       CodeSuccess,
		DebugStack: nil,
		Message:    MsgSuccess,
		RequestID:  c.GetString("request_id"),
		Response:   data,
		Version:    "2022.11.15.20:44",
	})
}

func Error(c *gin.Context, code string, msg string, data interface{}, httpStatus int) {
	c.JSON(httpStatus, BaseResponse{
		Code:       code,
		DebugStack: nil,
		Message:    msg,
		RequestID:  c.GetString("request_id"),
		Response:   data,
		Version:    "2022.11.15.20:44",
	})
}
