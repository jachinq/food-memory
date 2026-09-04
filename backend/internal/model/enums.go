package model

const (
	StatusWantToCook   = "want_to_cook"
	StatusCooked       = "cooked"
	StatusSuccess      = "success"
	StatusFailed       = "failed"
	StatusWantToRecook = "want_to_recook"
	StatusPaused       = "paused"

	ResultSuccess = "success"
	ResultNormal  = "normal"
	ResultFailed  = "failed"

	TagIngredient = "ingredient"
	TagTaste      = "taste"
	TagScene      = "scene"
	TagCuisine    = "cuisine"
	TagTool       = "tool"
	TagCustom     = "custom"

	BizDishCover         = "dish_cover"
	BizCookRecord        = "cook_record"
	BizSourceScreenshot  = "source_screenshot"

	RecookActive    = "active"
	RecookCompleted = "completed"
	RecookCancelled = "cancelled"
)

func ValidDishStatus(s string) bool {
	switch s {
	case StatusWantToCook, StatusCooked, StatusSuccess, StatusFailed, StatusWantToRecook, StatusPaused:
		return true
	}
	return false
}

func ValidResult(s string) bool {
	switch s {
	case ResultSuccess, ResultNormal, ResultFailed:
		return true
	}
	return false
}

func ValidTagType(s string) bool {
	switch s {
	case TagIngredient, TagTaste, TagScene, TagCuisine, TagTool, TagCustom:
		return true
	}
	return false
}

func StatusLabel(s string) string {
	switch s {
	case StatusWantToCook:
		return "想做"
	case StatusCooked:
		return "已做"
	case StatusSuccess:
		return "做成功"
	case StatusFailed:
		return "翻车"
	case StatusWantToRecook:
		return "想再做"
	case StatusPaused:
		return "暂不做"
	default:
		return s
	}
}
