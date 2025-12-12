# ACH Concourse - Curl Commands for Postman Import

## Console Service (Port 8080)

### Get All ACH Items
```bash
curl --location 'http://localhost:8080/api/v1/ach-items'
```

### Get All ACH Items - Filter by ODFI
```bash
curl --location 'http://localhost:8080/api/v1/ach-items?side=ODFI'
```

### Get All ACH Items - Filter by RDFI
```bash
curl --location 'http://localhost:8080/api/v1/ach-items?side=RDFI'
```

### Get All ACH Items - Filter by Status
```bash
curl --location 'http://localhost:8080/api/v1/ach-items?status=PENDING'
```

### Get All ACH Items - Filter by Trace Number
```bash
curl --location 'http://localhost:8080/api/v1/ach-items?trace_number=1000000000000001'
```

### Get Single ACH Item (ODFI)
```bash
curl --location 'http://localhost:8080/api/v1/ach-items/ODFI/{{odfi_entry_id}}'
```

### Get Single ACH Item (RDFI)
```bash
curl --location 'http://localhost:8080/api/v1/ach-items/RDFI/{{rdfi_entry_id}}'
```

### Return RDFI Entry via Console
```bash
curl --location 'http://localhost:8080/api/v1/ach-items/RDFI/{{rdfi_entry_id}}/return' \
--header 'Content-Type: application/json' \
--data '{
  "reason": "R01"
}'
```

### Console Health Check
```bash
curl --location 'http://localhost:8080/healthz'
```

---

## ODFI Service (Port 8081)

### Create ODFI Entry
```bash
curl --location 'http://localhost:8081/api/v1/entries' \
--header 'Content-Type: application/json' \
--data '{
  "trace_number": "9999999999999999",
  "company_name": "ACME Corp",
  "sec_code": "PPD",
  "amount_cents": 50000
}'
```

### List All ODFI Entries
```bash
curl --location 'http://localhost:8081/api/v1/entries'
```

### List ODFI Entries - Filter by Status
```bash
curl --location 'http://localhost:8081/api/v1/entries?status=PENDING'
```

### List ODFI Entries - Filter by Trace Number
```bash
curl --location 'http://localhost:8081/api/v1/entries?trace_number=1000000000000001'
```

### Get ODFI Entry by ID
```bash
curl --location 'http://localhost:8081/api/v1/entries/{{odfi_entry_id}}'
```

### Update ODFI Entry Status to SENT
```bash
curl --location --request PATCH 'http://localhost:8081/api/v1/entries/{{odfi_entry_id}}/status' \
--header 'Content-Type: application/json' \
--data '{
  "status": "SENT"
}'
```

### Update ODFI Entry Status to CANCELLED
```bash
curl --location --request PATCH 'http://localhost:8081/api/v1/entries/{{odfi_entry_id}}/status' \
--header 'Content-Type: application/json' \
--data '{
  "status": "CANCELLED"
}'
```

### ODFI Health Check
```bash
curl --location 'http://localhost:8081/healthz'
```

---

## RDFI Service (Port 8082)

### Create RDFI Entry
```bash
curl --location 'http://localhost:8082/api/v1/entries' \
--header 'Content-Type: application/json' \
--data '{
  "trace_number": "8888888888888888",
  "receiver_name": "John Smith",
  "amount_cents": 25000
}'
```

### List All RDFI Entries
```bash
curl --location 'http://localhost:8082/api/v1/entries'
```

### List RDFI Entries - Filter by Status RECEIVED
```bash
curl --location 'http://localhost:8082/api/v1/entries?status=RECEIVED'
```

### List RDFI Entries - Filter by Status RETURNED
```bash
curl --location 'http://localhost:8082/api/v1/entries?status=RETURNED'
```

### List RDFI Entries - Filter by Trace Number
```bash
curl --location 'http://localhost:8082/api/v1/entries?trace_number=2000000000000001'
```

### Get RDFI Entry by ID
```bash
curl --location 'http://localhost:8082/api/v1/entries/{{rdfi_entry_id}}'
```

### Return RDFI Entry - Insufficient Funds (R01)
```bash
curl --location 'http://localhost:8082/api/v1/entries/{{rdfi_entry_id}}/return' \
--header 'Content-Type: application/json' \
--data '{
  "reason": "R01"
}'
```

### Return RDFI Entry - Account Closed (R02)
```bash
curl --location 'http://localhost:8082/api/v1/entries/{{rdfi_entry_id}}/return' \
--header 'Content-Type: application/json' \
--data '{
  "reason": "R02"
}'
```

