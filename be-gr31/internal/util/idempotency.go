package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const idempotencyPrefix = "idempotency:"

// IdempotencyStore
type IdempotencyStore struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewIdempotencyStore
func NewIdempotencyStore(rdb *redis.Client, ttl time.Duration) *IdempotencyStore {
	return &IdempotencyStore{rdb: rdb, ttl: ttl}
}

type cachedResponse struct {
	StatusCode int             `json:"statusCode"`
	Body       json.RawMessage `json:"body"`
}

// Middleware
func (s *IdempotencyStore) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Idempotency-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "X-Idempotency-Key header diperlukan untuk request ini",
			})
			return
		}

		redisKey := idempotencyPrefix + key
		ctx := c.Request.Context()

		raw, err := s.rdb.Get(ctx, redisKey).Bytes()
		if err == nil {
			var cached cachedResponse
			if jsonErr := json.Unmarshal(raw, &cached); jsonErr == nil {
				c.Header("X-Idempotency-Replayed", "true")
				c.Data(cached.StatusCode, "application/json", cached.Body)
				c.Abort()
				return
			}
		}

		// Belum ada: wrap writer untuk capture response
		brw := &bodyResponseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = brw

		c.Next()

		// Simpan response jika sukses (2xx)
		if brw.statusCode >= 200 && brw.statusCode < 300 && brw.body.Len() > 0 {
			cached := cachedResponse{
				StatusCode: brw.statusCode,
				Body:       brw.body.Bytes(),
			}
			data, _ := json.Marshal(cached)
			_ = s.rdb.Set(ctx, redisKey, data, s.ttl).Err()
		}
	}
}

// bodyResponseWriter
type bodyResponseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *bodyResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
