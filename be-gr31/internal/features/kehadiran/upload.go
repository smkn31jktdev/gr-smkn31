package kehadiran

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"be-gr31/internal/config"
	"be-gr31/internal/model/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// validMimeTypes adalah MIME type yang diizinkan untuk upload.
var validMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/jpg":  ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/heic": ".heic",
	"image/heif": ".heic",
}

// HandleUploadIzin menangani upload foto izin/sakit.
// Endpoint: POST /v1/student/kehadiran/upload-izin
func HandleUploadIzin(cfg *config.Config, nis string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Batasi ukuran request
		maxBytes := cfg.UploadMaxSizeMB * 1024 * 1024
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		if err := c.Request.ParseMultipartForm(maxBytes); err != nil {
			c.JSON(http.StatusRequestEntityTooLarge, common.Fail(
				fmt.Sprintf("ukuran file melebihi batas %dMB", cfg.UploadMaxSizeMB),
			))
			return
		}

		file, header, err := c.Request.FormFile("foto")
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Fail("file 'foto' tidak ditemukan"))
			return
		}
		defer file.Close()

		// Validasi MIME type (baca 512 byte pertama)
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		mimeType := http.DetectContentType(buf[:n])
		ext, ok := validMimeTypes[mimeType]
		if !ok && mimeType == "application/octet-stream" {
			origExt := strings.ToLower(filepath.Ext(header.Filename))
			if origExt == ".heic" || origExt == ".heif" {
				ext = ".heic"
				ok = true
			}
		}
		if !ok {
			c.JSON(http.StatusBadRequest, common.Fail("tipe file tidak valid, hanya jpeg/png/webp/heic"))
			return
		}

		// Reset ke awal file setelah deteksi MIME
		if _, err := file.Seek(0, 0); err != nil {
			c.JSON(http.StatusInternalServerError, common.Fail("gagal memproses file"))
			return
		}

		// Validasi ukuran header
		if header.Size > maxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, common.Fail(
				fmt.Sprintf("ukuran file melebihi batas %dMB", cfg.UploadMaxSizeMB),
			))
			return
		}

		// Sanitize NIS dan generate nama file unik (jangan pakai nama asli dari client)
		safeNIS := sanitizePathComponent(nis)
		fileName := uuid.New().String() + ext
		dir := filepath.Join(cfg.UploadDir, "izin", safeNIS)

		if err := os.MkdirAll(dir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, common.Fail("gagal membuat direktori upload"))
			return
		}

		destPath := filepath.Join(dir, fileName)
		destFile, err := os.Create(destPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Fail("gagal menyimpan file"))
			return
		}
		defer destFile.Close()

		// Copy isi file
		written := int64(0)
		tmp := make([]byte, 32*1024)
		for {
			nr, er := file.Read(tmp)
			if nr > 0 {
				nw, ew := destFile.Write(tmp[:nr])
				written += int64(nw)
				if ew != nil {
					c.JSON(http.StatusInternalServerError, common.Fail("gagal menulis file"))
					return
				}
				if written > maxBytes {
					os.Remove(destPath)
					c.JSON(http.StatusRequestEntityTooLarge, common.Fail("file terlalu besar"))
					return
				}
			}
			if er != nil {
				break
			}
		}

		// Kembalikan path relatif
		relativePath := fmt.Sprintf("/uploads/izin/%s/%s", safeNIS, fileName)
		c.JSON(http.StatusOK, common.OK(gin.H{"url": relativePath}, "upload berhasil"))
	}
}

// sanitizePathComponent membersihkan string dari karakter berbahaya untuk path.
func sanitizePathComponent(s string) string {
	s = strings.ReplaceAll(s, "/", "")
	s = strings.ReplaceAll(s, "\\", "")
	s = strings.ReplaceAll(s, "..", "")
	s = strings.ReplaceAll(s, ":", "")
	return strings.TrimSpace(s)
}
