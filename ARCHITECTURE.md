# Arsitektur Sistem Backend SMK Negeri 31 Jakarta

Dokumen ini menjelaskan arsitektur teknis dari backend API SMK Negeri 31 Jakarta secara lengkap dan terperinci.

## 📐 Arsitektur Tingkat Tinggi

Berikut adalah alur request dari Client (baik web SvelteKit maupun aplikasi mobile) melalui layer API router, middleware, logika bisnis, hingga persistensi data:

```mermaid
graph TD
    Client["CLIENT (SvelteKit App / Mobile)"] -->|HTTP/REST| Router["API Gateway / Router (Gin Web Framework)"]

    subgraph "Request Filtering Stack"
        Router --> Middleware["Middleware Stack"]
        Middleware --> CORS["CORS & Request Logging"]
        CORS --> JWT["JWT Validation & Auth Check"]
        JWT --> RateLimiter["Rate Limiting (School-Aware)"]
    end

    RateLimiter --> Handlers["Feature Handlers (internal/features/*)"]

    subgraph "Core Clean Architecture"
        Handlers --> Services["Service Layer (Business Logic Validation)"]
        Services --> Repos["Repository Layer (Data Abstraction Pattern)"]
    end

    subgraph "Data Storage & Caching"
        Repos <--> RedisCache["Redis (Session & Rate Limit Cache)"]
        Repos <--> AstraDB["Astra DB / MongoDB (Primary Database)"]
    end

    style Client fill:#ff3e00,stroke:#333,stroke-width:2px,color:#fff
    style Router fill:#00add8,stroke:#333,stroke-width:2px,color:#fff
    style RedisCache fill:#dc382d,stroke:#333,stroke-width:2px,color:#fff
    style AstraDB fill:#1287b1,stroke:#333,stroke-width:2px,color:#fff
```

---

## 🏗️ Clean Architecture (Layered Architecture)

Backend menggunakan arsitektur bersih (**Clean Architecture**) yang membagi fungsionalitas ke dalam 4 lapisan utama dengan arah dependensi ke dalam:

```mermaid
graph TD
    subgraph "Clean Architecture Layers"
        Handler["1. Handler Layer (Presentation/API)
        - Menerima HTTP request
        - Validasi skema input
        - Return HTTP response"]

        Service["2. Service Layer (Business Logic)
        - Logika bisnis inti
        - Validasi aturan sekolah
        - Orkestrasi query & update"]

        Repo["3. Repository Layer (Data Access)
        - Query database & CRUD
        - Abstraksi vendor database"]

        Storage["4. Storage Layer (Database Client)
        - Inisialisasi koneksi DB
        - Query builder & execution"]

        Handler -->|Memanggil| Service
        Service -->|Memanggil| Repo
        Repo -->|Memanggil| Storage
    end
```

### 1. **Handler Layer** (API/Presentation)

- **Lokasi**: `internal/features/*/handler.go`
- **Tanggung Jawab**:
  - Menerima request HTTP.
  - Memvalidasi input skema dari client (_Request Binding_).
  - Memanggil service layer yang bersangkutan.
  - Mengembalikan HTTP response terstandarisasi.
- **File Utama**:
  - `internal/features/auth/handler.go` - Autentikasi akun.
  - `internal/features/student/handler.go` - Manajemen Siswa.
  - `internal/features/kehadiran/handler.go` - Logika kehadiran & Geolocation.
  - `internal/features/aduan/handler.go` - Layanan Pengaduan BK.

### 2. **Service Layer** (Business Logic)

- **Lokasi**: `internal/features/*/service.go`
- **Tanggung Jawab**:
  - Memproses logika bisnis utama aplikasi.
  - Melakukan validasi aturan sekolah (_business rules_).
  - Mengoordinasikan pemanggilan ke Repository layer.
  - Menangani _error handling_ bisnis.

### 3. **Repository Layer** (Data Access)

