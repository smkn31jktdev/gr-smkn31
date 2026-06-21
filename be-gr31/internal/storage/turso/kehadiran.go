package turso

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"be-gr31/internal/model/common"
	"be-gr31/internal/model/kehadiran"
)

type KehadiranStore struct {
	client *Client
}

func NewKehadiranStore(client *Client) *KehadiranStore {
	return &KehadiranStore{client: client}
}

var columnMapping = map[string]string{
	"id":               "id",
	"nis":              "nis",
	"namaSiswa":        "nama_siswa",
	"nama":             "nama_siswa",
	"kelas":            "kelas",
	"tanggal":          "tanggal",
	"hari":             "hari",
	"status":           "status",
	"waktuAbsen":       "waktu_absen",
	"alasan":           "alasan",
	"alasanTidakHadir": "alasan",
	"fotoIzin":         "foto_izin",
	"jarak":            "jarak",
	"akurasi":          "akurasi",
	"deviceInfo":       "device_info",
	"createdAt":        "created_at",
	"updatedAt":        "updated_at",
}

// Menyimpan atau memperbarui data kehadiran
func (s *KehadiranStore) Create(ctx context.Context, data *kehadiran.Kehadiran) error {
	var lat, lng any
	if data.Koordinat != nil {
		lat = data.Koordinat.Lat
		lng = data.Koordinat.Lng
	}

	var deviceInfoStr sql.NullString
	if data.DeviceInfo != nil {
		if bytes, err := json.Marshal(data.DeviceInfo); err == nil {
			deviceInfoStr = sql.NullString{String: string(bytes), Valid: true}
		}
	}

	query := `
	INSERT INTO kehadiran (
		id, nis, nama_siswa, kelas, tanggal, hari, status, waktu_absen,
		alasan, foto_izin, koordinat_lat, koordinat_lng, jarak, akurasi,
		device_info, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status = excluded.status,
		waktu_absen = excluded.waktu_absen,
		alasan = excluded.alasan,
		foto_izin = excluded.foto_izin,
		koordinat_lat = excluded.koordinat_lat,
		koordinat_lng = excluded.koordinat_lng,
		jarak = excluded.jarak,
		akurasi = excluded.akurasi,
		device_info = excluded.device_info,
		updated_at = excluded.updated_at
	`
	_, err := s.client.DB.ExecContext(ctx, query,
		data.ID, data.NIS, data.NamaSiswa, data.Kelas, data.Tanggal, data.Hari, data.Status, data.WaktuAbsen,
		data.Alasan, data.FotoIzin, lat, lng, data.Jarak, data.Akurasi, deviceInfoStr,
		string(data.CreatedAt), string(data.UpdatedAt),
	)
	return err
}

// FindByID
func (s *KehadiranStore) FindByID(ctx context.Context, id string) (*kehadiran.Kehadiran, error) {
	query := `
	SELECT id, nis, nama_siswa, kelas, tanggal, hari, status, waktu_absen,
	       alasan, foto_izin, koordinat_lat, koordinat_lng, jarak, akurasi,
	       device_info, created_at, updated_at
	FROM kehadiran
	WHERE id = ?
	`
	row := s.client.DB.QueryRowContext(ctx, query, id)
	return scanKehadiran(row)
}

// FindByNISNTanggal
func (s *KehadiranStore) FindByNISNTanggal(ctx context.Context, nis, tanggal string) (*kehadiran.Kehadiran, error) {
	query := `
	SELECT id, nis, nama_siswa, kelas, tanggal, hari, status, waktu_absen,
	       alasan, foto_izin, koordinat_lat, koordinat_lng, jarak, akurasi,
	       device_info, created_at, updated_at
	FROM kehadiran
	WHERE nis = ? AND tanggal = ?
	`
	row := s.client.DB.QueryRowContext(ctx, query, nis, tanggal)
	return scanKehadiran(row)
}

// ListPaged mengambil kehadiran dengan filter, pageSize, dan pageState (offset cursor).
func (s *KehadiranStore) ListPaged(ctx context.Context, filter kehadiran.KehadiranFilter, pageSize int, pageState string) ([]kehadiran.Kehadiran, string, error) {
	query := `
	SELECT id, nis, nama_siswa, kelas, tanggal, hari, status, waktu_absen,
	       alasan, foto_izin, koordinat_lat, koordinat_lng, jarak, akurasi,
	       device_info, created_at, updated_at
	FROM kehadiran
	WHERE 1=1
	`
	var args []any

	if filter.NIS != "" {
		query += " AND nis = ?"
		args = append(args, filter.NIS)
	}
	if filter.Kelas != "" {
		query += " AND kelas = ?"
		args = append(args, filter.Kelas)
	}
	if filter.Status != "" {
		if filter.Status == "izin_sakit" {
			query += " AND (status = 'izin' OR status = 'sakit')"
		} else {
			query += " AND status = ?"
			args = append(args, filter.Status)
		}
	}
	if filter.Tanggal != "" {
		query += " AND tanggal = ?"
		args = append(args, filter.Tanggal)
	} else if filter.TanggalDari != "" && filter.TanggalSampai != "" {
		query += " AND tanggal >= ? AND tanggal <= ?"
		args = append(args, filter.TanggalDari, filter.TanggalSampai)
	} else if filter.TanggalDari != "" {
		query += " AND tanggal >= ?"
		args = append(args, filter.TanggalDari)
	} else if filter.TanggalSampai != "" {
		query += " AND tanggal <= ?"
		args = append(args, filter.TanggalSampai)
	} else if filter.BulanDari != "" && filter.BulanKe != "" {
		dariVal := filter.BulanDari
		if len(dariVal) == 7 {
			dariVal += "-01"
		}
		keVal := filter.BulanKe
		if len(keVal) == 7 {
			keVal += "-31"
		}
		query += " AND tanggal >= ? AND tanggal <= ?"
		args = append(args, dariVal, keVal)
	} else if filter.BulanDari != "" {
		dariVal := filter.BulanDari
		if len(dariVal) == 7 {
			dariVal += "-01"
		}
		query += " AND tanggal >= ?"
		args = append(args, dariVal)
	} else if filter.BulanKe != "" {
		keVal := filter.BulanKe
		if len(keVal) == 7 {
			keVal += "-31"
		}
		query += " AND tanggal <= ?"
		args = append(args, keVal)
	}

	// Urutkan berdasarkan tanggal desc, waktu_absen desc
	query += " ORDER BY tanggal DESC, waktu_absen DESC"

	offset := 0
	if pageState != "" {
		if o, err := strconv.Atoi(pageState); err == nil {
			offset = o
		}
	}

	if pageSize > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, pageSize, offset)
	}

	rows, err := s.client.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("query list paged: %w", err)
	}
	defer rows.Close()

	items, err := scanKehadiranRows(rows)
	if err != nil {
		return nil, "", fmt.Errorf("scan rows: %w", err)
	}

	var nextState string
	if pageSize > 0 && len(items) == pageSize {
		nextState = strconv.Itoa(offset + len(items))
	}

	return items, nextState, nil
}

