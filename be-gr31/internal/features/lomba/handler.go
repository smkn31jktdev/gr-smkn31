package lomba

import (
	"errors"
	"net/http"
	"strconv"

	"be-gr31/internal/model/common"
	lombamodel "be-gr31/internal/model/lomba"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateLomba menangani POST /v1/student/lomba.
func (h *Handler) CreateLomba(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req lombamodel.KebersihanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Create(c.Request.Context(), claims.NIS, req)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			c.JSON(http.StatusConflict, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, common.OK(result, "data kebersihan berhasil ditambahkan"))
}

// UpdateLomba menangani PUT /v1/student/lomba.
func (h *Handler) UpdateLomba(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req lombamodel.KebersihanUpdateRequest
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

	c.JSON(http.StatusOK, common.OK(result, "data kebersihan berhasil diperbarui"))
}

// ListLombaSiswa menangani GET /v1/student/lomba.
func (h *Handler) ListLombaSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	student, err := h.service.GetStudentInfo(c.Request.Context(), claims.NIS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("gagal mengambil data siswa"))
		return
	}
	if student == nil {
		c.JSON(http.StatusNotFound, common.Fail("siswa tidak ditemukan"))
		return
	}

	filter := lombamodel.KebersihanFilter{
		Kelas:     student.Kelas, // Class-scoped: see all uploads for this class
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

// DeleteLombaSiswa menangani DELETE /v1/student/lomba.
func (h *Handler) DeleteLombaSiswa(c *gin.Context) {
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

// ListLombaAdmin menangani GET /v1/admin/lomba.
func (h *Handler) ListLombaAdmin(c *gin.Context) {
	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	filter := lombamodel.KebersihanFilter{
		Kelas:     c.Query("kelas"),
		Tanggal:   c.Query("tanggal"),
		DariTgl:   c.Query("dari"),
		SampaiTgl: c.Query("sampai"),
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
