package v1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/haogeng/DiaryGo/internal/models"
	"github.com/haogeng/DiaryGo/internal/repository"
	"github.com/haogeng/DiaryGo/pkg/response"
)

// BillHandler 处理账单相关的请求
// 提供账单管理功能，包括创建、查询、更新、删除、查看详情和统计
// URL:
// - POST /api/v1/bills - 创建账单
// - GET /api/v1/bills - 获取账单列表，支持分页和搜索
// - PUT /api/v1/bills/:id - 更新账单
// - DELETE /api/v1/bills/:id - 删除账单
// - GET /api/v1/bills/:id - 获取账单详情
// - GET /api/v1/bills/stats - 获取账单统计
type BillHandler struct {
	repo                repository.BillRepository
	accountBookRepo     repository.AccountBookRepository
	accountBookUserRepo repository.AccountBookUserRepository
	tagRepo             repository.TagRepository
}

// NewBillHandler 创建账单处理器
// 初始化处理器并注入所需的仓库依赖
func NewBillHandler() *BillHandler {
	return &BillHandler{
		repo:                repository.NewBillRepository(),
		accountBookRepo:     repository.NewAccountBookRepository(),
		accountBookUserRepo: repository.NewAccountBookUserRepository(),
		tagRepo:             repository.NewTagRepository(),
	}
}

// RegisterRoutes 注册账单相关的路由
// 注册以下路由:
// - POST /bills - 创建账单
// - GET /bills - 获取账单列表，支持分页和搜索
// - PUT /bills/:id - 更新账单
// - DELETE /bills/:id - 删除账单
// - GET /bills/:id - 获取账单详情
// - GET /bills/stats - 获取账单统计
func (h *BillHandler) RegisterRoutes(router *gin.RouterGroup) {
	bills := router.Group("/bills")
	{
		bills.POST("", h.CreateBill)
		bills.GET("", h.ListBills)
		bills.PUT("/:id", h.UpdateBill)
		bills.DELETE("/:id", h.DeleteBill)
		bills.GET("/:id", h.GetBill)
		bills.GET("/stats", h.GetBillStats)
	}
}

// BillRequest 账单请求参数
type BillRequest struct {
	AccountBookID uuid.UUID   `json:"account_book_id" binding:"required"`
	Amount        float64     `json:"amount" binding:"required"`
	Type          string      `json:"type" binding:"required"` // income 收入 / expense 支出
	TagIDs        []uuid.UUID `json:"tag_ids" binding:"required,min=1"`
	BillTime      time.Time   `json:"bill_time" binding:"required"`
	Remark        string      `json:"remark"`
	ImageUrl      string      `json:"image_url"`
}

// BillQuery 账单查询参数
type BillQuery struct {
	AccountBookID string    `form:"account_book_id" binding:"required"`
	Page          int       `form:"page" binding:"omitempty,min=1"`
	PageSize      int       `form:"page_size" binding:"omitempty,min=1,max=100"`
	Type          string    `form:"type"`
	TagIDs        []string  `form:"tag_ids"`
	StartTime     time.Time `form:"start_time" time_format:"2006-01-02"`
	EndTime       time.Time `form:"end_time" time_format:"2006-01-02"`
	MinAmount     float64   `form:"min_amount"`
	MaxAmount     float64   `form:"max_amount"`
	Keyword       string    `form:"keyword"`
}

// StatsQuery 统计查询参数
type StatsQuery struct {
	AccountBookID string    `form:"account_book_id" binding:"required"`
	StartTime     time.Time `form:"start_time" time_format:"2006-01-02"`
	EndTime       time.Time `form:"end_time" time_format:"2006-01-02"`
	GroupBy       string    `form:"group_by" binding:"omitempty,oneof=day week month year"`
}

// BillResponse 账单响应结构
type BillResponse struct {
	Bill models.Bill  `json:"bill"`
	Tags []models.Tag `json:"tags"`
}

