package admin

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

// CreateAdmin
func (h *Handler) CreateAdmin(c *gin.Context) {
	var req authmodel.AdminCreateRequest
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

	c.JSON(http.StatusCreated, common.OK(result, "admin berhasil ditambahkan"))
}

// DeleteAdmin
func (h *Handler) DeleteAdmin(c *gin.Context) {
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

// UpdateAdmin
func (h *Handler) UpdateAdmin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Fail("id diperlukan"))
		return
	}

	var req authmodel.AdminUpdateRequest
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

	c.JSON(http.StatusOK, common.OK(result, "admin berhasil diperbarui"))
}


// ListAdmin
func (h *Handler) ListAdmin(c *gin.Context) {
	filter := authmodel.AdminFilter{
		Role:  c.Query("role"),
		Query: c.Query("q"),
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

// BulkImportAdmin
func (h *Handler) BulkImportAdmin(c *gin.Context) {
	var items []authmodel.AdminBulkItem
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

// ImportFromSheets
func (h *Handler) ImportFromSheets(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, common.Fail("fitur import dari Google Sheets belum diimplementasikan"))
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
