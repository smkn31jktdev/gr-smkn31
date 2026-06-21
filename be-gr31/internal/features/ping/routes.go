package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes
func RegisterRoutes(r *gin.Engine, rdb *redis.Client) {
	handler := NewHandler(rdb)

	r.GET("/", handler.Root)
	r.GET("/health", handler.Health)
	r.GET("/v1/ping", handler.Ping)
}
