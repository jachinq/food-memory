package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"food-memory/internal/config"
	"food-memory/internal/handler"
	"food-memory/internal/middleware"
	"food-memory/internal/repository"
	"food-memory/internal/service"
	"food-memory/internal/storage"
	"food-memory/internal/unfurl"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB, fs *storage.Local) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), middleware.CORS())
	r.MaxMultipartMemory = cfg.MaxUploadBytes

	dishRepo := repository.NewDishRepo(db)
	recordRepo := repository.NewRecordRepo(db)
	tagRepo := repository.NewTagRepo(db)
	attRepo := repository.NewAttachmentRepo(db)
	recookRepo := repository.NewRecookRepo(db)
	homeRepo := repository.NewHomeRepo(db)

	dishSvc := service.NewDishService(db, dishRepo, tagRepo, attRepo, recordRepo)
	recordSvc := service.NewRecordService(db, recordRepo, dishRepo, attRepo)
	tagSvc := service.NewTagService(tagRepo)
	uploadSvc := service.NewUploadService(cfg, fs, attRepo)
	homeSvc := service.NewHomeService(homeRepo, attRepo)
	recookSvc := service.NewRecookService(db, recookRepo, dishRepo)

	dishH := handler.NewDishHandler(dishSvc)
	recordH := handler.NewRecordHandler(recordSvc)
	tagH := handler.NewTagHandler(tagSvc)
	uploadH := handler.NewUploadHandler(uploadSvc)
	homeH := handler.NewHomeHandler(homeSvc)
	recookH := handler.NewRecookHandler(recookSvc)
	sourceH := handler.NewSourceHandler(unfurl.New())

	api := r.Group("/api")
	api.Use(middleware.AccessToken(cfg.AccessToken))
	{
		api.GET("/health", handler.Health)
		api.GET("/dishes", dishH.List)
		api.POST("/dishes", dishH.Create)
		api.GET("/dishes/:id", dishH.Get)
		api.PUT("/dishes/:id", dishH.Update)
		api.DELETE("/dishes/:id", dishH.Delete)
		api.GET("/dishes/:id/records", recordH.List)
		api.POST("/dishes/:id/records", recordH.Create)
		api.PUT("/records/:id", recordH.Update)
		api.DELETE("/records/:id", recordH.Delete)
		api.GET("/tags", tagH.List)
		api.POST("/tags", tagH.Create)
		api.POST("/upload", uploadH.Upload)
		api.GET("/home/summary", homeH.Summary)
		api.GET("/recook-plans", recookH.List)
		api.POST("/recook-plans", recookH.Create)
		api.PUT("/recook-plans/:id", recookH.Update)
		api.POST("/recook-plans/:id/complete", recookH.Complete)
		api.POST("/recook-plans/:id/cancel", recookH.Cancel)
		api.DELETE("/recook-plans/:id", recookH.Delete)
		api.POST("/source-preview", sourceH.Preview)
	}

	r.Static("/uploads", cfg.UploadDir)
	mountSPA(r, cfg.PublicDir)
	return r
}

func mountSPA(r *gin.Engine, publicDir string) {
	index := filepath.Join(publicDir, "index.html")
	if _, err := os.Stat(index); err != nil {
		return
	}
	r.Static("/assets", filepath.Join(publicDir, "assets"))
	r.StaticFile("/favicon.svg", filepath.Join(publicDir, "favicon.svg"))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, handler.Body{Code: 404, Message: "接口不存在", Data: nil})
			return
		}
		if strings.HasPrefix(c.Request.URL.Path, "/uploads") {
			c.Status(http.StatusNotFound)
			return
		}
		c.File(index)
	})
}
