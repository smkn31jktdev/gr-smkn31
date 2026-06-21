package bukti

import (
	"net/http"
	"strconv"

	"be-gr31/internal/model/common"
	sekolahmodel "be-gr31/internal/model/sekolah"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Handler bukti
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateBukti
func (h *Handler) CreateBukti(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req sekolahmodel.BuktiCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Create(c.Request.Context(), claims.NIS, claims.Email, "", req)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, common.OK(result, "bukti berhasil disimpan"))
}

// ListBuktiSiswa
func (h *Handler) ListBuktiSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	filter := sekolahmodel.BuktiFilter{
		NIS:   claims.NIS,
		Bulan: c.Query("bulan"),
		Page:  parseIntQuery(c, "page", 1),
		Limit: parseIntQuery(c, "limit", util.DefaultPageSize),
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.Paginated(result, filter.Page, filter.Limit, total, hasMore, "ok"))
}

// ListBuktiAdmin
func (h *Handler) ListBuktiAdmin(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}

	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	filter := sekolahmodel.BuktiFilter{
		NIS:       nis,
		Kelas:     c.Query("kelas"),
		Bulan:     c.Query("bulan"),
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
