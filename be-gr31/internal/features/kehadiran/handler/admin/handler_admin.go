package admin

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"be-gr31/internal/features/kehadiran/monitoring"
	"be-gr31/internal/kalender"
	kehadiranmodel "be-gr31/internal/model/kehadiran"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Service defines the required interface for attendance service operations for admin
type Service interface {
	CreateByAdmin(ctx context.Context, req kehadiranmodel.AdminAbsenRequest, siswaInfo kehadiranmodel.SiswaInfo) (*kehadiranmodel.Kehadiran, error)
	List(ctx context.Context, filter kehadiranmodel.KehadiranFilter) ([]kehadiranmodel.Kehadiran, bool, int, error)
	Delete(ctx context.Context, id string) error
	UpdateByAdmin(ctx context.Context, id, statusBaru, alasan string) (*kehadiranmodel.Kehadiran, error)
}

// RekapService defines the required interface for attendance recap operations for admin
type RekapService interface {
	ListRekap(ctx context.Context, filter rekapmodel.RekapFilter) ([]rekapmodel.RekapBulanan, bool, int, error)
	GetRekapHarian(ctx context.Context, tanggal, kelas string) (*rekapmodel.RekapHarian, error)
	GetRekapSemesterKelas(ctx context.Context, kelas, semester string) (*rekapmodel.RekapSemesterKelas, error)
	GetRingkasanKelas(ctx context.Context, kelas, bulanTahun string, page, limit int) ([]rekapmodel.RingkasanSiswa, bool, int, error)
	GetRekapKelasLengkap(ctx context.Context, bulan, kelas string) (*rekapmodel.RekapKelasLengkap, error)
	GetPersentaseKehadiranKelas(ctx context.Context, bulan, kelas string) (*rekapmodel.RekapPersentaseKelas, error)
	GetKelasJurusan(ctx context.Context) (*rekapmodel.KelasJurusanResponse, error)
	GetRekapMingguanKelas(ctx context.Context, kelas, senin string) (*rekapmodel.RekapMingguanKelas, error)
	GetKelas(ctx context.Context) ([]string, error)
	GetKehadiranBulananSiswa(ctx context.Context, nis, bulanTahun string) (*rekapmodel.KehadiranBulananSiswa, error)
	GetKehadiranMingguanSiswa(ctx context.Context, nis, senin string) (*rekapmodel.KehadiranBulananSiswa, error)
}

// Handler handles admin attendance requests
type Handler struct {
	service  Service
	rekapSvc RekapService
}

// NewHandler creates a new admin Handler instance
func NewHandler(service Service, rekapSvc RekapService) *Handler {
	return &Handler{service: service, rekapSvc: rekapSvc}
}

// AbsenAdmin handles manual attendance creation by admin
func (h *Handler) AbsenAdmin(c *gin.Context) {
	var req kehadiranmodel.AdminAbsenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	siswa := kehadiranmodel.SiswaInfo{NIS: req.NIS, Nama: "", Kelas: ""}
	result, err := h.service.CreateByAdmin(c.Request.Context(), req, siswa)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "kehadiran berhasil ditambahkan", "data": result})
}

// ListKehadiranAdmin handles listing all students' attendance logs with monitoring summary
func (h *Handler) ListKehadiranAdmin(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}

	var adminID, adminRole string
	if claimsVal, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claimsVal.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	filter := kehadiranmodel.KehadiranFilter{
		NIS:       nis,
		Kelas:     c.Query("kelas"),
		Jurusan:   c.Query("jurusan"),
		Status:    c.Query("status"),
		Tanggal:   c.Query("tanggal"),
		BulanDari: c.Query("dari"),
		BulanKe:   c.Query("sampai"),
		Query:     c.Query("q"),
		Page:      parseIntQuery(c, "page", 1),
		Limit:     parseIntQueryAllowZero(c, "limit", util.DefaultPageSize),
		AdminID:   adminID,
		AdminRole: adminRole,
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

	summary := monitoring.BuildSummaryByStudent(result)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"items":            result,
			"page":             filter.Page,
			"limit":            filter.Limit,
			"total":            total,
			"hasMore":          hasMore,
			"summaryByStudent": summary,
		},
	})
}

