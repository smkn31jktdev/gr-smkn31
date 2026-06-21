package aduan

import (
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/mongodb"
	"be-gr31/internal/storage/supabase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// RegisterRoutes
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, mongoDB *mongo.Database, rdb *redis.Client, supabaseClient *supabase.Client) {
	store := mongodb.NewAduanStore(mongoDB)
	studentStore := astra.NewStudentStore(client)
	adminStore := astra.NewAdminStore(client)
	repo := NewRepo(store)
	service := NewService(repo, studentStore, adminStore, supabaseClient)
	handler := NewHandler(service)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	studentGroup := r.Group("/v1/student", middleware.StudentGuard(cfg))
	{
		studentGroup.POST("/aduan",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			handler.CreateAduan,
		)
		studentGroup.GET("/aduan", handler.ListAduanSiswa)
	}

	adminGroup := r.Group("/v1/admin", middleware.AdminGuard(cfg))
	{
		adminGroup.GET("/aduan", handler.ListAduanAdmin)
		adminGroup.GET("/aduan/room", handler.GetRoomAdmin)
		adminGroup.POST("/aduan/status", handler.UpdateStatusAdmin)
		adminGroup.POST("/aduan/respond",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			handler.RespondAdmin,
		)
		adminGroup.GET("/aduan/export/csv", handler.ExportCSVAduan)
		adminGroup.GET("/aduan/export/html", handler.ExportHTMLAduan)
	}
}
