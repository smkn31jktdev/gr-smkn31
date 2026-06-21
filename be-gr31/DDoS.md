# CLAUDE.md — Security Hardening: GR31 School System

# Target: gr31.tech | smkn31jkt.sch.id

# Stack: Go + Gin, AstraDB, Redis, Hostinger VPS

---

## Profil Sistem & Risiko

| Atribut       | Detail                               |
| ------------- | ------------------------------------ |
| Domain aktif  | `gr31.tech`, `smkn31jkt.sch.id`      |
| User harian   | ~700 siswa + guru + staf               |
| Fungsi kritis | Absensi harian, kegiatan siswa         |
| Backend       | Go + Gin (layered-by-feature)          |
| Database      | AstraDB (Cassandra) & MongoDB          |
| Cache / State | Redis (self-hosted, Hostinger VPS)     |
| Anti-judol    | ✅ Sudah terimplementasi              |
| Anti-DDoS     | ❌ Belum — dokumen ini fokus di sini  |

### Pola Traffic Khas Sekolah

```
05:30 - 07:30  │████████████████│ SPIKE — absensi masuk (burst 400+ req/menit)
08:00 - 12:00  │░░░████░░░████░░│ Sedang — aktivitas belajar
12:00 - 13:30  │████████████████│ SPIKE — absen ishoma / kegiatan siang
13:30 - 15:00  │░░░████░░░░░░░░░│ Sedang menurun
15:00 - 17:00  │████████████████│ SPIKE — absensi pulang
17:00 - dst    │░░░░░░░░░░░░░░░░│ Sangat rendah
Sabtu/Minggu   │░░░████░░░████░░│ SPIKE — Mengisi kegiatan G7KAIH
```

> **Implikasi penting:** Rate limiter harus TOLERAN terhadap burst jam 07.00, 12.00, dan 15.00.
> Threshold yang terlalu ketat akan memblokir siswa yang legitimate saat peak hour.

---

## Arsitektur Pertahanan (Defense in Depth)

```
Internet / Penyerang
        │
        ▼
┌───────────────────────────────┐
│  CLOUDFLARE (Edge)            │  Layer 1 — DDoS absorb, WAF, Bot Fight
│  gr31.tech + smkn31jkt.sch.id │
└───────────┬───────────────────┘
            │ Hanya IP Cloudflare
            ▼
┌───────────────────────────────┐
│  UFW FIREWALL (Hostinger VPS) │  Layer 2 — Drop semua non-CF traffic
│  Port 80/443: CF IPs only     │
└───────────┬───────────────────┘
            │
            ▼
┌───────────────────────────────┐
│  NGINX (Reverse Proxy)        │  Layer 3 — Rate limit koneksi, real IP resolve
│  limit_req + limit_conn       │
└───────────┬───────────────────┘
            │
            ▼
┌───────────────────────────────┐
│  GIN MIDDLEWARE STACK         │  Layer 4 — Smart rate limit (school-aware)
│  RealIP → CFOnly → BanCheck   │            IP ban, circuit breaker
│  → SchoolRateLimit → AntiFill │
└───────────┬───────────────────┘
            │
            ▼
┌───────────────────────────────┐
│  ASTRADB + REDIS              │  Layer 5 — Connection pool, query timeout
└───────────────────────────────┘
```

---

## Layer 1: Cloudflare — Konfigurasi Wajib

### DNS Setup (Semua harus Proxied / awan orange)

```
# gr31.tech
A    @           → [IP VPS]     ✅ Proxied
A    www         → [IP VPS]     ✅ Proxied
A    api         → [IP VPS]     ✅ Proxied
A    portal      → [IP VPS]     ✅ Proxied

# smkn31jkt.sch.id (jika dikelola via CF)
A    @           → [IP VPS]     ✅ Proxied
A    www         → [IP VPS]     ✅ Proxied

# ⚠️  JANGAN ada subdomain dengan DNS Only (awan abu-abu)
#     Itu yang dieksploitasi cf-hero untuk menemukan IP asli
```

### Security Settings Dashboard

```
SSL/TLS:
  ├── Mode: Full (Strict)
  ├── Always Use HTTPS: ON
  ├── Minimum TLS Version: TLS 1.2
  └── HSTS: Enabled, max-age 1 tahun

Firewall:
  ├── Bot Fight Mode: ON
  ├── Browser Integrity Check: ON
  └── Security Level: Medium (naikkan ke High saat serangan)

Speed:
  └── Auto Minify: JS + CSS + HTML (ringankan beban server)
```

