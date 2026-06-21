package astra

import (
	"context"
	"encoding/json"
	"fmt"

	lombamodel "be-gr31/internal/model/lomba"
)

const (
	colKebersihan = "kebersihan_kelas"
)

type KebersihanStore struct {
	client *Client
}

func NewKebersihanStore(client *Client) *KebersihanStore {
	return &KebersihanStore{client: client}
}

func (s *KebersihanStore) Create(ctx context.Context, data *lombamodel.KebersihanKelas) error {
	return s.client.UpsertOne(ctx, colKebersihan, data.ID, data)
}

func (s *KebersihanStore) Update(ctx context.Context, data *lombamodel.KebersihanKelas) error {
	return s.client.UpsertOne(ctx, colKebersihan, data.ID, data)
}

func (s *KebersihanStore) FindByID(ctx context.Context, id string) (*lombamodel.KebersihanKelas, error) {
	var result lombamodel.KebersihanKelas
	if err := s.client.GetDocument(ctx, colKebersihan, id, &result); err != nil {
		return nil, fmt.Errorf("kebersihan FindByID: %w", err)
	}
	return &result, nil
}

func (s *KebersihanStore) FindByKelasTanggal(ctx context.Context, kelas, tanggal string) (*lombamodel.KebersihanKelas, error) {
	where := map[string]any{
		"kelas":   map[string]any{"$eq": kelas},
		"tanggal": map[string]any{"$eq": tanggal},
	}
	items, _, err := s.client.Query(ctx, colKebersihan, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	var result lombamodel.KebersihanKelas
	if err := json.Unmarshal(items[0], &result); err != nil {
		return nil, fmt.Errorf("unmarshal kebersihan: %w", err)
	}
	return &result, nil
}

func (s *KebersihanStore) Delete(ctx context.Context, id string) error {
	url := s.client.CollectionDocumentURL(colKebersihan, id)
	return s.client.Delete(ctx, url)
}

func (s *KebersihanStore) ListPaged(ctx context.Context, filter lombamodel.KebersihanFilter, pageSize int, pageState string) ([]lombamodel.KebersihanKelas, string, error) {
	where := map[string]any{}
	if filter.Kelas != "" {
		where["kelas"] = map[string]any{"$eq": filter.Kelas}
	}
	if filter.Tanggal != "" {
		where["tanggal"] = map[string]any{"$eq": filter.Tanggal}
	} else if filter.DariTgl != "" && filter.SampaiTgl != "" {
		where["tanggal"] = map[string]any{"$gte": filter.DariTgl, "$lte": filter.SampaiTgl}
	}

	items, nextState, err := s.client.Query(ctx, colKebersihan, AstraQuery{
		Where: where, PageSize: pageSize, PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]lombamodel.KebersihanKelas, 0, len(items))
	for _, raw := range items {
		var item lombamodel.KebersihanKelas
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result, nextState, nil
}
