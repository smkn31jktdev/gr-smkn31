package student

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	kehadiranmodel "be-gr31/internal/model/kehadiran"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Service defines the required interface for attendance service operations
type Service interface {
	Create(ctx context.Context, siswa kehadiranmodel.SiswaInfo, req kehadiranmodel.AbsenRequest) (*kehadiranmodel.Kehadiran, error)
	List(ctx context.Context, filter kehadiranmodel.KehadiranFilter) ([]kehadiranmodel.Kehadiran, bool, int, error)
}

// RekapService defines the required interface for attendance recap operations
type RekapService interface {
	GetRingkasanSiswa(ctx context.Context, nis, dari, sampai string) ([]rekapmodel.RekapBulanan, error)
	GetRekap(ctx context.Context, nis, bulanTahun string) (*rekapmodel.RekapBulanan, error)
	GetKehadiranBulananSiswa(ctx context.Context, nis, bulanTahun string) (*rekapmodel.KehadiranBulananSiswa, error)
}

// Handler handles student attendance requests
type Handler struct {
	service  Service
	rekapSvc RekapService
}

// NewHandler creates a new student Handler instance
func NewHandler(service Service, rekapSvc RekapService) *Handler {
	return &Handler{service: service, rekapSvc: rekapSvc}
}

// AbsenSiswa handles student daily check-in/out
func (h *Handler) AbsenSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req kehadiranmodel.AbsenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	siswa := kehadiranmodel.SiswaInfo{NIS: claims.NIS}
	result, err := h.service.Create(c.Request.Context(), siswa, req)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "absensi berhasil", "data": result})
}

// RekapBulananSiswa handles retrieving monthly recap for a student
func (h *Handler) RekapBulananSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	dari := c.Query("dari")
	sampai := c.Query("sampai")
	if dari != "" && sampai != "" {
		result, err := h.rekapSvc.GetRingkasanSiswa(c.Request.Context(), claims.NIS, dari, sampai)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
		return
	}

	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}
	result, err := h.rekapSvc.GetRekap(c.Request.Context(), claims.NIS, bulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}
	if result == nil {
		result = &rekapmodel.RekapBulanan{
			NIS:        claims.NIS,
			BulanTahun: bulan,
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// KehadiranBulananSiswa handles retrieving detailed monthly attendance calendar/logs for a student
func (h *Handler) KehadiranBulananSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}

	result, err := h.rekapSvc.GetKehadiranBulananSiswa(c.Request.Context(), claims.NIS, bulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// ListKehadiranSiswa handles listing paginated attendance logs for a student
func (h *Handler) ListKehadiranSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	filter := kehadiranmodel.KehadiranFilter{
		NIS:       claims.NIS,
		Status:    c.Query("status"),
		Tanggal:   c.Query("tanggal"),
		BulanDari: c.Query("dari"),
		BulanKe:   c.Query("sampai"),
		Page:      parseIntQuery(c, "page", 1),
		Limit:     parseIntQuery(c, "limit", util.DefaultPageSize),
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

	// Format Paginated response manually to avoid circular dependencies on common package helper
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"items":   result,
			"page":    filter.Page,
			"limit":   filter.Limit,
			"total":   total,
			"hasMore": hasMore,
		},
	})
}

func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, kehadiranmodel.ErrDuplicate):
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, kehadiranmodel.ErrOutOfTime):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, kehadiranmodel.ErrWeekend):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, kehadiranmodel.ErrTooFar):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, kehadiranmodel.ErrBadRequest):
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	case errors.Is(err, kehadiranmodel.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
	}
}

func parseIntQuery(c *gin.Context, key string, fallback int) int {
	v := c.Query(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func clampLimit(limit *int) {
	if *limit == 0 {
		return
	}
	if *limit > util.MaxAPIPageSize {
		*limit = util.MaxAPIPageSize
	}
}
