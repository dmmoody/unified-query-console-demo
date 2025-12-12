package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Configurable seed counts - over 1000 of each type for demo
const (
	ODFICount   = 1200 // ODFI entries (1000+)
	RDFICount   = 1200 // RDFI entries (1000+)
	LedgerCount = 800  // Ledger postings
	EIPCount    = 500  // EIP cases
	Workers     = 5    // Concurrent workers (reduced to prevent Docker resource exhaustion)
)

var (
	odfiURL   = "http://localhost:8081"
	rdfiURL   = "http://localhost:8082"
	ledgerURL = "http://localhost:8083"
	eipURL    = "http://localhost:8084"

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}
)

type ODFIEntry struct {
	ID          string `json:"id"`
	TraceNumber string `json:"trace_number"`
	CompanyName string `json:"company_name"`
	SecCode     string `json:"sec_code"`
	AmountCents int64  `json:"amount_cents"`
}

type RDFIEntry struct {
	ID           string `json:"id"`
	TraceNumber  string `json:"trace_number"`
	ReceiverName string `json:"receiver_name"`
	AmountCents  int64  `json:"amount_cents"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("🌱 ACH Concourse - High-Volume Database Seeding")
	fmt.Println("================================================")
	fmt.Printf("  Target: %d ODFI + %d RDFI + %d Ledger + %d EIP = %d records\n",
		ODFICount, RDFICount, LedgerCount, EIPCount,
		ODFICount+RDFICount+LedgerCount+EIPCount)
	fmt.Println("  Mode: Interleaved ODFI/RDFI for realistic timestamp distribution")
	fmt.Println()

	// Check if services are running
	if !checkHealth() {
		fmt.Println("❌ Services are not running. Please run 'make up' first.")
		return
	}

	fmt.Println("✅ Services are running")
	fmt.Println()

	start := time.Now()

	// Seed ODFI and RDFI interleaved (for realistic timestamp mixing)
	// Ledger and EIP can run concurrently alongside
	var wg sync.WaitGroup

	wg.Add(3)

	// Interleaved ODFI/RDFI seeding - timestamps will be mixed
	go func() {
		defer wg.Done()
		seedODFIAndRDFIInterleaved(ODFICount, RDFICount)
	}()

	// Ledger seeding (concurrent with above)
	go func() {
		defer wg.Done()
		seedLedgerConcurrent(LedgerCount)
	}()

	// EIP seeding (concurrent with above)
	go func() {
		defer wg.Done()
		seedEIPConcurrent(EIPCount)
	}()

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println()
	fmt.Println("🎉 Database seeding completed!")
	fmt.Println()
	fmt.Println("📊 Summary:")
	fmt.Printf("  ODFI entries:      %d records\n", ODFICount)
	fmt.Printf("  RDFI entries:      %d records\n", RDFICount)
	fmt.Printf("  Ledger postings:   %d records\n", LedgerCount)
	fmt.Printf("  EIP cases:         %d records\n", EIPCount)
	fmt.Println("  ─────────────────────────────")
	fmt.Printf("  Total:             %d records\n", ODFICount+RDFICount+LedgerCount+EIPCount)
	fmt.Printf("  Time elapsed:      %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  Throughput:        %.0f records/sec\n",
		float64(ODFICount+RDFICount+LedgerCount+EIPCount)/elapsed.Seconds())
	fmt.Println()
	fmt.Println("🔍 Verify unified query (fan-out/fan-in):")
	fmt.Println("  curl http://localhost:8080/api/v1/ach-items | jq '.total_count, .service_info'")
	fmt.Println("  curl http://localhost:8083/api/v1/balances | jq .")
}

func checkHealth() bool {
	services := []string{odfiURL, rdfiURL, ledgerURL, eipURL}
	for _, service := range services {
		resp, err := httpClient.Get(service + "/healthz")
		if err != nil || resp.StatusCode != 200 {
			return false
		}
		resp.Body.Close()
	}
	return true
}

// seedODFIAndRDFIInterleaved creates ODFI and RDFI entries in alternating order
// This ensures timestamps are mixed for realistic demo data when sorting by created_at
func seedODFIAndRDFIInterleaved(odfiCount, rdfiCount int) {
	fmt.Println("📝 Seeding ODFI and RDFI entries (interleaved for realistic timestamps)...")

	companyNames := []string{"ACME Corp", "TechStart Inc", "Global Traders", "MegaCorp LLC", "SmallBiz Co",
		"Enterprise Solutions", "Digital Payments", "FinTech Group", "Payment Solutions", "Commerce Partners",
		"Alpha Industries", "Beta Systems", "Gamma Technologies", "Delta Services", "Epsilon Holdings"}
	secCodes := []string{"PPD", "CCD", "WEB", "TEL"}
	odfiStatuses := []string{"PENDING", "PENDING", "SENT", "SENT", "SENT", "CANCELLED"}

	receiverNames := []string{"John Smith", "Jane Doe", "Robert Johnson", "Mary Williams", "James Brown",
		"Patricia Davis", "Michael Miller", "Linda Wilson", "David Moore", "Barbara Taylor",
		"Sarah Connor", "Tom Anderson", "Alice Chen", "Bob Martinez", "Carol White"}
	returnCodes := []string{"R01", "R02", "R03", "R04", "R10", "R16", "R20"}

	var odfiCreated, rdfiCreated int64
	maxCount := odfiCount
	if rdfiCount > maxCount {
		maxCount = rdfiCount
	}

	// Use a small worker pool for concurrent but ordered execution
	// We send pairs of jobs (ODFI then RDFI) to maintain interleaving
	type job struct {
		index int
		side  string // "ODFI" or "RDFI"
	}

	jobs := make(chan job, maxCount*2)
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < Workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if j.side == "ODFI" {
					createODFIEntry(j.index, companyNames, secCodes, odfiStatuses)
					atomic.AddInt64(&odfiCreated, 1)
				} else {
					createRDFIEntry(j.index, receiverNames, returnCodes)
					atomic.AddInt64(&rdfiCreated, 1)
				}
			}
		}()
	}

	// Send jobs in interleaved order: ODFI, RDFI, ODFI, RDFI...
	// This gives us mixed timestamps even with concurrent workers
	for i := 1; i <= maxCount; i++ {
		if i <= odfiCount {
			jobs <- job{index: i, side: "ODFI"}
		}
		if i <= rdfiCount {
			jobs <- job{index: i, side: "RDFI"}
		}

		// Progress every 100 pairs
		if i%100 == 0 {
			fmt.Printf("  Progress: %d/%d ODFI, %d/%d RDFI...\n",
				atomic.LoadInt64(&odfiCreated), odfiCount,
				atomic.LoadInt64(&rdfiCreated), rdfiCount)
		}
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("✅ ODFI entries created: %d\n", atomic.LoadInt64(&odfiCreated))
	fmt.Printf("✅ RDFI entries created: %d\n", atomic.LoadInt64(&rdfiCreated))
}

func createODFIEntry(i int, companyNames, secCodes, odfiStatuses []string) {
	traceNum := fmt.Sprintf("%015d", 1000000000000+i)

	entry := map[string]interface{}{
		"trace_number": traceNum,
		"company_name": companyNames[i%len(companyNames)],
		"sec_code":     secCodes[i%len(secCodes)],
		"amount_cents": rand.Int63n(100000) + 1000,
	}

	body, _ := json.Marshal(entry)
	resp, err := httpClient.Post(odfiURL+"/api/v1/entries", "application/json", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Update status if needed
	status := odfiStatuses[i%len(odfiStatuses)]
	if status != "PENDING" && resp.StatusCode == 201 {
		var createdEntry ODFIEntry
		json.NewDecoder(resp.Body).Decode(&createdEntry)
		updateReq := map[string]string{"status": status}
		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest("PATCH", odfiURL+"/api/v1/entries/"+createdEntry.ID+"/status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		httpClient.Do(req)
	}
}

func createRDFIEntry(i int, receiverNames, returnCodes []string) {
	traceNum := fmt.Sprintf("%015d", 2000000000000+i)

	entry := map[string]interface{}{
		"trace_number":  traceNum,
		"receiver_name": receiverNames[i%len(receiverNames)],
		"amount_cents":  rand.Int63n(80000) + 500,
	}

	body, _ := json.Marshal(entry)
	resp, err := httpClient.Post(rdfiURL+"/api/v1/entries", "application/json", bytes.NewReader(body))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Return ~16% of entries (every 6th)
	if i%6 == 0 && resp.StatusCode == 201 {
		var createdEntry RDFIEntry
		json.NewDecoder(resp.Body).Decode(&createdEntry)
		returnReq := map[string]string{"reason": returnCodes[i%len(returnCodes)]}
		body, _ := json.Marshal(returnReq)
		httpClient.Post(rdfiURL+"/api/v1/entries/"+createdEntry.ID+"/return", "application/json", bytes.NewReader(body))
	}
}

// seedLedgerConcurrent seeds ledger postings using concurrent workers
func seedLedgerConcurrent(count int) {
	fmt.Println("📝 Seeding Ledger postings...")

	descriptions := []string{"Payment processing", "Vendor payment", "Payroll deposit", "Invoice payment",
		"Refund processing", "Settlement transfer", "Account funding", "Bill payment",
		"Wire transfer", "ACH settlement", "Batch reconciliation", "Fee collection"}

	var created int64
	jobs := make(chan int, count)
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < Workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				var achSide, traceNum string
				if i%2 == 0 {
					achSide = "ODFI"
					traceNum = fmt.Sprintf("%015d", 1000000000000+(i/2))
				} else {
					achSide = "RDFI"
					traceNum = fmt.Sprintf("%015d", 2000000000000+((i+1)/2))
				}

				direction := "DEBIT"
				if i%2 == 1 {
					direction = "CREDIT"
				}

				posting := map[string]interface{}{
					"ach_side":     achSide,
					"trace_number": traceNum,
					"amount_cents": rand.Int63n(90000) + 1000,
					"direction":    direction,
					"description":  descriptions[i%len(descriptions)],
				}

				body, _ := json.Marshal(posting)
				resp, err := httpClient.Post(ledgerURL+"/api/v1/postings", "application/json", bytes.NewReader(body))
				if err == nil {
					resp.Body.Close()
					atomic.AddInt64(&created, 1)
				}
			}
		}()
	}

	// Send jobs
	for i := 1; i <= count; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("✅ Ledger postings created: %d\n", atomic.LoadInt64(&created))
}

// seedEIPConcurrent seeds EIP cases using concurrent workers
func seedEIPConcurrent(count int) {
	fmt.Println("📝 Seeding EIP cases...")

	types := []string{"RETURN_REVIEW", "NOC_REVIEW", "CUSTOMER_DISPUTE"}
	notesOptions := []string{
		"Customer called to dispute charge",
		"Return received from bank, needs review",
		"NOC received, account number correction needed",
		"Duplicate transaction reported",
		"Unauthorized transaction claim",
		"Amount discrepancy reported",
		"Timing issue with settlement",
		"Customer requested investigation",
		"Bank returned entry with R03",
		"Notification of change received",
		"Account holder claims fraud",
		"Missing authorization documentation",
	}
	statuses := []string{"OPEN", "OPEN", "OPEN", "IN_PROGRESS", "IN_PROGRESS", "RESOLVED"}

	var created int64
	jobs := make(chan int, count)
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < Workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				var side, traceNum string
				if i%2 == 0 {
					side = "ODFI"
					traceNum = fmt.Sprintf("%015d", 1000000000000+(i/2))
				} else {
					side = "RDFI"
					traceNum = fmt.Sprintf("%015d", 2000000000000+((i+1)/2))
				}

				caseData := map[string]interface{}{
					"side":         side,
					"trace_number": traceNum,
					"type":         types[i%len(types)],
					"notes":        notesOptions[i%len(notesOptions)],
				}

				body, _ := json.Marshal(caseData)
				resp, err := httpClient.Post(eipURL+"/api/v1/cases", "application/json", bytes.NewReader(body))
				if err == nil {
					// Update some case statuses
					status := statuses[i%len(statuses)]
					if status != "OPEN" && resp.StatusCode == 201 {
						var createdCase struct {
							ID string `json:"id"`
						}
						json.NewDecoder(resp.Body).Decode(&createdCase)
						updateReq := map[string]string{"status": status}
						body, _ := json.Marshal(updateReq)
						req, _ := http.NewRequest("PATCH", eipURL+"/api/v1/cases/"+createdCase.ID+"/status", bytes.NewReader(body))
						req.Header.Set("Content-Type", "application/json")
						httpClient.Do(req)
					}
					resp.Body.Close()
					atomic.AddInt64(&created, 1)
				}
			}
		}()
	}

	// Send jobs
	for i := 1; i <= count; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("✅ EIP cases created: %d\n", atomic.LoadInt64(&created))
}
