package kegiatan

import (
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan semua route kegiatan.
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, rdb *redis.Client) {
	store := astra.NewKegiatanStore(client)
	repo := NewRepo(store)
	service := NewService(repo)
	handler := NewHandler(service)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	// Student routes
	studentGroup := r.Group("/v1/student", middleware.StudentGuard(cfg))
	{
		studentGroup.POST("/kegiatan",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			handler.CreateKegiatan,
		)
		studentGroup.PUT("/kegiatan", handler.UpdateKegiatan)
		studentGroup.GET("/kegiatan", handler.ListKegiatanSiswa)
		studentGroup.DELETE("/kegiatan", handler.DeleteKegiatanSiswa)
	}

	// Admin routes
	adminGroup := r.Group("/v1/admin", middleware.AdminGuard(cfg))
	{
		adminGroup.GET("/kegiatan", handler.ListKegiatanAdmin)
	}
}
