package security

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Rate Limiter
type RateLimiterConfig struct {
	Window time.Duration
	Max    int
	Prefix string
}

// RateLimiter middleware berbasis Redis sliding-window counter per IP
func RateLimiter(rdb *redis.Client, cfg RateLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := realIP(c)
		key := fmt.Sprintf("rl:%s:%s", cfg.Prefix, ip)
		ctx := context.Background()

		pipe := rdb.Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, cfg.Window)
		_, err := pipe.Exec(ctx)

		if err != nil {
			c.Next()
			return
		}

		count := incr.Val()
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Max))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max64(0, int64(cfg.Max)-count)))

		if count > int64(cfg.Max) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "terlalu banyak permintaan, coba lagi nanti",
			})
			return
		}
		c.Next()
	}
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// Mengambil IP asli klien dengan mempertimbangkan proxy header
func realIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	return c.ClientIP()
}

// Request Shield

const requestSizeLimit = 1 << 20

// Pola injeksi berbahaya
var injectionPatterns = []string{
	"' or ", "\" or ", "' and ", "\" and ",
	"union select", "union all select",
	"drop table", "drop database",
	"insert into", "delete from", "update set",
	"; exec", "; execute", "xp_cmdshell",
	"information_schema", "sys.tables",
	"$where", "$ne", "$gt", "$regex", "$or", "$and",
	"../", "..\\", "%2e%2e", "%252e%252e",
	"/etc/passwd", "/etc/shadow", "c:\\windows", "c:/windows",
	"win.ini", "boot.ini",
	"; ls", "; cat ", "; wget ", "; curl ", "| ls", "| cat ",
	"`ls`", "`id`", "$(id)", "$(ls)",
	"<script", "</script>", "javascript:", "vbscript:",
	"onload=", "onerror=", "onclick=", "onmouseover=",
	"<iframe", "<object", "<embed", "<svg",
	"alert(", "eval(", "document.cookie", "window.location",
	"{{", "}}", "${", "#{",
	"169.254.169.254",
	"metadata.google",
	"localhost", "127.0.0.1", "::1", "0.0.0.0",
}

// Request Shield
func RequestShield() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, requestSizeLimit)

		path := strings.ToLower(c.Request.URL.RawPath + c.Request.URL.Path)
		if containsInjection(path) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "request ditolak",
			})
			return
		}

		// Periksa query string
		query := strings.ToLower(c.Request.URL.RawQuery)
		if containsInjection(query) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "request ditolak",
			})
			return
		}

		// Periksa User-Agent suspicious
		ua := c.GetHeader("User-Agent")
		if isSuspiciousUA(ua) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "request ditolak",
			})
			return
		}

		// Validasi UTF-8
		if !utf8.ValidString(c.Request.URL.Path) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "request tidak valid",
			})
			return
		}

		c.Next()
	}
}

// Memeriksa string terhadap pola injeksi
func containsInjection(s string) bool {
	lower := strings.ToLower(s)
	for _, p := range injectionPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// suspiciousUAs
var suspiciousUAs = []string{
	"sqlmap", "nikto", "nmap", "masscan", "zgrab",
	"dirbuster", "gobuster", "wfuzz", "ffuf", "burpsuite",
	"nuclei", "metasploit", "hydra", "acunetix",
	"openvas", "nessus", "w3af", "skipfish",
	"python-requests", "go-http-client", "libwww-perl",
	"curl/", "wget/",
}

// isSuspiciousUA
func isSuspiciousUA(ua string) bool {
	if ua == "" {
		return true
	}
	lower := strings.ToLower(ua)
	for _, bad := range suspiciousUAs {
		if strings.Contains(lower, bad) {
			return true
		}
	}
	return false
}

// Brute Force Guard

// BruteForceConfig
type BruteForceConfig struct {
	MaxAttempts  int
	LockDuration time.Duration
	Prefix       string
}

// BruteForceGuard
func BruteForceGuard(rdb *redis.Client, cfg BruteForceConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := realIP(c)
		lockKey := fmt.Sprintf("bf:lock:%s:%s", cfg.Prefix, ip)
		ctx := context.Background()

		locked, err := rdb.Exists(ctx, lockKey).Result()
		if err == nil && locked > 0 {
			ttl, _ := rdb.TTL(ctx, lockKey).Result()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": fmt.Sprintf("terlalu banyak percobaan, coba lagi dalam %.0f detik", ttl.Seconds()),
			})
			return
		}

		c.Next()

		// Jika response adalah 401 atau 403
		status := c.Writer.Status()
		switch status {
		case http.StatusUnauthorized, http.StatusForbidden:
			attemptKey := fmt.Sprintf("bf:attempts:%s:%s", cfg.Prefix, ip)
			count, _ := rdb.Incr(ctx, attemptKey).Result()
			rdb.Expire(ctx, attemptKey, cfg.LockDuration)

			if count >= int64(cfg.MaxAttempts) {
				rdb.Set(ctx, lockKey, 1, cfg.LockDuration)
				rdb.Del(ctx, attemptKey)
			}
		case http.StatusOK, http.StatusCreated:
			rdb.Del(ctx, fmt.Sprintf("bf:attempts:%s:%s", cfg.Prefix, ip))
		}
	}
}

// RecordLoginFailure
func RecordLoginFailure(rdb *redis.Client, ip, prefix string, maxAttempts int, lockDuration time.Duration) {
	ctx := context.Background()
	attemptKey := fmt.Sprintf("bf:attempts:%s:%s", prefix, ip)
	lockKey := fmt.Sprintf("bf:lock:%s:%s", prefix, ip)

	count, _ := rdb.Incr(ctx, attemptKey).Result()
	rdb.Expire(ctx, attemptKey, lockDuration)

	if count >= int64(maxAttempts) {
		rdb.Set(ctx, lockKey, 1, lockDuration)
		rdb.Del(ctx, attemptKey)
	}
}

// ResetLoginAttempts
func ResetLoginAttempts(rdb *redis.Client, ip, prefix string) {
	ctx := context.Background()
	rdb.Del(ctx, fmt.Sprintf("bf:attempts:%s:%s", prefix, ip))
	rdb.Del(ctx, fmt.Sprintf("bf:lock:%s:%s", prefix, ip))
}

// IsIPLocked
func IsIPLocked(rdb *redis.Client, ip, prefix string) bool {
	ctx := context.Background()
	locked, err := rdb.Exists(ctx, fmt.Sprintf("bf:lock:%s:%s", prefix, ip)).Result()
	return err == nil && locked > 0
}
