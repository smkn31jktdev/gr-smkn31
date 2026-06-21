package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Menyimpan semua konfigurasi aplikasi dari environment variables
type Config struct {
	// Server
	ServerPort string
	AppMode    string
	AppEnv     string

	// Internal
	InternalKey string

	// JWT
	JWTSecret         string
	JWTIssuer         string
	JWTAccessTTLHours int
	JWTRefreshTTLDays int
	SuperAdminEmails  []string

	// AstraDB
	AstraEndpoint string
	AstraToken    string
	AstraKeyspace string

	// MongoDB
	MongoURI      string
	MongoDatabase string

	// Redis
	RedisAddr               string
	RedisPassword           string
	RedisDB                 int
	RedisIdempotencyTTLSecs int
	RedisCacheTTLSecs       int

	// Upload
	UploadDir       string
	UploadMaxSizeMB int64

	// Sekolah
	SekolahLat          float64
	SekolahLng          float64
	SekolahRadiusMeter  float64
	AbsensiStartHour    int
	AbsensiStartMinute  int
	AbsensiEndHour      int
	AbsensiEndMinute    int
	MinGPSAccuracyMeter float64
	MaxGPSAccuracyMeter float64

	// Turso
	TursoURL   string
	TursoToken string

	// Supabase
	SupabaseDBURL string
}

// Membaca env dan mengembalikan Config. Panic jika field wajib kosong
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: .env file not found, reading from system env")
	}

	cfg := &Config{
		ServerPort:              getEnv("SERVER_PORT", "8080"),
		AppMode:                 getEnv("APP_MODE", "http"),
		AppEnv:                  getEnv("APP_ENV", "development"),
		InternalKey:             mustEnv("INTERNAL_KEY"),
		JWTSecret:               mustEnv("JWT_SECRET"),
		JWTIssuer:               getEnv("JWT_ISSUER", "gr31-backend"),
		JWTAccessTTLHours:       getEnvInt("JWT_ACCESS_TTL_HOURS", 24),
		JWTRefreshTTLDays:       getEnvInt("JWT_REFRESH_TTL_DAYS", 30),
		SuperAdminEmails:        getEnvStringSlice("SUPER_ADMIN_EMAILS", ","),
		AstraEndpoint:           mustEnv("ASTRA_DB_ENDPOINT"),
		AstraToken:              mustEnv("ASTRA_DB_TOKEN"),
		AstraKeyspace:           getEnv("ASTRA_DB_KEYSPACE", "gr31"),
		MongoURI:                mustEnv("MONGODB_URI"),
		MongoDatabase:           getEnv("MONGODB_DATABASE", "smkn31jkt"),
		RedisAddr:               getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:           getEnv("REDIS_PASSWORD", ""),
		RedisDB:                 getEnvInt("REDIS_DB", 0),
		RedisIdempotencyTTLSecs: getEnvInt("REDIS_IDEMPOTENCY_TTL_SECONDS", 86400),
		RedisCacheTTLSecs:       getEnvInt("REDIS_CACHE_TTL_SECONDS", 300),
		UploadDir:               getEnv("UPLOAD_DIR", "./uploads"),
		UploadMaxSizeMB:         int64(getEnvInt("UPLOAD_MAX_SIZE_MB", 5)),
		SekolahLat:              getEnvFloat("SEKOLAH_LAT", -6.1819399),
		SekolahLng:              getEnvFloat("SEKOLAH_LNG", 106.8518572),
		SekolahRadiusMeter:      getEnvFloat("SEKOLAH_RADIUS_METER", 50),
		AbsensiStartHour:        getEnvInt("ABSENSI_START_HOUR", 5),
		AbsensiStartMinute:      getEnvInt("ABSENSI_START_MINUTE", 30),
		AbsensiEndHour:          getEnvInt("ABSENSI_END_HOUR", 7),
		AbsensiEndMinute:        getEnvInt("ABSENSI_END_MINUTE", 0),
		MinGPSAccuracyMeter:     getEnvFloat("GPS_MIN_ACCURACY_METER", 2.0),
		MaxGPSAccuracyMeter:     getEnvFloat("GPS_MAX_ACCURACY_METER", 80.0),
		TursoURL:                getEnv("TURSO_URL", ""),
		TursoToken:              getEnv("TURSO_TOKEN", ""),
		SupabaseDBURL:           getEnv("SUPABASE_DB_URL", ""),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("FATAL: environment variable %s is required", key)
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvFloat(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}

func getEnvStringSlice(key, sep string) []string {
	v := os.Getenv(key)
	if v == "" {
		return []string{}
	}
	parts := strings.Split(v, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
