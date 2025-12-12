# ACH Concourse - Demo Quick Reference

## 🚀 Quick Start

```bash
# 1. Start all services
make up

# 2. Seed with 620 test records
make seed

# 3. Verify everything is working
make verify
```

## 📊 Demo Data Overview

After seeding, you'll have:
- **150 ODFI entries** (Origination - outgoing ACH)
  - ~50 PENDING
  - ~75 SENT
  - ~25 CANCELLED
- **150 RDFI entries** (Receiving - incoming ACH)
  - ~75 RECEIVED
  - ~50 POSTED
  - ~25 RETURNED (with return codes R01-R10)
- **200 Ledger postings**
  - 100 DEBITS
  - 100 CREDITS
  - Mix of ODFI and RDFI sides
- **120 EIP cases** (Exception tracking)
  - ~60 OPEN
  - ~40 IN_PROGRESS
  - ~20 RESOLVED

## 🎯 Demo Scenarios

### Scenario 1: Query All ACH Items (Unified View)

```bash
# Get all entries via Console (unified gateway)
curl http://localhost:8080/api/v1/ach-items | jq .

# Filter by side
curl http://localhost:8080/api/v1/ach-items?side=ODFI | jq .
curl http://localhost:8080/api/v1/ach-items?side=RDFI | jq .

# Filter by status
curl "http://localhost:8080/api/v1/ach-items?status=PENDING" | jq .
curl "http://localhost:8080/api/v1/ach-items?status=RETURNED" | jq .
```

### Scenario 2: ODFI Operations (Origination)

```bash
# List all ODFI entries
curl http://localhost:8081/api/v1/entries | jq .

# Filter by status
curl "http://localhost:8081/api/v1/entries?status=PENDING" | jq .

# Create new entry
curl -X POST http://localhost:8081/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "999999999999999",
    "company_name": "Demo Corp",
    "sec_code": "WEB",
    "amount_cents": 75000
  }' | jq .

# Update status (save the ID from above)
ENTRY_ID="paste-id-here"
curl -X PATCH "http://localhost:8081/api/v1/entries/$ENTRY_ID/status" \
  -H "Content-Type: application/json" \
  -d '{"status": "SENT"}' | jq .
```

### Scenario 3: RDFI Operations (Receiving)

```bash
# List all RDFI entries
curl http://localhost:8082/api/v1/entries | jq .

# Filter returned entries
curl "http://localhost:8082/api/v1/entries?status=RETURNED" | jq .

# Create new incoming entry
curl -X POST http://localhost:8082/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "888888888888888",
    "receiver_name": "Demo User",
    "amount_cents": 50000
  }' | jq .

# Return an entry (save the ID from above)
ENTRY_ID="paste-id-here"
curl -X POST "http://localhost:8082/api/v1/entries/$ENTRY_ID/return" \
  -H "Content-Type: application/json" \
  -d '{"reason": "R01"}' | jq .

# Use Console to return (unified gateway)
curl -X POST "http://localhost:8080/api/v1/ach-items/RDFI/$ENTRY_ID/return" \
  -H "Content-Type: application/json" \
  -d '{"reason": "R02"}' | jq .
```

### Scenario 4: Ledger & Balance Tracking

```bash
# View all postings
curl http://localhost:8083/api/v1/postings | jq .

# Filter by side
curl "http://localhost:8083/api/v1/postings?ach_side=ODFI" | jq .
curl "http://localhost:8083/api/v1/postings?ach_side=RDFI" | jq .

# Check balances (should show net of all debits/credits)
curl http://localhost:8083/api/v1/balances | jq .

# Create new posting
curl -X POST http://localhost:8083/api/v1/postings \
  -H "Content-Type: application/json" \
  -d '{
    "ach_side": "ODFI",
    "trace_number": "999999999999999",
    "amount_cents": 75000,
    "direction": "DEBIT",
    "description": "Demo payment"
  }' | jq .
```

### Scenario 5: Exception Management (EIP)

```bash
# List all cases
curl http://localhost:8084/api/v1/cases | jq .

# Filter by status
curl "http://localhost:8084/api/v1/cases?status=OPEN" | jq .
curl "http://localhost:8084/api/v1/cases?status=IN_PROGRESS" | jq .

# Filter by type
curl "http://localhost:8084/api/v1/cases?type=CUSTOMER_DISPUTE" | jq .

# Create new case
curl -X POST http://localhost:8084/api/v1/cases \
  -H "Content-Type: application/json" \
  -d '{
    "side": "RDFI",
    "trace_number": "888888888888888",
    "type": "CUSTOMER_DISPUTE",
    "notes": "Customer called to dispute - investigating"
  }' | jq .

# Update case status (save the ID from above)
CASE_ID="paste-id-here"
curl -X PATCH "http://localhost:8084/api/v1/cases/$CASE_ID/status" \
  -H "Content-Type: application/json" \
  -d '{"status": "IN_PROGRESS"}' | jq .
```

## 🔍 Common Queries

### Find entries by trace number

