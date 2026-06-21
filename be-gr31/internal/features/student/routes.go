package student

import (
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Mendaftarkan semua route manajemen siswa.
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, supabaseClient *supabase.Client, rdb *redis.Client) {
	studentStore := astra.NewStudentStore(client)
	repo := NewRepo(client, studentStore, supabaseClient)
	service := NewService(repo)
	handler := NewHandler(service)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	adminGroup := r.Group("/v1/admin", middleware.AdminGuard(cfg))
	{
		adminGroup.POST("/students", handler.CreateSiswa)
		adminGroup.GET("/students", handler.ListSiswa)
		adminGroup.PUT("/students/:id", handler.UpdateSiswa)
		adminGroup.DELETE("/students/:id", handler.DeleteSiswa)
		adminGroup.POST("/tambah-siswa/bulk",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			handler.BulkImportSiswa,
		)
		adminGroup.POST("/tambah-siswa/sheets", handler.ImportFromSheets)
	}
}
