package console

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Service handles business logic for console operations
type Service struct {
	httpClient    *http.Client
	odfiBaseURL   string
	rdfiBaseURL   string
	ledgerBaseURL string
	eipBaseURL    string
}

// NewService creates a new console service
func NewService() *Service {
	return &Service{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		odfiBaseURL:   getEnv("ODFI_BASE_URL", "http://localhost:8081"),
		rdfiBaseURL:   getEnv("RDFI_BASE_URL", "http://localhost:8082"),
		ledgerBaseURL: getEnv("LEDGER_BASE_URL", "http://localhost:8083"),
		eipBaseURL:    getEnv("EIP_BASE_URL", "http://localhost:8084"),
	}
}

// ========== ODFI Operations ==========

// CreateODFIEntry creates an ODFI entry via the ODFI service
func (s *Service) CreateODFIEntry(ctx context.Context, req *CreateODFIEntryRequest) (*ODFIEntry, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.odfiBaseURL+"/api/v1/entries", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ODFI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var entry ODFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// ListODFIEntries lists ODFI entries with optional filters
func (s *Service) ListODFIEntries(ctx context.Context, status, traceNumber string) ([]*ODFIEntry, error) {
	queryParams := url.Values{}
	if status != "" {
		queryParams.Add("status", status)
	}
	if traceNumber != "" {
		queryParams.Add("trace_number", traceNumber)
	}

	url := fmt.Sprintf("%s/api/v1/entries?%s", s.odfiBaseURL, queryParams.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ODFI service returned status %d", resp.StatusCode)
	}

	var entries []*ODFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetODFIEntry gets a single ODFI entry by ID
func (s *Service) GetODFIEntry(ctx context.Context, id string) (*ODFIEntry, error) {
	url := fmt.Sprintf("%s/api/v1/entries/%s", s.odfiBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ODFI service returned status %d", resp.StatusCode)
	}

	var entry ODFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// UpdateODFIEntryStatus updates an ODFI entry status
func (s *Service) UpdateODFIEntryStatus(ctx context.Context, id, status string) (*ODFIEntry, error) {
	reqBody := UpdateODFIStatusRequest{Status: status}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/entries/%s/status", s.odfiBaseURL, id)
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ODFI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var entry ODFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// ========== RDFI Operations ==========

// CreateRDFIEntry creates an RDFI entry via the RDFI service
func (s *Service) CreateRDFIEntry(ctx context.Context, req *CreateRDFIEntryRequest) (*RDFIEntry, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.rdfiBaseURL+"/api/v1/entries", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("RDFI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var entry RDFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// ListRDFIEntries lists RDFI entries with optional filters
func (s *Service) ListRDFIEntries(ctx context.Context, status, traceNumber string) ([]*RDFIEntry, error) {
	queryParams := url.Values{}
	if status != "" {
		queryParams.Add("status", status)
	}
	if traceNumber != "" {
		queryParams.Add("trace_number", traceNumber)
	}

	url := fmt.Sprintf("%s/api/v1/entries?%s", s.rdfiBaseURL, queryParams.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RDFI service returned status %d", resp.StatusCode)
	}

	var entries []*RDFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetRDFIEntry gets a single RDFI entry by ID
func (s *Service) GetRDFIEntry(ctx context.Context, id string) (*RDFIEntry, error) {
	url := fmt.Sprintf("%s/api/v1/entries/%s", s.rdfiBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RDFI service returned status %d", resp.StatusCode)
	}

	var entry RDFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// ReturnEntry proxies a return request to the RDFI service
func (s *Service) ReturnEntry(ctx context.Context, id, reason string) (*RDFIEntry, error) {
	reqBody := ReturnRequest{Reason: reason}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/entries/%s/return", s.rdfiBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("RDFI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var entry RDFIEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// ========== Ledger Operations ==========

// CreateLedgerPosting creates a ledger posting
func (s *Service) CreateLedgerPosting(ctx context.Context, req *CreateLedgerPostingRequest) (*LedgerEntry, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.ledgerBaseURL+"/api/v1/postings", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ledger service returned status %d: %s", resp.StatusCode, string(body))
	}

	var entry LedgerEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

// ListLedgerPostings lists ledger postings with optional filters
func (s *Service) ListLedgerPostings(ctx context.Context, achSide, traceNumber string) ([]*LedgerEntry, error) {
	queryParams := url.Values{}
	if achSide != "" {
		queryParams.Add("ach_side", achSide)
	}
	if traceNumber != "" {
		queryParams.Add("trace_number", traceNumber)
	}

	url := fmt.Sprintf("%s/api/v1/postings?%s", s.ledgerBaseURL, queryParams.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ledger service returned status %d", resp.StatusCode)
	}

	var entries []*LedgerEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetBalances gets ledger balances
func (s *Service) GetBalances(ctx context.Context) (*BalanceResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.ledgerBaseURL+"/api/v1/balances", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ledger service returned status %d", resp.StatusCode)
	}

	var balances BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&balances); err != nil {
		return nil, err
	}

	return &balances, nil
}

// ========== EIP Operations ==========

// CreateEIPCase creates an EIP case
func (s *Service) CreateEIPCase(ctx context.Context, req *CreateEIPCaseRequest) (*EIPCase, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.eipBaseURL+"/api/v1/cases", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("EIP service returned status %d: %s", resp.StatusCode, string(body))
	}

	var eipCase EIPCase
	if err := json.NewDecoder(resp.Body).Decode(&eipCase); err != nil {
		return nil, err
	}

	return &eipCase, nil
}

// ListEIPCases lists EIP cases with optional filters
func (s *Service) ListEIPCases(ctx context.Context, status, side, traceNumber string) ([]*EIPCase, error) {
	queryParams := url.Values{}
	if status != "" {
		queryParams.Add("status", status)
	}
	if side != "" {
		queryParams.Add("side", side)
	}
	if traceNumber != "" {
		queryParams.Add("trace_number", traceNumber)
	}

	url := fmt.Sprintf("%s/api/v1/cases?%s", s.eipBaseURL, queryParams.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("EIP service returned status %d", resp.StatusCode)
	}

	var cases []*EIPCase
	if err := json.NewDecoder(resp.Body).Decode(&cases); err != nil {
		return nil, err
	}

	return cases, nil
}

// GetEIPCase gets a single EIP case by ID
func (s *Service) GetEIPCase(ctx context.Context, id string) (*EIPCase, error) {
	url := fmt.Sprintf("%s/api/v1/cases/%s", s.eipBaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("EIP service returned status %d", resp.StatusCode)
	}

	var eipCase EIPCase
	if err := json.NewDecoder(resp.Body).Decode(&eipCase); err != nil {
		return nil, err
	}

	return &eipCase, nil
}

// UpdateEIPCaseStatus updates an EIP case status
func (s *Service) UpdateEIPCaseStatus(ctx context.Context, id, status string) (*EIPCase, error) {
	reqBody := UpdateEIPCaseStatusRequest{Status: status}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/cases/%s/status", s.eipBaseURL, id)
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("EIP service returned status %d: %s", resp.StatusCode, string(body))
	}

	var eipCase EIPCase
	if err := json.NewDecoder(resp.Body).Decode(&eipCase); err != nil {
		return nil, err
	}

	return &eipCase, nil
}

// ========== Unified ACH Items (Fan-Out/Fan-In with Graceful Degradation) ==========

// serviceResult holds the result from a single service call
type serviceResult struct {
	serviceName string
	items       []*UnifiedAchItem
	err         error
	latency     time.Duration
}

// GetAchItems fetches and unifies entries from ODFI and RDFI services using fan-out/fan-in.
// If a service is unavailable, returns partial results from healthy services with degradation info.
func (s *Service) GetAchItems(ctx context.Context, side, status, traceNumber, sortBy, sortOrder string, limit, offset int) (*UnifiedAchResponse, error) {
	// Determine which services to query
	queryODFI := side == "" || strings.ToUpper(side) == "ODFI"
	queryRDFI := side == "" || strings.ToUpper(side) == "RDFI"

	// Channel to collect results - buffer for max expected services
	resultsChan := make(chan serviceResult, 2)

	// WaitGroup to track goroutines
	var wg sync.WaitGroup

	// Fan-out: Launch concurrent requests
	if queryODFI {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			items, err := s.fetchODFIEntries(ctx, status, traceNumber)
			resultsChan <- serviceResult{
				serviceName: "ODFI",
				items:       items,
				err:         err,
				latency:     time.Since(start),
			}
		}()
	}

	if queryRDFI {
		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			items, err := s.fetchRDFIEntries(ctx, status, traceNumber)
			resultsChan <- serviceResult{
				serviceName: "RDFI",
				items:       items,
				err:         err,
				latency:     time.Since(start),
			}
		}()
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Fan-in: Collect results as they arrive
	var allItems []*UnifiedAchItem
	var serviceInfo []ServiceHealth
	partial := false

	for result := range resultsChan {
		health := ServiceHealth{
			Service: result.serviceName,
			Latency: result.latency.Round(time.Millisecond).String(),
		}

		if result.err != nil {
			// Service failed - record degradation but continue
			health.Available = false
			health.Error = result.err.Error()
			partial = true
			fmt.Printf("[DEGRADED] %s service unavailable: %v (latency: %s)\n",
				result.serviceName, result.err, health.Latency)
		} else {
			// Service healthy - collect items
			health.Available = true
			allItems = append(allItems, result.items...)
			fmt.Printf("[OK] %s service returned %d items (latency: %s)\n",
				result.serviceName, len(result.items), health.Latency)
		}

		serviceInfo = append(serviceInfo, health)
	}

	// Sort combined results
	sortUnifiedAchItemsOptimized(allItems, sortBy, sortOrder)

	// Apply pagination
	totalCount := len(allItems)
	start := offset
	if start > totalCount {
		start = totalCount
	}

	end := start + limit
	if limit == 0 || end > totalCount {
		end = totalCount
	}

	return &UnifiedAchResponse{
		Items:       allItems[start:end],
		ServiceInfo: serviceInfo,
		Partial:     partial,
		TotalCount:  totalCount,
	}, nil
}

// GetAchItemsLegacy is the old synchronous version (deprecated)
func (s *Service) GetAchItemsLegacy(ctx context.Context, side, status, traceNumber, sortBy, sortOrder string, limit, offset int) ([]*UnifiedAchItem, error) {
	resp, err := s.GetAchItems(ctx, side, status, traceNumber, sortBy, sortOrder, limit, offset)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetAchItem fetches a single entry from the specified service
func (s *Service) GetAchItem(ctx context.Context, side, id string) (*UnifiedAchItem, error) {
	side = strings.ToUpper(side)

	switch side {
	case "ODFI":
		return s.fetchODFIEntry(ctx, id)
	case "RDFI":
		return s.fetchRDFIEntry(ctx, id)
	default:
		return nil, errors.New("invalid side: must be ODFI or RDFI")
	}
}

// fetchODFIEntries fetches entries from the ODFI service
func (s *Service) fetchODFIEntries(ctx context.Context, status, traceNumber string) ([]*UnifiedAchItem, error) {
	entries, err := s.ListODFIEntries(ctx, status, traceNumber)
	if err != nil {
		return nil, err
	}

	var items []*UnifiedAchItem
	for _, entry := range entries {
		items = append(items, &UnifiedAchItem{
			Side:        "ODFI",
			Source:      "odfi",
			EntryID:     entry.ID,
			TraceNumber: entry.TraceNumber,
			AmountCents: entry.AmountCents,
			Status:      entry.Status,
			CreatedAt:   entry.CreatedAt,
			Extra: map[string]interface{}{
				"company_name": entry.CompanyName,
				"sec_code":     entry.SecCode,
			},
		})
	}

	return items, nil
}

// fetchRDFIEntries fetches entries from the RDFI service
func (s *Service) fetchRDFIEntries(ctx context.Context, status, traceNumber string) ([]*UnifiedAchItem, error) {
	entries, err := s.ListRDFIEntries(ctx, status, traceNumber)
	if err != nil {
		return nil, err
	}

	var items []*UnifiedAchItem
	for _, entry := range entries {
		extra := map[string]interface{}{
			"receiver_name": entry.ReceiverName,
		}
		if entry.ReturnReason != "" {
			extra["return_reason"] = entry.ReturnReason
		}

		items = append(items, &UnifiedAchItem{
			Side:        "RDFI",
			Source:      "rdfi",
			EntryID:     entry.ID,
			TraceNumber: entry.TraceNumber,
			AmountCents: entry.AmountCents,
			Status:      entry.Status,
			CreatedAt:   entry.CreatedAt,
			Extra:       extra,
		})
	}

	return items, nil
}

// fetchODFIEntry fetches a single entry from ODFI service
func (s *Service) fetchODFIEntry(ctx context.Context, id string) (*UnifiedAchItem, error) {
	entry, err := s.GetODFIEntry(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &UnifiedAchItem{
		Side:        "ODFI",
		Source:      "odfi",
		EntryID:     entry.ID,
		TraceNumber: entry.TraceNumber,
		AmountCents: entry.AmountCents,
		Status:      entry.Status,
		CreatedAt:   entry.CreatedAt,
		Extra: map[string]interface{}{
			"company_name": entry.CompanyName,
			"sec_code":     entry.SecCode,
		},
	}, nil
}

// fetchRDFIEntry fetches a single entry from RDFI service
func (s *Service) fetchRDFIEntry(ctx context.Context, id string) (*UnifiedAchItem, error) {
	entry, err := s.GetRDFIEntry(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	extra := map[string]interface{}{
		"receiver_name": entry.ReceiverName,
	}
	if entry.ReturnReason != "" {
		extra["return_reason"] = entry.ReturnReason
	}

	return &UnifiedAchItem{
		Side:        "RDFI",
		Source:      "rdfi",
		EntryID:     entry.ID,
		TraceNumber: entry.TraceNumber,
		AmountCents: entry.AmountCents,
		Status:      entry.Status,
		CreatedAt:   entry.CreatedAt,
		Extra:       extra,
	}, nil
}

// sortUnifiedAchItemsOptimized sorts unified ACH items using sort.Slice (O(n log n))
func sortUnifiedAchItemsOptimized(items []*UnifiedAchItem, sortBy, sortOrder string) {
	// Default to created_at descending if not specified
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	ascending := sortOrder == "asc"

	sort.Slice(items, func(i, j int) bool {
		var less bool

		switch sortBy {
		case "created_at":
			less = items[i].CreatedAt < items[j].CreatedAt
		case "status":
			less = items[i].Status < items[j].Status
		case "amount", "amount_cents":
			less = items[i].AmountCents < items[j].AmountCents
		case "trace_number":
			less = items[i].TraceNumber < items[j].TraceNumber
		case "side":
			less = items[i].Side < items[j].Side
		default:
			less = items[i].CreatedAt < items[j].CreatedAt
		}

		// Flip for descending order
		if ascending {
			return less
		}
		return !less
	})
}

// sortUnifiedAchItems is deprecated - use sortUnifiedAchItemsOptimized
func sortUnifiedAchItems(items []*UnifiedAchItem, sortBy, sortOrder string) {
	sortUnifiedAchItemsOptimized(items, sortBy, sortOrder)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
