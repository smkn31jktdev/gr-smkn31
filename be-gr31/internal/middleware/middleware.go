package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/model/common"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Memvalidasi JWT dan memastikan role adalah "student"
func StudentGuard(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Fail("unauthorized"))
			return
		}
		claims, err := util.ParseJWT(token, cfg.JWTSecret)
		if err != nil || claims.Role != "student" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Fail("unauthorized"))
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// Memvalidasi JWT dan memastikan role mengandung prefix "admin" atau "super_admin"
func AdminGuard(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Fail("unauthorized"))
			return
		}
		claims, err := util.ParseJWT(token, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Fail("unauthorized"))
			return
		}
		if !isAdminRole(claims.Role) {
			log.Printf("AdminGuard DENY: role=%q email=%q path=%q", claims.Role, claims.Email, c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusForbidden, common.Fail("forbidden"))
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// Memvalidasi JWT dan memastikan role adalah "super_admin"
func SuperAdminGuard(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Fail("unauthorized"))
			return
		}
		claims, err := util.ParseJWT(token, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Fail("unauthorized"))
			return
		}
		if claims.Role != "super_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, common.Fail("forbidden: super_admin only"))
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// Memvalidasi header X-Internal-Key
func InternalGuard(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Internal-Key")
		if key != cfg.InternalKey {
			c.AbortWithStatusJSON(http.StatusForbidden, common.Fail("forbidden"))
			return
		}
		c.Next()
	}
}

// Middleware idempotency berbasis Redis
func IdempotencyGuard(rdb *redis.Client, ttl time.Duration) gin.HandlerFunc {
	store := util.NewIdempotencyStore(rdb, ttl)
	return store.Middleware()
}

// Mengizinkan semua origin (sesuaikan untuk production)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Idempotency-Key, X-Internal-Key")
		c.Header("Access-Control-Expose-Headers", "X-Idempotency-Replayed")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Mengambil token dari header Authorization: Bearer <token>
func extractBearer(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// Memeriksa apakah role termasuk kategori admin.
// Sesuai ABSEN.md §7B: ["admin", "super_admin", "piket", "guru_wali", "bk"].
// Disertakan juga "guru_bk" sebagai alias historis untuk kompatibilitas data lama.
func isAdminRole(role string) bool {
	adminRoles := []string{"admin", "super_admin", "guru_bk", "guru_wali", "piket", "bk", "admin_bk", "walas", "admin_piket"}
	for _, r := range adminRoles {
		if role == r {
			return true
		}
	}
	return false
}
