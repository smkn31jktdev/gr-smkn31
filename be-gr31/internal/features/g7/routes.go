package g7

import (
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/features/g7/handler/admin"
	"be-gr31/internal/features/g7/handler/student"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes G7
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, rdb *redis.Client, supabaseClient *supabase.Client) {
	store := astra.NewG7Store(client)
	rekapStore := astra.NewG7RekapStore(client)
	studentStore := astra.NewStudentStore(client)
	adminStore := astra.NewAdminStore(client)
	repo := NewRepo(store, rekapStore, studentStore, adminStore, supabaseClient)
	service := NewService(repo, rdb)
	
	studentHandler := student.NewHandler(service)
	adminHandler := admin.NewHandler(service)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	studentGroup := r.Group("/v1/student", middleware.StudentGuard(cfg))
	{
		studentGroup.GET("/g7/dashboard", studentHandler.DashboardSiswa)
		studentGroup.POST("/g7",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			studentHandler.UpsertG7Siswa,
		)
		studentGroup.GET("/g7", studentHandler.ListG7Siswa)
		studentGroup.GET("/g7/:tanggal", studentHandler.GetG7ByTanggal)
	}

	adminGroup := r.Group("/v1/admin", middleware.AdminGuard(cfg))
	{
		adminGroup.GET("/g7", adminHandler.ListG7Admin)
		adminGroup.GET("/g7/summary", adminHandler.SummaryG7Admin)
		adminGroup.DELETE("/g7", adminHandler.DeleteG7Admin)

		adminGroup.GET("/g7/rekap", adminHandler.ListRekap)
		adminGroup.GET("/g7/statistik", adminHandler.StatistikRekap)
		adminGroup.GET("/g7/rekap-kelas", adminHandler.RekapKelas)
		adminGroup.GET("/g7/rekap-semester", adminHandler.GetRekapSemesterKelas)
		adminGroup.GET("/g7/rekap/:nis/:bulan", adminHandler.GetRekapDetail)
		adminGroup.GET("/g7/rekap/:nis/:bulan/pdf", adminHandler.ExportPDF)
		adminGroup.GET("/g7/evaluate/:nis/:bulan", adminHandler.EvaluateJurnal)
		adminGroup.POST("/g7/rekap",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			adminHandler.UpsertRekap,
		)
		adminGroup.DELETE("/g7/rekap", adminHandler.DeleteRekap)
		adminGroup.GET("/g7/suggest/:nis/:bulan", adminHandler.SuggestSkor)
	}
}