### Return RDFI Entry - No Account (R03)
```bash
curl --location 'http://localhost:8082/api/v1/entries/{{rdfi_entry_id}}/return' \
--header 'Content-Type: application/json' \
--data '{
  "reason": "R03"
}'
```

### Return RDFI Entry - Invalid Account (R04)
```bash
curl --location 'http://localhost:8082/api/v1/entries/{{rdfi_entry_id}}/return' \
--header 'Content-Type: application/json' \
--data '{
  "reason": "R04"
}'
```

### Return RDFI Entry - Unauthorized (R10)
```bash
curl --location 'http://localhost:8082/api/v1/entries/{{rdfi_entry_id}}/return' \
--header 'Content-Type: application/json' \
--data '{
  "reason": "R10"
}'
```

### RDFI Health Check
```bash
curl --location 'http://localhost:8082/healthz'
```

---

## Ledger Service (Port 8083)

### Create Ledger Posting - ODFI Debit
```bash
curl --location 'http://localhost:8083/api/v1/postings' \
--header 'Content-Type: application/json' \
--data '{
  "ach_side": "ODFI",
  "trace_number": "1000000000000001",
  "amount_cents": 50000,
  "direction": "DEBIT",
  "description": "Payment to vendor"
}'
```

### Create Ledger Posting - RDFI Credit
```bash
curl --location 'http://localhost:8083/api/v1/postings' \
--header 'Content-Type: application/json' \
--data '{
  "ach_side": "RDFI",
  "trace_number": "2000000000000001",
  "amount_cents": 25000,
  "direction": "CREDIT",
  "description": "Payment received"
}'
```

### List All Ledger Postings
```bash
curl --location 'http://localhost:8083/api/v1/postings'
```

### List Ledger Postings - Filter by ODFI Side
```bash
curl --location 'http://localhost:8083/api/v1/postings?ach_side=ODFI'
```

### List Ledger Postings - Filter by RDFI Side
```bash
curl --location 'http://localhost:8083/api/v1/postings?ach_side=RDFI'
```

### List Ledger Postings - Filter by Trace Number
```bash
curl --location 'http://localhost:8083/api/v1/postings?trace_number=1000000000000001'
```

### Get Ledger Balances
```bash
curl --location 'http://localhost:8083/api/v1/balances'
```

### Ledger Health Check
```bash
curl --location 'http://localhost:8083/healthz'
```

---

## EIP Service (Port 8084)

### Create EIP Case - Customer Dispute
```bash
curl --location 'http://localhost:8084/api/v1/cases' \
--header 'Content-Type: application/json' \
--data '{
  "side": "RDFI",
  "trace_number": "2000000000000001",
  "type": "CUSTOMER_DISPUTE",
  "notes": "Customer called to dispute charge"
}'
```

### Create EIP Case - Return Review
```bash
curl --location 'http://localhost:8084/api/v1/cases' \
--header 'Content-Type: application/json' \
--data '{
  "side": "RDFI",
  "trace_number": "2000000000000002",
  "type": "RETURN_REVIEW",
  "notes": "Return received from bank, needs review"
}'
```

### Create EIP Case - NOC Review
```bash
curl --location 'http://localhost:8084/api/v1/cases' \
--header 'Content-Type: application/json' \
--data '{
  "side": "ODFI",
  "trace_number": "1000000000000001",
  "type": "NOC_REVIEW",
  "notes": "NOC received, account number correction needed"
}'
```

### List All EIP Cases
```bash
curl --location 'http://localhost:8084/api/v1/cases'
```

### List EIP Cases - Filter by Status OPEN
```bash
curl --location 'http://localhost:8084/api/v1/cases?status=OPEN'
```

### List EIP Cases - Filter by Status IN_PROGRESS
```bash
curl --location 'http://localhost:8084/api/v1/cases?status=IN_PROGRESS'
```

### List EIP Cases - Filter by Status RESOLVED
```bash
curl --location 'http://localhost:8084/api/v1/cases?status=RESOLVED'
```

### List EIP Cases - Filter by Side
```bash
curl --location 'http://localhost:8084/api/v1/cases?side=RDFI'
```

### List EIP Cases - Filter by Trace Number
```bash
curl --location 'http://localhost:8084/api/v1/cases?trace_number=2000000000000001'
```

### Get EIP Case by ID
```bash
curl --location 'http://localhost:8084/api/v1/cases/{{case_id}}'
```

