package response

import (
	"golang-rest-user/provider/configProvider"
	"strings"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	Success(c, GetSuccess, nil)
}

func Data(c *gin.Context, data *Response) {
	if data.Data == nil {
		data.Data = EmptyData{}
	}
	if data.Err != nil {
		Error(c, data.MessageCode, data.Err)
		return
	}
	Success(c, data.MessageCode, data.Data)
}

func Success(c *gin.Context, messageCode string, data interface{}) {
	config := configProvider.GetConfig()
	vn := config.GetStringMap("messages.vn")
	codeMap := config.GetStringMap("messages.code")
	code, _ := codeMap[strings.ToLower(messageCode)].(int)
	c.JSON(code, gin.H{
		"code":        messageCode,
		"debug_stack": nil,
		"message":     vn[strings.ToLower(messageCode)],
		"request_id":  c.GetString("request_id"),
		"response":    data,
		"version":     "2022.11.15.20:44",
	})
}

func Error(c *gin.Context, messageCode string, err error) {
	config := configProvider.GetConfig()
	vn := config.GetStringMap("messages.vn")
	codeMap := config.GetStringMap("messages.code")
	code, _ := codeMap[strings.ToLower(messageCode)].(int)
	c.JSON(code, gin.H{
		"code":        messageCode,
		"debug_stack": err.Error(),
		"message":     vn[strings.ToLower(messageCode)],
		"request_id":  c.GetString("request_id"),
		"response":    EmptyData{},
		"version":     "2022.11.15.20:44",
	})
}
