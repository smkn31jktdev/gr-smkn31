package aduan

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	aduanmodel "be-gr31/internal/model/aduan"
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

// CreateAduan
func (h *Handler) CreateAduan(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req aduanmodel.AduanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Create(c.Request.Context(), claims.NIS, req.Isi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusCreated, common.OK(result, "aduan berhasil dibuat"))
}

// ListAduanSiswa
func (h *Handler) ListAduanSiswa(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	filter := aduanmodel.AduanFilter{
		NISN:   claims.NIS,
		Status: c.Query("status"),
		Page:   parseIntQuery(c, "page", 1),
		Limit:  parseIntQuery(c, "limit", util.DefaultPageSize),
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.Paginated(result, filter.Page, filter.Limit, total, hasMore, "ok"))
}

// ListAduanAdmin
func (h *Handler) ListAduanAdmin(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}
	filter := aduanmodel.AduanFilter{
		NISN:   nis,
		Status: c.Query("status"),
		Page:   parseIntQuery(c, "page", 1),
		Limit:  parseIntQuery(c, "limit", util.DefaultPageSize),
	}
	clampLimit(&filter.Limit)

	result, hasMore, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.Paginated(result, filter.Page, filter.Limit, total, hasMore, "ok"))
}

// GetRoomAdmin
func (h *Handler) GetRoomAdmin(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Fail("id aduan diperlukan"))
		return
	}

	result, err := h.service.GetRoom(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "ok"))
}

// UpdateStatusAdmin
func (h *Handler) UpdateStatusAdmin(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req aduanmodel.AduanStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.UpdateStatus(c.Request.Context(), req.AduanID, req.Status, claims.Email, claims.Role)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "status aduan diperbarui"))
}

// RespondAdmin
func (h *Handler) RespondAdmin(c *gin.Context) {
	claims := c.MustGet("claims").(*util.JWTClaims)

	var req aduanmodel.AduanRespondRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	result, err := h.service.Respond(c.Request.Context(), claims.Email, req.AduanID, req.Isi)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.JSON(http.StatusOK, common.OK(result, "balasan berhasil dikirim"))
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

// ExportCSVAduan handles exporting all complaints matching filters to CSV format
func (h *Handler) ExportCSVAduan(c *gin.Context) {
	nis := c.Query("nis")
	if nis == "" {
		nis = c.Query("nisn")
	}
	filter := aduanmodel.AduanFilter{
		NISN:   nis,
		Status: c.Query("status"),
	}

	all, err := h.service.ListAll(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=arsip-aduan-siswa.csv")

	// UTF-8 BOM for Excel
	_, _ = c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Writer)
	header := []string{
		"ID Tiket",
		"NISN",
		"Nama Siswa",
		"Kelas",
		"Guru Wali",
		"Status",
		"Tanggal Pengaduan",
		"Terakhir Diperbarui",
		"Ditangani Oleh",
		"Pesan Pengaduan Pertama",
		"Jumlah Pesan",
	}
	if err := writer.Write(header); err != nil {
		return
	}

	for _, item := range all {
		firstMsg := ""
		if len(item.Messages) > 0 {
			firstMsg = item.Messages[0].Isi
		}

		statusLabel := item.Status
		switch item.Status {
		case "open":
			statusLabel = "Baru"
		case "in_progress":
			statusLabel = "Diproses"
		case "closed":
			statusLabel = "Selesai"
		case "pending":
			statusLabel = "Tertunda"
		}

		row := []string{
			item.ID,
			item.NISN,
			item.NamaSiswa,
			item.Kelas,
			item.Walas,
			statusLabel,
			item.CreatedAt,
			item.UpdatedAt,
			item.AdminNama,
			firstMsg,
			strconv.Itoa(len(item.Messages)),
		}
		if err := writer.Write(row); err != nil {
			return
		}
	}
	writer.Flush()
}

// ExportHTMLAduan
func (h *Handler) ExportHTMLAduan(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		id = c.Query("aduanId")
	}
	if id == "" {
		c.JSON(http.StatusBadRequest, common.Fail("id aduan diperlukan"))
		return
	}

	aduan, err := h.service.GetRoom(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, common.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Fail("internal server error"))
		return
	}

	// Format date function inside template
	formatDateIndo := func(isoStr string) string {
		if isoStr == "" {
			return "-"
		}
		t, err := time.Parse(time.RFC3339, isoStr)
		if err != nil {
			t, err = time.Parse("2006-01-02 15:04:05", isoStr)
			if err != nil {
				return isoStr
			}
		}
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err == nil {
			t = t.In(loc)
		}

		bulanMap := map[time.Month]string{
			time.January: "Januari", time.February: "Februari", time.March: "Maret",
			time.April: "April", time.May: "Mei", time.June: "Juni",
			time.July: "Juli", time.August: "Agustus", time.September: "September",
			time.October: "Oktober", time.November: "November", time.December: "Desember",
		}
		return fmt.Sprintf("%d %s %d %02d:%02d", t.Day(), bulanMap[t.Month()], t.Year(), t.Hour(), t.Minute())
	}

	currentDateIndo := func() string {
		t := util.NowJakarta()
		bulanMap := map[time.Month]string{
			time.January: "Januari", time.February: "Februari", time.March: "Maret",
			time.April: "April", time.May: "Mei", time.June: "Juni",
			time.July: "Juli", time.August: "Agustus", time.September: "September",
			time.October: "Oktober", time.November: "November", time.December: "Desember",
		}
		return fmt.Sprintf("%d %s %d", t.Day(), bulanMap[t.Month()], t.Year())
	}

	tmpl, err := template.New("aduan_report").Funcs(template.FuncMap{
		"formatDate":  formatDateIndo,
		"currentDate": currentDateIndo,
	}).Parse(htmlTemplate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("template compilation error"))
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, aduan); err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("template execution error"))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Header("Content-Disposition", "inline; filename=arsip-aduan-chat-"+id+".html")
	c.String(http.StatusOK, buf.String())
}
