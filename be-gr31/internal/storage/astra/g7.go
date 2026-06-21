package astra

import (
	"context"
	"encoding/json"
	"fmt"

	"be-gr31/internal/model/g7"
)

const colG7 = "g7"

type G7Store struct {
	client *Client
}

func NewG7Store(client *Client) *G7Store {
	return &G7Store{client: client}
}

type g7RawDoc struct {
	ID            string          `json:"_id"`
	NISN          string          `json:"nis"`
	Name          string          `json:"name"`
	NamaSiswa     string          `json:"namaSiswa"`
	Nama          string          `json:"nama"`
	Kelas         string          `json:"kelas"`
	Tanggal       string          `json:"tanggal"`
	BangunPagi    json.RawMessage `json:"bangunPagi"`
	Beribadah     json.RawMessage `json:"beribadah"`
	Ibadah        json.RawMessage `json:"ibadah"`
	Olahraga      json.RawMessage `json:"olahraga"`
	MakanSehat    json.RawMessage `json:"makanSehat"`
	Makan         json.RawMessage `json:"makan"`
	GemarBelajar  json.RawMessage `json:"gemarBelajar"`
	Belajar       json.RawMessage `json:"belajar"`
	Bermasyarakat json.RawMessage `json:"bermasyarakat"`
	TidurCepat    json.RawMessage `json:"tidurCepat"`
	Tidur         json.RawMessage `json:"tidur"`
	TotalDone     int             `json:"totalDone"`
	Bangun        json.RawMessage `json:"bangun"`
	CreatedAt     json.RawMessage `json:"createdAt"`
	UpdatedAt     json.RawMessage `json:"updatedAt"`
}

func rawToG7(raw g7RawDoc) g7.G7 {
	nama := raw.NamaSiswa
	if nama == "" {
		nama = raw.Nama
	}
	if nama == "" {
		nama = raw.Name
	}

	result := g7.G7{
		ID:        raw.ID,
		NISN:      raw.NISN,
		NamaSiswa: nama,
		Kelas:     raw.Kelas,
		Tanggal:   raw.Tanggal,
		TotalDone: raw.TotalDone,
	}

	// Bangun Pagi
	if raw.Bangun != nil && string(raw.Bangun) != "null" {
		result.Bangun = unmarshalAktivitas(raw.Bangun)
	} else if raw.BangunPagi != nil && string(raw.BangunPagi) != "null" {
		result.Bangun = parseBangunPagi(raw.BangunPagi)
	}

	// Beribadah
	if raw.Ibadah != nil && string(raw.Ibadah) != "null" {
		result.Ibadah = unmarshalAktivitas(raw.Ibadah)
	} else if raw.Beribadah != nil && string(raw.Beribadah) != "null" {
		result.Ibadah = parseBeribadah(raw.Beribadah)
	}

	// Makan Sehat
	if raw.Makan != nil && string(raw.Makan) != "null" {
		result.Makan = unmarshalAktivitas(raw.Makan)
	} else if raw.MakanSehat != nil && string(raw.MakanSehat) != "null" {
		result.Makan = unmarshalAktivitasGeneric(raw.MakanSehat)
	}

	// Olahraga
	if raw.Olahraga != nil && string(raw.Olahraga) != "null" {
		result.Olahraga = unmarshalAktivitas(raw.Olahraga)
	}

	// Gemar Belajar
	if raw.Belajar != nil && string(raw.Belajar) != "null" {
		result.Belajar = unmarshalAktivitas(raw.Belajar)
	} else if raw.GemarBelajar != nil && string(raw.GemarBelajar) != "null" {
		result.Belajar = unmarshalAktivitasGeneric(raw.GemarBelajar)
	}

	// Bermasyarakat
	if raw.Bermasyarakat != nil && string(raw.Bermasyarakat) != "null" {
		result.Bermasyarakat = unmarshalAktivitas(raw.Bermasyarakat)
	}

	// Tidur Cepat
	if raw.Tidur != nil && string(raw.Tidur) != "null" {
		result.Tidur = unmarshalAktivitas(raw.Tidur)
	} else if raw.TidurCepat != nil && string(raw.TidurCepat) != "null" {
		result.Tidur = parseTidurCepat(raw.TidurCepat)
	}

	if result.TotalDone == 0 {
		result.TotalDone = countAktivitasDone(&result)
	}

	return result
}

func unmarshalAktivitas(raw json.RawMessage) *g7.Aktivitas {
	var a g7.Aktivitas
	if err := json.Unmarshal(raw, &a); err != nil {
		return nil
	}
	return &a
}

