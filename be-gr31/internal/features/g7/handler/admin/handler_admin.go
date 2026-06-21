package admin

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"be-gr31/internal/features/g7/content/evaluate"
	"be-gr31/internal/features/g7/content/jurnal"
	"be-gr31/internal/features/g7/content/rekap"
	"be-gr31/internal/model/common"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Service defines the interface required by the admin handler
type Service interface {
	List(ctx context.Context, filter g7model.G7Filter) ([]g7model.G7, bool, int, error)
	Summary(ctx context.Context, bulanTahun, kelas, adminID, adminRole string) ([]g7model.G7Summary, error)
	Delete(ctx context.Context, id string) error
	VerifyAdminAccessToStudent(ctx context.Context, adminID, adminRole, studentNIS string) error
	UpsertRekap(ctx context.Context, req g7model.G7RekapUpsertRequest) (*g7model.G7Rekap, error)
	ListRekap(ctx context.Context, filter g7model.G7RekapFilter) ([]g7model.G7Rekap, bool, int, error)
	GetRekap(ctx context.Context, nisn, bulan string) (*g7model.G7Rekap, error)
	Statistik(ctx context.Context, bulan, kelas, adminID, adminRole string) (*g7model.G7RekapStatistik, error)
	RekapKelasLengkap(ctx context.Context, bulan, kelas, adminID, adminRole string) (*g7model.G7RekapKelasLengkap, error)
	Suggest(ctx context.Context, nisn, bulan string) (*g7model.G7SuggestResponse, error)
	EvaluateJurnalBulanan(ctx context.Context, nisn, bulan string) (*evaluate.EvalReport, error)
	BuildLaporanPDF(ctx context.Context, nisn, bulan string) (*g7model.PDFLaporan, string, error)
	DeleteRekap(ctx context.Context, id string) error
	GetRekapSemesterKelas(ctx context.Context, semester, kelas, adminID, adminRole string) ([]g7model.G7SemesterStudentItem, error)
	FindG7ByID(ctx context.Context, id string) (*g7model.G7, error)
	FindRekapByID(ctx context.Context, id string) (*g7model.G7Rekap, error)
}

// Handler G7 Admin
type Handler struct {
	service Service
}

// NewHandler creates a new admin handler instance
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListG7Admin(c *gin.Context) {
	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}
	filter := g7model.G7Filter{
		NISN:      nis,
		Kelas:     c.Query("kelas"),
		Tanggal:   c.Query("tanggal"),
		BulanDari: c.Query("dari"),
		BulanKe:   c.Query("sampai"),
		Query:     c.Query("q"),
		Page:      parseIntQuery(c, "page", 1),
		Limit:     parseIntQuery(c, "limit", util.DefaultPageSize),
		AdminID:   adminID,
		AdminRole: adminRole,
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.Paginated(result, filter.Page, filter.Limit, total, hasMore, "ok"))
}

