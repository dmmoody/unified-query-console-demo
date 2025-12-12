# Pagination Guide

## Overview

The unified gateway now supports pagination to handle large datasets efficiently. This document explains how to use pagination effectively.

## Query Parameters

| Parameter | Type | Default | Max | Description |
|-----------|------|---------|-----|-------------|
| `limit` | integer | 100 | 1000 | Number of items to return per page |
| `offset` | integer | 0 | - | Number of items to skip |

## Examples

### Basic Pagination

**First page (20 items):**
```bash
curl "http://localhost:8080/api/v1/ach-items?limit=20&offset=0"
```

**Second page:**
```bash
curl "http://localhost:8080/api/v1/ach-items?limit=20&offset=20"
```

**Third page:**
```bash
curl "http://localhost:8080/api/v1/ach-items?limit=20&offset=40"
```

### Combining with Sorting

**Top 10 highest value transactions:**
```bash
curl "http://localhost:8080/api/v1/ach-items?sort_by=amount&sort_order=desc&limit=10"
```

**Most recent 50 ODFI entries:**
```bash
curl "http://localhost:8080/api/v1/ach-items?side=ODFI&sort_by=created_at&sort_order=desc&limit=50"
```

### Processing Queue Pattern

**Get next batch of PENDING items (FIFO):**
```bash
# First batch
curl "http://localhost:8080/api/v1/ach-items?status=PENDING&sort_by=created_at&sort_order=asc&limit=10&offset=0"

# Second batch
curl "http://localhost:8080/api/v1/ach-items?status=PENDING&sort_by=created_at&sort_order=asc&limit=10&offset=10"
```

### Large Datasets

**Use the default page size (100) for balanced performance:**
```bash
curl "http://localhost:8080/api/v1/ach-items?limit=100"
```

**For reporting/exports, paginate through all data:**
```bash
# Page 1: items 0-99
curl "http://localhost:8080/api/v1/ach-items?limit=100&offset=0"

# Page 2: items 100-199
curl "http://localhost:8080/api/v1/ach-items?limit=100&offset=100"

# Page 3: items 200-299
curl "http://localhost:8080/api/v1/ach-items?limit=100&offset=200"
```

## Response Format

The response is a standard JSON array, regardless of pagination:

```json
[
  {
    "side": "ODFI",
    "source": "odfi",
    "entry_id": "uuid-1",
    "trace_number": "1000000000000001",
    "amount_cents": 50000,
    "status": "SENT",
    "created_at": "2024-01-15T10:00:00Z",
    "extra": {"company_name": "ACME Corp", "sec_code": "PPD"}
  },
  {
    "side": "RDFI",
    "source": "rdfi",
    "entry_id": "uuid-2",
    "trace_number": "2000000000000001",
    "amount_cents": 25000,
    "status": "RECEIVED",
    "created_at": "2024-01-15T09:55:00Z",
    "extra": {"receiver_name": "Jane Smith"}
  }
]
```

## Best Practices

### 1. Choose Appropriate Page Sizes

| Use Case | Recommended Limit |
|----------|-------------------|
| UI list views | 20-50 |
| API consumers | 100 (default) |
| Background jobs | 100-500 |
| High-frequency polling | 10-20 |

### 2. Always Use Pagination for Production

❌ **Don't do this:**
```bash
# Unbounded query - will get ALL records (could be millions!)
curl "http://localhost:8080/api/v1/ach-items"
```

✅ **Do this:**
```bash
# Explicit limit - safe for production
curl "http://localhost:8080/api/v1/ach-items?limit=100"
```

### 3. Combine Filters to Reduce Result Sets

```bash
# Better: Filter first, then paginate
curl "http://localhost:8080/api/v1/ach-items?side=ODFI&status=PENDING&limit=50"

# vs fetching everything and filtering client-side
```

### 4. Implement "Load More" UI Pattern

For web/mobile UIs, implement infinite scroll or "Load More" buttons:

```javascript
let offset = 0;
const limit = 20;

async function loadMore() {
  const response = await fetch(
    `http://localhost:8080/api/v1/ach-items?limit=${limit}&offset=${offset}`
  );
  const items = await response.json();
  
  // Append to UI
  renderItems(items);
  
  // Prepare for next page
  offset += limit;
  
  // Hide "Load More" if less than limit returned (end of data)
  if (items.length < limit) {
    hideLoadMoreButton();
  }
}
```

## Implementation Notes

### Current Implementation (POC)

- Fetches all records from both services
- Merges and sorts in memory
- **Then** applies pagination

### Production Considerations

For large datasets, consider:

1. **Push pagination down to services**
   - Each service handles its own pagination
   - Gateway merges only requested page ranges

2. **Cursor-based pagination** (instead of offset)
   - More efficient for large datasets
   - Prevents issues with concurrent inserts

3. **Aggregation database**
   - Materialized view of unified data
   - Direct pagination at DB level

See [PRODUCTION_CONSIDERATIONS.md](PRODUCTION_CONSIDERATIONS.md) for detailed scaling strategies.

## Postman Collection

The Postman collection includes 5 pagination examples:

1. **First Page (20 items)** - Basic pagination
2. **Second Page (20 items)** - Page navigation
3. **Top 10 by Amount** - Combining sort + pagination
4. **Small Batches for Processing (FIFO)** - Queue pattern
5. **Large Page (100 items)** - Default page size

Import `postman_collection.json` to try them out!

## Troubleshooting

### Empty Results

If you get `[]`, you may have paginated past the end of the dataset:

```bash
# Check total count first (without pagination)
curl "http://localhost:8080/api/v1/ach-items" | jq length

# 150 items total

# This will return empty (150 < 200)
curl "http://localhost:8080/api/v1/ach-items?offset=200"
```

### Performance Issues

If queries are slow:

1. Check your `limit` - smaller is faster
2. Add filters (`side`, `status`, etc.) to reduce data fetched
3. Consider if you need to paginate at the service level (see PRODUCTION_CONSIDERATIONS.md)

### Max Limit Exceeded

Requests with `limit > 1000` will be clamped to 1000:

```bash
# This will return max 1000 items, not 5000
curl "http://localhost:8080/api/v1/ach-items?limit=5000"
```

## See Also

- [SORTING_OPTIONS.md](SORTING_OPTIONS.md) - Flexible sorting guide
- [GATEWAY.md](GATEWAY.md) - Complete gateway API reference
- [PRODUCTION_CONSIDERATIONS.md](PRODUCTION_CONSIDERATIONS.md) - Scaling strategies
- [README.md](../README.md) - Main project documentation

