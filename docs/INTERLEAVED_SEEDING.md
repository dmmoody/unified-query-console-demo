# Interleaved Timestamp Seeding

## Problem Solved

Previously, the seed script created all ODFI entries first, then all RDFI entries. This meant:
- All ODFI entries had earlier timestamps
- All RDFI entries had later timestamps
- When sorted by `created_at`, results were still grouped by service
- **Sorting didn't demonstrate the true mixed/interleaved capability**

## Solution

The seed script now creates entries in an **alternating pattern**:

```
ODFI #1  → delay → RDFI #1  → delay →
ODFI #2  → delay → RDFI #2  → delay →
ODFI #3  → delay → RDFI #3  → delay →
...
```

### Key Changes

1. **New function**: `seedODFIAndRDFIInterleaved()`
2. **Alternating creation**: Creates one ODFI, waits 50ms, creates one RDFI, waits 50ms
3. **True timestamp mixing**: ODFI and RDFI timestamps are genuinely interspersed

## Result

Now when you query with default sort (`created_at desc`):

```bash
curl http://localhost:8080/api/v1/ach-items | jq '.[] | {side, trace_number, created_at}'
```

You'll see:
```json
[
  {"side": "RDFI", "trace_number": "2000000000000075", "created_at": "2024-01-15T15:45:23Z"},
  {"side": "ODFI", "trace_number": "1000000000000075", "created_at": "2024-01-15T15:45:22Z"},
  {"side": "RDFI", "trace_number": "2000000000000074", "created_at": "2024-01-15T15:45:21Z"},
  {"side": "ODFI", "trace_number": "1000000000000074", "created_at": "2024-01-15T15:45:20Z"}
]
```

Notice how ODFI and RDFI entries are **truly mixed** based on timestamps!

## Demo Impact

### Before
- "Here's the unified view" → Shows all ODFI, then all RDFI (not impressive)
- Sorting didn't really demonstrate the feature

### After  
- "Here's the unified view" → Shows **interleaved** ODFI and RDFI entries!
- **Clearly demonstrates** that the gateway merges from multiple services
- Sorting by different fields shows the flexibility

## Testing

```bash
# Reseed with interleaved timestamps
make down && make up
make seed

# View the mixed results
curl http://localhost:8080/api/v1/ach-items | jq '.[:10] | .[] | {side, created_at}'

# You should see alternating ODFI/RDFI entries!
```

## Technical Details

- **50ms delay** between each entry creation
- **100ms total delay** per pair (ODFI + RDFI)
- **75 ODFI + 75 RDFI** = 150 total entries (vs 300 before)
- Reduced count to speed up seeding (~15 seconds vs ~30 seconds)
- **Same total demo value** but better timestamp distribution

## Additional Benefits

1. **Faster seeding** - 150 entries vs 300, with deliberate delays
2. **More realistic** - Real-world ACH systems don't batch all ODFI then all RDFI
3. **Better demo** - Clearly shows unified gateway capability
4. **Sorting flexibility** - All sort options now show true interleaving

This makes your demo **much more impressive**! 🎯

