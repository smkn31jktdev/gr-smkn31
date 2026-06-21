package util

import (
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"strings"
	"unicode"
)

// URL Sanitizer
var allowedHosts = []string{
	"youtube.com",
	"youtu.be",
	"instagram.com",
	"tiktok.com",
	"facebook.com",
	"fb.com",
	"fb.watch",
	"drive.google.com",
}

// Memvalidasi URL hanya dari platform yang diizinkan
func SanitizeURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("url tidak boleh kosong")
	}

	// Wajib ada scheme https
	if !strings.HasPrefix(raw, "https://") {
		return "", errors.New("url harus menggunakan https")
	}

	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", errors.New("url tidak valid")
	}

	// Hanya izinkan https
	if parsed.Scheme != "https" {
		return "", errors.New("url harus menggunakan https")
	}

	// Normalisasi host: hilangkan port, lowercase
	host := strings.ToLower(parsed.Hostname())

	if !isHostAllowed(host) {
		return "", errors.New("url tidak diizinkan — hanya YouTube, Instagram, TikTok, Facebook, dan Google Drive")
	}

	return raw, nil
}

// Memeriksa host terhadap daftar domain yang diizinkan
func isHostAllowed(host string) bool {
	for _, allowed := range allowedHosts {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}

// File Sanitizer

type allowedFile struct {
	ext   string
	magic [][]byte
}

// Daftar ekstensi + magic bytes yang diizinkan
var allowedFiles = []allowedFile{
	{
		ext: ".jpg",
		magic: [][]byte{
			{0xFF, 0xD8, 0xFF},
		},
	},
	{
		ext: ".jpeg",
		magic: [][]byte{
			{0xFF, 0xD8, 0xFF},
		},
	},
	{
		ext: ".png",
		magic: [][]byte{
			{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		},
	},
	{
		ext: ".heic",
		magic: [][]byte{
			{'f', 't', 'y', 'p'},
		},
	},
}

const maxFileSize = 5 << 20

// Memvalidasi file upload: ekstensi, ukuran, dan magic bytes
func SanitizeUploadedFile(fh *multipart.FileHeader) error {
	// Validasi ukuran
	if fh.Size > maxFileSize {
		return errors.New("ukuran file melebihi batas 5 MB")
	}

	// Validasi ekstensi (lowercase, case-insensitive)
	name := strings.ToLower(strings.TrimSpace(fh.Filename))
	ext, allowed := resolveAllowedExt(name)
	if !allowed {
		return errors.New("tipe file tidak diizinkan — hanya JPG, PNG, dan HEIC")
	}

	// Buka file untuk baca magic bytes
	f, err := fh.Open()
	if err != nil {
		return errors.New("gagal membaca file")
	}
	defer f.Close()

	// Baca 16 byte pertama (cukup untuk semua magic bytes)
	header := make([]byte, 16)
	n, err := io.ReadFull(f, header)
	if err != nil && n == 0 {
		return errors.New("file kosong atau tidak dapat dibaca")
	}
	header = header[:n]

	// Validasi magic bytes sesuai ekstensi
	if !matchMagic(ext, header) {
		return errors.New("konten file tidak sesuai dengan ekstensi — kemungkinan file berbahaya")
	}

	return nil
}

// Mengembalikan ekstensi yang cocok dan status diizinkan
func resolveAllowedExt(filename string) (string, bool) {
	dotIdx := strings.LastIndex(filename, ".")
	if dotIdx < 0 {
		return "", false
	}
	ext := filename[dotIdx:]
	for _, af := range allowedFiles {
		if ext == af.ext {
			return ext, true
		}
	}
	return "", false
}

// Memeriksa header bytes sesuai daftar magic bytes ekstensi
func matchMagic(ext string, header []byte) bool {
	for _, af := range allowedFiles {
		if af.ext != ext {
			continue
		}
		for _, magic := range af.magic {
			if ext == ".heic" {
				// HEIC: "ftyp" berada di offset 4 (bukan offset 0)
				if len(header) >= 8 && matchAt(header, magic, 4) {
					return true
				}
				continue
			}
			if matchAt(header, magic, 0) {
				return true
			}
		}
	}
	return false
}

// Memeriksa apakah slice pattern cocok di posisi offset dalam data
func matchAt(data, pattern []byte, offset int) bool {
	if len(data) < offset+len(pattern) {
		return false
	}
	for i, b := range pattern {
		if data[offset+i] != b {
			return false
		}
	}
	return true
}

// Text Sanitizer

// Membersihkan input teks dari karakter kontrol, null byte, dan spasi berlebih
func SanitizeText(s string) string {
	// Hapus null byte dan karakter kontrol non-printable (kecuali newline/tab yang sah)
	cleaned := strings.Map(func(r rune) rune {
		if r == 0 {
			return -1
		}
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			return -1
		}
		return r
	}, s)

	return strings.TrimSpace(cleaned)
}

// Membersihkan input teks
func SanitizeTextStrict(s string) string {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
	return strings.TrimSpace(cleaned)
}