// CreateBill 创建账单
// URL: POST /api/v1/bills
// 功能: 创建新账单
// 权限: 需要用户登录，且必须有对应账本的访问权限
// 请求体:
//
//	{
//	  "account_book_id": "账本UUID",
//	  "amount": 100.00,
//	  "type": "income/expense",
//	  "tag_ids": ["标签UUID1", "标签UUID2"],
//	  "bill_time": "2023-01-01T00:00:00Z",
//	  "remark": "备注",
//	  "image_url": "图片URL"
//	}
//
// 返回:
//   - 成功: 200 状态码，包含创建的账单信息
//   - 失败: 400 (参数错误), 401 (未授权), 403 (无权限), 500 (服务器错误)
//
// @Summary 创建账单
// @Description 创建新账单，需要有对应账本的访问权限
// @Tags bills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BillRequest true "账单信息"
// @Success 200 {object} response.Response{data=models.Bill}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/bills [post]
func (h *BillHandler) CreateBill(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析请求
	var req BillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 检查用户是否有权限访问该账本
	_, err := h.accountBookUserRepo.GetByAccountBookIDAndUserID(req.AccountBookID, userID.(uuid.UUID))
	if err != nil {
		response.Forbidden(c, "您没有权限访问此账本")
		return
	}

	// 检查所有标签是否存在
	for _, tagID := range req.TagIDs {
		_, err = h.tagRepo.GetByID(tagID)
		if err != nil {
			response.NotFound(c, "标签不存在: "+tagID.String())
			return
		}
	}

	// 创建账单
	bill := &models.Bill{
		AccountBookId: req.AccountBookID,
		UserId:        userID.(uuid.UUID),
		Amount:        req.Amount,
		Type:          req.Type,
		Remark:        req.Remark,
		ImageUrl:      req.ImageUrl,
	}

	if err := h.repo.Create(bill, req.TagIDs); err != nil {
		response.ServerError(c, "创建账单失败")
		return
	}

	// 获取完整的账单信息（包括标签）
	billWithTags, tags, err := h.repo.GetBillWithTags(bill.Id)
	if err != nil {
		response.SuccessWithMessage(c, "创建账单成功，但获取详情失败", bill)
		return
	}

	response.SuccessWithMessage(c, "创建账单成功", BillResponse{
		Bill: *billWithTags,
		Tags: tags,
	})
}

