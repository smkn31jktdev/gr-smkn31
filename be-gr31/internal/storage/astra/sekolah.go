package astra

import (
	"context"
	"encoding/json"
	"fmt"

	"be-gr31/internal/model/aduan"
	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/sekolah"
)

const (
	colStudents = "akun_siswa"
	colAdmins   = "akun_admin"
	colClients  = "clients"
	colBukti    = "bukti"
	colKegiatan = "kebiasaan_hebat"
	colAduan    = "aduan"
)

// Students
type StudentStore struct {
	client *Client
}

func NewStudentStore(client *Client) *StudentStore {
	return &StudentStore{client: client}
}

func (s *StudentStore) FindByNISN(ctx context.Context, nisn string) (*authmodel.Siswa, error) {
	where := map[string]any{"nis": map[string]any{"$eq": nisn}}
	items, _, err := s.client.Query(ctx, colStudents, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		where2 := map[string]any{"nisn": map[string]any{"$eq": nisn}}
		items, _, err = s.client.Query(ctx, colStudents, AstraQuery{Where: where2, PageSize: 1})
		if err != nil {
			return nil, err
		}
	}
	if len(items) == 0 {
		return nil, nil
	}
	var result authmodel.Siswa
	if err := json.Unmarshal(items[0], &result); err != nil {
		return nil, fmt.Errorf("unmarshal siswa: %w", err)
	}
	return &result, nil
}

func (s *StudentStore) FindByID(ctx context.Context, id string) (*authmodel.Siswa, error) {
	var result authmodel.Siswa
	if err := s.client.GetDocument(ctx, colStudents, id, &result); err != nil {
		return nil, fmt.Errorf("siswa FindByID: %w", err)
	}
	return &result, nil
}

func (s *StudentStore) Upsert(ctx context.Context, data *authmodel.Siswa) error {
	return s.client.UpsertOne(ctx, colStudents, data.ID, data)
}

