package auth

import (
	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client, supabaseClient *supabase.Client) {
	studentStore := astra.NewStudentStore(client)
	adminStore := astra.NewAdminStore(client)
	repo := NewRepo(client, studentStore, adminStore, supabaseClient)

	jwtConfig := util.JWTConfig{
		Secret:         cfg.JWTSecret,
		Issuer:         cfg.JWTIssuer,
		AccessTTLHours: cfg.JWTAccessTTLHours,
		RefreshTTLDays: cfg.JWTRefreshTTLDays,
	}
	service := NewService(repo, jwtConfig, cfg.SuperAdminEmails)
	handler := NewHandler(service)

	studentGroup := r.Group("/v1/student")
	{
		studentGroup.POST("/login", handler.LoginSiswa)
		studentGroup.POST("/refresh-token", handler.RefreshTokenSiswa)
		studentGroup.GET("/me", middleware.StudentGuard(cfg), handler.MeSiswa)
		studentGroup.POST("/me", middleware.StudentGuard(cfg), handler.MeSiswa) // backward compat
	}

	adminGroup := r.Group("/v1/admin")
	{
		adminGroup.POST("/login", handler.LoginAdmin)
		adminGroup.POST("/refresh-token", handler.RefreshTokenAdmin)
		adminGroup.GET("/me", middleware.AdminGuard(cfg), handler.MeAdmin)
		adminGroup.POST("/me", middleware.AdminGuard(cfg), handler.MeAdmin) // backward compat
	}
}
