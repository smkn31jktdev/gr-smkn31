package util

import (
	"context"
	"fmt"
	"time"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 20    // AstraDB Data API v1 hard-limit: 20 dokumen per request
	MaxAPIPageSize  = 1000  // batas limit yang boleh diminta client via query param
	MaxFetchAll     = 5000  // batas atas total dokumen yang bisa di-fetch sekaligus
	MaxPages        = 500
	RequestTimeout  = 30 * time.Second // diperbesar untuk koleksi besar (607+ dokumen = ~31 request)
)

type PagedFetcher[T any] func(ctx context.Context, pageSize int, pageState string) (items []T, nextState string, err error)

// FetchAll 
func FetchAll[T any](ctx context.Context, fetcher PagedFetcher[T], pageSize int) ([]T, error) {
	if pageSize <= 0 || pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	pageState := ""
	result := make([]T, 0, 64)
	seenStates := map[string]struct{}{"": {}} // hindari loop tak terbatas

	for page := 0; page < MaxPages; page++ {
		reqCtx, cancel := context.WithTimeout(ctx, RequestTimeout)
		items, nextState, err := fetcher(reqCtx, pageSize, pageState)
		cancel()

		if err != nil {
			return nil, fmt.Errorf("fetch page %d: %w", page+1, err)
		}

		result = append(result, items...)

		if nextState == "" {
			break
		}
		if _, seen := seenStates[nextState]; seen {
			break
		}
		seenStates[nextState] = struct{}{}
		pageState = nextState

		if len(result) >= MaxFetchAll {
			break
		}
	}

	return result, nil
}

// Mengambil satu halaman saja
func FetchPage[T any](ctx context.Context, fetcher PagedFetcher[T], page, limit int) ([]T, bool, error) {
	if limit <= 0 || limit > MaxPageSize {
		limit = DefaultPageSize
	}
	if page <= 0 {
		page = 1
	}

	pageState := ""
	for skip := 1; skip < page; skip++ {
		reqCtx, cancel := context.WithTimeout(ctx, RequestTimeout)
		_, nextState, err := fetcher(reqCtx, limit, pageState)
		cancel()
		if err != nil {
			return nil, false, err
		}
		if nextState == "" {
			return []T{}, false, nil
		}
		pageState = nextState
	}

	reqCtx, cancel := context.WithTimeout(ctx, RequestTimeout)
	defer cancel()
	items, nextState, err := fetcher(reqCtx, limit, pageState)
	if err != nil {
		return nil, false, err
	}

	return items, nextState != "", nil
}

// WithRetry membungkus operasi dengan retry backoff linear
func WithRetry(attempts int, fn func() error) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		if lastErr = fn(); lastErr == nil {
			return nil
		}
		time.Sleep(time.Duration(i+1) * 300 * time.Millisecond)
	}
	return fmt.Errorf("after %d attempts: %w", attempts, lastErr)
}

// Melakukan pagination manual pada slice
func PaginateSlice[T any](items []T, page, limit int) ([]T, bool, int) {
	total := len(items)
	if page <= 0 {
		page = 1
	}

	// limit=0 → return semua data (mode "show all")
	if limit <= 0 {
		return items, false, total
	}

	start := (page - 1) * limit
	if start >= total {
		return []T{}, false, total
	}

	end := start + limit
	if end > total {
		end = total
	}

	return items[start:end], end < total, total
}
