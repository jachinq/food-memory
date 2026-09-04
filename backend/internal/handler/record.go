package handler

import (
	"errors"

	"food-memory/internal/model"
	"food-memory/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecordHandler struct {
	svc *service.RecordService
}

func NewRecordHandler(svc *service.RecordService) *RecordHandler {
	return &RecordHandler{svc: svc}
}

func (h *RecordHandler) List(c *gin.Context) {
	dishID, ok := parseID(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.List(dishID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "菜品不存在")
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, items)
}

func (h *RecordHandler) Create(c *gin.Context) {
	dishID, ok := parseID(c, "id")
	if !ok {
		return
	}
	var in model.RecordInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	data, err := h.svc.Create(dishID, in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "菜品不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}
	OK(c, data)
}

func (h *RecordHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var in model.RecordInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	data, err := h.svc.Update(id, in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "制作记录不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}
	OK(c, data)
}

func (h *RecordHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "制作记录不存在")
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, gin.H{"ok": true})
}