func (h *Handler) SummaryG7Admin(c *gin.Context) {
	bulanTahun := c.Query("bulan")
	kelas := c.Query("kelas")

	if bulanTahun == "" {
		bulanTahun = util.BulanTahun(util.NowJakarta())
	}

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	result, err := h.service.Summary(c.Request.Context(), bulanTahun, kelas, adminID, adminRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

func (h *Handler) DeleteG7Admin(c *gin.Context) {
	var req g7model.DeleteG7Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	g7Doc, err := h.service.FindG7ByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}
	if g7Doc == nil {
		c.JSON(http.StatusNotFound, common.Fail("data G7 tidak ditemukan"))
		return
	}

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	if err := h.service.VerifyAdminAccessToStudent(c.Request.Context(), adminID, adminRole, g7Doc.NISN); err != nil {
		c.JSON(http.StatusForbidden, common.Fail(err.Error()))
		return
	}

	if err := h.service.Delete(c.Request.Context(), req.ID); err != nil {
		if errors.Is(err, jurnal.ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail("data G7 tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
	
func (h *Handler) UpsertRekap(c *gin.Context) {
	var req g7model.G7RekapUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	if err := h.service.VerifyAdminAccessToStudent(c.Request.Context(), adminID, adminRole, req.NISN); err != nil {
		c.JSON(http.StatusForbidden, common.Fail(err.Error()))
		return
	}

	result, err := h.service.UpsertRekap(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "rekap G7 berhasil disimpan"))
}

func (h *Handler) ListRekap(c *gin.Context) {
	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}
	filter := g7model.G7RekapFilter{
		NISN:       nis,
		Kelas:      c.Query("kelas"),
		BulanTahun: c.Query("bulan"),
		Predikat:   c.Query("predikat"),
		Status:     c.Query("status"),
		Query:      c.Query("q"),
		Page:       parseIntQuery(c, "page", 1),
		Limit:      parseIntQuery(c, "limit", util.DefaultPageSize),
		AdminID:    adminID,
		AdminRole:  adminRole,
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.ListRekap(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.Paginated(result, filter.Page, filter.Limit, total, hasMore, "ok"))
}

// Hander Recap
func (h *Handler) GetRekapDetail(c *gin.Context) {
	nis := c.Param("nis")
	if nis == "" {
		nis = c.Param("nisn")
	}
	bulan := c.Param("bulan")

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	if err := h.service.VerifyAdminAccessToStudent(c.Request.Context(), adminID, adminRole, nis); err != nil {
		c.JSON(http.StatusForbidden, common.Fail(err.Error()))
		return
	}

	result, err := h.service.GetRekap(c.Request.Context(), nis, bulan)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

// StatistikRekap handles retrieving G7 statistics for admin
func (h *Handler) StatistikRekap(c *gin.Context) {
	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}
	kelas := c.Query("kelas")

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	result, err := h.service.Statistik(c.Request.Context(), bulan, kelas, adminID, adminRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

func (h *Handler) RekapKelas(c *gin.Context) {
	bulan := c.Query("bulan")
	if bulan == "" {
		bulan = util.BulanTahun(util.NowJakarta())
	}
	kelas := c.Query("kelas")

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	result, err := h.service.RekapKelasLengkap(c.Request.Context(), bulan, kelas, adminID, adminRole)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

func (h *Handler) SuggestSkor(c *gin.Context) {
	nis := c.Param("nis")
	if nis == "" {
		nis = c.Param("nisn")
	}
	bulan := c.Param("bulan")

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	if err := h.service.VerifyAdminAccessToStudent(c.Request.Context(), adminID, adminRole, nis); err != nil {
		c.JSON(http.StatusForbidden, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Suggest(c.Request.Context(), nis, bulan)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

func (h *Handler) EvaluateJurnal(c *gin.Context) {
	nis := c.Param("nis")
	if nis == "" {
		nis = c.Param("nisn")
	}
	bulan := c.Param("bulan")

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	if err := h.service.VerifyAdminAccessToStudent(c.Request.Context(), adminID, adminRole, nis); err != nil {
		c.JSON(http.StatusForbidden, common.Fail(err.Error()))
		return
	}

	report, err := h.service.EvaluateJurnalBulanan(c.Request.Context(), nis, bulan)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.OK(report, "ok"))
}

// ExportPDF handles exporting monthly G7 report to HTML format
func (h *Handler) ExportPDF(c *gin.Context) {
	nis := c.Param("nis")
	if nis == "" {
		nis = c.Param("nisn")
	}
	bulan := c.Param("bulan")

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	if err := h.service.VerifyAdminAccessToStudent(c.Request.Context(), adminID, adminRole, nis); err != nil {
		c.JSON(http.StatusForbidden, common.Fail(err.Error()))
		return
	}

	_, html, err := h.service.BuildLaporanPDF(c.Request.Context(), nis, bulan)
	if err != nil {
		if errors.Is(err, rekap.ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail("rekap G7 tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Content-Disposition", "inline; filename=laporan-g7-"+nis+"-"+bulan+".html")
	c.String(http.StatusOK, html)
}

// DeleteRekap handles deleting G7 monthly recap
func (h *Handler) DeleteRekap(c *gin.Context) {
	var req g7model.DeleteG7RekapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	rekapDoc, err := h.service.FindRekapByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}
	if rekapDoc == nil {
		c.JSON(http.StatusNotFound, common.Fail("rekap tidak ditemukan"))
		return
	}

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	if err := h.service.VerifyAdminAccessToStudent(c.Request.Context(), adminID, adminRole, rekapDoc.NISN); err != nil {
		c.JSON(http.StatusForbidden, common.Fail(err.Error()))
		return
	}

	if err := h.service.DeleteRekap(c.Request.Context(), req.ID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetRekapSemesterKelas handles GET /v1/admin/g7/rekap-semester
func (h *Handler) GetRekapSemesterKelas(c *gin.Context) {
	semester := c.Query("semester")
	kelas := c.Query("kelas")

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	result, err := h.service.GetRekapSemesterKelas(c.Request.Context(), semester, kelas, adminID, adminRole)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, rekap.ErrNotFound):
		c.JSON(http.StatusNotFound, common.Fail(err.Error()))
	case errors.Is(err, rekap.ErrSiswaNotFound):
		c.JSON(http.StatusNotFound, common.Fail(err.Error()))
	case errors.Is(err, jurnal.ErrFutureDate):
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
	case errors.Is(err, g7model.ErrBadRequest):
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
	case errors.Is(err, g7model.ErrSkorRange):
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
	case errors.Is(err, rekap.ErrInvalidStatus):
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
	case errors.Is(err, rekap.ErrRekapFinal):
		c.JSON(http.StatusConflict, common.Fail(err.Error()))
	case errors.Is(err, rekap.ErrFinalizeRule):
		c.JSON(http.StatusUnprocessableEntity, common.Fail(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
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
	if *limit > util.MaxAPIPageSize {
		*limit = util.MaxAPIPageSize
	}
}
