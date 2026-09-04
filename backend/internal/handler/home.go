package handler

import (
	"food-memory/internal/service"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	svc *service.HomeService
}

func NewHomeHandler(svc *service.HomeService) *HomeHandler {
	return &HomeHandler{svc: svc}
}

func (h *HomeHandler) Summary(c *gin.Context) {
	data, err := h.svc.Summary()
	if err != nil {
		ServerError(c, err.Error())
		return
	}
	OK(c, data)
}
