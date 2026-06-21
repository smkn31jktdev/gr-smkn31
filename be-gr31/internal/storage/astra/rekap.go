package astra

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"be-gr31/internal/kalender"
	"be-gr31/internal/model/rekap"
)

const colRekap = "rekap_absensi"

type RekapStore struct {
	client *Client
}

func NewRekapStore(client *Client) *RekapStore {
	return &RekapStore{client: client}
}

type astraRekap struct {
	ID               string          `json:"_id"`           
	AppID            string          `json:"id,omitempty"`  
	RekapKey         string          `json:"rekap_key"`
	NIS              string          `json:"nis"`
	Nama             string          `json:"nama"`
	Kelas            string          `json:"kelas"`
	Bulan            string          `json:"bulan"`
	Semester         string          `json:"semester,omitempty"`
	TotalHadir       int             `json:"total_hadir"`
	TotalIzin        int             `json:"total_izin"`
	TotalSakit       int             `json:"total_sakit"`
	TotalTidakHadir  int             `json:"total_tidak_hadir"`
	TotalMagang      int             `json:"total_magang"`
	TotalHariEfektif int             `json:"total_hari_efektif"`
	PersentaseHadir  float64         `json:"persentase_hadir"`
	CreatedAt        json.RawMessage `json:"created_at,omitempty"`
	UpdatedAt        json.RawMessage `json:"updated_at,omitempty"`
}

func (a *astraRekap) toModel() rekap.RekapBulanan {
	docID := a.ID
	if docID == "" {
		docID = a.AppID
	}

	hariEfektif := kalender.HariEfektif(a.Bulan)
	persentase := a.PersentaseHadir
	if hariEfektif > 0 {
		persentase = float64(a.TotalHadir+a.TotalMagang) / float64(hariEfektif) * 100
	} else {
		total := a.TotalHadir + a.TotalIzin + a.TotalSakit + a.TotalTidakHadir + a.TotalMagang
		if total > 0 {
			persentase = float64(a.TotalHadir+a.TotalMagang) / float64(total) * 100
			hariEfektif = total
		}
	}

	if persentase > 100 {
		persentase = 100
	} else if persentase < 0 {
		persentase = 0
	}

	return rekap.RekapBulanan{
		ID:               docID,
		RekapKey:         a.RekapKey,
		NIS:              a.NIS,
		NamaSiswa:        a.Nama,
		Kelas:            a.Kelas,
		BulanTahun:       a.Bulan,
		Semester:         a.Semester,
		TotalHadir:       a.TotalHadir,
		TotalIzin:        a.TotalIzin,
		TotalSakit:       a.TotalSakit,
		TotalTidakHadir:  a.TotalTidakHadir,
		TotalMagang:      a.TotalMagang,
		TotalHariEfektif: hariEfektif,
		PersentaseHadir:  math.Round(persentase*10) / 10,
		CreatedAt:        string(a.CreatedAt),
		UpdatedAt:        string(a.UpdatedAt),
	}
}

func stringToRawJSON(s string) json.RawMessage {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		return json.RawMessage([]byte(s))
	}
	return json.RawMessage([]byte(fmt.Sprintf("%q", s)))
}

func fromModelRekap(r rekap.RekapBulanan) astraRekap {
	return astraRekap{
		ID:               r.ID, 
		AppID:            r.ID,
		RekapKey:         r.RekapKey,
		NIS:              r.NIS,
		Nama:             r.NamaSiswa,
		Kelas:            r.Kelas,
		Bulan:            r.BulanTahun,
		Semester:         r.Semester,
		TotalHadir:       r.TotalHadir,
		TotalIzin:        r.TotalIzin,
		TotalSakit:       r.TotalSakit,
		TotalTidakHadir:  r.TotalTidakHadir,
		TotalMagang:      r.TotalMagang,
		TotalHariEfektif: r.TotalHariEfektif,
		PersentaseHadir:  r.PersentaseHadir,
		CreatedAt:        stringToRawJSON(r.CreatedAt),
		UpdatedAt:        stringToRawJSON(r.UpdatedAt),
	}
}

func (s *RekapStore) Upsert(ctx context.Context, data *rekap.RekapBulanan) error {
	doc := fromModelRekap(*data)
	return s.client.UpsertOne(ctx, colRekap, data.ID, doc)
}

func (s *RekapStore) IncrementCounters(ctx context.Context, rekapKey string, incMap map[string]int, updatedAt string) error {
	incFields := make(map[string]any, len(incMap))
	for field, delta := range incMap {
		incFields[field] = delta
	}

	payload := map[string]any{
		"updateOne": map[string]any{
			"filter": map[string]any{"rekap_key": rekapKey},
			"update": map[string]any{
				"$inc": incFields,
				"$set": map[string]any{"updated_at": updatedAt},
			},
		},
	}
	_, err := s.client.do(ctx, colRekap, payload)
	return err
}

func (s *RekapStore) FindByKey(ctx context.Context, rekapKey string) (*rekap.RekapBulanan, error) {
	where := map[string]any{
		"rekap_key": map[string]any{"$eq": rekapKey},
	}
	items, _, err := s.client.Query(ctx, colRekap, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	var doc astraRekap
	if err := json.Unmarshal(items[0], &doc); err != nil {
		return nil, fmt.Errorf("unmarshal rekap: %w", err)
	}
	r := doc.toModel()
	return &r, nil
}

func (s *RekapStore) FindByID(ctx context.Context, id string) (*rekap.RekapBulanan, error) {
	var doc astraRekap
	if err := s.client.GetDocument(ctx, colRekap, id, &doc); err != nil {
		return nil, fmt.Errorf("rekap FindByID: %w", err)
	}
	r := doc.toModel()
	return &r, nil
}

func (s *RekapStore) ListByKeys(ctx context.Context, keys []string) ([]rekap.RekapBulanan, error) {
	if len(keys) == 0 {
		return []rekap.RekapBulanan{}, nil
	}
	result := make([]rekap.RekapBulanan, 0, len(keys))
	for _, key := range keys {
		r, err := s.FindByKey(ctx, key)
		if err != nil || r == nil {
			continue
		}
		result = append(result, *r)
	}
	return result, nil
}

func (s *RekapStore) ListPaged(ctx context.Context, filter rekap.RekapFilter, pageSize int, pageState string) ([]rekap.RekapBulanan, string, error) {
	where := map[string]any{}
	if filter.NIS != "" {
		where["nis"] = map[string]any{"$eq": filter.NIS}
	}
	if filter.Kelas != "" {
		where["kelas"] = map[string]any{"$eq": filter.Kelas}
	}
	if filter.BulanTahun != "" {
		where["bulan"] = map[string]any{"$eq": filter.BulanTahun}
	}
	if filter.Semester != "" {
		where["semester"] = map[string]any{"$eq": filter.Semester}
	}

	items, nextState, err := s.client.Query(ctx, colRekap, AstraQuery{
		Where:     where,
		PageSize:  pageSize,
		PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}

	result := make([]rekap.RekapBulanan, 0, len(items))
	for _, raw := range items {
		var doc astraRekap
		if err := json.Unmarshal(raw, &doc); err != nil {
			continue
		}
		result = append(result, doc.toModel())
	}
	return result, nextState, nil
}