### WAF Custom Rules (urutan eksekusi penting)

```
# Rule 1 — Blokir akses via IP langsung (bypass hostname)
# Ini menutup celah cf-hero yang sudah temukan IP asli
Expression:
  (not http.host in {"gr31.tech" "www.gr31.tech" "api.gr31.tech"
                     "smkn31jkt.sch.id" "www.smkn31jkt.sch.id"})
Action: Block

# Rule 2 — Challenge bot di endpoint absensi
# Endpoint absensi tidak seharusnya diakses dari non-browser
Expression:
  (http.request.uri.path contains "/absensi" and
   not http.user_agent contains "Mozilla" and
   not http.user_agent contains "Chrome" and
   not http.user_agent contains "Safari" and
   not http.user_agent contains "Android")
Action: Managed Challenge

# Rule 3 — Blokir scanning tools (update berkala)
Expression:
  (http.user_agent contains "cf-hero" or
   http.user_agent contains "sqlmap" or
   http.user_agent contains "nikto" or
   http.user_agent contains "subfinder" or
   http.user_agent contains "nuclei" or
   http.user_agent contains "masscan" or
   http.user_agent contains "zgrab" or
   http.user_agent eq "")
Action: Block

# Rule 4 — Rate limit endpoint auth (brute-force login siswa)
Expression:
  (http.request.uri.path contains "/auth/login" and
   http.request.method eq "POST")
Action: Rate Limit — 10 req/menit per IP

# Rule 5 — Preserve anti-judol rules yang sudah ada
# (pastikan rules ini tidak konflik dengan yang sudah ada)
```

### Page Rules / Cache Rules

```
# Cache aset statis — kurangi beban VPS
Match: gr31.tech/static/*
Settings:
  Cache Level: Cache Everything
  Browser Cache TTL: 4 hours
  Edge Cache TTL: 1 day

# Bypass cache untuk API dinamis
Match: gr31.tech/api/*
Settings:
  Cache Level: Bypass

# Bypass cache untuk halaman absensi (data real-time)
Match: gr31.tech/absensi*
Settings:
  Cache Level: Bypass
```

---

## Layer 2: UFW Firewall (Hostinger VPS)

### Setup UFW — Hanya Izinkan IP Cloudflare

```bash
#!/bin/bash
# /usr/local/bin/setup-ufw-cloudflare.sh
# Jalankan SEKALI sebagai root

set -e

echo "[GR31] Resetting UFW rules..."
ufw --force reset

# Default policy
ufw default deny incoming
ufw default allow outgoing

# SSH — WAJIB sebelum enable UFW (ganti 22 jika port SSH sudah diubah)
ufw allow 22/tcp comment "SSH"

# Hanya izinkan IP Cloudflare untuk HTTP/HTTPS
# Source: https://www.cloudflare.com/ips/
CF_IPV4=(
  "173.245.48.0/20"
  "103.21.244.0/22"
  "103.22.200.0/22"
  "103.31.4.0/22"
  "141.101.64.0/18"
  "108.162.192.0/18"
  "190.93.240.0/20"
  "188.114.96.0/20"
  "197.234.240.0/22"
  "198.41.128.0/17"
  "162.158.0.0/15"
  "104.16.0.0/13"
  "104.24.0.0/14"
  "172.64.0.0/13"
  "131.0.72.0/22"
)

CF_IPV6=(
  "2400:cb00::/32"
  "2606:4700::/32"
  "2803:f800::/32"
  "2405:b500::/32"
  "2405:8100::/32"
  "2a06:98c0::/29"
  "2c0f:f248::/32"
)

for ip in "${CF_IPV4[@]}" "${CF_IPV6[@]}"; do
  ufw allow from "$ip" to any port 80,443 proto tcp comment "Cloudflare"
done

ufw --force enable
echo "[GR31] UFW setup complete. Status:"
ufw status numbered
```

### Cron: Auto-Update IP Cloudflare (Mingguan)

