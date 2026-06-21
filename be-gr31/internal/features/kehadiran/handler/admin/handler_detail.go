package admin

import (
	"context"
	"errors"
	"math"
	"net/http"

	kehadiranmodel "be-gr31/internal/model/kehadiran"

	"github.com/gin-gonic/gin"
)

// DetailService tambahan interface untuk mengambil satu record kehadiran berdasarkan ID.
type DetailService interface {
	GetByID(ctx context.Context, id string) (*kehadiranmodel.Kehadiran, error)
}

// AbsensiDetailResponse adalah respons detail satu record absensi untuk admin.
// Menyertakan koordinat, radius sekolah, jarak siswa, dan informasi device.
type AbsensiDetailResponse struct {
	ID         string                   `json:"id"`
	NIS        string                   `json:"nis"`
	NamaSiswa  string                   `json:"namaSiswa"`
	Kelas      string                   `json:"kelas"`
	Tanggal    string                   `json:"tanggal"`
	Hari       string                   `json:"hari"`
	Status     string                   `json:"status"`
	WaktuAbsen string                   `json:"waktuAbsen"`
	Alasan     string                   `json:"alasan,omitempty"`
	FotoIzin   string                   `json:"fotoIzin,omitempty"`
	Lokasi     *LokasiDetail            `json:"lokasi,omitempty"`
	Device     *kehadiranmodel.DeviceInfo `json:"device,omitempty"`
}

// LokasiDetail menyimpan detail koordinat + radius untuk rekap absensi admin
type LokasiDetail struct {
	Koordinat     *kehadiranmodel.LatLng `json:"koordinat,omitempty"`
	JarakMeter    float64                `json:"jarakMeter"`            // jarak siswa ke sekolah (m)
	AkurasiMeter  float64                `json:"akurasiMeter"`          // akurasi GPS (m)
	RadiusMeter   float64                `json:"radiusMeter"`           // radius yang ditetapkan sekolah (m)
	DalamRadius   bool                   `json:"dalamRadius"`           // apakah jarak <= radius
}

// ConfigGetter
type ConfigGetter interface {
	SekolahRadiusMeter() float64
}

// DetailHandler
type DetailHandler struct {
	detailSvc DetailService
	radiusMeter float64 // nilai radius sekolah (meter) dari config
}

// NewDetailHandler membuat DetailHandler baru
func NewDetailHandler(svc DetailService, radiusMeter float64) *DetailHandler {
	return &DetailHandler{detailSvc: svc, radiusMeter: radiusMeter}
}

// GetDetail menangani GET /v1/admin/kehadiran/:id
func (h *DetailHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id wajib diisi"})
		return
	}

	rec, err := h.detailSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, kehadiranmodel.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "data kehadiran tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "internal server error"})
		return
	}
	if rec == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "data kehadiran tidak ditemukan"})
		return
	}

	resp := buildDetailResponse(rec, h.radiusMeter)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": resp})
}

// buildDetailResponse
func buildDetailResponse(rec *kehadiranmodel.Kehadiran, radiusMeter float64) AbsensiDetailResponse {
	resp := AbsensiDetailResponse{
		ID:         rec.ID,
		NIS:        rec.NIS,
		NamaSiswa:  rec.NamaSiswa,
		Kelas:      rec.Kelas,
		Tanggal:    rec.Tanggal,
		Hari:       rec.Hari,
		Status:     rec.Status,
		WaktuAbsen: rec.WaktuAbsen,
		Alasan:     rec.Alasan,
		FotoIzin:   rec.FotoIzin,
		Device:     rec.DeviceInfo,
	}

	if rec.Koordinat != nil {
		jarakBulat := math.Round(rec.Jarak*100) / 100
		resp.Lokasi = &LokasiDetail{
			Koordinat:    rec.Koordinat,
			JarakMeter:   jarakBulat,
			AkurasiMeter: rec.Akurasi,
			RadiusMeter:  radiusMeter,
			DalamRadius:  rec.Jarak <= radiusMeter,
		}
	}

	return resp
}
