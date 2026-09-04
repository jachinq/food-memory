package handler

import (
	"errors"
	"strconv"

	"food-memory/internal/model"
	"food-memory/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DishHandler struct {
	svc *service.DishService
}

func NewDishHandler(svc *service.DishService) *DishHandler {
	return &DishHandler{svc: svc}
}

func (h *DishHandler) List(c *gin.Context) {
	q := model.DishListQuery{
		Keyword:  c.Query("keyword"),
		Status:   c.Query("status"),
		Tag:      c.Query("tag"),
		Page:     atoiDefault(c.Query("page"), 1),
		PageSize: atoiDefault(c.Query("pageSize"), 20),
	}
	if v, ok := parseBoolQuery(c.Query("cooked")); ok {
		q.Cooked = &v
	}
	if v, ok := parseBoolQuery(c.Query("success")); ok {
		q.Success = &v
	}
	data, err := h.svc.List(q)
	if err != nil {
		ServerError(c, err.Error())
		return
	}
	OK(c, data)
}

func (h *DishHandler) Create(c *gin.Context) {
	var in model.DishCreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	dish, err := h.svc.Create(in)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	OK(c, dish)
}

func (h *DishHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	dish, err := h.svc.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "菜品不存在")
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, dish)
}

func (h *DishHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var in model.DishCreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	dish, err := h.svc.Update(id, in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "菜品不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}
	OK(c, dish)
}

func (h *DishHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "菜品不存在")
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, gin.H{"ok": true})
}

func parseID(c *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		BadRequest(c, "无效的 ID")
		return 0, false
	}
	return id, true
}

func atoiDefault(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

func parseBoolQuery(s string) (bool, bool) {
	switch s {
	case "true", "1":
		return true, true
	case "false", "0":
		return false, true
	default:
		return false, false
	}
}
