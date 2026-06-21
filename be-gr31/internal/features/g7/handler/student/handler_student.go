package student

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"be-gr31/internal/features/g7/content/jurnal"
	"be-gr31/internal/features/g7/content/rekap"
	"be-gr31/internal/model/common"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Service defines the interface required by the student handler
type Service interface {
	DashboardSiswa(ctx context.Context, nisn string) (*g7model.G7DashboardSiswa, error)
	Upsert(ctx context.Context, nisn string, req g7model.G7UpsertRequest) (*g7model.G7, error)
	List(ctx context.Context, filter g7model.G7Filter) ([]g7model.G7, bool, int, error)
	GetByTanggal(ctx context.Context, nisn, tanggal string) (*g7model.G7, error)
}

// Handler G7 Student
type Handler struct {
	service Service
}

// NewHandler creates a new student handler instance
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// DashboardSiswa handles retrieving the dashboard data for a student
func (h *Handler) DashboardSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	result, err := h.service.DashboardSiswa(c.Request.Context(), claims.NIS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

// UpsertG7Siswa handles updating/inserting a student's daily G7 log
func (h *Handler) UpsertG7Siswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req g7model.G7UpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Upsert(c.Request.Context(), claims.NIS, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "G7 berhasil disimpan"))
}

// ListG7Siswa handles listing G7 records for a student
func (h *Handler) ListG7Siswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	filter := g7model.G7Filter{
		NISN:      claims.NIS,
		BulanDari: c.Query("dari"),
		BulanKe:   c.Query("sampai"),
		Page:      parseIntQuery(c, "page", 1),
		Limit:     parseIntQuery(c, "limit", util.DefaultPageSize),
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.Paginated(result, filter.Page, filter.Limit, total, hasMore, "ok"))
}

// GetG7ByTanggal handles retrieving a student's daily G7 log by date
func (h *Handler) GetG7ByTanggal(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)
	tanggal := c.Param("tanggal")

	result, err := h.service.GetByTanggal(c.Request.Context(), claims.NIS, tanggal)
	if err != nil {
		if errors.Is(err, jurnal.ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail("data G7 tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, jurnal.ErrNotFound):
		c.JSON(http.StatusNotFound, common.Fail(err.Error()))
	case errors.Is(err, rekap.ErrSiswaNotFound):
		c.JSON(http.StatusNotFound, common.Fail(err.Error()))
	case errors.Is(err, jurnal.ErrFutureDate):
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
	case errors.Is(err, jurnal.ErrNotTodayDate):
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
