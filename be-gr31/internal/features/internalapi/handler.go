package internalapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"be-gr31/internal/model/common"
	"be-gr31/internal/storage/astra"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Client merepresentasikan internal client registry.
type Client struct {
	ID        string `json:"id" bson:"_id"`
	AppName   string `json:"appName" bson:"appName"`
	AppKey    string `json:"appKey" bson:"appKey"`
	IsActive  bool   `json:"isActive" bson:"isActive"`
	CreatedAt string `json:"createdAt" bson:"createdAt"`
	UpdatedAt string `json:"updatedAt" bson:"updatedAt"`
}

// ClientUpsertRequest adalah body request upsert client.
type ClientUpsertRequest struct {
	AppName  string `json:"appName" binding:"required"`
	AppKey   string `json:"appKey" binding:"required"`
	IsActive *bool  `json:"isActive"`
}

// ClientFindRequest adalah body request find client.
type ClientFindRequest struct {
	AppName string `json:"appName"`
	AppKey  string `json:"appKey"`
}

// ClientDeleteRequest adalah body request delete client.
type ClientDeleteRequest struct {
	ID string `json:"id" binding:"required"`
}

const colClients = "clients"

// Handler menangani HTTP request untuk internal client registry.
type Handler struct {
	astraClient *astra.Client
}

// NewHandler membuat instance Handler baru.
func NewHandler(astraClient *astra.Client) *Handler {
	return &Handler{astraClient: astraClient}
}

// UpsertClient menangani POST /internal/client-upsert.
func (h *Handler) UpsertClient(c *gin.Context) {
	var req ClientUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	ctx := c.Request.Context()
	// Cari existing client by appName
	existing := h.findByAppName(ctx, req.AppName)

	now := time.Now().Format(time.RFC3339)
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var client Client
	if existing != nil {
		client = *existing
		client.AppKey = req.AppKey
		client.IsActive = isActive
		client.UpdatedAt = now
	} else {
		client = Client{
			ID:        uuid.New().String(),
			AppName:   req.AppName,
			AppKey:    req.AppKey,
			IsActive:  isActive,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	url := h.astraClient.CollectionDocumentURL(colClients, client.ID)
	if err := h.astraClient.Put(ctx, url, client); err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("gagal menyimpan client"))
		return
	}

	c.JSON(http.StatusOK, common.OK(client, "client berhasil disimpan"))
}

// FindClient menangani POST /internal/client-find.
func (h *Handler) FindClient(c *gin.Context) {
	var req ClientFindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	ctx := c.Request.Context()
	where := make(map[string]interface{})
	if req.AppName != "" {
		where["appName"] = map[string]interface{}{"$eq": req.AppName}
	}
	if req.AppKey != "" {
		where["appKey"] = map[string]interface{}{"$eq": req.AppKey}
	}

	items, _, err := h.astraClient.Query(ctx, colClients, astra.AstraQuery{
		Where:    where,
		PageSize: 10,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("gagal mengambil client"))
		return
	}

	clients := make([]Client, 0, len(items))
	for _, raw := range items {
		var cl Client
		if err := json.Unmarshal(raw, &cl); err != nil {
			continue
		}
		clients = append(clients, cl)
	}

	c.JSON(http.StatusOK, common.OK(clients, "ok"))
}

// DeleteClient menangani DELETE /internal/client-delete.
func (h *Handler) DeleteClient(c *gin.Context) {
	var req ClientDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Fail(err.Error()))
		return
	}

	url := h.astraClient.CollectionDocumentURL(colClients, req.ID)
	if err := h.astraClient.Delete(c.Request.Context(), url); err != nil {
		c.JSON(http.StatusInternalServerError, common.Fail("gagal menghapus client"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// findByAppName mencari client berdasarkan appName.
func (h *Handler) findByAppName(ctx context.Context, appName string) *Client {
	where := map[string]interface{}{
		"appName": map[string]interface{}{"$eq": appName},
	}
	items, _, err := h.astraClient.Query(ctx, colClients, astra.AstraQuery{
		Where:    where,
		PageSize: 1,
	})
	if err != nil || len(items) == 0 {
		return nil
	}
	var result Client
	if err := json.Unmarshal(items[0], &result); err != nil {
		return nil
	}
	return &result
}
