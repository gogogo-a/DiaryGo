package v1

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/jwt"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// 微信登录请求参数
type WechatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// 微信登录响应
type WechatLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// 登录成功响应
type LoginSuccessResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

type UserLoginHandler struct {
	repo repository.UserLoginRepository
}

func NewUserLoginHandler() *UserLoginHandler {
	return &UserLoginHandler{
		repo: repository.NewUserLoginRepository(),
	}
}

// RegisterRoutes 注册用户登录相关的路由
func (h *UserLoginHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/wx-login", h.WechatLogin) // 微信登录
	}
}

// WechatLogin 处理微信登录
// @Summary 微信登录
// @Description 通过微信code获取用户信息并登录
// @Tags auth
// @Accept json
// @Produce json
// @Param code body WechatLoginRequest true "微信code"
// @Success 200 {object} response.Response{data=LoginSuccessResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/wx-login [post]
func (h *UserLoginHandler) WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "无效的请求参数")
		return
	}

	// 获取微信小程序配置
	appID := os.Getenv("WX_APP_ID")
	appSecret := os.Getenv("WX_APP_SECRET")
	if appID == "" || appSecret == "" {
		response.ServerError(c, "微信小程序配置缺失")
		return
	}
	fmt.Println("req.Code", req.Code)
	// 调用微信API获取OpenID和SessionKey
	wxResp, err := getWechatSession(appID, appSecret, req.Code)
	if err != nil {
		response.ServerError(c, fmt.Sprintf("调用微信API失败: %v", err))
		return
	}

	if wxResp.ErrCode != 0 {
		response.ServerError(c, fmt.Sprintf("微信API返回错误: %s", wxResp.ErrMsg))
		return
	}

	// 查找或创建用户
	user, err := h.findOrCreateUser(wxResp.OpenID)
	if err != nil {
		response.ServerError(c, "处理用户信息失败")
		return
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.Id)
	if err != nil {
		response.ServerError(c, "生成令牌失败")
		return
	}

	// 返回用户信息和令牌
	loginResp := LoginSuccessResponse{
		User:  user,
		Token: token,
	}

	response.Success(c, loginResp)
}

// getWechatSession 调用微信API获取会话信息
func getWechatSession(appID, appSecret, code string) (*WechatLoginResponse, error) {
	client := resty.New()
	url := "https://api.weixin.qq.com/sns/jscode2session"

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"appid":      appID,
			"secret":     appSecret,
			"js_code":    code,
			"grant_type": "authorization_code",
		}).
		Get(url)

	if err != nil {
		return nil, err
	}

	var wxResp WechatLoginResponse
	if err := json.Unmarshal(resp.Body(), &wxResp); err != nil {
		return nil, err
	}

	return &wxResp, nil
}

// findOrCreateUser 根据OpenID查找或创建用户
func (h *UserLoginHandler) findOrCreateUser(openID string) (*models.User, error) {
	// 查找用户
	user, err := h.repo.FindByOpenID(openID)
	if err == nil && user != nil {
		// 用户已存在
		return user, nil
	}

	// 创建新用户
	newUser := &models.User{
		PlantId:   openID,
		PlantForm: "wechat",
		UserName:  "微信用户", // 可以后续更新为微信昵称
	}

	// 保存用户
	if err := h.repo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
