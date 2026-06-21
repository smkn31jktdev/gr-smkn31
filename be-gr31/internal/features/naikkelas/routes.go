package naikkelas

import (
	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan route naik kelas dan memulai cron scheduler.
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, supabaseClient *supabase.Client) {
	studentStore := astra.NewStudentStore(client)
	svc := NewService(client, studentStore, supabaseClient)
	handler := NewHandler(svc)

	// Mulai cron otomatis di background
	StartCronNaikKelas(svc)

	// Route SuperAdmin only
	superAdminGroup := r.Group("/v1/admin", middleware.SuperAdminGuard(cfg))
	{
		superAdminGroup.POST("/naik-kelas", handler.TriggerNaikKelas)

		superAdminGroup.GET("/naik-kelas/status", handler.StatusNaikKelas)
	}
}
