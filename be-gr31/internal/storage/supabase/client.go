package supabase

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Client struct {
	DB *sql.DB
}

// NewClient membuat client Supabase (PostgreSQL) baru dan menjalankan migrasi serta seeder.
func NewClient(dbURL string) (*Client, error) {
	if dbURL == "" || strings.Contains(dbURL, "[YOUR-PASSWORD]") {
		return nil, nil
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open supabase db connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping supabase db: %w", err)
	}

	client := &Client{DB: db}
	if err := client.Migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate supabase db: %w", err)
	}
	return client, nil
}

func (c *Client) Migrate() error {
	queryAdmins := `
	CREATE TABLE IF NOT EXISTS akun_admin (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nama TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_walas BOOLEAN DEFAULT FALSE,
		kelas TEXT DEFAULT '',
		role TEXT DEFAULT 'admin',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
	);
	`
	if _, err := c.DB.Exec(queryAdmins); err != nil {
		return fmt.Errorf("execute migrate query for akun_admin: %w", err)
	}

	queryKelas := `
	CREATE TABLE IF NOT EXISTS kelas (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nama TEXT NOT NULL UNIQUE,
		walas_id UUID REFERENCES akun_admin(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
	);
	`
	if _, err := c.DB.Exec(queryKelas); err != nil {
		return fmt.Errorf("execute migrate query for kelas: %w", err)
	}

	// Seed default classes
	defaultClasses := []string{
		"X Akuntansi", "X Animasi", "X Bisnis Ritel", "X DKV", "X Layanan Perbankan", "X Manajemen Perkantoran",
		"XI Akuntansi", "XI Animasi", "XI Bisnis Ritel", "XI DKV", "XI Layanan Perbankan", "XI Manajemen Perkantoran",
		"XII Akuntansi", "XII Animasi", "XII Bisnis Ritel", "XII DKV", "XII Layanan Perbankan", "XII Manajemen Perkantoran",
	}
	for _, className := range defaultClasses {
		_, err := c.DB.Exec("INSERT INTO kelas (nama) VALUES ($1) ON CONFLICT (nama) DO NOTHING", className)
		if err != nil {
			return fmt.Errorf("seed default class %s: %w", className, err)
		}
	}

	// Ensure 'kelas' and 'role' columns exist in 'akun_admin' table
	if _, err := c.DB.Exec(`
		ALTER TABLE akun_admin 
		ADD COLUMN IF NOT EXISTS kelas TEXT DEFAULT '',
		ADD COLUMN IF NOT EXISTS role TEXT DEFAULT 'admin'
	`); err != nil {
		fmt.Printf("  [Supabase Migration] Warning: failed to alter 'akun_admin' table: %v\n", err)
	}

	// Ensure smkn31jktdev@gmail.com has super_admin role
	if _, err := c.DB.Exec(`
		UPDATE akun_admin 
		SET role = 'super_admin' 
		WHERE email = 'smkn31jktdev@gmail.com'
	`); err != nil {
		fmt.Printf("  [Supabase Migration] Warning: failed to update super_admin role: %v\n", err)
	}

	queryStudents := `
	CREATE TABLE IF NOT EXISTS akun_siswa (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		nis TEXT NOT NULL UNIQUE,
		nama TEXT NOT NULL,
		kelas TEXT NOT NULL,
		wali_id UUID REFERENCES akun_admin(id) ON DELETE SET NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
	);
	`
	if _, err := c.DB.Exec(queryStudents); err != nil {
		return fmt.Errorf("execute migrate query for akun_siswa: %w", err)
	}

	// Rename column in Supabase if nisn still exists
	var exists bool
	err := c.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name='akun_siswa' AND column_name='nisn'
		)
	`).Scan(&exists)
	if err == nil && exists {
		_, err = c.DB.Exec("ALTER TABLE akun_siswa RENAME COLUMN nisn TO nis")
		if err != nil {
			fmt.Printf("  [Supabase Migration] Warning: failed to rename column nisn to nis: %v\n", err)
		} else {
			fmt.Println("  [Supabase Migration] Success: Renamed column nisn to nis in akun_siswa table")
		}
	}

	return nil
}

type adminRecord struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type studentRecord struct {
	ID    string `json:"id"`
	NIS   string `json:"nisn"`
	Nama  string `json:"nama"`
	Kelas string `json:"kelas"`
	Wali  string `json:"wali"`
}

// SeedAdmins membaca file JSON seeder, mengenkripsi password dengan bcrypt, dan menyimpannya ke database Supabase.
func (c *Client) SeedAdmins(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		// Jika file seeder tidak ditemukan, lewati dengan warning di log
		fmt.Printf("  [Supabase Seeder] Warning: %s tidak ditemukan, skip seeding admins\n", filePath)
		return nil
	}

	var records []adminRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("unmarshal admin seeder json: %w", err)
	}

	// Mulai transaksi database agar proses sekuensial cepat dan aman
	tx, err := c.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO akun_admin (nama, email, password, is_walas)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer stmt.Close()

	seededCount := 0
	for _, rec := range records {
		if rec.Email == "" || rec.Password == "" {
			continue
		}

		// Enkripsi password menggunakan bcrypt
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(rec.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("bcrypt hash password for %s: %w", rec.Email, err)
		}

		res, err := stmt.Exec(rec.Nama, rec.Email, string(hashedPass), false)
		if err != nil {
			return fmt.Errorf("exec insert statement for %s: %w", rec.Email, err)
		}

		if rows, err := res.RowsAffected(); err == nil && rows > 0 {
			seededCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if seededCount > 0 {
		fmt.Printf("  [Supabase Seeder] Success: Berhasil memasukkan %d data admin baru ke akun_admin\n", seededCount)
	} else {
		fmt.Println("  [Supabase Seeder] Info: Tidak ada data admin baru yang dimasukkan (seluruh email sudah terdaftar)")
	}

	return nil
}

// SeedStudents membaca file JSON seeder siswa, memetakan wali ke akun_admin, mengenkripsi password dengan bcrypt, dan menyimpannya ke database Supabase.
func (c *Client) SeedStudents(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("  [Supabase Seeder] Warning: %s tidak ditemukan, skip seeding siswa\n", filePath)
		return nil
	}

	var records []studentRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("unmarshal student seeder json: %w", err)
	}

	// Ambil semua admin dari database untuk pemetaan wali
	rows, err := c.DB.Query("SELECT id, nama FROM akun_admin")
	if err != nil {
		return fmt.Errorf("query admins for mapping: %w", err)
	}
	defer rows.Close()

	adminMap := make(map[string]string)
	for rows.Next() {
		var id, nama string
		if err := rows.Scan(&id, &nama); err != nil {
			return fmt.Errorf("scan admin row: %w", err)
		}
		adminMap[strings.ToLower(strings.TrimSpace(nama))] = id
	}

	findWaliID := func(waliName string) (string, bool) {
		if waliName == "" {
			return "", false
		}
		normalized := strings.ToLower(strings.TrimSpace(waliName))
		// 1. Direct match
		if id, ok := adminMap[normalized]; ok {
			return id, true
		}

		// 2. Custom typo/spelling fallbacks
		var corrected string
		switch normalized {
		case "meythannisa salsabila":
			corrected = "meythannisa salsabilla"
		case "sriyani":
			corrected = "sriyani spd"
		case "syafira khairuninisa":
			corrected = "syafira khairunnisa"
		case "melki felix m":
			corrected = "melki felix mapan"
		}
		if corrected != "" {
			if id, ok := adminMap[corrected]; ok {
				return id, true
			}
		}

		// 3. Prefix/substring fallback match
		for adminName, id := range adminMap {
			if len(normalized) >= 4 && strings.HasPrefix(adminName, normalized) {
				return id, true
			}
			if len(adminName) >= 4 && strings.HasPrefix(normalized, adminName) {
				return id, true
			}
		}

		return "", false
	}

	tx, err := c.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction for students: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO akun_siswa (id, nis, nama, kelas, wali_id, password)
		VALUES (COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6)
		ON CONFLICT (nis) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare student statement: %w", err)
	}
	defer stmt.Close()

	seededCount := 0
	unmappedWalis := make(map[string]bool)

	for _, rec := range records {
		if rec.NIS == "" {
			continue
		}

		// Cari wali_id
		var waliID sql.NullString
		if id, ok := findWaliID(rec.Wali); ok {
			waliID = sql.NullString{String: id, Valid: true}
		} else {
			waliID = sql.NullString{Valid: false}
			if rec.Wali != "" {
				unmappedWalis[rec.Wali] = true
			}
		}

		// Enkripsi password menggunakan bcrypt (default password: "123456")
		hashedPass, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("bcrypt hash password for student %s: %w", rec.NIS, err)
		}

		res, err := stmt.Exec(rec.ID, rec.NIS, rec.Nama, rec.Kelas, waliID, string(hashedPass))
		if err != nil {
			return fmt.Errorf("exec insert student statement for %s: %w", rec.NIS, err)
		}

		if rows, err := res.RowsAffected(); err == nil && rows > 0 {
			seededCount++
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit student transaction: %w", err)
	}

	if len(unmappedWalis) > 0 {
		var names []string
		for name := range unmappedWalis {
			names = append(names, name)
		}
		fmt.Printf("  [Supabase Seeder] Warning: Wali berikut tidak ditemukan di akun_admin: %s\n", strings.Join(names, ", "))
	}

	if seededCount > 0 {
		fmt.Printf("  [Supabase Seeder] Success: Berhasil memasukkan %d data siswa baru ke akun_siswa\n", seededCount)
	} else {
		fmt.Println("  [Supabase Seeder] Info: Tidak ada data siswa baru yang dimasukkan (seluruh NIS sudah terdaftar)")
	}

	return nil
}

func (c *Client) ResetAllStudentPasswords(plainTextPassword string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = c.DB.Exec("UPDATE akun_siswa SET password = $1", string(hashedPass))
	if err != nil {
		return err
	}
	fmt.Printf("  [Supabase Admin] Success: Reset seluruh password siswa ke default '%s'\n", plainTextPassword)
	return nil
}

func (c *Client) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
