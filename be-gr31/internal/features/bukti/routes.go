package bukti

import (
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan semua route bukti
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, rdb *redis.Client, supabaseClient *supabase.Client) {
	store := astra.NewBuktiStore(client)
	studentStore := astra.NewStudentStore(client)
	adminStore := astra.NewAdminStore(client)
	repo := NewRepo(store, studentStore, adminStore, supabaseClient)
	service := NewService(repo)
	handler := NewHandler(service)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	studentGroup := r.Group("/v1/student", middleware.StudentGuard(cfg))
	{
		studentGroup.POST("/bukti",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			handler.CreateBukti,
		)
		studentGroup.GET("/bukti", handler.ListBuktiSiswa)
	}

	adminGroup := r.Group("/v1/admin", middleware.AdminGuard(cfg))
	{
		adminGroup.GET("/bukti", handler.ListBuktiAdmin)
	}
}
