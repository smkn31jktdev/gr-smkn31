package astra

import (
	"context"
	"encoding/json"
	"fmt"

	"be-gr31/internal/model/kehadiran"
)

const colKehadiran = "kehadiran"

type KehadiranStore struct {
	client *Client
}

func NewKehadiranStore(client *Client) *KehadiranStore {
	return &KehadiranStore{client: client}
}

type kehadiranRawDoc struct {
	ID               string            `json:"_id"`
	NIS              string            `json:"nis"`
	Nama             string            `json:"nama"`
	NamaSiswa        string            `json:"namaSiswa"`
	Kelas            string            `json:"kelas"`
	Tanggal          string            `json:"tanggal"`
	Hari             string            `json:"hari"`
	Status           string            `json:"status"`
	WaktuAbsen       string            `json:"waktuAbsen"`
	AlasanTidakHadir string            `json:"alasanTidakHadir"`
	Alasan           string            `json:"alasan"`
	FotoIzin         string            `json:"fotoIzin,omitempty"`
	Koordinat        *kehadiran.LatLng `json:"koordinat,omitempty"`
	Jarak            float64           `json:"jarak"`
	Akurasi          float64           `json:"akurasi"`
	VerifiedAt       json.RawMessage   `json:"verifiedAt,omitempty"`
	UpdatedBy        string            `json:"updatedBy,omitempty"`
	CreatedAt        json.RawMessage   `json:"createdAt,omitempty"`
	UpdatedAt        json.RawMessage   `json:"updatedAt,omitempty"`
}

func rawToKehadiran(raw kehadiranRawDoc) kehadiran.Kehadiran {
	nama := raw.NamaSiswa
	if nama == "" {
		nama = raw.Nama
	}
	alasan := raw.Alasan
	if alasan == "" {
		alasan = raw.AlasanTidakHadir
	}
	return kehadiran.Kehadiran{
		ID:         raw.ID,
		NIS:        raw.NIS,
		NamaSiswa:  nama,
		Kelas:      raw.Kelas,
		Tanggal:    raw.Tanggal,
		Hari:       raw.Hari,
		Status:     raw.Status,
		WaktuAbsen: raw.WaktuAbsen,
		Alasan:     alasan,
		FotoIzin:   raw.FotoIzin,
		Koordinat:  raw.Koordinat,
		Jarak:      raw.Jarak,
		Akurasi:    raw.Akurasi,
	}
}

func (s *KehadiranStore) Create(ctx context.Context, data *kehadiran.Kehadiran) error {
	doc := map[string]interface{}{
		"id":               data.ID,
		"nis":              data.NIS,
		"nama":             data.NamaSiswa,
		"namaSiswa":        data.NamaSiswa,
		"kelas":            data.Kelas,
		"tanggal":          data.Tanggal,
		"hari":             data.Hari,
		"status":           data.Status,
		"waktuAbsen":       data.WaktuAbsen,
		"alasanTidakHadir": data.Alasan,
		"alasan":           data.Alasan,
		"fotoIzin":         data.FotoIzin,
		"koordinat":        data.Koordinat,
		"jarak":            data.Jarak,
		"akurasi":          data.Akurasi,
	}
	return s.client.UpsertOne(ctx, colKehadiran, data.ID, doc)
}

func (s *KehadiranStore) FindByID(ctx context.Context, id string) (*kehadiran.Kehadiran, error) {
	where := map[string]any{"id": map[string]any{"$eq": id}}
	items, _, err := s.client.Query(ctx, colKehadiran, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, fmt.Errorf("kehadiran FindByID: %w", err)
	}
	if len(items) == 0 {
		where2 := map[string]any{"_id": map[string]any{"$eq": id}}
		items, _, err = s.client.Query(ctx, colKehadiran, AstraQuery{Where: where2, PageSize: 1})
		if err != nil {
			return nil, fmt.Errorf("kehadiran FindByID (_id): %w", err)
		}
	}
	if len(items) == 0 {
		return nil, nil
	}
	var raw kehadiranRawDoc
	if err := json.Unmarshal(items[0], &raw); err != nil {
		return nil, fmt.Errorf("unmarshal kehadiran: %w", err)
	}
	result := rawToKehadiran(raw)
	return &result, nil
}

