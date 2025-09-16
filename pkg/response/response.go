package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义常用的状态码
const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
	NOT_FOUND      = 404
	UNAUTHORIZED   = 401
	FORBIDDEN      = 403
)

// 状态码对应的消息
var codeMessages = map[int]string{
	SUCCESS:        "成功",
	ERROR:          "服务器内部错误",
	INVALID_PARAMS: "请求参数错误",
	NOT_FOUND:      "资源不存在",
	UNAUTHORIZED:   "未授权",
	FORBIDDEN:      "禁止访问",
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PagedData 分页数据结构
type PagedData struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页记录数
}

// Result 返回自定义响应
func Result(c *gin.Context, code int, message string, data interface{}) {
	// 如果没有提供消息，则使用默认消息
	if message == "" {
		message = codeMessages[code]
	}

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, "", data)
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	Result(c, SUCCESS, message, data)
}

// Fail 失败响应
func Fail(c *gin.Context, code int, data interface{}) {
	Result(c, code, "", data)
}

// FailWithMessage 带消息的失败响应
func FailWithMessage(c *gin.Context, code int, message string, data interface{}) {
	Result(c, code, message, data)
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	Result(c, INVALID_PARAMS, message, nil)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	Result(c, ERROR, message, nil)
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	Result(c, NOT_FOUND, message, nil)
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	Result(c, UNAUTHORIZED, message, nil)
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	Result(c, FORBIDDEN, message, nil)
}
