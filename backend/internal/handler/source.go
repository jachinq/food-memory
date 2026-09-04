package handler

import (
	"errors"

	"food-memory/internal/model"
	"food-memory/internal/unfurl"

	"github.com/gin-gonic/gin"
)

type SourceHandler struct {
	svc *unfurl.Service
}

func NewSourceHandler(svc *unfurl.Service) *SourceHandler {
	return &SourceHandler{svc: svc}
}

func (h *SourceHandler) Preview(c *gin.Context) {
	var in model.SourcePreviewInput
	if err := c.ShouldBindJSON(&in); err != nil {
		BadRequest(c, "请求参数无效")
		return
	}
	got, err := h.svc.Preview(c.Request.Context(), in.URL)
	if err != nil {
		if errors.Is(err, unfurl.ErrInvalidURL) {
			BadRequest(c, "无效的链接")
			return
		}
		ServerError(c, err.Error())
		return
	}
	OK(c, model.SourcePreview{
		URL:             got.URL,
		SourcePlatform:  got.SourcePlatform,
		Name:            got.Name,
		CoverImageURL:   got.CoverImageURL,
		MainIngredients: got.MainIngredients,
		CookTimeMinutes: got.CookTimeMinutes,
		Partial:         got.Partial,
	})
}
