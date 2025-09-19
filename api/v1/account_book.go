package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

type AccountBookHandler struct {
	repo                repository.AccountBookRepository
	userRepo            repository.UserLoginRepository
	accountBookUserRepo repository.AccountBookUserRepository
}

func NewAccountBookHandler() *AccountBookHandler {
	return &AccountBookHandler{
		repo:                repository.NewAccountBookRepository(),
		userRepo:            repository.NewUserLoginRepository(),
		accountBookUserRepo: repository.NewAccountBookUserRepository(),
	}
}
func (h *AccountBookHandler) RegisterRoutes(router *gin.RouterGroup) {
	accountBooks := router.Group("/account-books")
	{
		accountBooks.POST("", h.Create)       //创建账本
		accountBooks.GET("", h.List)          //获取该用户所有账本
		accountBooks.GET("/:id", h.Get)       //根据该账本id获取账本
		accountBooks.PUT("/:id", h.Update)    //更新账本的名称
		accountBooks.DELETE("/:id", h.Delete) //删除账本
	}
}

func (h *AccountBookHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	var accountBook models.AccountBook
	if err := c.ShouldBindJSON(&accountBook); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	user, err := h.userRepo.GetByID(userID.(uuid.UUID))
	if err != nil {
		response.ServerError(c, "获取用户信息失败")
		return
	}

	// 先创建账本
	if err := h.repo.Create(&accountBook); err != nil {
		response.ServerError(c, "创建账本失败")
		return
	}

	// 再创建账本与用户的关联
	accountBookUser := models.AccountBookUser{
		AccountBookId: accountBook.Id,
		UserId:        user.Id,
	}
	if err := h.accountBookUserRepo.Create(&accountBookUser); err != nil {
		response.ServerError(c, "创建账本用户关联失败")
		return
	}

	response.SuccessWithMessage(c, "创建账本成功", accountBook)
}

func (h *AccountBookHandler) List(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}
	accountBooks, err := h.repo.GetAll(userID.(uuid.UUID))
	if err != nil {
		response.ServerError(c, "获取账本失败")
		return
	}
	response.Success(c, accountBooks)
}

func (h *AccountBookHandler) Get(c *gin.Context) {
	id := c.Param("id")
	accountBook, err := h.repo.GetByID(uuid.MustParse(id))
	if err != nil {
		response.ServerError(c, "获取账本失败")
		return
	}
	response.Success(c, accountBook)
}

func (h *AccountBookHandler) Update(c *gin.Context) {
	// 先获取id
	// 在获取用户id
	// 在account_book_user中检验权限
	// 有权限改名字，没权限返回错误信息
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析ID
	accountBookID, err := uuid.Parse(id)
	if err != nil {
		response.ParamError(c, "无效的账本ID格式")
		return
	}

	// 获取账本用户关系，检验权限
	accountBookUser, err := h.accountBookUserRepo.GetByAccountBookID(accountBookID)
	if err != nil {
		response.ServerError(c, "获取账本用户关系失败")
		return
	}

	// 验证用户是否有权限更新此账本
	if accountBookUser.UserId != userID.(uuid.UUID) {
		response.Forbidden(c, "您没有权限更新此账本")
		return
	}

	// 从请求体获取新的名称
	var nameUpdate struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&nameUpdate); err != nil {
		response.ParamError(c, "无效的名称参数")
		return
	}

	// 获取原始账本
	originalAccountBook := accountBookUser.AccountBook

	// 只更新名称
	originalAccountBook.Name = nameUpdate.Name

	// 保存更新
	if err := h.repo.Update(&originalAccountBook); err != nil {
		response.ServerError(c, "更新账本失败")
		return
	}

	response.SuccessWithMessage(c, "更新账本名称成功", originalAccountBook)
}

func (h *AccountBookHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(uuid.MustParse(id)); err != nil {
		response.ServerError(c, "删除账本失败")
		return
	}
	response.SuccessWithMessage(c, "删除账本成功", nil)
}