```bash
#!/bin/bash
# /usr/local/bin/update-cf-ips.sh
# Cron: 0 2 * * 0 /usr/local/bin/update-cf-ips.sh >> /var/log/gr31-cf-update.log 2>&1

LOG_PREFIX="[$(date '+%Y-%m-%d %H:%M:%S')] [GR31-CF-UPDATE]"

# Ambil IP terbaru dari Cloudflare
CF_IPV4=$(curl -sf https://www.cloudflare.com/ips-v4 || echo "")
CF_IPV6=$(curl -sf https://www.cloudflare.com/ips-v6 || echo "")

if [ -z "$CF_IPV4" ]; then
  echo "$LOG_PREFIX ERROR: Gagal fetch CF IPs. Aborted."
  exit 1
fi

# Hapus rules Cloudflare lama
ufw status numbered | grep "Cloudflare" | awk '{print $1}' | \
  tr -d '[]' | sort -rn | xargs -I{} ufw --force delete {}

# Tambah rules baru
for ip in $CF_IPV4 $CF_IPV6; do
  ufw allow from "$ip" to any port 80,443 proto tcp comment "Cloudflare" > /dev/null
done

echo "$LOG_PREFIX CF IPs updated successfully."
```

---

## Layer 3: Nginx (Reverse Proxy)

```nginx
# /etc/nginx/conf.d/gr31-cloudflare.conf

# Restore IP asli dari header Cloudflare
set_real_ip_from 173.245.48.0/20;
set_real_ip_from 103.21.244.0/22;
set_real_ip_from 103.22.200.0/22;
set_real_ip_from 103.31.4.0/22;
set_real_ip_from 141.101.64.0/18;
set_real_ip_from 108.162.192.0/18;
set_real_ip_from 190.93.240.0/20;
set_real_ip_from 188.114.96.0/20;
set_real_ip_from 197.234.240.0/22;
set_real_ip_from 198.41.128.0/17;
set_real_ip_from 162.158.0.0/15;
set_real_ip_from 104.16.0.0/13;
set_real_ip_from 104.24.0.0/14;
set_real_ip_from 172.64.0.0/13;
set_real_ip_from 131.0.72.0/22;
real_ip_header CF-Connecting-IP;

# Zone rate limiting
# "school_general" — toleran untuk burst jam absensi
# Rate 100r/s = cukup longgar untuk 700 siswa dalam 1 menit
limit_req_zone  $binary_remote_addr zone=school_general:20m rate=100r/s;
limit_req_zone  $binary_remote_addr zone=school_auth:10m    rate=5r/s;
limit_req_zone  $binary_remote_addr zone=school_absen:10m   rate=10r/s;
limit_conn_zone $binary_remote_addr zone=school_conn:10m;

server {
    listen 80;
    server_name gr31.tech www.gr31.tech smkn31jkt.sch.id www.smkn31jkt.sch.id;

    # Blokir akses via IP langsung (double-check setelah Cloudflare WAF)
    if ($host ~* "^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$") {
        return 444;  # Drop koneksi tanpa response
    }

    # Batasi ukuran body — cukup untuk upload foto profil
    client_max_body_size 5m;
    client_body_timeout  10s;
    client_header_timeout 10s;

    # Endpoint auth — rate limit ketat (anti brute-force)
    location ~* ^/api/v[0-9]+/(auth|login|register) {
        limit_req  zone=school_auth burst=10 nodelay;
        limit_conn school_conn 5;
        proxy_pass http://127.0.0.1:8080;
        include    /etc/nginx/snippets/proxy-headers.conf;
    }

    # Endpoint absensi — rate limit sedang, burst toleran
    location ~* ^/api/v[0-9]+/absensi {
        limit_req  zone=school_absen burst=30 nodelay;
        limit_conn school_conn 10;
        proxy_pass http://127.0.0.1:8080;
        include    /etc/nginx/snippets/proxy-headers.conf;
    }

    # Semua endpoint lain
    location / {
        limit_req  zone=school_general burst=200 nodelay;
        limit_conn school_conn 20;
        proxy_pass http://127.0.0.1:8080;
        include    /etc/nginx/snippets/proxy-headers.conf;
    }
}

# /etc/nginx/snippets/proxy-headers.conf
# proxy_set_header Host              $host;
# proxy_set_header X-Real-IP         $remote_addr;
# proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
# proxy_set_header CF-Connecting-IP  $http_cf_connecting_ip;
# proxy_connect_timeout              5s;
# proxy_send_timeout                 15s;
# proxy_read_timeout                 30s;
```

---

## Layer 4: Go/Gin Middleware Stack