// ListBills 获取账单列表
// URL: GET /api/v1/bills
// 功能: 获取账单列表，支持分页和多条件搜索
// 权限: 需要用户登录，且必须有对应账本的访问权限
// 查询参数:
//   - account_book_id: 账本ID (必填)
//   - page: 页码，默认为1
//   - page_size: 每页数量，默认为10，最大100
//   - type: 账单类型 (income/expense)
//   - tag_ids: 标签ID列表
//   - start_time: 开始时间 (YYYY-MM-DD)
//   - end_time: 结束时间 (YYYY-MM-DD)
//   - min_amount: 最小金额
//   - max_amount: 最大金额
//   - keyword: 搜索关键词 (搜索备注)
//
// 返回:
//   - 成功: 200 状态码，包含账单列表和分页信息
//   - 失败: 400 (参数错误), 401 (未授权), 403 (无权限), 500 (服务器错误)
//
// @Summary 获取账单列表
// @Description 获取账单列表，支持分页和多条件搜索
// @Tags bills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account_book_id query string true "账本ID"
// @Param page query int false "页码，默认为1"
// @Param page_size query int false "每页数量，默认为10，最大100"
// @Param type query string false "账单类型 (income/expense)"
// @Param tag_ids query []string false "标签ID列表"
// @Param start_time query string false "开始时间 (YYYY-MM-DD)"
// @Param end_time query string false "结束时间 (YYYY-MM-DD)"
// @Param min_amount query number false "最小金额"
// @Param max_amount query number false "最大金额"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PagedData{list=[]models.Bill}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/bills [get]
func (h *BillHandler) ListBills(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析查询参数
	var query BillQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 解析账本ID
	accountBookID, err := uuid.Parse(query.AccountBookID)
	if err != nil {
		response.ParamError(c, "无效的账本ID格式")
		return
	}

	// 检查用户是否有权限访问该账本
	_, err = h.accountBookUserRepo.GetByAccountBookIDAndUserID(accountBookID, userID.(uuid.UUID))
	if err != nil {
		response.Forbidden(c, "您没有权限访问此账本")
		return
	}

	// 解析标签ID（如果有）
	var tagIDs []uuid.UUID
	for _, tagIDStr := range query.TagIDs {
		tagID, err := uuid.Parse(tagIDStr)
		if err != nil {
			response.ParamError(c, "无效的标签ID格式: "+tagIDStr)
			return
		}
		tagIDs = append(tagIDs, tagID)
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	} else if query.PageSize > 100 {
		query.PageSize = 100
	}

	// 获取账单列表
	bills, total, err := h.repo.GetBills(
		accountBookID,
		query.Page,
		query.PageSize,
		query.Type,
		tagIDs,
		query.StartTime,
		query.EndTime,
		query.MinAmount,
		query.MaxAmount,
		query.Keyword,
	)
	if err != nil {
		response.ServerError(c, "获取账单列表失败")
		return
	}

	// 构造分页响应
	pagedData := response.PagedData{
		List:     bills,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	response.Success(c, pagedData)
}

// UpdateBill 更新账单
// URL: PUT /api/v1/bills/:id
// 功能: 更新账单信息
// 权限: 需要用户登录，且必须有对应账本的访问权限
// 参数:
//   - id: 账单ID (UUID格式，通过URL路径参数传递)
//
// 请求体:
//
//	{
//	  "account_book_id": "账本UUID",
//	  "amount": 100.00,
//	  "type": "income/expense",
//	  "tag_ids": ["标签UUID1", "标签UUID2"],
//	  "bill_time": "2023-01-01T00:00:00Z",
//	  "remark": "备注",
//	  "image_url": "图片URL"
//	}
//
// 返回:
//   - 成功: 200 状态码，包含更新后的账单信息
//   - 失败: 400 (参数错误), 401 (未授权), 403 (无权限), 404 (账单不存在), 500 (服务器错误)
//
// @Summary 更新账单
// @Description 更新账单信息，需要有对应账本的访问权限
// @Tags bills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "账单ID"
// @Param request body BillRequest true "账单信息"
// @Success 200 {object} response.Response{data=models.Bill}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/bills/{id} [put]
func (h *BillHandler) UpdateBill(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析账单ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的账单ID格式")
		return
	}

	// 获取原始账单信息
	originalBill, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(c, "账单不存在")
		return
	}

	// 检查用户是否有权限访问该账本
	_, err = h.accountBookUserRepo.GetByAccountBookIDAndUserID(originalBill.AccountBookId, userID.(uuid.UUID))
	if err != nil {
		response.Forbidden(c, "您没有权限访问此账本")
		return
	}

	// 解析请求
	var req BillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 检查新账本ID是否与原账本ID不同，如果不同，需要检查用户是否有权限访问新账本
	if req.AccountBookID != originalBill.AccountBookId {
		_, err = h.accountBookUserRepo.GetByAccountBookIDAndUserID(req.AccountBookID, userID.(uuid.UUID))
		if err != nil {
			response.Forbidden(c, "您没有权限访问目标账本")
			return
		}
	}

	// 检查所有标签是否存在
	for _, tagID := range req.TagIDs {
		_, err = h.tagRepo.GetByID(tagID)
		if err != nil {
			response.NotFound(c, "标签不存在: "+tagID.String())
			return
		}
	}

	// 更新账单信息
	originalBill.AccountBookId = req.AccountBookID
	originalBill.Amount = req.Amount
	originalBill.Type = req.Type
	originalBill.Remark = req.Remark
	originalBill.ImageUrl = req.ImageUrl

	if err := h.repo.Update(originalBill, req.TagIDs); err != nil {
		response.ServerError(c, "更新账单失败")
		return
	}

	// 获取完整的账单信息（包括标签）
	billWithTags, tags, err := h.repo.GetBillWithTags(originalBill.Id)
	if err != nil {
		response.SuccessWithMessage(c, "更新账单成功，但获取详情失败", originalBill)
		return
	}

	response.SuccessWithMessage(c, "更新账单成功", BillResponse{
		Bill: *billWithTags,
		Tags: tags,
	})
}

