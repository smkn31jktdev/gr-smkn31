package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/router"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/storage/turso"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorWhite  = "\033[97m"
	colorGray   = "\033[90m"
	bold        = "\033[1m"
)

func logInfo(icon, label, value string) {
	fmt.Printf("  %s%-14s%s %s%s%s\n", colorGray, label, colorReset, colorWhite, value, colorReset)
	_ = icon
}

func printBanner(cfg *config.Config) {
	fmt.Println()
	fmt.Printf("%s%s╔══════════════════════════════════════════════╗%s\n", bold, colorCyan, colorReset)
	fmt.Printf("%s%s║           GR31 BACKEND  v1.0.0               ║%s\n", bold, colorCyan, colorReset)
	fmt.Printf("%s%s╚══════════════════════════════════════════════╝%s\n", bold, colorCyan, colorReset)
	fmt.Println()
	fmt.Printf("  %s●%s Server\n", colorGreen, colorReset)
	logInfo("", "Port        :", cfg.ServerPort)
	logInfo("", "Environment :", cfg.AppEnv)
	logInfo("", "Mode        :", cfg.AppMode)
	fmt.Println()
	fmt.Printf("  %s●%s Database\n", colorGreen, colorReset)
	logInfo("", "AstraDB     :", maskSecret(cfg.AstraEndpoint))
	logInfo("", "Keyspace    :", cfg.AstraKeyspace)
	logInfo("", "MongoDB     :", maskSecret(cfg.MongoURI))
	logInfo("", "Database    :", cfg.MongoDatabase)
	logInfo("", "Turso       :", maskSecret(cfg.TursoURL))
	logInfo("", "Supabase    :", maskSecret(cfg.SupabaseDBURL))
	fmt.Println()
	fmt.Printf("  %s●%s Cache\n", colorGreen, colorReset)
	logInfo("", "Redis       :", cfg.RedisAddr)
	fmt.Println()
	fmt.Printf("  %s●%s Upload\n", colorGreen, colorReset)
	logInfo("", "Directory   :", cfg.UploadDir)
	logInfo("", "Max Size    :", fmt.Sprintf("%d MB", cfg.UploadMaxSizeMB))
	fmt.Println()
}

func printStarted(port string) {
	fmt.Printf("%s%s┌────────────────────────────────────────────────┐%s\n", bold, colorGreen, colorReset)
	fmt.Printf("%s%s│  ✓  Server running                             │%s\n", bold, colorGreen, colorReset)
	fmt.Printf("%s%s│     http://localhost:%-25s │%s\n", bold, colorGreen, port, colorReset)
	fmt.Printf("%s%s└────────────────────────────────────────────────┘%s\n", bold, colorGreen, colorReset)
	fmt.Println()
}

func printStatus(label, status string, ok bool) {
	icon := colorGreen + "✓" + colorReset
	statusColor := colorGreen
	if !ok {
		icon = colorRed + "✗" + colorReset
		statusColor = colorRed
	}
	fmt.Printf("  %s  %-14s %s%s%s\n", icon, label, statusColor, status, colorReset)
}

func checkServices(cfg *config.Config, rdb *redis.Client, mongoClient *mongo.Client, tursoClient *turso.Client, supabaseClient *supabase.Client) {
	fmt.Printf("  %s●%s Checking services...\n", colorYellow, colorReset)

	// Cek Redis
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		printStatus("Redis:", "unreachable — idempotency & cache disabled", false)
	} else {
		printStatus("Redis:", "connected", true)
	}

	// AstraDB
	if cfg.AstraEndpoint != "" && cfg.AstraToken != "" {
		printStatus("AstraDB:", "configured", true)
	} else {
		printStatus("AstraDB:", "not configured!", false)
	}

	// MongoDB
	if mongoClient != nil {
		ctxM, cancelM := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancelM()
		if err := mongoClient.Ping(ctxM, nil); err != nil {
			printStatus("MongoDB:", "unreachable", false)
		} else {
			printStatus("MongoDB:", "connected", true)
		}
	} else {
		printStatus("MongoDB:", "not configured", false)
	}

	// Turso
	if tursoClient != nil {
		printStatus("Turso:", "connected", true)
	} else {
		printStatus("Turso:", "not configured!", false)
	}

	// Supabase
	if supabaseClient != nil {
		printStatus("Supabase:", "connected", true)
	} else {
		printStatus("Supabase:", "not configured (placeholder password)", false)
	}

	fmt.Println()
}

