package handler

import (
	"food-memory/internal/model"
	"food-memory/internal/service"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	svc *service.TagService
}

func NewTagHandler(svc *service.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

func (h *TagHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Query("type"), c.Query("keyword"))
	if err != nil {
		ServerError(c, err.Error())
		return
	}
	OK(c, items)
}

func (h *TagHandler) Create(c *gin.Context) {
	var in model.TagInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	tag, err := h.svc.Create(in.Name, in.Type)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	OK(c, tag)
}
