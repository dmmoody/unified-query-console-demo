# Unified ACH Items - Sorted & Merged

## Problem Solved

Previously, the unified ACH items endpoint (`GET /api/v1/ach-items`) was simply appending ODFI results to RDFI results, resulting in entries grouped by service rather than sorted chronologically.

## Solution

The unified view now properly **merges and sorts** entries from both ODFI and RDFI services by their `created_at` timestamp (most recent first).

### What Changed

1. **Added `created_at` field** to `UnifiedAchItem` model
2. **Implemented sorting function** that merges entries chronologically
3. **Preserved timestamps** when mapping from service responses

### Example Output

#### Before (Grouped by Service)
```json
[
  {"side": "ODFI", "trace_number": "100...001", "created_at": "2024-01-15T10:00:00Z"},
  {"side": "ODFI", "trace_number": "100...002", "created_at": "2024-01-15T09:00:00Z"},
  {"side": "ODFI", "trace_number": "100...003", "created_at": "2024-01-15T08:00:00Z"},
  {"side": "RDFI", "trace_number": "200...001", "created_at": "2024-01-15T11:00:00Z"},
  {"side": "RDFI", "trace_number": "200...002", "created_at": "2024-01-15T07:00:00Z"},
  {"side": "RDFI", "trace_number": "200...003", "created_at": "2024-01-15T06:00:00Z"}
]
```

#### After (Sorted by Timestamp)
```json
[
  {"side": "RDFI", "trace_number": "200...001", "created_at": "2024-01-15T11:00:00Z"},
  {"side": "ODFI", "trace_number": "100...001", "created_at": "2024-01-15T10:00:00Z"},
  {"side": "ODFI", "trace_number": "100...002", "created_at": "2024-01-15T09:00:00Z"},
  {"side": "ODFI", "trace_number": "100...003", "created_at": "2024-01-15T08:00:00Z"},
  {"side": "RDFI", "trace_number": "200...002", "created_at": "2024-01-15T07:00:00Z"},
  {"side": "RDFI", "trace_number": "200...003", "created_at": "2024-01-15T06:00:00Z"}
]
```

Notice how entries are now **interleaved** based on their creation time, giving a true unified chronological view.

## Testing

```bash
# Get unified items (now sorted by timestamp)
curl http://localhost:8080/api/v1/ach-items | jq '.[] | {side, trace_number, created_at}'

# Should show most recent entries first, regardless of side
```

## Benefits

1. **Chronological View** - See ACH activity in the order it happened
2. **Better UX** - Users see most recent activity first
3. **True Unification** - ODFI and RDFI entries are truly merged, not just concatenated
4. **Dashboard Ready** - Perfect for activity feeds and timelines

## Implementation Details

The sorting is done in-memory after fetching from both services:

```go
func sortUnifiedAchItems(items []*UnifiedAchItem) {
    // Sort by created_at descending (most recent first)
    for i := 0; i < len(items); i++ {
        for j := i + 1; j < len(items); j++ {
            if items[i].CreatedAt < items[j].CreatedAt {
                items[i], items[j] = items[j], items[i]
            }
        }
    }
}
```

For production with larger datasets, consider:
- Using `sort.Slice` with proper time parsing
- Pagination with cursor-based navigation
- Database-level sorting if needed

