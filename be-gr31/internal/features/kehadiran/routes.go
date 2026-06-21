package kehadiran

import (
	"net/http"
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/features/kehadiran/handler/admin"
	"be-gr31/internal/features/kehadiran/handler/student"
	"be-gr31/internal/middleware"
	"be-gr31/internal/model/common"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/storage/turso"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan semua route kehadiran
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, tursoClient *turso.Client, rdb *redis.Client, supabaseClient *supabase.Client) {
	kehadiranStore := turso.NewKehadiranStore(tursoClient)
	rekapStore := astra.NewRekapStore(client)
	studentStore := astra.NewStudentStore(client)
	adminStore := astra.NewAdminStore(client)
	repo := NewRepo(kehadiranStore)
	rekapSvc := NewRekapService(rekapStore, rdb, supabaseClient)
	rekapSvc.SetKehadiranStore(kehadiranStore)
	rekapSvc.SetStudentStore(studentStore)
	service := NewService(repo, rekapSvc, cfg, rdb, supabaseClient, adminStore)
	
	studentHandler := student.NewHandler(service, rekapSvc)
	adminHandler := admin.NewHandler(service, rekapSvc)
	detailHandler := admin.NewDetailHandler(service, cfg.SekolahRadiusMeter)

	idempotencyTTL := time.Duration(cfg.RedisIdempotencyTTLSecs) * time.Second

	// Student routes
	studentGroup := r.Group("/v1/student", middleware.StudentGuard(cfg))
	{
		studentGroup.GET("/kehadiran/rekap", studentHandler.RekapBulananSiswa)
		studentGroup.GET("/kehadiran/bulanan", studentHandler.KehadiranBulananSiswa)
		studentGroup.POST("/kehadiran",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			studentHandler.AbsenSiswa,
		)
		studentGroup.GET("/kehadiran", studentHandler.ListKehadiranSiswa)
		studentGroup.POST("/kehadiran/upload-izin",
			uploadIzinStudentHandler(cfg),
		)
	}

	// Admin routes
	adminGroup := r.Group("/v1/admin", middleware.AdminGuard(cfg))
	{
		adminGroup.POST("/kehadiran",
			middleware.IdempotencyGuard(rdb, idempotencyTTL),
			adminHandler.AbsenAdmin,
		)
		adminGroup.GET("/kehadiran", adminHandler.ListKehadiranAdmin)
		adminGroup.GET("/kehadiran/:id", detailHandler.GetDetail)
		adminGroup.PUT("/kehadiran", adminHandler.UpdateKehadiran)
		adminGroup.DELETE("/kehadiran", adminHandler.DeleteKehadiran)
		adminGroup.POST("/kehadiran/upload-izin",
			uploadIzinAdminHandler(cfg),
		)
		// Rekap endpoints
		adminGroup.GET("/rekap-bulanan", adminHandler.ListRekapBulanan)
		adminGroup.GET("/rekap-harian", adminHandler.GetRekapHarian)
		adminGroup.GET("/rekap-kelas", adminHandler.GetRingkasanKelas)
		adminGroup.GET("/rekap-lengkap", adminHandler.GetRekapLengkap)
		adminGroup.GET("/rekap-persentase-kelas", adminHandler.GetPersentaseKelas)
		adminGroup.GET("/rekap-semester", adminHandler.GetRekapSemesterKelas)
		adminGroup.GET("/rekap-mingguan", adminHandler.GetRekapMingguan)
		adminGroup.GET("/rekap-siswa-detail", adminHandler.GetRekapSiswaDetail)
		adminGroup.GET("/kalender", adminHandler.GetKalender)
		adminGroup.GET("/kelas", adminHandler.GetKelas)
		adminGroup.GET("/kelas-jurusan", adminHandler.GetKelasJurusan)
	}
}

func uploadIzinStudentHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validasi hari sekolah
		now := util.NowJakarta()
		if util.IsWeekend(now) {
			c.JSON(http.StatusBadRequest, common.Fail("tidak ada sekolah pada hari Sabtu dan Minggu"))
			return
		}

		// Validasi jam absensi untuk siswa
		if err := util.ValidateAbsenTime(now, cfg.AbsensiStartHour, cfg.AbsensiStartMinute, cfg.AbsensiEndHour, cfg.AbsensiEndMinute); err != nil {
			c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
			return
		}

		claims := c.MustGet("claims").(*util.JWTClaims)
		HandleUploadIzin(cfg, claims.NIS)(c)
	}
}

func uploadIzinAdminHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		nis := c.Query("nis")
		if nis == "" {
			nis = c.Query("nisn")
		}
		if nis == "" {
			nis = "unknown"
		}
		HandleUploadIzin(cfg, nis)(c)
	}
}
