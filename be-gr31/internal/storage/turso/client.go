package turso

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Client struct {
	DB *sql.DB
}

func NewClient(url, token string) (*Client, error) {
	dbURL := url
	if token != "" {
		if strings.Contains(dbURL, "?") {
			dbURL = dbURL + "&authToken=" + token
		} else {
			dbURL = dbURL + "?authToken=" + token
		}
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open turso db connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping turso db: %w", err)
	}

	client := &Client{DB: db}
	if err := client.Migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate turso db: %w", err)
	}

	return client, nil
}

// Migrate
func (c *Client) Migrate() error {
	// Jalankan rename dari 'nisn' ke 'nis' bila kolom lama terdeteksi (migrasi skema)
	var hasNisn bool
	var hasDeviceInfo bool
	rows, err := c.DB.Query("PRAGMA table_info(kehadiran)")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dfltVal any
			if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltVal, &pk); err == nil {
				if name == "nisn" {
					hasNisn = true
				}
				if name == "device_info" {
					hasDeviceInfo = true
				}
			}
		}
	}
	if hasNisn {
		_, _ = c.DB.Exec("ALTER TABLE kehadiran RENAME COLUMN nisn TO nis")
		_, _ = c.DB.Exec("DROP INDEX IF EXISTS idx_kehadiran_nisn_tanggal")
	}
	if !hasDeviceInfo {
		_, _ = c.DB.Exec("ALTER TABLE kehadiran ADD COLUMN device_info TEXT")
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS kehadiran (
			id TEXT PRIMARY KEY,
			nis TEXT NOT NULL,
			nama_siswa TEXT NOT NULL,
			kelas TEXT NOT NULL,
			tanggal TEXT NOT NULL,
			hari TEXT NOT NULL,
			status TEXT NOT NULL,
			waktu_absen TEXT NOT NULL,
			alasan TEXT,
			foto_izin TEXT,
			koordinat_lat REAL,
			koordinat_lng REAL,
			jarak REAL,
			akurasi REAL,
			device_info TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_kehadiran_nis_tanggal ON kehadiran (nis, tanggal);`,
		`CREATE INDEX IF NOT EXISTS idx_kehadiran_kelas ON kehadiran (kelas);`,
		`CREATE INDEX IF NOT EXISTS idx_kehadiran_tanggal ON kehadiran (tanggal);`,
	}

	for _, query := range queries {
		if _, err := c.DB.Exec(query); err != nil {
			return fmt.Errorf("execute migration query: %w", err)
		}
	}

	return nil
}

func (c *Client) Close() error {
	return c.DB.Close()
}
