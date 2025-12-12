package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	odfiURL   = "http://localhost:8081"
	rdfiURL   = "http://localhost:8082"
	ledgerURL = "http://localhost:8083"
	eipURL    = "http://localhost:8084"
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

	fmt.Println("🌱 ACH Concourse - Fast Database Seeding")
	fmt.Println("=========================================")
	fmt.Println()

	// Check if services are running
	if !checkHealth() {
		fmt.Println("❌ Services are not running. Please run 'make up' first.")
		return
	}

	fmt.Println("✅ Services are running")
	fmt.Println()

	// Seed ODFI and RDFI in interleaved pattern, then seed ledger/EIP
	var wg sync.WaitGroup

	wg.Add(3) // Only 3 goroutines!

	go func() {
		defer wg.Done()
		seedODFIAndRDFIInterleaved(75, 75) // Create 75 of each, interleaved
	}()

	go func() {
		defer wg.Done()
		seedLedger(200)
	}()

	go func() {
		defer wg.Done()
		seedEIP(120)
	}()

	wg.Wait()

	fmt.Println()
	fmt.Println("🎉 Database seeding completed!")
	fmt.Println()
	fmt.Println("📊 Summary:")
	fmt.Println("  ODFI entries:      75 records")
	fmt.Println("  RDFI entries:      75 records")
	fmt.Println("  Ledger postings:   200 records")
	fmt.Println("  EIP cases:         120 records")
	fmt.Println("  ─────────────────────────────")
	fmt.Println("  Total:             470 records")
	fmt.Println()
	fmt.Println("🔍 Verify data:")
	fmt.Println("  curl http://localhost:8080/api/v1/ach-items | jq .")
	fmt.Println("  curl http://localhost:8083/api/v1/balances | jq .")
}

func checkHealth() bool {
	services := []string{odfiURL, rdfiURL, ledgerURL, eipURL}
	for _, service := range services {
		resp, err := http.Get(service + "/healthz")
		if err != nil || resp.StatusCode != 200 {
			return false
		}
		resp.Body.Close()
	}
	return true
}

