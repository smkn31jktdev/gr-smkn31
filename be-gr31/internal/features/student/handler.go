package student

import (
	"errors"
	"net/http"
	"strconv"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	"be-gr31/internal/util"

	"github.com/gin-gonic/gin"
)

// Handler
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateSiswa
func (h *Handler) CreateSiswa(c *gin.Context) {
	var req authmodel.SiswaBulkItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			c.JSON(http.StatusConflict, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusCreated, common.OK(result, "siswa berhasil ditambahkan"))
}

// UpdateSiswa
func (h *Handler) UpdateSiswa(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Fail("id diperlukan"))
		return
	}

	var req authmodel.SiswaUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "data siswa berhasil diperbarui"))
}

// DeleteSiswa
func (h *Handler) DeleteSiswa(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Fail("id diperlukan"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListSiswa
func (h *Handler) ListSiswa(c *gin.Context) {
	var adminID, adminRole string
	if claims, exists := c.Get("claims"); exists {
		if jwtClaims, ok := claims.(*util.JWTClaims); ok {
			adminID = jwtClaims.ID
			adminRole = jwtClaims.Role
		}
	}

	limitStr := c.Query("limit")
	var limit int
	if limitStr == "all" || limitStr == "-1" {
		limit = -1
	} else {
		limit = parseIntQuery(c, "limit", util.DefaultPageSize)
		clampLimit(&limit)
	}

	filter := authmodel.SiswaFilter{
		Kelas:     c.Query("kelas"),
		Query:     c.Query("q"),
		Page:      parseIntQuery(c, "page", 1),
		Limit:     limit,
		AdminID:   adminID,
		AdminRole: adminRole,
	}

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.Paginated(result, filter.Page, filter.Limit, total, hasMore, "ok"))
}

// BulkImportSiswa
func (h *Handler) BulkImportSiswa(c *gin.Context) {
	var items []authmodel.SiswaBulkItem
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, common.Fail("data tidak boleh kosong"))
		return
	}

	success, errs, _ := h.service.BulkCreate(c.Request.Context(), items)
	c.JSON(http.StatusOK, common.OK(gin.H{
		"berhasil": success,
		"gagal":    len(errs),
		"errors":   errs,
	}, "import selesai"))
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

// ImportFromSheets
func (h *Handler) ImportFromSheets(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, common.Fail("fitur import dari Google Sheets belum diimplementasikan"))
}