// DeleteBill 删除账单
// URL: DELETE /api/v1/bills/:id
// 功能: 删除账单
// 权限: 需要用户登录，且必须有对应账本的访问权限
// 参数:
//   - id: 账单ID (UUID格式，通过URL路径参数传递)
//
// 返回:
//   - 成功: 200 状态码，删除成功消息
//   - 失败: 400 (参数错误), 401 (未授权), 403 (无权限), 404 (账单不存在), 500 (服务器错误)
//
// @Summary 删除账单
// @Description 删除账单，需要有对应账本的访问权限
// @Tags bills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "账单ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/bills/{id} [delete]
func (h *BillHandler) DeleteBill(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析账单ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的账单ID格式")
		return
	}

	// 获取账单信息
	bill, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(c, "账单不存在")
		return
	}

	// 检查用户是否有权限访问该账本
	_, err = h.accountBookUserRepo.GetByAccountBookIDAndUserID(bill.AccountBookId, userID.(uuid.UUID))
	if err != nil {
		response.Forbidden(c, "您没有权限访问此账本")
		return
	}

	// 删除账单
	if err := h.repo.Delete(id); err != nil {
		response.ServerError(c, "删除账单失败")
		return
	}

	response.SuccessWithMessage(c, "删除账单成功", nil)
}

// GetBill 获取账单详情
// URL: GET /api/v1/bills/:id
// 功能: 获取账单详细信息
// 权限: 需要用户登录，且必须有对应账本的访问权限
// 参数:
//   - id: 账单ID (UUID格式，通过URL路径参数传递)
//
// 返回:
//   - 成功: 200 状态码，包含账单详情的JSON数据
//   - 失败: 400 (参数错误), 401 (未授权), 403 (无权限), 404 (账单不存在)
//
// @Summary 获取账单详情
// @Description 获取账单详细信息，需要有对应账本的访问权限
// @Tags bills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "账单ID"
// @Success 200 {object} response.Response{data=models.Bill}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/bills/{id} [get]
func (h *BillHandler) GetBill(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析账单ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.ParamError(c, "无效的账单ID格式")
		return
	}

	// 获取账单信息
	bill, err := h.repo.GetByID(id)
	if err != nil {
		response.NotFound(c, "账单不存在")
		return
	}

	// 检查用户是否有权限访问该账本
	_, err = h.accountBookUserRepo.GetByAccountBookIDAndUserID(bill.AccountBookId, userID.(uuid.UUID))
	if err != nil {
		response.Forbidden(c, "您没有权限访问此账本")
		return
	}

	// 获取账单及其标签
	billWithTags, tags, err := h.repo.GetBillWithTags(id)
	if err != nil {
		response.ServerError(c, "获取账单详情失败")
		return
	}

	response.Success(c, BillResponse{
		Bill: *billWithTags,
		Tags: tags,
	})
}

// GetBillStats 获取账单统计
// URL: GET /api/v1/bills/stats
// 功能: 获取账单统计信息
// 权限: 需要用户登录，且必须有对应账本的访问权限
// 查询参数:
//   - account_book_id: 账本ID (必填)
//   - start_time: 开始时间 (YYYY-MM-DD)
//   - end_time: 结束时间 (YYYY-MM-DD)
//   - group_by: 分组方式 (day/week/month/year)
//
// 返回:
//   - 成功: 200 状态码，包含账单统计信息
//   - 失败: 400 (参数错误), 401 (未授权), 403 (无权限), 500 (服务器错误)
//
// @Summary 获取账单统计
// @Description 获取账单统计信息，需要有对应账本的访问权限
// @Tags bills
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account_book_id query string true "账本ID"
// @Param start_time query string false "开始时间 (YYYY-MM-DD)"
// @Param end_time query string false "结束时间 (YYYY-MM-DD)"
// @Param group_by query string false "分组方式 (day/week/month/year)"
// @Success 200 {object} response.Response{data=repository.BillStats}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/bills/stats [get]
func (h *BillHandler) GetBillStats(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未找到用户信息")
		return
	}

	// 解析查询参数
	var query StatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 解析账本ID
	accountBookID, err := uuid.Parse(query.AccountBookID)
	if err != nil {
		response.ParamError(c, "无效的账本ID格式")
		return
	}

	// 检查用户是否有权限访问该账本
	_, err = h.accountBookUserRepo.GetByAccountBookIDAndUserID(accountBookID, userID.(uuid.UUID))
	if err != nil {
		response.Forbidden(c, "您没有权限访问此账本")
		return
	}

	// 获取账单统计
	stats, err := h.repo.GetStats(accountBookID, query.StartTime, query.EndTime, query.GroupBy)
	if err != nil {
		response.ServerError(c, "获取账单统计失败")
		return
	}

	response.Success(c, stats)
}