// Memperbarui data berdasarkan fields map
func (s *KehadiranStore) Update(ctx context.Context, id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	query := "UPDATE kehadiran SET "
	var args []any
	var setParts []string

	for k, v := range fields {
		colName, exists := columnMapping[k]
		if !exists {
			continue
		}
		if colName == "id" {
			continue
		}

		if colName == "koordinat" {
			if coords, ok := v.(*kehadiran.LatLng); ok && coords != nil {
				setParts = append(setParts, "koordinat_lat = ?", "koordinat_lng = ?")
				args = append(args, coords.Lat, coords.Lng)
			} else if coords, ok := v.(kehadiran.LatLng); ok {
				setParts = append(setParts, "koordinat_lat = ?", "koordinat_lng = ?")
				args = append(args, coords.Lat, coords.Lng)
			}
			continue
		}

		if colName == "device_info" {
			if dev, ok := v.(*kehadiran.DeviceInfo); ok && dev != nil {
				if bytes, err := json.Marshal(dev); err == nil {
					setParts = append(setParts, "device_info = ?")
					args = append(args, string(bytes))
				}
			} else if dev, ok := v.(kehadiran.DeviceInfo); ok {
				if bytes, err := json.Marshal(dev); err == nil {
					setParts = append(setParts, "device_info = ?")
					args = append(args, string(bytes))
				}
			}
			continue
		}

		setParts = append(setParts, fmt.Sprintf("%s = ?", colName))
		args = append(args, v)
	}

	if len(setParts) == 0 {
		return nil
	}

	if _, hasUpdatedAt := fields["updatedAt"]; !hasUpdatedAt {
		setParts = append(setParts, "updated_at = ?")
		args = append(args, time.Now().Format(time.RFC3339))
	}

	query += strings.Join(setParts, ", ")
	query += " WHERE id = ?"
	args = append(args, id)

	_, err := s.client.DB.ExecContext(ctx, query, args...)
	return err
}

// Menghapus data kehadiran berdasarkan ID
func (s *KehadiranStore) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM kehadiran WHERE id = ?"
	_, err := s.client.DB.ExecContext(ctx, query, id)
	return err
}

// Helper scanner tunggal
func scanKehadiran(scanner interface{ Scan(dest ...any) error }) (*kehadiran.Kehadiran, error) {
	var k kehadiran.Kehadiran
	var latVal, lngVal sql.NullFloat64
	var alasanVal, fotoVal sql.NullString
	var deviceInfoVal sql.NullString
	var createdAtStr, updatedAtStr string

	err := scanner.Scan(
		&k.ID, &k.NIS, &k.NamaSiswa, &k.Kelas, &k.Tanggal, &k.Hari, &k.Status, &k.WaktuAbsen,
		&alasanVal, &fotoVal, &latVal, &lngVal, &k.Jarak, &k.Akurasi, &deviceInfoVal,
		&createdAtStr, &updatedAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if latVal.Valid && lngVal.Valid {
		k.Koordinat = &kehadiran.LatLng{
			Lat: latVal.Float64,
			Lng: lngVal.Float64,
		}
	}
	if alasanVal.Valid {
		k.Alasan = alasanVal.String
	}
	if fotoVal.Valid {
		k.FotoIzin = fotoVal.String
	}
	if deviceInfoVal.Valid && deviceInfoVal.String != "" {
		var dev kehadiran.DeviceInfo
		if err := json.Unmarshal([]byte(deviceInfoVal.String), &dev); err == nil {
			k.DeviceInfo = &dev
		}
	}
	k.CreatedAt = common.FlexTime(createdAtStr)
	k.UpdatedAt = common.FlexTime(updatedAtStr)

	return &k, nil
}

// Helper scanner list
func scanKehadiranRows(rows *sql.Rows) ([]kehadiran.Kehadiran, error) {
	var result []kehadiran.Kehadiran
	for rows.Next() {
		kPtr, err := scanKehadiran(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *kPtr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