### 4.1 Registrasi Middleware (urutan kritis)

```go
// internal/server/middleware.go

package server

import (
    "github.com/gin-gonic/gin"
    "gr31.tech/internal/middleware"
    "gr31.tech/internal/config"
    "github.com/redis/go-redis/v9"
)

func RegisterMiddlewares(r *gin.Engine, cfg *config.Config, rdb *redis.Client) {
    // Urutan ini kritis — JANGAN diubah urutannya
    r.Use(middleware.RealIPResolver())              // 1. Ambil IP asli dari CF header
    r.Use(middleware.CloudflareGuard(cfg))          // 2. Tolak non-Cloudflare (production only)
    r.Use(middleware.IPBanCheck(rdb))               // 3. Cek IP banned di Redis
    r.Use(middleware.SchoolAwareRateLimiter(rdb))   // 4. Rate limit kontekstual
    r.Use(middleware.RequestSizeGuard(5 << 20))     // 5. Max 5MB body
    r.Use(middleware.SecurityHeaders())             // 6. HTTP security headers
    r.Use(middleware.AntiJudolFilter())             // 7. Filter judol (existing — preserve)
    r.Use(gin.Recovery())                           // 8. Panic recovery terakhir
}
```

### 4.2 Real IP Resolver

```go
// internal/middleware/realip.go
package middleware

import (
    "net"
    "strings"

    "github.com/gin-gonic/gin"
)

// CtxKeyClientIP adalah key untuk menyimpan IP client di Gin context
const CtxKeyClientIP = "gr31.client_ip"

var cloudflareIPNets []*net.IPNet

func init() {
    cidrs := []string{
        "173.245.48.0/20", "103.21.244.0/22", "103.22.200.0/22",
        "103.31.4.0/22", "141.101.64.0/18", "108.162.192.0/18",
        "190.93.240.0/20", "188.114.96.0/20", "197.234.240.0/22",
        "198.41.128.0/17", "162.158.0.0/15", "104.16.0.0/13",
        "104.24.0.0/14", "172.64.0.0/13", "131.0.72.0/22",
        "127.0.0.0/8", // localhost untuk development
    }
    for _, cidr := range cidrs {
        _, network, err := net.ParseCIDR(cidr)
        if err == nil {
            cloudflareIPNets = append(cloudflareIPNets, network)
        }
    }
}

// IsCloudflareIP memeriksa apakah IP berasal dari Cloudflare
func IsCloudflareIP(ipStr string) bool {
    ip := net.ParseIP(strings.TrimSpace(ipStr))
    if ip == nil {
        return false
    }
    for _, network := range cloudflareIPNets {
        if network.Contains(ip) {
            return true
        }
    }
    return false
}

// RealIPResolver mengambil IP client asli dari header CF-Connecting-IP
func RealIPResolver() gin.HandlerFunc {
    return func(c *gin.Context) {
        remoteIP := c.RemoteIP()

        if IsCloudflareIP(remoteIP) {
            if cfIP := c.GetHeader("CF-Connecting-IP"); cfIP != "" {
                c.Set(CtxKeyClientIP, cfIP)
                c.Next()
                return
            }
        }

        c.Set(CtxKeyClientIP, remoteIP)
        c.Next()
    }
}

// GetClientIP helper — gunakan ini di semua handler, bukan c.ClientIP()
func GetClientIP(c *gin.Context) string {
    if ip, exists := c.Get(CtxKeyClientIP); exists {
        return ip.(string)
    }
    return c.RemoteIP()
}
```

### 4.3 Cloudflare Guard

```go
// internal/middleware/cloudflare_guard.go
package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gr31.tech/internal/config"
)

// CloudflareGuard — di production, hanya izinkan traffic dari Cloudflare.
// Development/staging bisa bypass dengan APP_ENV != "production".
func CloudflareGuard(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !cfg.IsProduction() {
            c.Next()
            return
        }

        remoteIP := c.RemoteIP()
        if !IsCloudflareIP(remoteIP) {
            // Simulasikan server tidak ada — jangan berikan info error
            c.AbortWithStatus(http.StatusServiceUnavailable)
            return
        }

        c.Next()
    }
}
```

### 4.4 School-Aware Rate Limiter

