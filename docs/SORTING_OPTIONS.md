# Unified ACH Items - Flexible Sorting Options

## 🎯 New Sorting Capabilities

The unified ACH items endpoint now supports **flexible sorting** with query parameters!

### Query Parameters

| Parameter | Values | Default | Description |
|-----------|--------|---------|-------------|
| `sort_by` | `created_at`, `status`, `amount`, `trace_number`, `side` | `created_at` | Field to sort by |
| `sort_order` | `asc`, `desc` | `desc` | Sort direction |
| `side` | `ODFI`, `RDFI` | (all) | Filter by side |
| `status` | varies | (all) | Filter by status |
| `trace_number` | string | (all) | Filter by trace number |

## 📋 Postman Examples

### 1. Sort by Created Date (Most Recent First) - DEFAULT
```bash
GET http://localhost:8080/api/v1/ach-items
# or explicitly:
GET http://localhost:8080/api/v1/ach-items?sort_by=created_at&sort_order=desc
```

**Use Case:** Activity feed showing latest ACH entries

### 2. Sort by Created Date (Oldest First)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=created_at&sort_order=asc
```

**Use Case:** Historical view, processing queue (FIFO)

### 3. Sort by Amount (Highest First)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=amount&sort_order=desc
```

**Use Case:** Show largest transactions, prioritize high-value items

### 4. Sort by Amount (Lowest First)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=amount&sort_order=asc
```

**Use Case:** Show small transactions, micro-payment analysis

### 5. Sort by Status (Alphabetically)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=status&sort_order=asc
```

**Use Case:** Group by status (CANCELLED, PENDING, POSTED, RECEIVED, RETURNED, SENT)

### 6. Sort by Status (Reverse Alphabetical)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=status&sort_order=desc
```

**Use Case:** Show SENT/RETURNED items first

### 7. Sort by Trace Number (Ascending)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=trace_number&sort_order=asc
```

**Use Case:** Sequential processing, trace number order

### 8. Sort by Side (ODFI First, then RDFI)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=side&sort_order=asc
```

**Use Case:** Separate origination from receiving, grouped view

### 9. Sort by Side (RDFI First, then ODFI)
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=side&sort_order=desc
```

**Use Case:** Prioritize incoming transactions

## 🔥 Combined Filters + Sorting

### Show PENDING ODFI entries by amount (largest first)
```bash
GET http://localhost:8080/api/v1/ach-items?side=ODFI&status=PENDING&sort_by=amount&sort_order=desc
```

### Show all RDFI entries by date (oldest first)
```bash
GET http://localhost:8080/api/v1/ach-items?side=RDFI&sort_by=created_at&sort_order=asc
```

### Show RETURNED entries by amount (to identify large returns)
```bash
GET http://localhost:8080/api/v1/ach-items?status=RETURNED&sort_by=amount&sort_order=desc
```

### Show all entries grouped by status, then by amount
```bash
GET http://localhost:8080/api/v1/ach-items?sort_by=status&sort_order=asc
# Note: For multiple sort fields, you'd make a second request or handle client-side
```

## 📊 Demo Scenarios

### Scenario 1: Transaction Dashboard
**Goal:** Show recent activity across all ACH

```bash
# Latest entries (default)
curl "http://localhost:8080/api/v1/ach-items" | jq '.[] | {side, trace_number, amount_cents, status, created_at}'
```

### Scenario 2: High-Value Transaction Monitor
**Goal:** Identify large transactions for review

```bash
# Entries over $500 equivalent, sorted by amount
curl "http://localhost:8080/api/v1/ach-items?sort_by=amount&sort_order=desc" | \
  jq '[.[] | select(.amount_cents > 50000)] | .[] | {side, trace_number, amount: .amount_cents, status}'
```

### Scenario 3: Processing Queue
**Goal:** Process entries in FIFO order

```bash
# Oldest PENDING entries first
curl "http://localhost:8080/api/v1/ach-items?status=PENDING&sort_by=created_at&sort_order=asc" | \
  jq '.[] | {side, entry_id, trace_number, created_at}'
```

### Scenario 4: Returns Analysis
**Goal:** Analyze returned transactions

```bash
# All returns, largest first
curl "http://localhost:8080/api/v1/ach-items?status=RETURNED&sort_by=amount&sort_order=desc" | \
  jq '.[] | {side, trace_number, amount_cents, extra}'
```

### Scenario 5: Side-by-Side Comparison
**Goal:** Compare ODFI vs RDFI volumes

```bash
# Group by side
curl "http://localhost:8080/api/v1/ach-items?sort_by=side" | \
  jq 'group_by(.side) | map({side: .[0].side, count: length, total: (map(.amount_cents) | add)})'
```

## 🎨 Postman Collection Examples

### Collection Variables
```json
{
  "gateway_url": "http://localhost:8080",
  "sort_by": "created_at",
  "sort_order": "desc"
}
```

### Request Examples

#### 1. Latest Activity (Default)
**GET** `{{gateway_url}}/api/v1/ach-items`

#### 2. Highest Amounts
**GET** `{{gateway_url}}/api/v1/ach-items?sort_by=amount&sort_order=desc`

#### 3. Alphabetical by Status
**GET** `{{gateway_url}}/api/v1/ach-items?sort_by=status&sort_order=asc`

#### 4. ODFI by Amount
**GET** `{{gateway_url}}/api/v1/ach-items?side=ODFI&sort_by=amount&sort_order=desc`

#### 5. Processing Queue (FIFO)
**GET** `{{gateway_url}}/api/v1/ach-items?status=PENDING&sort_by=created_at&sort_order=asc`

## 🔍 Validation

The gateway validates sort parameters:

**Invalid sort_by:**
```bash
curl "http://localhost:8080/api/v1/ach-items?sort_by=invalid"
# Response: {"error":"sort_by must be one of: created_at, status, amount, trace_number, side"}
```

**Invalid sort_order:**
```bash
curl "http://localhost:8080/api/v1/ach-items?sort_order=invalid"
# Response: {"error":"sort_order must be 'asc' or 'desc'"}
```

## 💡 Best Practices

1. **Use defaults when possible** - Default `created_at desc` covers most cases
2. **Combine filters with sorting** - `?status=PENDING&sort_by=amount&sort_order=desc`
3. **Client-side secondary sort** - Gateway does primary sort, handle complex multi-field sorting client-side
4. **Pagination consideration** - For production, add pagination with sorting
5. **Performance** - Sorting happens in-memory after fetching; consider limits for large datasets

## 🚀 Quick Reference

| Sort Goal | Query String |
|-----------|-------------|
| Latest first (default) | `?sort_by=created_at&sort_order=desc` |
| Oldest first | `?sort_by=created_at&sort_order=asc` |
| Highest amounts | `?sort_by=amount&sort_order=desc` |
| Lowest amounts | `?sort_by=amount&sort_order=asc` |
| By status A-Z | `?sort_by=status&sort_order=asc` |
| By trace number | `?sort_by=trace_number&sort_order=asc` |
| ODFI first | `?sort_by=side&sort_order=asc` |
| RDFI first | `?sort_by=side&sort_order=desc` |

## 📝 Response Format

All sorted responses maintain the same structure:

```json
[
  {
    "side": "ODFI",
    "source": "odfi",
    "entry_id": "uuid",
    "trace_number": "1000000000000001",
    "amount_cents": 75000,
    "status": "SENT",
    "created_at": "2024-01-15T10:30:00Z",
    "extra": {
      "company_name": "ACME Corp",
      "sec_code": "PPD"
    }
  }
]
```

The order changes based on your sort parameters!