func (s *StudentStore) Delete(ctx context.Context, id string) error {
	url := s.client.CollectionDocumentURL(colStudents, id)
	return s.client.Delete(ctx, url)
}
func (s *StudentStore) ListPaged(ctx context.Context, kelas string, pageSize int, pageState string) ([]authmodel.Siswa, string, error) {
	where := map[string]any{}
	if kelas != "" {
		where["kelas"] = map[string]any{"$eq": kelas}
	}
	items, nextState, err := s.client.Query(ctx, colStudents, AstraQuery{
		Where: where, PageSize: pageSize, PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]authmodel.Siswa, 0, len(items))
	for _, raw := range items {
		var item authmodel.Siswa
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result, nextState, nil
}

func (s *StudentStore) ListByWalas(ctx context.Context, walas string) ([]authmodel.Siswa, error) {
	where := map[string]any{"walas": map[string]any{"$eq": walas}}
	items, _, err := s.client.Query(ctx, colStudents, AstraQuery{Where: where, PageSize: 100})
	if err != nil {
		return nil, err
	}
	result := make([]authmodel.Siswa, 0, len(items))
	for _, raw := range items {
		var item authmodel.Siswa
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *StudentStore) Count(ctx context.Context, kelas string) (int, error) {
	where := map[string]any{}
	if kelas != "" {
		where["kelas"] = map[string]any{"$eq": kelas}
	}
	n, err := s.client.CountDocuments(ctx, colStudents, where)
	if err != nil {
		return 0, err
	}
	if kelas == "" && n == 0 {
		if est, errEst := s.client.EstimatedDocumentCount(ctx, colStudents); errEst == nil && est > 0 {
			return est, nil
		}
	}
	return n, nil
}

// Admins
type AdminStore struct {
	client *Client
}

func NewAdminStore(client *Client) *AdminStore {
	return &AdminStore{client: client}
}

func (s *AdminStore) FindByEmail(ctx context.Context, email string) (*authmodel.Admin, error) {
	where := map[string]any{"email": map[string]any{"$eq": email}}
	items, _, err := s.client.Query(ctx, colAdmins, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	var result authmodel.Admin
	if err := json.Unmarshal(items[0], &result); err != nil {
		return nil, fmt.Errorf("unmarshal admin: %w", err)
	}
	return &result, nil
}

func (s *AdminStore) FindByID(ctx context.Context, id string) (*authmodel.Admin, error) {
	where := map[string]any{"id": map[string]any{"$eq": id}}
	items, _, err := s.client.Query(ctx, colAdmins, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, fmt.Errorf("admin FindByID (by id): %w", err)
	}
	if len(items) > 0 {
		var result authmodel.Admin
		if err := json.Unmarshal(items[0], &result); err != nil {
			return nil, fmt.Errorf("unmarshal admin: %w", err)
		}
		return &result, nil
	}

	where2 := map[string]any{"_id": map[string]any{"$eq": id}}
	items2, _, err2 := s.client.Query(ctx, colAdmins, AstraQuery{Where: where2, PageSize: 1})
	if err2 != nil {
		return nil, fmt.Errorf("admin FindByID (by _id): %w", err2)
	}
	if len(items2) == 0 {
		return nil, nil
	}
	var result authmodel.Admin
	if err := json.Unmarshal(items2[0], &result); err != nil {
		return nil, fmt.Errorf("unmarshal admin: %w", err)
	}
	return &result, nil
}

func (s *AdminStore) Upsert(ctx context.Context, data *authmodel.Admin) error {
	return s.client.UpsertOne(ctx, colAdmins, data.ID, data)
}
func (s *AdminStore) Delete(ctx context.Context, id string) error {
	url := s.client.CollectionDocumentURL(colAdmins, id)
	return s.client.Delete(ctx, url)
}

func (s *AdminStore) ListPaged(ctx context.Context, role string, pageSize int, pageState string) ([]authmodel.Admin, string, error) {
	where := map[string]any{}
	if role != "" {
		where["role"] = map[string]any{"$eq": role}
	}
	items, nextState, err := s.client.Query(ctx, colAdmins, AstraQuery{
		Where: where, PageSize: pageSize, PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]authmodel.Admin, 0, len(items))
	for _, raw := range items {
		var item authmodel.Admin
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result, nextState, nil
}

// Bukti
type buktiRawDoc struct {
	ID          string          `json:"_id"`
	NIS         string          `json:"nis"`
	Nama        string          `json:"nama"`
	NamaSiswa   string          `json:"namaSiswa"`
	Kelas       string          `json:"kelas"`
	Bulan       string          `json:"bulan"`
	Foto        []string        `json:"foto"`
	LinkYouTube []string        `json:"linkYouTube"`
	LinkYT      []string        `json:"linkYT"`
	CreatedAt   json.RawMessage `json:"createdAt,omitempty"`
	UpdatedAt   json.RawMessage `json:"updatedAt,omitempty"`
}

func rawToBukti(raw buktiRawDoc) sekolah.Bukti {
	nama := raw.NamaSiswa
	if nama == "" {
		nama = raw.Nama
	}
	linkYT := raw.LinkYT
	if len(linkYT) == 0 {
		linkYT = raw.LinkYouTube
	}
	return sekolah.Bukti{
		ID:        raw.ID,
		NIS:       raw.NIS,
		NamaSiswa: nama,
		Kelas:     raw.Kelas,
		Bulan:     raw.Bulan,
		Foto:      raw.Foto,
		LinkYT:    linkYT,
	}
}

type BuktiStore struct {
	client *Client
}

func NewBuktiStore(client *Client) *BuktiStore {
	return &BuktiStore{client: client}
}

func (s *BuktiStore) Upsert(ctx context.Context, data *sekolah.Bukti) error {
	doc := map[string]interface{}{
		"id":          data.ID,
		"nis":         data.NIS,
		"nama":        data.NamaSiswa,
		"namaSiswa":   data.NamaSiswa,
		"kelas":       data.Kelas,
		"bulan":       data.Bulan,
		"foto":        data.Foto,
		"linkYouTube": data.LinkYT,
		"linkYT":      data.LinkYT,
	}
	return s.client.UpsertOne(ctx, colBukti, data.ID, doc)
}

func (s *BuktiStore) ListPaged(ctx context.Context, filter sekolah.BuktiFilter, pageSize int, pageState string) ([]sekolah.Bukti, string, error) {
	where := map[string]any{}
	if filter.NIS != "" {
		where["nis"] = map[string]any{"$eq": filter.NIS}
	}
	if filter.Kelas != "" {
		where["kelas"] = map[string]any{"$eq": filter.Kelas}
	}
	if filter.Bulan != "" {
		where["bulan"] = map[string]any{"$eq": filter.Bulan}
	}
	items, nextState, err := s.client.Query(ctx, colBukti, AstraQuery{
		Where: where, PageSize: pageSize, PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]sekolah.Bukti, 0, len(items))
	for _, rawBytes := range items {
		var raw buktiRawDoc
		if err := json.Unmarshal(rawBytes, &raw); err != nil {
			continue
		}
		result = append(result, rawToBukti(raw))
	}
	return result, nextState, nil
}

// Kegiatan

type KegiatanStore struct {
	client *Client
}

func NewKegiatanStore(client *Client) *KegiatanStore {
	return &KegiatanStore{client: client}
}

func (s *KegiatanStore) Create(ctx context.Context, data *sekolah.Kegiatan) error {
	return s.client.UpsertOne(ctx, colKegiatan, data.ID, data)
}

func (s *KegiatanStore) Update(ctx context.Context, data *sekolah.Kegiatan) error {
	return s.client.UpsertOne(ctx, colKegiatan, data.ID, data)
}

func (s *KegiatanStore) FindByID(ctx context.Context, id string) (*sekolah.Kegiatan, error) {
	var result sekolah.Kegiatan
	if err := s.client.GetDocument(ctx, colKegiatan, id, &result); err != nil {
		return nil, fmt.Errorf("kegiatan FindByID: %w", err)
	}
	return &result, nil
}

func (s *KegiatanStore) Delete(ctx context.Context, id string) error {
	url := s.client.CollectionDocumentURL(colKegiatan, id)
	return s.client.Delete(ctx, url)
}

func (s *KegiatanStore) ListPaged(ctx context.Context, filter sekolah.KegiatanFilter, pageSize int, pageState string) ([]sekolah.Kegiatan, string, error) {
	where := map[string]any{}
	if filter.NISN != "" {
		where["nisn"] = map[string]any{"$eq": filter.NISN}
	}
	if filter.Kelas != "" {
		where["kelas"] = map[string]any{"$eq": filter.Kelas}
	}
	if filter.Section != "" {
		where["section"] = map[string]any{"$eq": filter.Section}
	}
	if filter.Tanggal != "" {
		where["tanggal"] = map[string]any{"$eq": filter.Tanggal}
	} else if filter.DariTgl != "" && filter.SampaiTgl != "" {
		where["tanggal"] = map[string]any{"$gte": filter.DariTgl, "$lte": filter.SampaiTgl}
	}

	items, nextState, err := s.client.Query(ctx, colKegiatan, AstraQuery{
		Where: where, PageSize: pageSize, PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]sekolah.Kegiatan, 0, len(items))
	for _, raw := range items {
		var item sekolah.Kegiatan
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result, nextState, nil
}

// Aduan
type AduanStore struct {
	client *Client
}

func NewAduanStore(client *Client) *AduanStore {
	return &AduanStore{client: client}
}

func (s *AduanStore) Create(ctx context.Context, data *aduan.Aduan) error {
	return s.client.UpsertOne(ctx, colAduan, data.ID, data)
}

func (s *AduanStore) Update(ctx context.Context, data *aduan.Aduan) error {
	return s.client.UpsertOne(ctx, colAduan, data.ID, data)
}

func (s *AduanStore) FindByID(ctx context.Context, id string) (*aduan.Aduan, error) {
	var result aduan.Aduan
	if err := s.client.GetDocument(ctx, colAduan, id, &result); err != nil {
		return nil, fmt.Errorf("aduan FindByID: %w", err)
	}
	return &result, nil
}

func (s *AduanStore) ListPaged(ctx context.Context, filter aduan.AduanFilter, pageSize int, pageState string) ([]aduan.Aduan, string, error) {
	where := map[string]any{}
	if filter.NISN != "" {
		where["nisn"] = map[string]any{"$eq": filter.NISN}
	}
	if filter.Status != "" {
		where["status"] = map[string]any{"$eq": filter.Status}
	}
	if filter.AdminNama != "" {
		where["adminNama"] = map[string]any{"$eq": filter.AdminNama}
	}
	items, nextState, err := s.client.Query(ctx, colAduan, AstraQuery{
		Where: where, PageSize: pageSize, PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]aduan.Aduan, 0, len(items))
	for _, raw := range items {
		var item aduan.Aduan
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result, nextState, nil
}