```go
// internal/middleware/ratelimit.go
package middleware

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
)

// SchoolTimeZone — WIB (UTC+7)
const SchoolTimeZone = "Asia/Jakarta"

// isPeakHour mendeteksi jam sibuk absensi (07:00-08:30, 11:30-13:00, 14:30-16:00)
func isPeakHour() bool {
    loc, _ := time.LoadLocation(SchoolTimeZone)
    now := time.Now().In(loc)

    hour := now.Hour()
    minute := now.Minute()
    totalMin := hour*60 + minute

    weekday := now.Weekday()
    // Sabtu/Minggu tidak ada absensi
    if weekday == time.Saturday || weekday == time.Sunday {
        return false
    }

    // Peak windows: 06:45-08:30, 11:30-13:00, 14:30-16:00
    peakWindows := [][2]int{
        {6*60 + 45, 8*60 + 30},
        {11*60 + 30, 13*60 + 0},
        {14*60 + 30, 16*60 + 0},
    }

    for _, w := range peakWindows {
        if totalMin >= w[0] && totalMin <= w[1] {
            return true
        }
    }
    return false
}

type RateRule struct {
    NormalLimit int           // Request per window di luar peak hour
    PeakLimit   int           // Request per window saat peak hour (lebih longgar)
    Window      time.Duration
    AutoBanAt   int           // Multiplier — auto-ban jika melebihi NormalLimit * AutoBanAt
}

// rateRules — dikonfigurasi berdasarkan endpoint GR31
var rateRules = map[string]RateRule{
    // Absensi: bisa 700 siswa dalam 30 menit → ~24 req/menit/siswa tapi per IP
    // Satu siswa cukup 1-3 request untuk absensi, beri burst 20
    "absensi": {NormalLimit: 10, PeakLimit: 30, Window: time.Minute, AutoBanAt: 20},

    // Auth/login: ketat untuk cegah brute-force
    "auth": {NormalLimit: 5, PeakLimit: 10, Window: time.Minute, AutoBanAt: 10},

    // Kegiatan siswa: akses normal
    "kegiatan": {NormalLimit: 30, PeakLimit: 60, Window: time.Minute, AutoBanAt: 15},

    // API umum: default
    "default": {NormalLimit: 60, PeakLimit: 120, Window: time.Minute, AutoBanAt: 10},
}

// resolveRule mencocokkan path ke rule yang sesuai
func resolveRule(path string) (string, RateRule) {
    switch {
    case containsAny(path, "/absensi", "/attendance"):
        return "absensi", rateRules["absensi"]
    case containsAny(path, "/auth", "/login", "/register", "/logout"):
        return "auth", rateRules["auth"]
    case containsAny(path, "/kegiatan", "/activity", "/event"):
        return "kegiatan", rateRules["kegiatan"]
    default:
        return "default", rateRules["default"]
    }
}

func containsAny(s string, subs ...string) bool {
    for _, sub := range subs {
        if len(s) >= len(sub) && (s == sub || len(s) > len(sub) &&
            (s[:len(sub)] == sub || contains(s, sub))) {
            return true
        }
    }
    return false
}

func contains(s, sub string) bool {
    for i := 0; i <= len(s)-len(sub); i++ {
        if s[i:i+len(sub)] == sub {
            return true
        }
    }
    return false
}

// SchoolAwareRateLimiter adalah rate limiter yang sadar konteks jam sekolah
func SchoolAwareRateLimiter(rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := GetClientIP(c)
        path := c.FullPath()
        if path == "" {
            path = c.Request.URL.Path
        }

        ruleName, rule := resolveRule(path)
        peak := isPeakHour()

        // Tentukan limit berdasarkan jam
        limit := rule.NormalLimit
        if peak {
            limit = rule.PeakLimit
        }

        key := fmt.Sprintf("rl:gr31:%s:%s:%s", ruleName, ip,
            time.Now().Truncate(rule.Window).Format("150405"))

        ctx := context.Background()

        // Pipeline Redis — atomic increment + expire
        pipe := rdb.TxPipeline()
        incr := pipe.Incr(ctx, key)
        pipe.Expire(ctx, key, rule.Window+5*time.Second)
        if _, err := pipe.Exec(ctx); err != nil {
            // Redis error → fail open (jangan blokir siswa)
            c.Next()
            return
        }

        count := incr.Val()
        remaining := int64(limit) - count
        if remaining < 0 {
            remaining = 0
        }

        c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
        c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
        c.Header("X-RateLimit-Peak-Mode", fmt.Sprintf("%v", peak))

        // Auto-ban jika sangat agresif (indikasi bot/attack)
        if int(count) > limit*rule.AutoBanAt {
            duration := 6 * time.Hour
            BanIPWithReason(rdb, ip, duration, fmt.Sprintf("auto-ban: %s flood", ruleName))
            c.AbortWithStatus(http.StatusForbidden)
            return
        }

        if int(count) > limit {
            c.Header("Retry-After", fmt.Sprintf("%d", int(rule.Window.Seconds())))
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error":    "terlalu banyak permintaan, coba lagi sebentar",
                "retry_in": int(rule.Window.Seconds()),
            })
            return
        }

        c.Next()
    }
}
```

