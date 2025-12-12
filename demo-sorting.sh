#!/bin/bash

# Unified ACH Items - Sorting Demo
# Demonstrates all available sorting options

set -e

GATEWAY="http://localhost:8080"

echo "🎯 Unified ACH Items - Sorting Options Demo"
echo "============================================"
echo ""

# Check if gateway is running
if ! curl -sf "$GATEWAY/healthz" > /dev/null 2>&1; then
    echo "❌ Gateway is not running. Please run 'make up' first."
    exit 1
fi

echo "✅ Gateway is healthy"
echo ""

# Get total count
TOTAL=$(curl -s "$GATEWAY/api/v1/ach-items" | grep -o '"entry_id"' | wc -l | tr -d ' ')
echo "📊 Total ACH Items: $TOTAL"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  Default Sort: Most Recent First (created_at desc)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "$GATEWAY/api/v1/ach-items" | jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | $\(.amount_cents/100) | \(.status) | \(.created_at)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2️⃣  Sort by Amount: Highest First"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "$GATEWAY/api/v1/ach-items?sort_by=amount&sort_order=desc" | \
  jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | $\(.amount_cents/100) | \(.status)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3️⃣  Sort by Amount: Lowest First"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "$GATEWAY/api/v1/ach-items?sort_by=amount&sort_order=asc" | \
  jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | $\(.amount_cents/100) | \(.status)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4️⃣  Sort by Status: Alphabetical"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "$GATEWAY/api/v1/ach-items?sort_by=status&sort_order=asc" | \
  jq -r '.[:5] | .[] | "\(.side) | \(.status) | \(.trace_number) | $\(.amount_cents/100)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5️⃣  Sort by Side: ODFI First"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "$GATEWAY/api/v1/ach-items?sort_by=side&sort_order=asc" | \
  jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | \(.status) | $\(.amount_cents/100)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "6️⃣  Sort by Trace Number"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "$GATEWAY/api/v1/ach-items?sort_by=trace_number&sort_order=asc" | \
  jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | $\(.amount_cents/100)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "7️⃣  Oldest First (FIFO Queue)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s "$GATEWAY/api/v1/ach-items?sort_by=created_at&sort_order=asc" | \
  jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | \(.created_at)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "8️⃣  Combined: PENDING ODFI by Amount"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
PENDING_ODFI=$(curl -s "$GATEWAY/api/v1/ach-items?side=ODFI&status=PENDING&sort_by=amount&sort_order=desc")
PENDING_COUNT=$(echo "$PENDING_ODFI" | grep -o '"entry_id"' | wc -l | tr -d ' ')
echo "Found $PENDING_COUNT PENDING ODFI entries"
echo "$PENDING_ODFI" | jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | $\(.amount_cents/100) | \(.status)"'
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "9️⃣  Combined: RETURNED entries by Amount"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
RETURNED=$(curl -s "$GATEWAY/api/v1/ach-items?status=RETURNED&sort_by=amount&sort_order=desc")
RETURNED_COUNT=$(echo "$RETURNED" | grep -o '"entry_id"' | wc -l | tr -d ' ')
echo "Found $RETURNED_COUNT RETURNED entries"
echo "$RETURNED" | jq -r '.[:5] | .[] | "\(.side) | \(.trace_number) | $\(.amount_cents/100) | \(.extra.return_reason // \"N/A\")"'
echo ""

echo "🎉 Sorting Demo Complete!"
echo ""
echo "📊 Summary of Sorting Options:"
echo "   ✅ created_at (asc/desc) - Chronological order"
echo "   ✅ amount (asc/desc) - By transaction amount"
echo "   ✅ status (asc/desc) - Alphabetical by status"
echo "   ✅ side (asc/desc) - Group by ODFI/RDFI"
echo "   ✅ trace_number (asc/desc) - Sequential order"
echo ""
echo "💡 Combine with filters for powerful queries:"
echo "   curl '$GATEWAY/api/v1/ach-items?side=ODFI&status=PENDING&sort_by=amount&sort_order=desc'"
echo ""
echo "📖 See docs/SORTING_OPTIONS.md for complete documentation"

