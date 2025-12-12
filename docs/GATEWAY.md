# Console API Gateway - Complete API Reference

The Console service now acts as a complete **Unified API Gateway** for all ACH Concourse operations. All CRUD operations for every service can be performed through the Console at `http://localhost:8080`.

## 🎯 Why Use the Gateway?

✅ **Single Entry Point** - One API for all operations  
✅ **Simplified Client Code** - No need to know individual service URLs  
✅ **Consistent Interface** - Uniform error handling and responses  
✅ **Service Discovery** - Gateway handles routing to backend services  
✅ **Future-Proof** - Add authentication, rate limiting, etc. in one place  

---

## 📡 Gateway Endpoints

### Legacy Unified View (Backward Compatible)

#### GET /api/v1/ach-items
Query entries from both ODFI and RDFI in a unified format with **flexible sorting options**.

**Query Parameters:**
- `side` (optional): `ODFI` or `RDFI` - Filter by side
- `status` (optional): Filter by status
- `trace_number` (optional): Filter by trace number
- **`sort_by`** (optional): Field to sort by - `created_at` (default), `status`, `amount`, `trace_number`, `side`
- **`sort_order`** (optional): Sort direction - `desc` (default) or `asc`
- **`limit`** (optional): Number of results to return (default: 100, max: 1000)
- **`offset`** (optional): Number of results to skip (default: 0)

**Key Features:** 
- Results are **merged and sorted** (not just appended)
- ODFI and RDFI entries are interleaved based on sort criteria
- **Pagination support** for handling large datasets
- Multiple sort fields available for different use cases

⚠️ **Production Note:** Current implementation loads all data then paginates. See [PRODUCTION_CONSIDERATIONS.md](PRODUCTION_CONSIDERATIONS.md) for scaling strategies.

```bash
# Default: Most recent first, limit 100
curl http://localhost:8080/api/v1/ach-items

# First page (50 items)
curl "http://localhost:8080/api/v1/ach-items?limit=50&offset=0"

# Second page
curl "http://localhost:8080/api/v1/ach-items?limit=50&offset=50"

# Sort by amount (highest first), first 20
curl "http://localhost:8080/api/v1/ach-items?sort_by=amount&sort_order=desc&limit=20"

# Filter ODFI, sort by amount, paginate
curl "http://localhost:8080/api/v1/ach-items?side=ODFI&sort_by=amount&sort_order=desc&limit=25&offset=0"

# Oldest first (FIFO processing queue), small batches
curl "http://localhost:8080/api/v1/ach-items?status=PENDING&sort_by=created_at&sort_order=asc&limit=10"
```

**📖 See [SORTING_OPTIONS.md](SORTING_OPTIONS.md) for complete sorting documentation and [PRODUCTION_CONSIDERATIONS.md](PRODUCTION_CONSIDERATIONS.md) for scaling strategies.**

**Response includes all fields for sorting:**
```json
[
  {
    "side": "RDFI",
    "source": "rdfi",
    "entry_id": "uuid-1",
    "trace_number": "2000000000000001",
    "amount_cents": 50000,
    "status": "RECEIVED",
    "created_at": "2024-01-15T11:00:00Z",
    "extra": {"receiver_name": "John Doe"}
  }
]
```

#### GET /api/v1/ach-items/{side}/{id}
Get a single entry (unified format).

```bash
curl http://localhost:8080/api/v1/ach-items/ODFI/{id}
curl http://localhost:8080/api/v1/ach-items/RDFI/{id}
```

#### POST /api/v1/ach-items/RDFI/{id}/return
Return an RDFI entry (unified endpoint).

```bash
curl -X POST http://localhost:8080/api/v1/ach-items/RDFI/{id}/return \
  -H "Content-Type: application/json" \
  -d '{"reason": "R01"}'
```

---

## 🏦 ODFI Operations (via Gateway)

All ODFI operations available at `/api/v1/odfi/entries`

### POST /api/v1/odfi/entries
Create an ODFI entry through the gateway.

```bash
curl -X POST http://localhost:8080/api/v1/odfi/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "1234567890123456",
    "company_name": "ACME Corp",
    "sec_code": "PPD",
    "amount_cents": 50000
  }'
```

### GET /api/v1/odfi/entries
List all ODFI entries through the gateway.

```bash
curl http://localhost:8080/api/v1/odfi/entries
curl "http://localhost:8080/api/v1/odfi/entries?status=PENDING"
curl "http://localhost:8080/api/v1/odfi/entries?trace_number=1234567890123456"
```

### GET /api/v1/odfi/entries/{id}
Get a single ODFI entry by ID.

```bash
curl http://localhost:8080/api/v1/odfi/entries/{id}
```