### 4.5 IP Ban System

```go
// internal/middleware/ipban.go
package middleware

import (
    "context"
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "net/http"
)

const (
    banKeyPrefix    = "gr31:ipban:"
    banReasonPrefix = "gr31:ipban_reason:"
)

// BanIPWithReason menyimpan ban beserta alasannya (untuk audit log)
func BanIPWithReason(rdb *redis.Client, ip string, duration time.Duration, reason string) {
    ctx := context.Background()
    pipe := rdb.Pipeline()
    pipe.Set(ctx, banKeyPrefix+ip, "1", duration)
    pipe.Set(ctx, banReasonPrefix+ip, reason, duration)
    pipe.Exec(ctx)
}

// ManualBanIP untuk digunakan dari admin endpoint
func ManualBanIP(rdb *redis.Client, ip string, hours int, reason string) {
    BanIPWithReason(rdb, ip, time.Duration(hours)*time.Hour, "manual: "+reason)
}

// UnbanIP untuk admin endpoint
func UnbanIP(rdb *redis.Client, ip string) {
    ctx := context.Background()
    rdb.Del(ctx, banKeyPrefix+ip, banReasonPrefix+ip)
}

// GetBanReason mengambil alasan ban (untuk admin panel)
func GetBanReason(rdb *redis.Client, ip string) string {
    val, err := rdb.Get(context.Background(), banReasonPrefix+ip).Result()
    if err != nil {
        return ""
    }
    return val
}

// IPBanCheck middleware
func IPBanCheck(rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := GetClientIP(c)
        key := banKeyPrefix + ip

        banned, err := rdb.Exists(context.Background(), key).Result()
        if err == nil && banned > 0 {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }

        c.Next()
    }
}
```

### 4.6 Request Size Guard & Security Headers

```go
// internal/middleware/security.go
package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// RequestSizeGuard membatasi ukuran body request
func RequestSizeGuard(maxBytes int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
        c.Next()
    }
}

// SecurityHeaders menambahkan HTTP security headers standar
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "SAMEORIGIN")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
        c.Header("Content-Security-Policy",
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data: https:; "+
            "connect-src 'self'")
        // Hapus header yang mengekspos tech stack
        c.Header("Server", "")
        c.Header("X-Powered-By", "")
        c.Next()
    }
}
```

### 4.7 Admin Endpoint untuk IP Management

```go
// internal/handler/admin/security.go
package admin

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "gr31.tech/internal/middleware"
)

type BanRequest struct {
    IP     string `json:"ip"     binding:"required"`
    Hours  int    `json:"hours"  binding:"required,min=1,max=720"`
    Reason string `json:"reason" binding:"required"`
}

func SecurityHandler(rdb *redis.Client) *SecurityHandlerImpl {
    return &SecurityHandlerImpl{rdb: rdb}
}

type SecurityHandlerImpl struct {
    rdb *redis.Client
}

// POST /admin/security/ban
func (h *SecurityHandlerImpl) BanIP(c *gin.Context) {
    var req BanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    middleware.ManualBanIP(h.rdb, req.IP, req.Hours, req.Reason)
    c.JSON(http.StatusOK, gin.H{
        "message": "IP berhasil di-ban",
        "ip":      req.IP,
        "hours":   req.Hours,
    })
}

// DELETE /admin/security/ban/:ip
func (h *SecurityHandlerImpl) UnbanIP(c *gin.Context) {
    ip := c.Param("ip")
    middleware.UnbanIP(h.rdb, ip)
    c.JSON(http.StatusOK, gin.H{"message": "IP berhasil di-unban", "ip": ip})
}

// GET /admin/security/ban/:ip
func (h *SecurityHandlerImpl) GetBanInfo(c *gin.Context) {
    ip := c.Param("ip")
    reason := middleware.GetBanReason(h.rdb, ip)
    if reason == "" {
        c.JSON(http.StatusNotFound, gin.H{"message": "IP tidak dalam daftar ban"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"ip": ip, "reason": reason})
}
```

