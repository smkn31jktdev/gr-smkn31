package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: blob:; "+
				"font-src 'self'; "+
				"connect-src 'self'; "+
				"media-src 'self'; "+
				"object-src 'none'; "+
				"frame-src 'none'; "+
				"base-uri 'self'; "+
				"form-action 'self'; "+
				"frame-ancestors 'none'; "+
				"upgrade-insecure-requests",
		)

		c.Header("Server", "")
		c.Header("X-Powered-By", "")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
		c.Header("Pragma", "no-cache")

		c.Next()
	}
}

// AntiJudolHeaders
func AntiJudolHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		existing := c.Writer.Header().Get("Content-Security-Policy")
		if existing == "" {
			c.Header("Content-Security-Policy", "frame-ancestors 'none'")
		}

		c.Next()
	}
}

// ScannerTrap
func ScannerTrap() gin.HandlerFunc {
	trapPaths := map[string]bool{
		"/admin":          true,
		"/wp-admin":       true,
		"/wp-login.php":   true,
		"/phpmyadmin":     true,
		"/.env":           true,
		"/config":         true,
		"/backup":         true,
		"/shell":          true,
		"/.git":           true,
		"/actuator":       true,
		"/api/v1/console": true,
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if trapPaths[path] {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
