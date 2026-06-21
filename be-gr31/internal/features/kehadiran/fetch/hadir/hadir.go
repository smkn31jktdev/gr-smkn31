package hadir

import (
	"errors"
	"fmt"

	"be-gr31/internal/util"
	kehadiranmodel "be-gr31/internal/model/kehadiran"
)

var (
	ErrTooFar           = errors.New("lokasi terlalu jauh dari sekolah")
	ErrAccuracySuspect  = errors.New("GPS accuracy mencurigakan, kemungkinan menggunakan mock location")
	ErrAccuracyTooWeak  = errors.New("sinyal GPS terlalu lemah")
)

// ValidateGPS memvalidasi radius GPS siswa untuk status hadir.
// Jika accuracy disediakan (> 0), validasi accuracy dilakukan sebelum cek radius.
func ValidateGPS(
	koordinat *kehadiranmodel.LatLng,
	status string,
	sekolahLat, sekolahLng, sekolahRadiusMeter float64,
	accuracy, minAccuracy, maxAccuracy float64,
) (float64, error) {
	if koordinat == nil {
		return 0, nil
	}

	// Validasi accuracy hanya jika status hadir dan accuracy dikirim client (> 0)
	if status == "hadir" && accuracy > 0 {
		if accuracy < minAccuracy {
			return 0, fmt.Errorf(
				"%w: akurasi %.1fm terlalu sempurna (minimal %.0fm). Matikan aplikasi mock location",
				ErrAccuracySuspect, accuracy, minAccuracy,
			)
		}
		if accuracy >= maxAccuracy {
			return 0, fmt.Errorf(
				"%w: akurasi %.0fm terlalu rendah (maks < %.0fm). Pindah ke area terbuka",
				ErrAccuracyTooWeak, accuracy, maxAccuracy,
			)
		}
	}

	if status == "hadir" {
		dist, err := util.ValidateRadius(
			koordinat.Lat, koordinat.Lng,
			sekolahLat, sekolahLng,
			sekolahRadiusMeter,
		)
		if err != nil {
			return 0, fmt.Errorf("lokasi terlalu jauh dari sekolah: %s", err.Error())
		}
		return dist, nil
	}

	return util.Haversine(
		koordinat.Lat, koordinat.Lng,
		sekolahLat, sekolahLng,
	), nil
}