---

## Layer 5: AstraDB + Redis — Connection Hardening

### AstraDB — Timeout & Pool Config

```go
// internal/database/astradb.go

package database

import (
    "context"
    "time"

    "github.com/datastax/astra-client-go/v2/astra"
)

type AstraConfig struct {
    Token     string
    Endpoint  string
    // Timeout untuk query — penting agar DDoS tidak exhaust goroutines
    QueryTimeout    time.Duration // default: 5s
    ConnectTimeout  time.Duration // default: 10s
    MaxConns        int           // default: 20 (sesuaikan dengan RAM VPS)
}

func DefaultAstraConfig(token, endpoint string) AstraConfig {
    return AstraConfig{
        Token:          token,
        Endpoint:       endpoint,
        QueryTimeout:   5 * time.Second,
        ConnectTimeout: 10 * time.Second,
        MaxConns:       20,
    }
}

// WrapWithTimeout membungkus context dengan timeout query
func WrapWithTimeout(cfg AstraConfig) (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), cfg.QueryTimeout)
}
```

### Redis — Config Self-Hosted (Hostinger VPS)

```go
// internal/cache/redis.go

package cache

import (
    "context"
    "time"
    "github.com/redis/go-redis/v9"
)

func NewRedisClient(host, port, password string) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:         host + ":" + port,
        Password:     password,
        DB:           0,

        // Pool sizing untuk 700 user
        PoolSize:     50,               // max concurrent connections
        MinIdleConns: 10,               // jaga koneksi siap pakai
        MaxIdleConns: 20,

        // Timeout ketat — jangan biarkan DDoS exhaust Redis pool
        DialTimeout:  3 * time.Second,
        ReadTimeout:  2 * time.Second,
        WriteTimeout: 2 * time.Second,
        PoolTimeout:  4 * time.Second,

        // Health check
        ConnMaxIdleTime: 5 * time.Minute,
    })

    // Validasi koneksi saat startup
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := rdb.Ping(ctx).Err(); err != nil {
        panic("GR31: Redis connection failed: " + err.Error())
    }

    return rdb
}
```

### Redis Security (/etc/redis/redis.conf)

```conf
# Bind hanya ke localhost — jangan expose ke internet
bind 127.0.0.1

# Wajib set password
requirepass YOUR_STRONG_REDIS_PASSWORD_HERE

# Nonaktifkan perintah berbahaya
rename-command FLUSHALL  ""
rename-command FLUSHDB   ""
rename-command CONFIG    "CONFIG_ADMIN_ONLY_XYZ"
rename-command DEBUG     ""
rename-command KEYS      ""  # Gunakan SCAN, bukan KEYS

# Batasi memory — cegah OOM saat DDoS flood key
maxmemory 256mb
maxmemory-policy allkeys-lru

# Nonaktifkan akses tanpa auth
protected-mode yes
```

---

## Layer 6: Fail2Ban (Last Resort)

```ini
# /etc/fail2ban/filter.d/gr31-nginx.conf
[Definition]
failregex = ^<HOST> -.* "(GET|POST|HEAD|PUT|DELETE|PATCH).*" (400|403|404|429|444|500) \d+
ignoreregex = ^<HOST> -.* "(GET|POST) /health.*" 200

# /etc/fail2ban/jail.d/gr31.conf
[gr31-nginx-ddos]
enabled   = true
port      = http,https
filter    = gr31-nginx
logpath   = /var/log/nginx/access.log
maxretry  = 200
findtime  = 60
bantime   = 3600
action    = iptables-multiport[name=GR31_DDOS, port="http,https", protocol=tcp]

[gr31-nginx-bruteforce]
enabled   = true
port      = http,https
filter    = gr31-nginx
logpath   = /var/log/nginx/access.log
maxretry  = 20
findtime  = 60
bantime   = 86400
# Filter khusus untuk endpoint auth
# Sesuaikan failregex untuk path /api/*/auth
```

---

## Integrasi dengan Anti-Judol (Existing)

