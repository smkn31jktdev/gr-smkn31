package lomba

import (
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan semua route lomba kebersihan kelas.
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, rdb *redis.Client, supabaseClient *supabase.Client) {
	store := astra.NewKebersihanStore(client)
	studentStore := astra.NewStudentStore(client)
	adminStore := astra.NewAdminStore(client)
	repo := NewRepo(store, studentStore, adminStore, supabaseClient)
	service := NewService(repo)
	handler := NewHandler(service)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	// Student routes
	studentGroup := r.Group("/v1/student", middleware.StudentGuard(cfg))
	{
		studentGroup.POST("/lomba",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			handler.CreateLomba,
		)
		studentGroup.PUT("/lomba", handler.UpdateLomba)
		studentGroup.GET("/lomba", handler.ListLombaSiswa)
		studentGroup.DELETE("/lomba", handler.DeleteLombaSiswa)
		studentGroup.POST("/lomba/upload",
			HandleUploadFotoLomba(cfg),
		)
	}

	// Admin routes (Wali Kelas and Super Admin)
	adminGroup := r.Group("/v1/admin", middleware.AdminGuard(cfg))
	{
		adminGroup.GET("/lomba", handler.ListLombaAdmin)
	}
}
