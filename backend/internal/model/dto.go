package model

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

type TagInput struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DishCreateInput struct {
	Name              string     `json:"name"`
	CoverImageURL     string     `json:"cover_image_url"`
	SourceURL         string     `json:"source_url"`
	SourcePlatform    string     `json:"source_platform"`
	Description       string     `json:"description"`
	Status            string     `json:"status"`
	Rating            *float64   `json:"rating"`
	Difficulty        *int16     `json:"difficulty"`
	CookTimeMinutes   *int       `json:"cook_time_minutes"`
	MainIngredients   string     `json:"main_ingredients"`
	Taste             string     `json:"taste"`
	Scene             string     `json:"scene"`
	Note              string     `json:"note"`
	IsFavorite        *bool      `json:"is_favorite"`
	Tags              []TagInput `json:"tags"`
	CoverAttachmentID *uint64    `json:"cover_attachment_id"`
}

type DishUpdateInput struct {
	DishCreateInput
}

type DishListQuery struct {
	Keyword  string
	Status   string
	Tag      string
	Untagged *bool
	Cooked   *bool
	Success  *bool
	Page     int
	PageSize int
}

type RecordInput struct {
	CookedAt        string   `json:"cooked_at"`
	Result          string   `json:"result"`
	Rating          *float64 `json:"rating"`
	Notes           string   `json:"notes"`
	Changes         string   `json:"changes"`
	FailureReason   string   `json:"failure_reason"`
	NextImprovement string   `json:"next_improvement"`
	PhotoURLs       []string `json:"photo_urls"`
	AttachmentIDs   []uint64 `json:"attachment_ids"`
	UpdateDishStatus string  `json:"update_dish_status"`
}

type RecookPlanInput struct {
	DishID      uint64 `json:"dish_id"`
	PlannedDate string `json:"planned_date"`
	Reason      string `json:"reason"`
}

type RecookPlanUpdateInput struct {
	PlannedDate string `json:"planned_date"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
}

type SourcePreviewInput struct {
	URL string `json:"url"`
}

type SourcePreview struct {
	URL             string `json:"url"`
	SourcePlatform  string `json:"source_platform"`
	Name            string `json:"name"`
	CoverImageURL   string `json:"cover_image_url"`
	MainIngredients string `json:"main_ingredients"`
	CookTimeMinutes *int   `json:"cook_time_minutes"`
	Partial         bool   `json:"partial"`
}

type HomeSummary struct {
	Stats             HomeStats    `json:"stats"`
	Recent            []Dish       `json:"recent"`
	OverdueHighRating []Dish       `json:"overdue_high_rating"`
	RandomOld         *Dish        `json:"random_old"`
	WantToCook        []Dish       `json:"want_to_cook"`
}

type HomeStats struct {
	TotalDishes    int64 `json:"total_dishes"`
	CookedCount    int64 `json:"cooked_count"`
	SuccessCount   int64 `json:"success_count"`
	MonthCookCount int64 `json:"month_cook_count"`
}

type RecordSaveResult struct {
	Record        CookRecord `json:"record"`
	SuggestStatus string     `json:"suggest_status,omitempty"`
}
