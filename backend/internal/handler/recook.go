package handler

import (
	"errors"

	"food-memory/internal/model"
	"food-memory/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecookHandler struct {
	svc *service.RecookService
}

func NewRecookHandler(svc *service.RecookService) *RecookHandler {
	return &RecookHandler{svc: svc}
}

func (h *RecookHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Query("status"))
	if err != nil {
		ServerError(c, err.Error())
		return
	}
	OK(c, items)
}

func (h *RecookHandler) Create(c *gin.Context) {
	var in model.RecookPlanInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	plan, err := h.svc.Create(in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "菜品不存在")
			return
		}
		if errors.Is(err, service.ErrRecookAlreadyActive) {
			BadRequest(c, err.Error())
			return
		}
		BadRequest(c, err.Error())
		return
	}
	OK(c, plan)
}

func (h *RecookHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var in model.RecookPlanUpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	plan, err := h.svc.Update(id, in)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "复做计划不存在")
			return
		}
		BadRequest(c, err.Error())
		return
	}
	OK(c, plan)
}

func (h *RecookHandler) Complete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	plan, err := h.svc.Complete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "复做计划不存在")
			return
		}
		if errors.Is(err, service.ErrEmptyRecookComplete) {
			BadRequest(c, err.Error())
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, plan)
}

func (h *RecookHandler) Cancel(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Cancel(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "复做计划不存在")
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, gin.H{"ok": true})
}

func (h *RecookHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			NotFound(c, "复做计划不存在")
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, gin.H{"ok": true})
}
