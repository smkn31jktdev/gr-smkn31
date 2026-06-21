package kegiatan

import (
	"errors"
	"net/http"
	"strconv"

	"be-gr31/internal/model/common"
	sekolahmodel "be-gr31/internal/model/sekolah"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Handler menangani HTTP request untuk kegiatan.
type Handler struct {
	service *Service
}

// NewHandler membuat instance Handler baru.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateKegiatan menangani POST /v1/student/kegiatan.
func (h *Handler) CreateKegiatan(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req sekolahmodel.KegiatanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Create(c.Request.Context(), claims.NIS, claims.Email, "", req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, common.OK(result, "kegiatan berhasil ditambahkan"))
}

// UpdateKegiatan menangani PUT /v1/student/kegiatan.
func (h *Handler) UpdateKegiatan(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req sekolahmodel.KegiatanUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Update(c.Request.Context(), claims.NIS, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail(err.Error()))
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "kegiatan berhasil diperbarui"))
}

// ListKegiatanSiswa menangani GET /v1/student/kegiatan.
func (h *Handler) ListKegiatanSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	filter := sekolahmodel.KegiatanFilter{
		NISN:      claims.NIS,
		Section:   c.Query("section"),
		Tanggal:   c.Query("tanggal"),
		DariTgl:   c.Query("dari"),
		SampaiTgl: c.Query("sampai"),
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

// DeleteKegiatanSiswa menangani DELETE /v1/student/kegiatan.
func (h *Handler) DeleteKegiatanSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Fail("id diperlukan"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), claims.NIS, id); err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail(err.Error()))
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListKegiatanAdmin menangani GET /v1/admin/kegiatan.
func (h *Handler) ListKegiatanAdmin(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}
	filter := sekolahmodel.KegiatanFilter{
		NISN:      nis,
		Kelas:     c.Query("kelas"),
		Section:   c.Query("section"),
		Tanggal:   c.Query("tanggal"),
		DariTgl:   c.Query("dari"),
		SampaiTgl: c.Query("sampai"),
		Query:     c.Query("q"),
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