- **Lokasi**: `internal/features/*/repo.go`
- **Tanggung Jawab**:
  - Berinteraksi langsung dengan database driver.
  - Melakukan query data & persistensi data.
  - Menerapkan pattern Repository untuk memisahkan logika bisnis dari detail database.

### 4. **Storage Layer** (Database Client)

- **Lokasi**: `internal/storage/` (termasuk `internal/storage/astra/`)
- **Tanggung Jawab**:
  - Mengelola koneksi database (pool koneksi).
  - Abstraksi database AstraDB (Cassandra) dan MongoDB.

---

## 📦 Feature-Based Structure

Kode backend dikelompokkan berdasarkan fitur (**Feature-Based Packaging**) untuk mempermudah pemeliharaan:

```
internal/features/<feature-name>/
├── handler.go      # HTTP handlers & binding
├── service.go      # Business logic & validations
├── repo.go         # Data access & queries
└── routes.go       # Route registration
```

### Fitur-Fitur Utama

#### 1. **Auth** - Autentikasi & Otorisasi

```
/internal/features/auth/
├── handler.go      - Endpoint Login, Register, Logout
├── service.go      - Logika validasi kredensial & pembuatan JWT
├── repo.go         - Query & persistensi user (Siswa/Admin)
└── routes.go       - Registrasi rute autentikasi
```

- **Models** (`internal/model/auth/`):
  - `admin.go` - Skema data Administrator/Guru.
  - `siswa.go` - Skema data Siswa.

#### 2. **Student** - Manajemen Siswa

```
/internal/features/student/
├── handler.go      - CRUD endpoint data siswa
├── service.go      - Logika data siswa
├── repo.go         - Query data siswa
└── routes.go       - Rute manajemen siswa
```

#### 3. **Kehadiran** - Manajemen Kehadiran

```
/internal/features/kehadiran/
├── handler.go      - Endpoint kehadiran harian
├── service.go      - Validasi koordinat GPS & radius sekolah
├── repo.go         - Query kehadiran & log
├── routes.go       - Rute kehadiran
├── upload.go       - Handler upload CSV/Excel data kehadiran
└── rekap.go        - Logika pembuatan laporan rekap absensi
```

#### 4. **Aduan** - Sistem Pengaduan & Bimbingan Konseling (BK)

```
/internal/features/aduan/
├── handler.go      - Endpoint pelaporan aduan siswa
├── service.go      - Logika distribusi aduan & status
├── repo.go         - Query & update aduan
└── routes.go       - Rute pengaduan
```

#### 5. **Admin** - Manajemen Admin

```
/internal/features/admin/
├── handler.go      - Endpoint pengelolaan data staff/admin
├── service.go      - Logika manajemen admin
├── repo.go         - Query database staff
└── routes.go       - Rute admin
```

---

## 🔄 Request Flow

### Contoh: Alur Login Siswa

Proses autentikasi dari input kredensial di frontend hingga token JWT disimpan dan dikembalikan:

```mermaid
sequenceDiagram
    autonumber
    actor Client as Client (SvelteKit)
    participant Router as Router (Gin)
    participant Handler as Auth Handler
    participant Service as Auth Service
    participant Repo as Auth Repository
    participant DB as Astra DB / Redis

    Client->>Router: POST /api/auth/login
    Router->>Handler: Rujuk ke Handler.Login()
    Handler->>Handler: Validasi Request Body (JSON Binding)
    Handler->>Service: Panggil Service.Login()
    Service->>Repo: Panggil Repo.FindByEmail(email)
    Repo->>DB: Query data user dari Astra DB
    DB-->>Repo: Data Kredensial User
    Repo-->>Service: Return User Record
    Service->>Service: Bandingkan password hash (bcrypt)
    Service->>Service: Generate JWT Token (HS256)
    Service->>DB: Simpan sesi di Redis (TTL 24 jam)
    Service-->>Handler: Return Token & Profil
    Handler-->>Client: HTTP 200 OK + JSON Token & User Data
```