func unmarshalAktivitasGeneric(raw json.RawMessage) *g7.Aktivitas {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	done, _ := m["done"].(bool)
	waktu, _ := m["waktu"].(string)
	ket, _ := m["keterangan"].(string)
	return &g7.Aktivitas{Done: done, Waktu: waktu, Keterangan: ket}
}

func parseBangunPagi(raw json.RawMessage) *g7.Aktivitas {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	jam, _ := m["jam"].(string)
	berdoa, _ := m["membacaDoaBangunTidur"].(bool)
	done := jam != ""
	ket := ""
	if berdoa {
		ket = "Membaca Doa: Ya"
	} else {
		ket = "Membaca Doa: Tidak"
	}
	return &g7.Aktivitas{Done: done, Waktu: jam, Keterangan: ket}
}

func parseTidurCepat(raw json.RawMessage) *g7.Aktivitas {
	return parseBangunPagi(raw)
}
func parseBeribadah(raw json.RawMessage) *g7.Aktivitas {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	done := false
	parts := []string{}
	boolFields := []string{"berdoa", "sholatFajar", "sholat5Waktu", "zikir", "sholatDhuha", "sholatSunah"}
	for _, f := range boolFields {
		if v, ok := m[f].(bool); ok && v {
			done = true
			parts = append(parts, f)
		}
	}
	if nominal, ok := m["infaq"].(float64); ok && nominal > 0 {
		done = true
	}

	ket := ""
	if len(parts) > 0 {
		ket = "Ibadah: " + joinStrings(parts)
	}
	return &g7.Aktivitas{Done: done, Keterangan: ket}
}

func joinStrings(ss []string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func countAktivitasDone(g7data *g7.G7) int {
	count := 0
	for _, a := range []*g7.Aktivitas{
		g7data.Bangun, g7data.Ibadah, g7data.Makan,
		g7data.Olahraga, g7data.Belajar, g7data.Bermasyarakat, g7data.Tidur,
	} {
		if a != nil && a.Done {
			count++
		}
	}
	return count
}

// Store methods
func (s *G7Store) Upsert(ctx context.Context, data *g7.G7) error {
	return s.client.UpsertOne(ctx, colG7, data.ID, data)
}

func (s *G7Store) FindByID(ctx context.Context, id string) (*g7.G7, error) {
	where := map[string]any{"id": map[string]any{"$eq": id}}
	items, _, err := s.client.Query(ctx, colG7, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, fmt.Errorf("g7 FindByID: %w", err)
	}
	if len(items) == 0 {
		where2 := map[string]any{"_id": map[string]any{"$eq": id}}
		items, _, err = s.client.Query(ctx, colG7, AstraQuery{Where: where2, PageSize: 1})
		if err != nil {
			return nil, fmt.Errorf("g7 FindByID (_id): %w", err)
		}
	}
	if len(items) == 0 {
		return nil, nil
	}
	var raw g7RawDoc
	if err := json.Unmarshal(items[0], &raw); err != nil {
		return nil, fmt.Errorf("unmarshal g7 raw: %w", err)
	}
	result := rawToG7(raw)
	return &result, nil
}

func (s *G7Store) FindByNISNTanggal(ctx context.Context, nisn, tanggal string) (*g7.G7, error) {
	where := map[string]any{
		"nis":     map[string]any{"$eq": nisn},
		"tanggal": map[string]any{"$eq": tanggal},
	}
	items, _, err := s.client.Query(ctx, colG7, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	var raw g7RawDoc
	if err := json.Unmarshal(items[0], &raw); err != nil {
		return nil, fmt.Errorf("unmarshal g7: %w", err)
	}
	result := rawToG7(raw)
	return &result, nil
}

func (s *G7Store) ListPaged(ctx context.Context, filter g7.G7Filter, pageSize int, pageState string) ([]g7.G7, string, error) {
	items, nextState, err := s.client.Query(ctx, colG7, AstraQuery{
		Where:     buildG7Filter(filter),
		PageSize:  pageSize,
		PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]g7.G7, 0, len(items))
	for _, rawBytes := range items {
		var raw g7RawDoc
		if err := json.Unmarshal(rawBytes, &raw); err != nil {
			continue
		}
		result = append(result, rawToG7(raw))
	}
	return result, nextState, nil
}

func (s *G7Store) Delete(ctx context.Context, id string) error {
	url := s.client.CollectionDocumentURL(colG7, id)
	return s.client.Delete(ctx, url)
}

func buildG7Filter(f g7.G7Filter) map[string]any {
	where := map[string]any{}
	if f.NISN != "" {
		where["nis"] = map[string]any{"$eq": f.NISN}
	}
	if f.Kelas != "" {
		where["kelas"] = map[string]any{"$eq": f.Kelas}
	}
	if f.Tanggal != "" {
		where["tanggal"] = map[string]any{"$eq": f.Tanggal}
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