### PATCH /api/v1/odfi/entries/{id}/status
Update ODFI entry status through the gateway.

```bash
curl -X PATCH http://localhost:8080/api/v1/odfi/entries/{id}/status \
  -H "Content-Type: application/json" \
  -d '{"status": "SENT"}'
```

Valid statuses: `PENDING`, `SENT`, `CANCELLED`

---

## 🏛️ RDFI Operations (via Gateway)

All RDFI operations available at `/api/v1/rdfi/entries`

### POST /api/v1/rdfi/entries
Create an RDFI entry through the gateway.

```bash
curl -X POST http://localhost:8080/api/v1/rdfi/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "9876543210987654",
    "receiver_name": "John Doe",
    "amount_cents": 25000
  }'
```

### GET /api/v1/rdfi/entries
List all RDFI entries through the gateway.

```bash
curl http://localhost:8080/api/v1/rdfi/entries
curl "http://localhost:8080/api/v1/rdfi/entries?status=RECEIVED"
curl "http://localhost:8080/api/v1/rdfi/entries?trace_number=9876543210987654"
```

### GET /api/v1/rdfi/entries/{id}
Get a single RDFI entry by ID.

```bash
curl http://localhost:8080/api/v1/rdfi/entries/{id}
```

### POST /api/v1/rdfi/entries/{id}/return
Return an RDFI entry through the gateway.

```bash
curl -X POST http://localhost:8080/api/v1/rdfi/entries/{id}/return \
  -H "Content-Type: application/json" \
  -d '{"reason": "R01"}'
```

Valid return codes: `R01`, `R02`, `R03`, `R04`, `R10`, etc.

---

## 📊 Ledger Operations (via Gateway)

All ledger operations available at `/api/v1/ledger`

### POST /api/v1/ledger/postings
Create a ledger posting through the gateway.

```bash
curl -X POST http://localhost:8080/api/v1/ledger/postings \
  -H "Content-Type: application/json" \
  -d '{
    "ach_side": "ODFI",
    "trace_number": "1234567890123456",
    "amount_cents": 50000,
    "direction": "DEBIT",
    "description": "Payment to vendor"
  }'
```

Valid values:
- `ach_side`: `ODFI` or `RDFI`
- `direction`: `DEBIT` or `CREDIT`

### GET /api/v1/ledger/postings
List all ledger postings through the gateway.

```bash
curl http://localhost:8080/api/v1/ledger/postings
curl "http://localhost:8080/api/v1/ledger/postings?ach_side=ODFI"
curl "http://localhost:8080/api/v1/ledger/postings?trace_number=1234567890123456"
```

### GET /api/v1/ledger/balances
Get ledger balances through the gateway.

```bash
curl http://localhost:8080/api/v1/ledger/balances
```

Response:
```json
{
  "total_debits": 1000000,
  "total_credits": 750000,
  "net_balance": -250000
}
```

---

## 🚨 EIP Operations (via Gateway)

All EIP operations available at `/api/v1/eip/cases`

### POST /api/v1/eip/cases
Create an EIP case through the gateway.

```bash
curl -X POST http://localhost:8080/api/v1/eip/cases \
  -H "Content-Type: application/json" \
  -d '{
    "side": "RDFI",
    "trace_number": "9876543210987654",
    "type": "CUSTOMER_DISPUTE",
    "notes": "Customer called to dispute charge"
  }'
```

Valid values:
- `side`: `ODFI` or `RDFI`
- `type`: `RETURN_REVIEW`, `NOC_REVIEW`, `CUSTOMER_DISPUTE`

### GET /api/v1/eip/cases
List all EIP cases through the gateway.

```bash
curl http://localhost:8080/api/v1/eip/cases
curl "http://localhost:8080/api/v1/eip/cases?status=OPEN"
curl "http://localhost:8080/api/v1/eip/cases?side=RDFI"
curl "http://localhost:8080/api/v1/eip/cases?trace_number=9876543210987654"
```

### GET /api/v1/eip/cases/{id}
Get a single EIP case by ID.

```bash
curl http://localhost:8080/api/v1/eip/cases/{id}
```

### PATCH /api/v1/eip/cases/{id}/status
Update EIP case status through the gateway.

```bash
curl -X PATCH http://localhost:8080/api/v1/eip/cases/{id}/status \
  -H "Content-Type: application/json" \
  -d '{"status": "IN_PROGRESS"}'
```

Valid statuses: `OPEN`, `IN_PROGRESS`, `RESOLVED`

---

## 🏥 Health Check

### GET /healthz
Check if the gateway is healthy.

```bash
curl http://localhost:8080/healthz
```

---

