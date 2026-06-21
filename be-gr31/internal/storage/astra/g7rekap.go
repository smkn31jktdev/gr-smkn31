package astra

import (
	"context"
	"encoding/json"
	"fmt"

	"be-gr31/internal/model/g7"
)

const colG7Rekap = "g7_rekap"

type G7RekapStore struct {
	client *Client
}

func NewG7RekapStore(client *Client) *G7RekapStore {
	return &G7RekapStore{client: client}
}

func (s *G7RekapStore) Upsert(ctx context.Context, data *g7.G7Rekap) error {
	return s.client.UpsertOne(ctx, colG7Rekap, data.ID, data)
}

func (s *G7RekapStore) FindByID(ctx context.Context, id string) (*g7.G7Rekap, error) {
	var result g7.G7Rekap
	if err := s.client.GetDocument(ctx, colG7Rekap, id, &result); err != nil {
		return nil, fmt.Errorf("g7rekap FindByID: %w", err)
	}
	return &result, nil
}

func (s *G7RekapStore) FindByNISNBulan(ctx context.Context, nisn, bulan string) (*g7.G7Rekap, error) {
	where := map[string]any{
		"nis":        map[string]any{"$eq": nisn},
		"bulanTahun": map[string]any{"$eq": bulan},
	}
	items, _, err := s.client.Query(ctx, colG7Rekap, AstraQuery{Where: where, PageSize: 1})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	var result g7.G7Rekap
	if err := json.Unmarshal(items[0], &result); err != nil {
		return nil, fmt.Errorf("unmarshal g7rekap: %w", err)
	}
	return &result, nil
}

func (s *G7RekapStore) ListPaged(ctx context.Context, filter g7.G7RekapFilter, pageSize int, pageState string) ([]g7.G7Rekap, string, error) {
	items, nextState, err := s.client.Query(ctx, colG7Rekap, AstraQuery{
		Where:     buildG7RekapFilter(filter),
		PageSize:  pageSize,
		PageState: pageState,
	})
	if err != nil {
		return nil, "", err
	}
	result := make([]g7.G7Rekap, 0, len(items))
	for _, raw := range items {
		var item g7.G7Rekap
		if err := json.Unmarshal(raw, &item); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result, nextState, nil
}

func (s *G7RekapStore) Delete(ctx context.Context, id string) error {
	url := s.client.CollectionDocumentURL(colG7Rekap, id)
	return s.client.Delete(ctx, url)
}

func buildG7RekapFilter(f g7.G7RekapFilter) map[string]any {
	where := map[string]any{}
	if f.NISN != "" {
		where["nis"] = map[string]any{"$eq": f.NISN}
	}
	if f.Kelas != "" {
		where["kelas"] = map[string]any{"$eq": f.Kelas}
	}
	if f.BulanTahun != "" {
		where["bulanTahun"] = map[string]any{"$eq": f.BulanTahun}
	}
	if f.Predikat != "" {
		where["predikat"] = map[string]any{"$eq": f.Predikat}
	}
	if f.Status != "" {
		where["status"] = map[string]any{"$eq": f.Status}
	}
	return where
}
