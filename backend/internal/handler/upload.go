package handler

import (
	"strconv"

	"food-memory/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	svc *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		BadRequest(c, "请选择要上传的图片")
		return
	}
	defer file.Close()

	bizType := c.PostForm("biz_type")
	var bizID uint64
	if v := c.PostForm("biz_id"); v != "" {
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			BadRequest(c, "无效的业务 ID")
			return
		}
		bizID = n
	}
	att, err := h.svc.Save(file, header, bizType, bizID)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	OK(c, att)
}
