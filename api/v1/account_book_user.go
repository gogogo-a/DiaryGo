package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// AccountBookUserHandler 处理账本用户权限相关的请求
type AccountBookUserHandler struct {
	repo            repository.AccountBookUserRepository
	accountBookRepo repository.AccountBookRepository
	userRepo        repository.UserLoginRepository
}

// NewAccountBookUserHandler 创建账本用户处理器
func NewAccountBookUserHandler() *AccountBookUserHandler {
	return &AccountBookUserHandler{
		repo:            repository.NewAccountBookUserRepository(),
		accountBookRepo: repository.NewAccountBookRepository(),
		userRepo:        repository.NewUserLoginRepository(),
	}
}

// RegisterRoutes 注册账本用户相关的路由
func (h *AccountBookUserHandler) RegisterRoutes(router *gin.RouterGroup) {
	accountBookUsers := router.Group("/account-book-users")
	{
		// 传入账本id查看该账本的用户
		accountBookUsers.GET("/book/:bookId", h.GetUsersByBookID)

		// 赋予该账本的权限给用户 第一个查出的用户id为管理员才能赋予权限
		accountBookUsers.POST("/grant", h.GrantPermission)

		// 删除该账本的权限给用户 第一个查出的用户id为管理员才能删除权限
		accountBookUsers.DELETE("/revoke", h.RevokePermission)
	}
}

// GetUsersByBookID 获取账本的所有用户
// @Summary 获取账本的所有用户
// @Description 根据账本ID获取所有有权限的用户
// @Tags account-book-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bookId path string true "账本ID"
// @Success 200 {object} response.Response{data=[]models.User}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/account-book-users/book/{bookId} [get]
func (h *AccountBookUserHandler) GetUsersByBookID(c *gin.Context) {
	// 获取账本ID
	bookIdStr := c.Param("bookId")
	bookId, err := uuid.Parse(bookIdStr)
	if err != nil {
		response.ParamError(c, "无效的账本ID格式")
		return
	}

	// 获取当前用户ID
	_, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 检查账本是否存在
	_, err = h.repo.GetByAccountBookID(bookId)
	if err != nil {
		response.NotFound(c, "账本不存在或您没有权限")
		return
	}

	// 获取账本的所有用户
	users, err := h.repo.GetAllUsersByAccountBookID(bookId)
	if err != nil {
		response.ServerError(c, "获取账本用户列表失败")
		return
	}

	response.Success(c, users)
}

// GrantPermissionRequest 授予权限请求
type GrantPermissionRequest struct {
	AccountBookID uuid.UUID `json:"account_book_id" binding:"required"`
	UserID        uuid.UUID `json:"user_id" binding:"required"`
}

// GrantPermission 授予用户账本权限
// @Summary 授予用户账本权限
// @Description 管理员为用户授予账本访问权限
// @Tags account-book-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GrantPermissionRequest true "授权请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/account-book-users/grant [post]
func (h *AccountBookUserHandler) GrantPermission(c *gin.Context) {
	// 获取当前用户ID
	currentUserID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析请求
	var req GrantPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 检查当前用户是否是账本的管理员
	// 获取账本的第一个用户（管理员）
	accountBookUser, err := h.repo.GetByAccountBookID(req.AccountBookID)
	if err != nil {
		response.NotFound(c, "账本不存在")
		return
	}

	// 验证当前用户是否为管理员
	if accountBookUser.UserId != currentUserID.(uuid.UUID) {
		response.Forbidden(c, "您没有权限授予此账本的访问权限")
		return
	}

	// 检查要授权的用户是否存在
	_, err = h.userRepo.GetByID(req.UserID)
	if err != nil {
		response.NotFound(c, "目标用户不存在")
		return
	}

	// 检查用户是否已经有权限
	existingPermission, err := h.repo.GetByAccountBookIDAndUserID(req.AccountBookID, req.UserID)
	if err == nil && existingPermission != nil {
		response.SuccessWithMessage(c, "用户已有此账本的权限", nil)
		return
	}

	// 创建新的账本用户关联
	newAccountBookUser := &models.AccountBookUser{
		AccountBookId: req.AccountBookID,
		UserId:        req.UserID,
	}

	if err := h.repo.Create(newAccountBookUser); err != nil {
		response.ServerError(c, "授予权限失败")
		return
	}

	response.SuccessWithMessage(c, "成功授予用户访问权限", nil)
}

// RevokePermissionRequest 撤销权限请求
type RevokePermissionRequest struct {
	AccountBookID uuid.UUID `json:"account_book_id" binding:"required"`
	UserID        uuid.UUID `json:"user_id" binding:"required"`
}

// RevokePermission 撤销用户账本权限
// @Summary 撤销用户账本权限
// @Description 管理员撤销用户的账本访问权限
// @Tags account-book-users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RevokePermissionRequest true "撤销权限请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/account-book-users/revoke [delete]
func (h *AccountBookUserHandler) RevokePermission(c *gin.Context) {
	// 获取当前用户ID
	currentUserID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析请求
	var req RevokePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 检查当前用户是否是账本的管理员
	// 获取账本的第一个用户（管理员）
	accountBookUser, err := h.repo.GetByAccountBookID(req.AccountBookID)
	if err != nil {
		response.NotFound(c, "账本不存在")
		return
	}

	// 验证当前用户是否为管理员
	if accountBookUser.UserId != currentUserID.(uuid.UUID) {
		response.Forbidden(c, "您没有权限撤销此账本的访问权限")
		return
	}

	// 检查用户是否有此账本的权限
	_, err = h.repo.GetByAccountBookIDAndUserID(req.AccountBookID, req.UserID)
	if err != nil {
		response.NotFound(c, "该用户没有此账本的权限")
		return
	}

	// 删除特定用户的特定账本权限
	if err := h.repo.DeleteByAccountBookIDAndUserID(req.AccountBookID, req.UserID); err != nil {
		response.ServerError(c, "撤销权限失败")
		return
	}

	response.SuccessWithMessage(c, "成功撤销用户访问权限", nil)
}