### Update EIP Case Status to IN_PROGRESS
```bash
curl --location --request PATCH 'http://localhost:8084/api/v1/cases/{{case_id}}/status' \
--header 'Content-Type: application/json' \
--data '{
  "status": "IN_PROGRESS"
}'
```

### Update EIP Case Status to RESOLVED
```bash
curl --location --request PATCH 'http://localhost:8084/api/v1/cases/{{case_id}}/status' \
--header 'Content-Type: application/json' \
--data '{
  "status": "RESOLVED"
}'
```

### EIP Health Check
```bash
curl --location 'http://localhost:8084/healthz'
```

---

## Complete Demo Workflow

### Step 1: Create ODFI Entry
```bash
curl --location 'http://localhost:8081/api/v1/entries' \
--header 'Content-Type: application/json' \
--data '{
  "trace_number": "DEMO123456789",
  "company_name": "Demo Corporation",
  "sec_code": "WEB",
  "amount_cents": 100000
}'
```

### Step 2: Get Entry ID and Update Status
```bash
# First get the entry
curl --location 'http://localhost:8081/api/v1/entries?trace_number=DEMO123456789'

# Then update status (replace {{odfi_entry_id}} with actual ID)
curl --location --request PATCH 'http://localhost:8081/api/v1/entries/{{odfi_entry_id}}/status' \
--header 'Content-Type: application/json' \
--data '{
  "status": "SENT"
}'
```

### Step 3: Create Ledger Posting
```bash
curl --location 'http://localhost:8083/api/v1/postings' \
--header 'Content-Type: application/json' \
--data '{
  "ach_side": "ODFI",
  "trace_number": "DEMO123456789",
  "amount_cents": 100000,
  "direction": "DEBIT",
  "description": "Demo payment outgoing"
}'
```

### Step 4: Create RDFI Entry
```bash
curl --location 'http://localhost:8082/api/v1/entries' \
--header 'Content-Type: application/json' \
--data '{
  "trace_number": "DEMO987654321",
  "receiver_name": "John Demo User",
  "amount_cents": 50000
}'
```

### Step 5: Return Entry via Console
```bash
# First get the entry
curl --location 'http://localhost:8082/api/v1/entries?trace_number=DEMO987654321'

# Then return it via Console (replace {{rdfi_entry_id}} with actual ID)
curl --location 'http://localhost:8080/api/v1/ach-items/RDFI/{{rdfi_entry_id}}/return' \
--header 'Content-Type: application/json' \
--data '{
  "reason": "R01"
}'
```

### Step 6: Create EIP Case for Return
```bash
curl --location 'http://localhost:8084/api/v1/cases' \
--header 'Content-Type: application/json' \
--data '{
  "side": "RDFI",
  "trace_number": "DEMO987654321",
  "type": "RETURN_REVIEW",
  "notes": "R01 return received - insufficient funds - needs review"
}'
```

### Step 7: Query Unified View
```bash
curl --location 'http://localhost:8080/api/v1/ach-items'
```

### Step 8: Check Balances
```bash
curl --location 'http://localhost:8083/api/v1/balances'
```

---

## How to Import into Postman

### Option 1: Import Individual Curls
1. Copy any curl command above
2. In Postman, click **Import** button
3. Select **Raw text** tab
4. Paste the curl command
5. Click **Continue** then **Import**

### Option 2: Import the JSON Collection
1. In Postman, click **Import** button
2. Select **File** tab
3. Choose `postman_collection.json` from the project root
4. Click **Import**

### Option 3: Use Postman Variables
After importing, set these variables in Postman:
- `odfi_entry_id` - Get from any ODFI entry creation
- `rdfi_entry_id` - Get from any RDFI entry creation
- `case_id` - Get from any EIP case creation

To set variables:
1. Click on your collection
2. Go to **Variables** tab
3. Add the variable name and current value
4. Save

---

## Testing with Seeded Data

After running `make seed`, you'll have data with these trace numbers:

**ODFI Trace Numbers:** 1000000000000001 through 1000000000000150
**RDFI Trace Numbers:** 2000000000000001 through 2000000000000150

Example queries:
```bash
# Get a specific ODFI entry
curl --location 'http://localhost:8081/api/v1/entries?trace_number=1000000000000001'

# Get a specific RDFI entry
curl --location 'http://localhost:8082/api/v1/entries?trace_number=2000000000000001'

# Query via Console
curl --location 'http://localhost:8080/api/v1/ach-items?trace_number=1000000000000001'
```