---

## 🔐 Security Layer (Middleware)

Middleware disematkan pada router Gin untuk menyaring seluruh request sebelum masuk ke aplikasi utama:

```mermaid
graph TD
    Request["Incoming HTTP Request"] --> CORS["CORS Check (Cross-Origin Setup)"]
    CORS --> JWT["JWT Validation Middleware (Token Verification)"]
    JWT --> RateLimit["Rate Limiting Middleware (School-Aware)"]
    RateLimit --> Handler["Feature Request Handler"]
    Handler --> Response["Outgoing Standardized Response"]
```

### Mekanisme Autentikasi JWT & Caching Sesi:

```
sequenceDiagram
    Client ->> Client: POST /login
    Client ->> Client: Service memvalidasi password
    Client ->> Client: Generate JWT token
    Client ->> Client: Simpan sesi di Redis dengan Key: "session:<user_id>" (TTL 24 jam)
    Client ->> Client: Kembalikan token ke Client

    rect rgb(211, 211, 211)
        Note over Client: Rute Terproteksi:
        Client ->> Client: Mengirim request dengan Header "Authorization: Bearer <token>"
        Client ->> Client: Middleware memvalidasi tanda tangan JWT
        Client ->> Client: Middleware mencocokkan ID sesi di Redis
        alt Jika Valid
            Client ->> Client: Lanjutkan ke Handler
        else Jika Invalid
            Client ->> Client: Kembalikan HTTP 401 Unauthorized
        end
    end
```

---

## 💾 Database Schema

### Koleksi di Astra DB / MongoDB

#### 1. **users**

Koleksi untuk data pengguna sistem (Siswa & Administrator):

