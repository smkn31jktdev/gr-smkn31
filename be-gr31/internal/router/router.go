package router

import (
	"be-gr31/internal/config"
	"be-gr31/internal/features/admin"
	"be-gr31/internal/features/aduan"
	"be-gr31/internal/features/auth"
	"be-gr31/internal/features/bukti"
	"be-gr31/internal/features/g7"
	featureinternal "be-gr31/internal/features/internalapi"
	"be-gr31/internal/kalender/puasa"
	"be-gr31/internal/features/kegiatan"
	"be-gr31/internal/features/kehadiran"
	"be-gr31/internal/features/lomba"
	"be-gr31/internal/features/naikkelas"
	"be-gr31/internal/features/ping"
	"be-gr31/internal/features/student"
	"be-gr31/internal/middleware"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/storage/turso"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Register mendaftarkan semua route ke Gin engine
func Register(r *gin.Engine, cfg *config.Config, client *astra.Client, mongoDB *mongo.Database, rdb *redis.Client, tursoClient *turso.Client, supabaseClient *supabase.Client) {
	// CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Fase 2: Auth + Ping + Internal
	ping.RegisterRoutes(r, rdb)
	auth.RegisterRoutes(r, cfg, client, supabaseClient)
	featureinternal.RegisterRoutes(r, cfg, client)

	// Fase 3: Kehadiran
	kehadiran.RegisterRoutes(r, cfg, client, tursoClient, rdb, supabaseClient)

	// Fase 4: G7
	g7.RegisterRoutes(r, cfg, client, rdb, supabaseClient)
	puasa.RegisterRoutes(r)

	// Fase 5: Sekolah
	kegiatan.RegisterRoutes(r, cfg, client, rdb)
	bukti.RegisterRoutes(r, cfg, client, rdb, supabaseClient)
	aduan.RegisterRoutes(r, cfg, client, mongoDB, rdb, supabaseClient)
	lomba.RegisterRoutes(r, cfg, client, rdb, supabaseClient)

	// Fase 6: Admin Management
	admin.RegisterRoutes(r, cfg, client, supabaseClient, rdb)
	student.RegisterRoutes(r, cfg, client, supabaseClient, rdb)

	// Fase 7: Naik Kelas otomatis 1 Juli
	naikkelas.RegisterRoutes(r, cfg, client, supabaseClient)
}