```bash
TRACE="1000000000000001"

# Via Console (searches both ODFI and RDFI)
curl "http://localhost:8080/api/v1/ach-items?trace_number=$TRACE" | jq .

# Via ODFI service
curl "http://localhost:8081/api/v1/entries?trace_number=$TRACE" | jq .

# Check if there's a case for it
curl "http://localhost:8084/api/v1/cases?trace_number=$TRACE" | jq .

# Check ledger postings
curl "http://localhost:8083/api/v1/postings?trace_number=$TRACE" | jq .
```

### Count records

```bash
# Count ODFI entries
curl -s http://localhost:8081/api/v1/entries | jq 'length'

# Count RDFI entries
curl -s http://localhost:8082/api/v1/entries | jq 'length'

# Count by status
curl -s "http://localhost:8081/api/v1/entries?status=PENDING" | jq 'length'
curl -s "http://localhost:8082/api/v1/entries?status=RETURNED" | jq 'length'
```

### Get specific entry details

```bash
# Get from Console (unified view)
curl "http://localhost:8080/api/v1/ach-items/ODFI/<entry-id>" | jq .
curl "http://localhost:8080/api/v1/ach-items/RDFI/<entry-id>" | jq .

# Get directly from service
curl "http://localhost:8081/api/v1/entries/<entry-id>" | jq .
curl "http://localhost:8082/api/v1/entries/<entry-id>" | jq .
```

## 📋 Status Values Reference

### ODFI Statuses
- `PENDING` - Entry created, awaiting transmission
- `SENT` - Entry transmitted to network
- `CANCELLED` - Entry cancelled before sending

### RDFI Statuses
- `RECEIVED` - Entry received from network
- `POSTED` - Entry posted to account
- `RETURNED` - Entry returned with reason code

### EIP Case Statuses
- `OPEN` - Case created, needs attention
- `IN_PROGRESS` - Case being investigated
- `RESOLVED` - Case closed/resolved

### Common Return Codes
- `R01` - Insufficient funds
- `R02` - Account closed
- `R03` - No account/unable to locate
- `R04` - Invalid account number
- `R10` - Customer advises unauthorized

## 🎬 Complete Demo Flow

```bash
# 1. Show unified view
curl http://localhost:8080/api/v1/ach-items | jq '.[0:5]'

# 2. Create ODFI entry
curl -X POST http://localhost:8081/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "DEMO123456789",
    "company_name": "Demo Corp",
    "sec_code": "PPD",
    "amount_cents": 100000
  }' | jq . | tee /tmp/odfi_entry.json

# 3. Extract ID and update status
ODFI_ID=$(cat /tmp/odfi_entry.json | jq -r '.id')
curl -X PATCH "http://localhost:8081/api/v1/entries/$ODFI_ID/status" \
  -H "Content-Type: application/json" \
  -d '{"status": "SENT"}' | jq .

# 4. Create ledger posting for it
curl -X POST http://localhost:8083/api/v1/postings \
  -H "Content-Type: application/json" \
  -d '{
    "ach_side": "ODFI",
    "trace_number": "DEMO123456789",
    "amount_cents": 100000,
    "direction": "DEBIT",
    "description": "Demo payment outgoing"
  }' | jq .

# 5. Create RDFI entry (receiving side)
curl -X POST http://localhost:8082/api/v1/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "DEMO987654321",
    "receiver_name": "John Doe",
    "amount_cents": 50000
  }' | jq . | tee /tmp/rdfi_entry.json

# 6. Return it via Console
RDFI_ID=$(cat /tmp/rdfi_entry.json | jq -r '.id')
curl -X POST "http://localhost:8080/api/v1/ach-items/RDFI/$RDFI_ID/return" \
  -H "Content-Type: application/json" \
  -d '{"reason": "R01"}' | jq .

# 7. Create EIP case for the return
curl -X POST http://localhost:8084/api/v1/cases \
  -H "Content-Type: application/json" \
  -d '{
    "side": "RDFI",
    "trace_number": "DEMO987654321",
    "type": "RETURN_REVIEW",
    "notes": "R01 return - needs review"
  }' | jq .

# 8. Check final balances
curl http://localhost:8083/api/v1/balances | jq .

# 9. Query everything via Console
curl "http://localhost:8080/api/v1/ach-items?trace_number=DEMO" | jq .
```

## 🛠️ Troubleshooting

```bash
# Check service health
curl http://localhost:8080/healthz  # Console
curl http://localhost:8081/healthz  # ODFI
curl http://localhost:8082/healthz  # RDFI
curl http://localhost:8083/healthz  # Ledger
curl http://localhost:8084/healthz  # EIP

# View logs
make logs
make logs-odfi
make logs-console

# Restart a service
docker-compose restart odfi
docker-compose restart console

# Clear data and reseed
make down
make up
make seed
```

## 📱 Using Postman

Import `postman_collection.json` for a complete collection with:
- All endpoints pre-configured
- Example request bodies
- Variables for IDs

Variables you can set:
- `odfi_entry_id`
- `rdfi_entry_id`
- `case_id`