// DeleteKehadiran handles deleting an attendance record
func (h *Handler) DeleteKehadiran(c *gin.Context) {
	var req kehadiranmodel.DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Delete(c.Request.Context(), req.ID); err != nil {
		if errors.Is(err, kehadiranmodel.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "kehadiran tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateKehadiran handles updating an attendance status by admin
func (h *Handler) UpdateKehadiran(c *gin.Context) {
	var req rekapmodel.UpdateKehadiranRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	result, err := h.service.UpdateByAdmin(c.Request.Context(), req.ID, req.Status, req.Alasan)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "kehadiran berhasil diperbarui", "data": result})
}

// ListRekapBulanan handles listing monthly recaps with pagination
func (h *Handler) ListRekapBulanan(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}
	filter := rekapmodel.RekapFilter{
		NIS:        nis,
		Kelas:      c.Query("kelas"),
		BulanTahun: c.Query("bulan"),
		Semester:   c.Query("semester"),
		Query:      c.Query("q"),
		Page:       parseIntQuery(c, "page", 1),
		Limit:      parseIntQueryAllowZero(c, "limit", util.DefaultPageSize),
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.rekapSvc.ListRekap(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

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

// GetRekapHarian handles retrieving daily recap summary for a class
func (h *Handler) GetRekapHarian(c *gin.Context) {
	tanggal := c.Query("tanggal")
	if tanggal == "" {
		tanggal = util.FormatTanggal(util.NowJakarta())
	}
	kelas := c.Query("kelas")

	result, err := h.rekapSvc.GetRekapHarian(c.Request.Context(), tanggal, kelas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// GetRekapSemesterKelas handles retrieving attendance trends for a semester
func (h *Handler) GetRekapSemesterKelas(c *gin.Context) {
	kelas := c.Query("kelas")
	semester := c.Query("semester")

	result, err := h.rekapSvc.GetRekapSemesterKelas(c.Request.Context(), kelas, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// GetRingkasanKelas handles retrieving class summaries for a month
func (h *Handler) GetRingkasanKelas(c *gin.Context) {
	kelas := c.Query("kelas")
	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", util.DefaultPageSize)
	clampLimit(&limit)

	result, hasMore, total, err := h.rekapSvc.GetRingkasanKelas(c.Request.Context(), kelas, bulan, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ok",
		"data": gin.H{
			"items":   result,
			"page":    page,
			"limit":   limit,
			"total":   total,
			"hasMore": hasMore,
		},
	})
}

// GetRekapLengkap handles retrieving complete monthly class recap with pagination
func (h *Handler) GetRekapLengkap(c *gin.Context) {
	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}
	kelas := c.Query("kelas")

	result, err := h.rekapSvc.GetRekapKelasLengkap(c.Request.Context(), bulan, kelas)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	page := parseIntQuery(c, "page", 1)
	limit := parseIntQueryAllowZero(c, "limit", 0)
	items, hasMore, total := util.PaginateSlice(result.SummaryByStudent, page, limit)
	result.SummaryByStudent = items
	result.TotalSiswa = total
	result.Page = page
	result.Limit = limit
	result.HasMore = hasMore

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// GetPersentaseKelas handles retrieving average attendance percentages per class
func (h *Handler) GetPersentaseKelas(c *gin.Context) {
	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}
	kelas := c.Query("kelas")

	result, err := h.rekapSvc.GetPersentaseKehadiranKelas(c.Request.Context(), bulan, kelas)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// GetKelasJurusan handles retrieving all unique classes and departments
func (h *Handler) GetKelasJurusan(c *gin.Context) {
	result, err := h.rekapSvc.GetKelasJurusan(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// GetRekapMingguan handles retrieving weekly class attendance summary
func (h *Handler) GetRekapMingguan(c *gin.Context) {
	kelas := c.Query("kelas")
	senin := c.Query("senin")

	result, err := h.rekapSvc.GetRekapMingguanKelas(c.Request.Context(), kelas, senin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	page := parseIntQuery(c, "page", 1)
	limit := parseIntQueryAllowZero(c, "limit", 0)
	items, hasMore, total := util.PaginateSlice(result.SummaryByStudent, page, limit)
	result.SummaryByStudent = items
	result.TotalSiswa = total
	result.Page = page
	result.Limit = limit
	result.HasMore = hasMore

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// GetKalender handles retrieving effective school days and holidays for a month
func (h *Handler) GetKalender(c *gin.Context) {
	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}

	resp := gin.H{
		"bulan":       bulan,
		"hariEfektif": kalender.HariEfektif(bulan),
		"daftarLibur": kalender.DaftarLibur(bulan),
	}
	if info, ok := kalender.Info(bulan); ok {
		resp["adaDiKalender"] = true
		resp["hariKerjaTersedia"] = info.HariKerjaTersedia
		resp["kegiatanAkademik"] = info.KegiatanAkademik
	} else {
		resp["adaDiKalender"] = false
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": resp})
}

// GetKelas handles retrieving list of all classes
func (h *Handler) GetKelas(c *gin.Context) {
	result, err := h.rekapSvc.GetKelas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
}

// GetRekapSiswaDetail handles retrieving detailed attendance for a single student (weekly or monthly)
func (h *Handler) GetRekapSiswaDetail(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}
	tipe := c.Query("tipe")
	bulan := c.Query("bulan")
	senin := c.Query("senin")

	if nis == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "nis atau nisn wajib diisi"})
		return
	}

	var result *rekapmodel.KehadiranBulananSiswa
	var err error

	if tipe == "mingguan" {
		result, err = h.rekapSvc.GetKehadiranMingguanSiswa(c.Request.Context(), nis, senin)
	} else {
		if bulan == "" {
			bulan = util.BulanTahun(util.NowJakarta())
		}
		result, err = h.rekapSvc.GetKehadiranBulananSiswa(c.Request.Context(), nis, bulan)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": result})
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

func parseIntQueryAllowZero(c *gin.Context, key string, fallback int) int {
	v := c.Query(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}
