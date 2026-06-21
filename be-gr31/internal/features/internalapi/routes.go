package internalapi

import (
	"be-gr31/internal/config"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan route internal.
func RegisterRoutes(r *gin.Engine, cfg *config.Config, client *astra.Client) {
	handler := NewHandler(client)

	internalGroup := r.Group("/internal", middleware.InternalGuard(cfg))
	{
		internalGroup.POST("/client-upsert", handler.UpsertClient)
		internalGroup.POST("/client-find", handler.FindClient)
		internalGroup.DELETE("/client-delete", handler.DeleteClient)
	}
}
