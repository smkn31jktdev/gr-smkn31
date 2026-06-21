package naikkelas

import (
	"context"
	"log"
	"net/http"
	"time"

	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Handler HTTP untuk naik kelas
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// TriggerNaikKelas
func (h *Handler) TriggerNaikKelas(c *gin.Context) {
	force := c.Query("force") == "true"
	tahunParam := c.Query("tahun_ajaran")

	now := util.NowJakarta()
	tahunAjaran := tahunParam
	if tahunAjaran == "" {
		tahunAjaran = TahunAjaranSekarang(now)
	}

	hasil, err := h.svc.JalankanNaikKelas(c.Request.Context(), tahunAjaran, force)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Proses naik kelas selesai",
		"data":    hasil,
	})
}

// StatusNaikKelas
func (h *Handler) StatusNaikKelas(c *gin.Context) {
	tahunParam := c.Query("tahun_ajaran")
	now := util.NowJakarta()
	tahunAjaran := tahunParam
	if tahunAjaran == "" {
		tahunAjaran = TahunAjaranSekarang(now)
	}

	sudah, err := h.svc.SudahDieksekusi(c.Request.Context(), tahunAjaran)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	msg := "Naik kelas belum dieksekusi untuk tahun ajaran " + tahunAjaran
	if sudah {
		msg = "Naik kelas sudah dieksekusi untuk tahun ajaran " + tahunAjaran
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"tahunAjaran":   tahunAjaran,
		"sudahEksekusi": sudah,
		"message":       msg,
	})
}

// StartCronNaikKelas
func StartCronNaikKelas(svc *Service) {
	go func() {
		log.Println("[naik-kelas] Cron scheduler dimulai — akan aktif otomatis setiap 1 Juli")

		for {
			now := util.NowJakarta()

			besok := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 1, 0, 0, now.Location())
			time.Sleep(time.Until(besok))

			now = util.NowJakarta()

			if now.Month() != time.July || now.Day() != 1 {
				continue
			}

			tahunAjaran := TahunAjaranSekarang(now)
			log.Printf("[naik-kelas] 1 Juli terdeteksi — memulai naik kelas tahun ajaran %s", tahunAjaran)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			hasil, err := svc.JalankanNaikKelas(ctx, tahunAjaran, false)
			cancel()

			if err != nil {
				log.Printf("[naik-kelas] Cron skip/gagal: %v", err)
			} else {
				log.Printf("[naik-kelas] Cron sukses: %d berhasil, %d gagal dari %d siswa",
					hasil.TotalBerhasil, hasil.TotalGagal, hasil.TotalDiproses)
			}

			time.Sleep(25 * time.Hour)
		}
	}()
}