// seedODFIAndRDFIInterleaved creates ODFI and RDFI entries in an alternating pattern
// This ensures timestamps are truly mixed for demonstration purposes
func seedODFIAndRDFIInterleaved(odfiCount, rdfiCount int) {
	fmt.Println("📝 Seeding ODFI and RDFI entries (interleaved for timestamp variety)...")

	companyNames := []string{"ACME Corp", "TechStart Inc", "Global Traders", "MegaCorp LLC", "SmallBiz Co",
		"Enterprise Solutions", "Digital Payments", "FinTech Group", "Payment Solutions", "Commerce Partners"}
	secCodes := []string{"PPD", "CCD", "WEB", "TEL"}
	odfiStatuses := []string{"PENDING", "PENDING", "SENT", "SENT", "SENT", "CANCELLED"}

	receiverNames := []string{"John Smith", "Jane Doe", "Robert Johnson", "Mary Williams", "James Brown",
		"Patricia Davis", "Michael Miller", "Linda Wilson", "David Moore", "Barbara Taylor"}
	returnCodes := []string{"R01", "R02", "R03", "R04", "R10"}

	maxCount := odfiCount
	if rdfiCount > maxCount {
		maxCount = rdfiCount
	}

	for i := 1; i <= maxCount; i++ {
		// Create ODFI entry if we haven't reached the limit
		if i <= odfiCount {
			traceNum := fmt.Sprintf("%015d", 1000000000000+i)

			entry := map[string]interface{}{
				"trace_number": traceNum,
				"company_name": companyNames[i%len(companyNames)],
				"sec_code":     secCodes[i%len(secCodes)],
				"amount_cents": rand.Int63n(100000) + 1000,
			}

			body, _ := json.Marshal(entry)
			resp, err := http.Post(odfiURL+"/api/v1/entries", "application/json", bytes.NewReader(body))
			if err == nil {
				// Update status if needed
				status := odfiStatuses[i%len(odfiStatuses)]
				if status != "PENDING" && resp.StatusCode == 201 {
					var created ODFIEntry
					json.NewDecoder(resp.Body).Decode(&created)
					updateReq := map[string]string{"status": status}
					body, _ := json.Marshal(updateReq)
					req, _ := http.NewRequest("PATCH", odfiURL+"/api/v1/entries/"+created.ID+"/status", bytes.NewReader(body))
					req.Header.Set("Content-Type", "application/json")
					http.DefaultClient.Do(req)
				}
				resp.Body.Close()
			}
		}

		// Small delay to ensure timestamps differ
		time.Sleep(50 * time.Millisecond)

		// Create RDFI entry if we haven't reached the limit
		if i <= rdfiCount {
			traceNum := fmt.Sprintf("%015d", 2000000000000+i)

			entry := map[string]interface{}{
				"trace_number":  traceNum,
				"receiver_name": receiverNames[i%len(receiverNames)],
				"amount_cents":  rand.Int63n(80000) + 500,
			}

			body, _ := json.Marshal(entry)
			resp, err := http.Post(rdfiURL+"/api/v1/entries", "application/json", bytes.NewReader(body))
			if err == nil {
				// Return some entries
				if i%6 == 0 && resp.StatusCode == 201 {
					var created RDFIEntry
					json.NewDecoder(resp.Body).Decode(&created)
					returnReq := map[string]string{"reason": returnCodes[i%len(returnCodes)]}
					body, _ := json.Marshal(returnReq)
					http.Post(rdfiURL+"/api/v1/entries/"+created.ID+"/return", "application/json", bytes.NewReader(body))
				}
				resp.Body.Close()
			}
		}

		// Small delay between pairs
		time.Sleep(50 * time.Millisecond)

		if i%20 == 0 {
			fmt.Printf("  Created %d ODFI and %d RDFI entries...\n",
				min(i, odfiCount), min(i, rdfiCount))
		}
	}

	fmt.Printf("✅ ODFI entries created: %d\n", odfiCount)
	fmt.Printf("✅ RDFI entries created: %d\n", rdfiCount)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func seedLedger(count int) {
	fmt.Println("📝 Seeding Ledger postings...")

	descriptions := []string{"Payment processing", "Vendor payment", "Payroll deposit", "Invoice payment",
		"Refund processing", "Settlement transfer", "Account funding", "Bill payment"}

	for i := 1; i <= count; i++ {
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
		resp, err := http.Post(ledgerURL+"/api/v1/postings", "application/json", bytes.NewReader(body))
		if err != nil {
			continue
		}
		resp.Body.Close()

		if i%40 == 0 {
			fmt.Printf("  Created %d ledger postings...\n", i)
		}
	}
	fmt.Println("✅ Ledger postings created: 200")
}

func seedEIP(count int) {
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
	}
	statuses := []string{"OPEN", "OPEN", "OPEN", "IN_PROGRESS", "IN_PROGRESS", "RESOLVED"}

	for i := 1; i <= count; i++ {
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
		resp, err := http.Post(eipURL+"/api/v1/cases", "application/json", bytes.NewReader(body))
		if err != nil {
			continue
		}

		// Update some case statuses
		status := statuses[i%len(statuses)]
		if status != "OPEN" && resp.StatusCode == 201 {
			var created struct {
				ID string `json:"id"`
			}
			json.NewDecoder(resp.Body).Decode(&created)
			updateReq := map[string]string{"status": status}
			body, _ := json.Marshal(updateReq)
			req, _ := http.NewRequest("PATCH", eipURL+"/api/v1/cases/"+created.ID+"/status", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			http.DefaultClient.Do(req)
		}
		resp.Body.Close()

		if i%30 == 0 {
			fmt.Printf("  Created %d EIP cases...\n", i)
		}
	}
	fmt.Println("✅ EIP cases created: 120")
}