## 🎯 Complete Demo Flow (All via Gateway)

```bash
# 1. Create ODFI entry via gateway
curl -X POST http://localhost:8080/api/v1/odfi/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "DEMO111111111111",
    "company_name": "Demo Corp",
    "sec_code": "WEB",
    "amount_cents": 100000
  }' | jq . | tee /tmp/odfi.json

# 2. Extract ID and update status via gateway
ODFI_ID=$(cat /tmp/odfi.json | jq -r '.id')
curl -X PATCH "http://localhost:8080/api/v1/odfi/entries/$ODFI_ID/status" \
  -H "Content-Type: application/json" \
  -d '{"status": "SENT"}' | jq .

# 3. Create ledger posting via gateway
curl -X POST http://localhost:8080/api/v1/ledger/postings \
  -H "Content-Type: application/json" \
  -d '{
    "ach_side": "ODFI",
    "trace_number": "DEMO111111111111",
    "amount_cents": 100000,
    "direction": "DEBIT",
    "description": "Demo payment outgoing"
  }' | jq .

# 4. Create RDFI entry via gateway
curl -X POST http://localhost:8080/api/v1/rdfi/entries \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "DEMO999999999999",
    "receiver_name": "Demo User",
    "amount_cents": 50000
  }' | jq . | tee /tmp/rdfi.json

# 5. Return entry via gateway
RDFI_ID=$(cat /tmp/rdfi.json | jq -r '.id')
curl -X POST "http://localhost:8080/api/v1/rdfi/entries/$RDFI_ID/return" \
  -H "Content-Type: application/json" \
  -d '{"reason": "R01"}' | jq .

# 6. Create EIP case via gateway
curl -X POST http://localhost:8080/api/v1/eip/cases \
  -H "Content-Type: application/json" \
  -d '{
    "side": "RDFI",
    "trace_number": "DEMO999999999999",
    "type": "RETURN_REVIEW",
    "notes": "R01 return - insufficient funds"
  }' | jq .

# 7. Check balances via gateway
curl http://localhost:8080/api/v1/ledger/balances | jq .

# 8. Query everything via unified view
curl "http://localhost:8080/api/v1/ach-items?trace_number=DEMO" | jq .
```

---

## 📊 Endpoint Summary

| Service | Direct Port | Gateway Path | Operations |
|---------|-------------|--------------|------------|
| **Console** | 8080 | `/api/v1/ach-items` | Unified view (legacy) |
| **ODFI** | 8081 | `/api/v1/odfi/entries` | Create, List, Get, Update Status |
| **RDFI** | 8082 | `/api/v1/rdfi/entries` | Create, List, Get, Return |
| **Ledger** | 8083 | `/api/v1/ledger/*` | Create Posting, List, Balances |
| **EIP** | 8084 | `/api/v1/eip/cases` | Create, List, Get, Update Status |

**Total Gateway Endpoints: 22 endpoints** (all operations for all services!)

---

## 💡 Best Practices

### Use the Gateway for Production
```bash
# ✅ Good - Use gateway
curl http://localhost:8080/api/v1/odfi/entries

# ❌ Avoid - Direct service access
curl http://localhost:8081/api/v1/entries
```

### Direct Service Access for Debugging
```bash
# Use direct access only for debugging
curl http://localhost:8081/healthz
curl http://localhost:8082/api/v1/entries
```

### Error Handling
The gateway returns consistent error responses:
```json
{
  "error": "descriptive error message"
}
```

### Gateway Benefits
1. **Single authentication point** (when added)
2. **Centralized logging** (when added)
3. **Rate limiting** (when added)
4. **Service discovery** - clients don't need service URLs
5. **API versioning** - easy to version the gateway API
6. **Load balancing** - gateway can distribute requests

---

## 🔄 Migration Path

### Phase 1: Current (Both Supported)
- Direct service access: `http://localhost:808X/api/v1/*`
- Gateway access: `http://localhost:8080/api/v1/{service}/*`

### Phase 2: Gateway Preferred
- Encourage all clients to use gateway
- Direct access still works for debugging

### Phase 3: Gateway Only (Future)
- All client traffic goes through gateway
- Direct service access blocked (except internal)
- Services only accessible within Docker network

---

## 🎉 Summary

The Console Gateway now provides **complete CRUD operations** for all services:

✅ **15 ODFI & RDFI operations** (create, read, update, return)  
✅ **3 Ledger operations** (create postings, list, get balances)  
✅ **4 EIP operations** (create cases, read, update status)  
✅ **3 Legacy unified endpoints** (backward compatible)  
✅ **1 Health check**  

**Total: 26 endpoints** all accessible through a single unified gateway! 🚀

