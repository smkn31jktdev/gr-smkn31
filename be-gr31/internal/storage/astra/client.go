package astra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	endpoint   string
	token      string
	keyspace   string
	httpClient *http.Client
}

type AstraError struct {
	StatusCode int
	Body       string
}

func (e *AstraError) Error() string {
	return fmt.Sprintf("astradb error status=%d body=%s", e.StatusCode, e.Body)
}

func NewClient(endpoint, token, keyspace string) *Client {
	return &Client{
		endpoint: endpoint,
		token:    token,
		keyspace: keyspace,
		httpClient: &http.Client{
			// Timeout per-request ke AstraDB.
			// Koleksi besar (607+ dokumen) membutuhkan ~31 request serial untuk FetchAll.
			// 30 detik per-request sudah lebih dari cukup.
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) collectionURL(collection string) string {
	return fmt.Sprintf("%s/api/json/v1/%s/%s", c.endpoint, c.keyspace, collection)
}

func (c *Client) CollectionURL(collection string) string {
	return c.collectionURL(collection)
}
func (c *Client) CollectionDocumentURL(collection, docID string) string {
	return fmt.Sprintf("%s/%s", c.collectionURL(collection), docID)
}

func (c *Client) do(ctx context.Context, collection string, payload interface{}) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.collectionURL(collection), bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &AstraError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var errResp struct {
		Errors []struct {
			Message   string `json:"message"`
			ErrorCode string `json:"errorCode"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil && len(errResp.Errors) > 0 {
		return nil, &AstraError{StatusCode: resp.StatusCode, Body: errResp.Errors[0].Message}
	}

	return body, nil
}

func (c *Client) Do(ctx context.Context, method, url string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, &AstraError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}
	return respBody, resp.StatusCode, nil
}

func (c *Client) InsertOne(ctx context.Context, collection string, data interface{}) error {
	payload := map[string]interface{}{
		"insertOne": map[string]interface{}{
			"document": data,
		},
	}
	_, err := c.do(ctx, collection, payload)
	return err
}

func (c *Client) UpsertOne(ctx context.Context, collection string, id string, data interface{}) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	var docMap map[string]interface{}
	if err := json.Unmarshal(raw, &docMap); err != nil {
		return fmt.Errorf("unmarshal to map: %w", err)
	}
	if _, hasID := docMap["_id"]; !hasID {
		docMap["_id"] = id
	}
	updatePayload := map[string]interface{}{
		"findOneAndReplace": map[string]interface{}{
			"filter":      map[string]interface{}{"_id": id},
			"replacement": docMap,
			"options":     map[string]interface{}{"upsert": true},
		},
	}
	_, err = c.do(ctx, collection, updatePayload)
	return err
}

func (c *Client) Put(ctx context.Context, url string, data interface{}) error {
	collection, id := parseCollectionAndID(url, c.endpoint, c.keyspace)
	if collection == "" {
		return fmt.Errorf("invalid url for Put: %s", url)
	}
	return c.UpsertOne(ctx, collection, id, data)
}

func (c *Client) Patch(ctx context.Context, url string, data interface{}) error {
	collection, id := parseCollectionAndID(url, c.endpoint, c.keyspace)
	if collection == "" {
		return fmt.Errorf("invalid url for Patch: %s", url)
	}
	payload := map[string]interface{}{
		"updateOne": map[string]interface{}{
			"filter": map[string]interface{}{"_id": id},
			"update": map[string]interface{}{"$set": data},
		},
	}
	_, err := c.do(ctx, collection, payload)
	return err
}

func (c *Client) Delete(ctx context.Context, url string) error {
	collection, id := parseCollectionAndID(url, c.endpoint, c.keyspace)
	if collection == "" {
		return fmt.Errorf("invalid url for Delete: %s", url)
	}
	payload := map[string]interface{}{
		"deleteOne": map[string]interface{}{
			"filter": map[string]interface{}{"_id": id},
		},
	}
	_, err := c.do(ctx, collection, payload)
	return err
}

func (c *Client) GetDocument(ctx context.Context, collection, appID string, out interface{}) error {
	payload := map[string]interface{}{
		"findOne": map[string]interface{}{
			"filter": map[string]interface{}{"_id": appID},
		},
	}
	body, err := c.do(ctx, collection, payload)
	if err != nil {
		return err
	}
	var resp struct {
		Data struct {
			Document json.RawMessage `json:"document"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("unmarshal findOne: %w", err)
	}
	if resp.Data.Document == nil {
		return &AstraError{StatusCode: 404, Body: "document not found"}
	}
	return json.Unmarshal(resp.Data.Document, out)
}

type AstraListResponse struct {
	Data      map[string]json.RawMessage `json:"data"`
	PageState string                     `json:"pageState"`
}

type AstraQuery struct {
	Where     map[string]interface{} `json:"where,omitempty"`
	PageSize  int                    `json:"page-size,omitempty"`
	PageState string                 `json:"page-state,omitempty"`
	Fields    string                 `json:"fields,omitempty"`
}

func (c *Client) Query(ctx context.Context, collection string, q AstraQuery) ([]json.RawMessage, string, error) {
	filter := convertFilter(q.Where)

	// PENTING — Semantik AstraDB Data API:
	//   options.limit = BATAS TOTAL dokumen lintas semua halaman (hard cap),
	//                   BUKAN ukuran per halaman.
	//   Ukuran per halaman selalu 20 (default API), dikontrol via pageState.
	//
	// Jika limit diset (mis. 20), AstraDB berhenti di 20 dan TIDAK mengembalikan
	// nextPageState → pagination terputus. Karena itu:
	//   - PageSize == 1  → lookup tunggal, kirim limit:1 (ambil tepat 1 dokumen).
	//   - PageSize lain  → JANGAN kirim limit; biarkan paging 20/halaman via
	//                       pageState hingga seluruh dokumen terkumpul (FetchAll).
	options := map[string]interface{}{}
	if q.PageSize == 1 {
		options["limit"] = 1
	}
	if q.PageState != "" {
		options["pageState"] = q.PageState
	}

	payload := map[string]interface{}{
		"find": map[string]interface{}{
			"filter":  filter,
			"options": options,
		},
	}

	body, err := c.do(ctx, collection, payload)
	if err != nil {
		return nil, "", err
	}

	var resp struct {
		Data struct {
			Documents     []json.RawMessage `json:"documents"`
			NextPageState string            `json:"nextPageState"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, "", fmt.Errorf("unmarshal find response: %w", err)
	}

	return resp.Data.Documents, resp.Data.NextPageState, nil
}

func IsNotFound(err error) bool {
	if aErr, ok := err.(*AstraError); ok {
		return aErr.StatusCode == 404
	}
	return false
}

// CountDocuments mengembalikan jumlah dokumen pada koleksi sesuai filter.
// Menggunakan perintah AstraDB Data API "countDocuments" — satu request,
// tanpa pagination. Cocok untuk menghitung total siswa (607+) secara efisien.
//
// Catatan: countDocuments AstraDB memiliki batas atas (default 1000). Jika
// koleksi melebihi itu, gunakan EstimatedDocumentCount.
func (c *Client) CountDocuments(ctx context.Context, collection string, where map[string]interface{}) (int, error) {
	payload := map[string]interface{}{
		"countDocuments": map[string]interface{}{
			"filter": convertFilter(where),
		},
	}
	body, err := c.do(ctx, collection, payload)
	if err != nil {
		return 0, err
	}
	var resp struct {
		Status struct {
			Count           int  `json:"count"`
			MoreData        bool `json:"moreData"`
		} `json:"status"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("unmarshal countDocuments: %w", err)
	}
	return resp.Status.Count, nil
}

// EstimatedDocumentCount mengembalikan perkiraan jumlah dokumen pada koleksi
// (tanpa filter). Tidak memiliki batas atas seperti CountDocuments — cocok
// untuk koleksi sangat besar. Hasilnya berupa estimasi.
func (c *Client) EstimatedDocumentCount(ctx context.Context, collection string) (int, error) {
	payload := map[string]interface{}{
		"estimatedDocumentCount": map[string]interface{}{},
	}
	body, err := c.do(ctx, collection, payload)
	if err != nil {
		return 0, err
	}
	var resp struct {
		Status struct {
			Count int `json:"count"`
		} `json:"status"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("unmarshal estimatedDocumentCount: %w", err)
	}
	return resp.Status.Count, nil
}

func convertFilter(where map[string]interface{}) map[string]interface{} {
	if where == nil {
		return map[string]interface{}{}
	}
	return where
}

func parseCollectionAndID(url, endpoint, keyspace string) (collection, id string) {
	prefix := fmt.Sprintf("%s/api/json/v1/%s/", endpoint, keyspace)
	if len(url) <= len(prefix) {
		return "", ""
	}
	rest := url[len(prefix):]
	for i, ch := range rest {
		if ch == '/' {
			return rest[:i], rest[i+1:]
		}
	}
	return "", ""
}