func (s *KehadiranStore) FindByNISNTanggal(ctx context.Context, nis, tanggal string) (*kehadiran.Kehadiran, error) {
	where := map[string]any{
		"nis":     map[string]any{"$eq": nis},
		"tanggal": map[string]any{"$eq": tanggal},
	}
	items, _, err := s.client.Query(ctx, colKehadiran, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	var raw kehadiranRawDoc
	if err := json.Unmarshal(items[0], &raw); err != nil {
		return nil, fmt.Errorf("unmarshal kehadiran: %w", err)
	}
	result := rawToKehadiran(raw)
	return &result, nil
}

func (s *KehadiranStore) ListPaged(ctx context.Context, filter kehadiran.KehadiranFilter, pageSize int, pageState string) ([]kehadiran.Kehadiran, string, error) {
	items, nextState, err := s.client.Query(ctx, colKehadiran, AstraQuery{
		Where:     buildKehadiranFilter(filter),
		PageSize:  pageSize,
		PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]kehadiran.Kehadiran, 0, len(items))
	for _, rawBytes := range items {
		var raw kehadiranRawDoc
		if err := json.Unmarshal(rawBytes, &raw); err != nil {
			continue
		}
		result = append(result, rawToKehadiran(raw))
	}
	return result, nextState, nil
}

func (s *KehadiranStore) Update(ctx context.Context, id string, fields map[string]any) error {
	payload := map[string]any{
		"updateOne": map[string]any{
			"filter": map[string]any{"id": map[string]any{"$eq": id}},
			"update": map[string]any{"$set": fields},
		},
	}
	_, err := s.client.do(ctx, colKehadiran, payload)
	return err
}

func (s *KehadiranStore) Delete(ctx context.Context, id string) error {
	url := s.client.CollectionDocumentURL(colKehadiran, id)
	return s.client.Delete(ctx, url)
}

func buildKehadiranFilter(f kehadiran.KehadiranFilter) map[string]any {
	where := map[string]any{}
	if f.NIS != "" {
		where["nis"] = map[string]any{"$eq": f.NIS}
	}
	if f.Kelas != "" {
		where["kelas"] = map[string]any{"$eq": f.Kelas}
	}
	if f.Status != "" {
		where["status"] = map[string]any{"$eq": f.Status}
	}
	if f.Tanggal != "" {
		where["tanggal"] = map[string]any{"$eq": f.Tanggal}
	} else if f.TanggalDari != "" && f.TanggalSampai != "" {
		where["tanggal"] = map[string]any{
			"$gte": f.TanggalDari,
			"$lte": f.TanggalSampai,
		}
	} else if f.TanggalDari != "" {
		where["tanggal"] = map[string]any{"$gte": f.TanggalDari}
	} else if f.TanggalSampai != "" {
		where["tanggal"] = map[string]any{"$lte": f.TanggalSampai}
	} else if f.BulanDari != "" && f.BulanKe != "" {
		dariVal := f.BulanDari
		if len(dariVal) == 7 {
			dariVal += "-01"
		}
		keVal := f.BulanKe
		if len(keVal) == 7 {
			keVal += "-31"
		}
		where["tanggal"] = map[string]any{
			"$gte": dariVal,
			"$lte": keVal,
		}
	} else if f.BulanDari != "" {
		dariVal := f.BulanDari
		if len(dariVal) == 7 {
			dariVal += "-01"
		}
		where["tanggal"] = map[string]any{"$gte": dariVal}
	} else if f.BulanKe != "" {
		keVal := f.BulanKe
		if len(keVal) == 7 {
			keVal += "-31"
		}
		where["tanggal"] = map[string]any{"$lte": keVal}
	}
	return where
}