// customLogger
func customLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		statusColor := colorGreen
		switch {
		case param.StatusCode >= 500:
			statusColor = colorRed
		case param.StatusCode >= 400:
			statusColor = colorYellow
		case param.StatusCode >= 300:
			statusColor = colorCyan
		}

		return fmt.Sprintf("%s%s  %s  %s%-4s%s  %s%-22s%s  %s%d%s  %s%s\n",
			colorGray,
			param.TimeStamp.Format("15:04:05"),
			colorReset,
			colorCyan,
			param.Method,
			colorReset,
			colorWhite,
			param.Path,
			colorReset,
			statusColor,
			param.StatusCode,
			colorReset,
			colorGray,
			param.Latency,
		)
	})
}

func maskSecret(s string) string {
	if len(s) <= 20 {
		return s
	}
	return s[:20] + "..."
}

func main() {
	// Load konfigurasi
	cfg := config.Load()

	// Print banner startup
	printBanner(cfg)

	// Init AstraDB client
	astraClient := astra.NewClient(cfg.AstraEndpoint, cfg.AstraToken, cfg.AstraKeyspace)

	// Init Turso client
	var tursoClient *turso.Client
	if cfg.TursoURL != "" {
		var err error
		tursoClient, err = turso.NewClient(cfg.TursoURL, cfg.TursoToken)
		if err != nil {
			fmt.Printf("  %s✗%s  Failed to connect to Turso: %v\n", colorRed, colorReset, err)
			os.Exit(1)
		}
	}

	// Init Supabase client
	var supabaseClient *supabase.Client
	if cfg.SupabaseDBURL != "" {
		var err error
		supabaseClient, err = supabase.NewClient(cfg.SupabaseDBURL)
		if err != nil {
			fmt.Printf("  %s⚠%s  Supabase unavailable (skipped): %v\n", colorYellow, colorReset, err)
		}
	}
	defer func() {
		if supabaseClient != nil {
			_ = supabaseClient.Close()
		}
	}()

	// Init MongoDB client
	ctxMongo, cancelMongo := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		fmt.Printf("  %s✗%s  Failed to connect to MongoDB: %v\n", colorRed, colorReset, err)
		cancelMongo()
		os.Exit(1)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		_ = mongoClient.Disconnect(disconnectCtx)
	}()

	// Ping to verify
	if err := mongoClient.Ping(ctxMongo, nil); err != nil {
		fmt.Printf("  %s✗%s  Failed to ping MongoDB: %v\n", colorRed, colorReset, err)
		cancelMongo()
		os.Exit(1)
	}
	cancelMongo()
	mongoDB := mongoClient.Database(cfg.MongoDatabase)

	// Init Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Cek koneksi services
	checkServices(cfg, rdb, mongoClient, tursoClient, supabaseClient)

	// Set Gin mode — matikan output default Gin
	gin.SetMode(gin.ReleaseMode)

	// Init Gin engine dengan logger custom
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(customLogger())

	// Buat direktori upload
	if err := os.MkdirAll(cfg.UploadDir+"/izin", 0755); err != nil {
		fmt.Printf("  %s✗%s  Upload dir: gagal membuat %s\n", colorRed, colorReset, cfg.UploadDir)
	}
	r.Static("/uploads", cfg.UploadDir)

	// Register semua routes
	router.Register(r, cfg, astraClient, mongoDB, rdb, tursoClient, supabaseClient)

	// Tampilkan info server siap
	printStarted(cfg.ServerPort)

	// Print header log request
	fmt.Printf("  %s%sTime      Method  Path                    Status  Latency%s\n", colorGray, bold, colorReset)
	fmt.Printf("  %s%s──────────────────────────────────────────────────────────%s\n", colorGray, bold, colorReset)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		fmt.Printf("\n  %s✗  FATAL: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
}
