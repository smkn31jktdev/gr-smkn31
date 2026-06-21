package ping

import (
	"context"
	"net/http"
	"time"

	"be-gr31/internal/model/common"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Handler
type Handler struct {
	rdb *redis.Client
}

// NewHandler
func NewHandler(rdb *redis.Client) *Handler {
	return &Handler{rdb: rdb}
}

// Root handle GET
func (h *Handler) Root(c *gin.Context) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, common.OK(gin.H{
			"name":    "GR31 Backend",
			"version": "1.0.0",
			"status":  "running",
		}, "GR31 Backend API is running"))
		return
	}

	now := time.Now().Format("02 Jan 2006 · 15:04:05 WIB")

	html := `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>GR31 Backend API</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      background: #0d1117;
      color: #e6edf3;
      font-family: 'Courier New', Courier, monospace;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 40px 20px;
    }
    .container {
      max-width: 860px;
      width: 100%;
    }
    .ascii {
      color: #58a6ff;
      font-size: clamp(7px, 1.1vw, 13px);
      line-height: 1.2;
      white-space: pre;
      text-align: center;
      margin-bottom: 8px;
      text-shadow: 0 0 20px rgba(88,166,255,0.4);
    }
    .tagline {
      text-align: center;
      color: #8b949e;
      font-size: 13px;
      margin-bottom: 36px;
      letter-spacing: 3px;
      text-transform: uppercase;
    }
    .badge-row {
      display: flex;
      justify-content: center;
      gap: 10px;
      flex-wrap: wrap;
      margin-bottom: 36px;
    }
    .badge {
      padding: 4px 12px;
      border-radius: 20px;
      font-size: 11px;
      font-weight: bold;
      letter-spacing: 1px;
    }
    .badge-green  { background: rgba(35,134,54,0.25);  color: #3fb950; border: 1px solid #238636; }
    .badge-blue   { background: rgba(31,111,235,0.2);  color: #58a6ff; border: 1px solid #1f6feb; }
    .badge-purple { background: rgba(188,140,255,0.15); color: #bc8cff; border: 1px solid #6e40c9; }
    .badge-orange { background: rgba(210,153,34,0.2);  color: #d29922; border: 1px solid #9e6a03; }
    .divider {
      border: none;
      border-top: 1px solid #21262d;
      margin: 0 0 28px 0;
    }
    .section-title {
      color: #8b949e;
      font-size: 10px;
      letter-spacing: 2px;
      text-transform: uppercase;
      margin-bottom: 12px;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
      gap: 12px;
      margin-bottom: 28px;
    }
    .card {
      background: #161b22;
      border: 1px solid #21262d;
      border-radius: 8px;
      padding: 16px 20px;
      transition: border-color .2s;
    }
    .card:hover { border-color: #58a6ff44; }
    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 10px;
    }
    .dot {
      width: 8px; height: 8px;
      border-radius: 50%;
    }
    .dot-green  { background: #3fb950; box-shadow: 0 0 6px #3fb950; }
    .dot-blue   { background: #58a6ff; box-shadow: 0 0 6px #58a6ff; }
    .dot-yellow { background: #d29922; box-shadow: 0 0 6px #d29922; }
    .dot-purple { background: #bc8cff; box-shadow: 0 0 6px #bc8cff; }
    .card-title { color: #e6edf3; font-size: 13px; font-weight: bold; }
    .card-body { color: #8b949e; font-size: 12px; line-height: 1.8; }
    .card-body span { color: #e6edf3; }
    .endpoint-list { list-style: none; }
    .endpoint-list li {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 4px 0;
      font-size: 12px;
      color: #8b949e;
      border-bottom: 1px solid #21262d22;
    }
    .endpoint-list li:last-child { border-bottom: none; }
    .method {
      font-size: 10px;
      font-weight: bold;
      padding: 2px 7px;
      border-radius: 4px;
      min-width: 38px;
      text-align: center;
    }
    .get    { background: rgba(35,134,54,0.2);   color: #3fb950; }
    .post   { background: rgba(31,111,235,0.2);  color: #58a6ff; }
    .put    { background: rgba(210,153,34,0.2);  color: #d29922; }
    .delete { background: rgba(248,81,73,0.2);   color: #f85149; }
    .path { color: #e6edf3; }
    .path-desc { color: #484f58; font-size: 11px; margin-left: auto; }
    .footer {
      text-align: center;
      color: #484f58;
      font-size: 11px;
      margin-top: 8px;
    }
    .footer a { color: #58a6ff; text-decoration: none; }
  </style>
</head>
<body>
<div class="container">

  <pre class="ascii"> 
 ________  ________  ________    _____     
|\   ____\|\   __  \|\_____  \  / __  \    
\ \  \___|\ \  \|\  \|____|\ /_|\/_|\  \   
 \ \  \  __\ \   _  _\    \|\  \|/ \ \  \  
  \ \  \|\  \ \  \\  \|  __\_\  \   \ \  \ 
   \ \_______\ \__\\ _\ |\_______\   \ \__\
    \|_______|\|__|\|__|\|_______|    \|__|
                                                                
	</pre>

  <p class="tagline">Backend API &nbsp;·&nbsp; SMK Negeri 31 Jakarta</p>

  <div class="badge-row">
    <span class="badge badge-green">● RUNNING</span>
    <span class="badge badge-blue">Go 1.22+</span>
    <span class="badge badge-blue">Gin Framework</span>
    <span class="badge badge-purple">AstraDB</span>
    <span class="badge badge-green">Mongodb</span>
    <span class="badge badge-orange">Redis</span>
    <span class="badge badge-blue">JWT HS256</span>
  </div>

  <hr class="divider"/>

  <p class="section-title">Server Info</p>
  <div class="grid" style="grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); margin-bottom: 28px;">
    <div class="card">
      <div class="card-header"><div class="dot dot-green"></div><span class="card-title">Status</span></div>
      <div class="card-body">
        Version &nbsp;<span>v1.0.0</span><br/>
        Env &nbsp;&nbsp;&nbsp;&nbsp;<span>development</span><br/>
        Time &nbsp;&nbsp;&nbsp;<span>` + now + `</span>
      </div>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-blue"></div><span class="card-title">Base URL</span></div>
      <div class="card-body">
        HTTP &nbsp;&nbsp;<span>http://localhost:8080</span><br/>
        Prefix &nbsp;<span>/v1/student/*</span><br/>
        Prefix &nbsp;<span>/v1/admin/*</span>
      </div>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-yellow"></div><span class="card-title">Auth</span></div>
      <div class="card-body">
        Type &nbsp;&nbsp;&nbsp;<span>Bearer JWT</span><br/>
        Header &nbsp;<span>Authorization</span><br/>
        Algo &nbsp;&nbsp;&nbsp;<span>HS256</span>
      </div>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-purple"></div><span class="card-title">Idempotency</span></div>
      <div class="card-body">
        Header &nbsp;<span>X-Idempotency-Key</span><br/>
        Storage &nbsp;<span>Redis TTL 24h</span><br/>
        Required <span>POST mutating</span>
      </div>
    </div>
  </div>

  <p class="section-title">Endpoints</p>
  <div class="grid">
    <div class="card">
      <div class="card-header"><div class="dot dot-blue"></div><span class="card-title">Auth</span></div>
      <ul class="endpoint-list">
        <li><span class="method post">POST</span><span class="path">/v1/student/login</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/student/me</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/student/refresh-token</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/admin/login</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/admin/me</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/admin/refresh-token</span></li>
      </ul>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-green"></div><span class="card-title">Kehadiran</span></div>
      <ul class="endpoint-list">
        <li><span class="method post">POST</span><span class="path">/v1/student/kehadiran</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/student/kehadiran</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/student/kehadiran/upload-izin</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/admin/kehadiran</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/kehadiran</span></li>
        <li><span class="method delete">DEL</span><span class="path">/v1/admin/kehadiran</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/rekap-bulanan</span></li>
      </ul>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-purple"></div><span class="card-title">Gerakan 7</span></div>
      <ul class="endpoint-list">
        <li><span class="method post">POST</span><span class="path">/v1/student/g7</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/student/g7</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/student/g7/:tanggal</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/g7</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/g7/summary</span></li>
        <li><span class="method delete">DEL</span><span class="path">/v1/admin/g7</span></li>
      </ul>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-yellow"></div><span class="card-title">Kegiatan & Bukti</span></div>
      <ul class="endpoint-list">
        <li><span class="method post">POST</span><span class="path">/v1/student/kegiatan</span></li>
        <li><span class="method put">PUT</span><span class="path">/v1/student/kegiatan</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/student/kegiatan</span></li>
        <li><span class="method delete">DEL</span><span class="path">/v1/student/kegiatan</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/student/bukti</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/student/bukti</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/kegiatan</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/bukti</span></li>
      </ul>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-blue"></div><span class="card-title">Aduan</span></div>
      <ul class="endpoint-list">
        <li><span class="method post">POST</span><span class="path">/v1/student/aduan</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/student/aduan</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/aduan</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/aduan/room</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/admin/aduan/status</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/admin/aduan/respond</span></li>
      </ul>
    </div>
    <div class="card">
      <div class="card-header"><div class="dot dot-green"></div><span class="card-title">Manajemen User</span></div>
      <ul class="endpoint-list">
        <li><span class="method post">POST</span><span class="path">/v1/admin/students</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/students</span></li>
        <li><span class="method put">PUT</span><span class="path">/v1/admin/students/:id</span></li>
        <li><span class="method delete">DEL</span><span class="path">/v1/admin/students/:id</span></li>
        <li><span class="method post">POST</span><span class="path">/v1/admin/admins</span></li>
        <li><span class="method get">GET</span><span class="path">/v1/admin/admins</span></li>
        <li><span class="method delete">DEL</span><span class="path">/v1/admin/admins/:id</span></li>
      </ul>
    </div>
  </div>

  <hr class="divider"/>
  <p class="footer">
    GR31 Backend &nbsp;·&nbsp; SMK Negeri 31 Jakarta &nbsp;·&nbsp;
    <a href="/health">Health Check</a> &nbsp;·&nbsp;
    <a href="/v1/ping">Ping</a>
  </p>

</div>
</body>
</html>`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// Health menangani GET
func (h *Handler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	redisStatus := "ok"
	overallStatus := "ok"

	if err := h.rdb.Ping(ctx).Err(); err != nil {
		redisStatus = "unavailable"
		overallStatus = "degraded"
	}

	httpStatus := http.StatusOK
	if overallStatus == "degraded" {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, common.OK(gin.H{
		"status": overallStatus,
		"time":   time.Now().Format(time.RFC3339),
		"services": gin.H{
			"redis": redisStatus,
		},
	}, "health check"))
}

// Ping
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, common.OK(gin.H{
		"pong": true,
		"time": time.Now().Format(time.RFC3339),
	}, "pong"))
}