Pastikan middleware anti-judol yang sudah ada tidak konflik dengan stack baru ini.
Posisi yang tepat ada di **urutan ke-7**, setelah rate limiter.

```go
// Contoh integrasi di RegisterMiddlewares:
r.Use(middleware.RealIPResolver())
r.Use(middleware.CloudflareGuard(cfg))
r.Use(middleware.IPBanCheck(rdb))
r.Use(middleware.SchoolAwareRateLimiter(rdb))
r.Use(middleware.RequestSizeGuard(5 << 20))
r.Use(middleware.SecurityHeaders())
r.Use(middleware.AntiJudolFilter()) // ← existing, preserve posisinya di sini
r.Use(gin.Recovery())

// AntiJudolFilter yang sudah ada tidak perlu diubah.
// Rate limiter akan mencegah flood sebelum request sampai ke filter ini,
// sehingga tidak ada performa penalty dari pengecekan judol saat DDoS.
```

---

## Environment Variables

```env
# ===== GR31 APP CONFIG =====
APP_ENV=production
APP_NAME=gr31-school-system
DOMAINS=gr31.tech,smkn31jkt.sch.id

# ===== REDIS (self-hosted) =====
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=your_strong_redis_password_here

# ===== ASTRADB =====
ASTRA_TOKEN=AstraCS:xxxxxxxxxxxx
ASTRA_ENDPOINT=https://xxxx-region.apps.astra.datastax.com

# ===== RATE LIMIT OVERRIDE (opsional, default sudah bagus) =====
# Uncomment untuk override saat event khusus (misal MPLS, penerimaan siswa baru)
# RL_PEAK_MULTIPLIER=3     # 3x lebih longgar dari normal saat event besar
# RL_ABSENSI_PEAK=50       # Override limit absensi saat jam padat

# ===== SECURITY =====
# IP whitelist admin (pisahkan koma, bisa IP rumah/kantor)
ADMIN_IP_WHITELIST=203.0.113.10
```

---

## Checklist Implementasi

### 🔴 Prioritas Tinggi — Lakukan Sekarang

- [ ] Jalankan `setup-ufw-cloudflare.sh` di VPS
- [ ] Pastikan **semua** subdomain `gr31.tech` dan `smkn31jkt.sch.id` pakai **CF Proxy** (awan orange)
- [ ] Tambah WAF Rule Cloudflare: blokir akses via IP langsung
- [ ] Tambah WAF Rule: blokir User-Agent scanner (`cf-hero`, `sqlmap`, dll)
- [ ] Aktifkan **Bot Fight Mode** di Cloudflare dashboard

### 🟡 Prioritas Menengah — Minggu Ini

- [ ] Implementasi `RealIPResolver`, `CloudflareGuard`, `IPBanCheck` di Gin
- [ ] Implementasi `SchoolAwareRateLimiter` dengan peak hour windows
- [ ] Konfigurasi Redis hardening (`/etc/redis/redis.conf`)
- [ ] Setup Nginx dengan `limit_req_zone` dan `limit_conn_zone`
- [ ] Test rate limiter tidak blokir siswa saat simulasi jam 07.00

### 🟢 Prioritas Normal — Bulan Ini

- [ ] Setup Fail2Ban dengan filter gr31-nginx
- [ ] Tambah admin endpoint `/admin/security/ban` dan `/admin/security/unban`
- [ ] Pasang cron weekly untuk update IP Cloudflare
- [ ] Load test dengan k6/locust simulasi 700 user concurrent (terutama endpoint absensi)
- [ ] Verifikasi integrasi anti-judol tidak terganggu

---

## Response Saat Terjadi Serangan

```
🚨 DETEKSI SERANGAN:

1. Cek Cloudflare Analytics → Traffic spike di tab "Security"
2. Naikkan Cloudflare Security Level ke "Under Attack Mode"
   (Dashboard → Security → Settings → Security Level)

3. Identifikasi IP sumber:
   sudo tail -f /var/log/nginx/access.log | grep " 429 "

4. Ban IP manual via admin endpoint:
   POST /admin/security/ban
   {"ip": "x.x.x.x", "hours": 24, "reason": "ddos-manual"}

5. Jika skala besar, aktifkan Cloudflare "I'm Under Attack" mode
   → Semua visitor dapat JS challenge otomatis
   → Biasanya cukup untuk menghentikan bot DDoS

6. Setelah serangan mereda, turunkan kembali ke Security Level "Medium"
```