```json
{
  "_id": "uuid",
  "email": "string",
  "nisn": "string (khusus siswa)",
  "password_hash": "string",
  "name": "string",
  "role": "student|admin",
  "class": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

#### 2. **kehadiran**

Koleksi untuk log absensi harian siswa:

```json
{
  "_id": "uuid",
  "nisn": "string",
  "date": "date",
  "status": "hadir|alfa|sakit|izin",
  "latitude": "float64",
  "longitude": "float64",
  "accuracy": "float64",
  "distance": "float64",
  "created_at": "timestamp"
}
```

#### 3. **aduan**

Koleksi untuk sistem pengaduan bimbingan konseling:

```json
{
  "_id": "uuid",
  "nisn": "string",
  "judul": "string",
  "isi": "string",
  "status": "open|responded|closed",
  "responses": [
    {
      "admin_name": "string",
      "isi": "string",
      "created_at": "timestamp"
    }
  ],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

---

## 🚀 Deployment Architecture

Sistem menggunakan skema beban terdistribusi (_multi-instance_) di belakang load balancer untuk menjamin ketersediaan tinggi:

```mermaid
graph TD
    Traffic["Client Traffic (HTTPS)"] --> LB["Nginx Reverse Proxy & Load Balancer"]

    subgraph "Application Cluster"
        LB --> Instance1["Go API Instance 1"]
        LB --> Instance2["Go API Instance 2"]
        LB --> Instance3["Go API Instance 3"]
    end

    Instance1 & Instance2 & Instance3 <--> Redis["Redis Cluster (Cache & Sessions)"]
    Instance1 & Instance2 & Instance3 <--> DB["Astra DB / MongoDB"]

    style Traffic fill:#f99,stroke:#333,stroke-width:1px
    style LB fill:#00add8,stroke:#333,stroke-width:2px,color:#fff
    style Redis fill:#dc382d,stroke:#333,stroke-width:2px,color:#fff
    style DB fill:#1287b1,stroke:#333,stroke-width:2px,color:#fff
```

---

## 📊 Data Flow Diagrams (DFD)

### 1. Alur Upload Absensi Massal (Admin/Piket)

Proses import data kehadiran via file Excel/CSV oleh Guru Piket atau Admin:

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Admin / Guru Piket
    participant Handler as Kehadiran Handler
    participant Service as Kehadiran Service
    participant Repo as Kehadiran Repository
    participant DB as Astra DB / MongoDB
    participant Redis as Redis Cache

    Admin->>Handler: Upload File Excel/CSV Kehadiran
    Handler->>Handler: Parse File menjadi Row Structs
    Handler->>Service: BulkInsertKehadiran(rows)
    loop Validasi Tiap Baris
        Service->>Service: Validasi NISN & Format Tanggal
    end
    Service->>Repo: BulkInsert(valid_rows)
    Repo->>DB: Kirim Bulk Insert Query
    DB-->>Repo: Hasil penyimpanan
    Repo-->>Service: Return Sukses
    Service->>Redis: Invalidate cache kehadiran terkait
    Service-->>Handler: Return ringkasan data (Sukses/Gagal)
    Handler-->>Admin: Kirim HTTP 200 (Menampilkan Statistik Upload)
```

### 2. Alur Penanganan Tanggapan Aduan (Guru BK)

Proses penanganan pengaduan dari Admin/Guru BK hingga notifikasi terkirim ke siswa:

```mermaid
sequenceDiagram
    autonumber
    actor Admin as Guru BK / Admin
    participant Handler as Aduan Handler
    participant Service as Aduan Service
    participant Repo as Aduan Repository
    participant DB as Astra DB / MongoDB
    participant Broker as Socket/Notification
    actor Student as Siswa

    Admin->>Handler: POST /api/aduan/{id}/respond {isi_tanggapan}
    Handler->>Service: Respond(id, admin_name, isi_tanggapan)
    Service->>Repo: FindAduanByID(id)
    Repo-->>Service: Dokumen Aduan Aktif
    Service->>Service: Tambah tanggapan & Ubah status (closed/responded)
    Service->>Repo: Update(aduan)
    Repo->>DB: Update dokumen di database
    Service->>Broker: Kirim payload notifikasi baru
    Broker-->>Student: Push Notifikasi secara Real-time (Web Socket)
    Service-->>Handler: Return sukses
    Handler-->>Admin: HTTP 200 OK (Tanggapan berhasil disimpan)
```

---

## 🛠️ Configuration Management

Sistem memisahkan kode dari konfigurasi dengan menggunakan variabel lingkungan (_Environment Variables_) yang dimuat pada saat startup:

```mermaid
graph TD
    Env[".env File / Server Environment Variables"] -->|Dimuat oleh| Loader["config.go (godotenv / viper)"]
    Loader -->|Parsing & Validasi Tipe Data| Config["Config Struct (Global)"]

    subgraph "Config Fields"
        Config --> P1["App Port (:8080)"]
        Config --> P2["Redis URI & Credentials"]
        Config --> P3["MongoDB URI & DB Name"]
        Config --> P4["JWT Secret Key"]
        Config --> P5["Environment (production/development)"]
    end
```

---

## 📝 Error Handling

Handler mengadopsi standar respon error JSON yang konsisten di semua rute API untuk memudahkan penanganan di sisi Client:

```mermaid
graph TD
    Err["Handler Menangkap Error"] --> Type{"Evaluasi Jenis Error"}

    Type -->|Validation Error| Err400["HTTP 400 Bad Request"]
    Type -->|Auth Error| Err401["HTTP 401 Unauthorized"]
    Type -->|Forbidden Access| Err403["HTTP 403 Forbidden"]
    Type -->|Resource Not Found| Err404["HTTP 404 Not Found"]
    Type -->|System/DB Failure| Err500["HTTP 500 Internal Server Error"]

    Err400 & Err401 & Err403 & Err404 & Err500 --> Res["Standardized JSON Response:
    {
      'code': 'ERROR_CODE',
      'message': 'Pesan ramah pengguna',
      'status': HTTP_STATUS
    }"]
```

---

© 2024-2026 Alvin Putra & SMK Negeri 31 Jakarta. Semua hak dilindungi.
