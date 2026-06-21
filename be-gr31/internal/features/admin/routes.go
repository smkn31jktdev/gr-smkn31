package admin

import (
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan semua route manajemen admin.
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, supabaseClient *supabase.Client, rdb *redis.Client) {
	adminStore := astra.NewAdminStore(client)
	repo := NewRepo(adminStore, supabaseClient)
	service := NewService(repo)
	handler := NewHandler(service)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	// SuperAdmin-only routes
	superAdminGroup := r.Group("/v1/admin", middleware.SuperAdminGuard(cfg))
	{
		superAdminGroup.POST("/admins", handler.CreateAdmin)
		superAdminGroup.GET("/admins", handler.ListAdmin)
		superAdminGroup.PUT("/admins/:id", handler.UpdateAdmin)
		superAdminGroup.DELETE("/admins/:id", handler.DeleteAdmin)
		superAdminGroup.POST("/tambah-admin/sheets", handler.ImportFromSheets)
		superAdminGroup.POST("/tambah-admin/bulk",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			handler.BulkImportAdmin,
		)
	}
}
