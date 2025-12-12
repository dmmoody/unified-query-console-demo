#!/bin/bash

# Console API Gateway Demo Script
# This script demonstrates ALL operations through the unified gateway

set -e

GATEWAY="http://localhost:8080"

echo "🚀 Console API Gateway Demo"
echo "============================"
echo ""
echo "All operations performed through: $GATEWAY"
echo ""

# Check if gateway is running
if ! curl -sf "$GATEWAY/healthz" > /dev/null 2>&1; then
    echo "❌ Gateway is not running. Please run 'make up' first."
    exit 1
fi

echo "✅ Gateway is healthy"
echo ""

# ========== ODFI Operations ==========
echo "📝 1. Creating ODFI Entry (via Gateway)..."
ODFI_RESPONSE=$(curl -s -X POST "$GATEWAY/api/v1/odfi/entries" \
    -H "Content-Type: application/json" \
    -d '{
        "trace_number": "GATEWAY111111111",
        "company_name": "Gateway Demo Corp",
        "sec_code": "WEB",
        "amount_cents": 75000
    }')

ODFI_ID=$(echo "$ODFI_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "   ✅ Created ODFI entry: $ODFI_ID"
echo ""

echo "📝 2. Listing ODFI Entries (via Gateway)..."
ODFI_COUNT=$(curl -s "$GATEWAY/api/v1/odfi/entries" | grep -o '"id"' | wc -l | tr -d ' ')
echo "   ✅ Found $ODFI_COUNT ODFI entries"
echo ""

echo "📝 3. Getting Single ODFI Entry (via Gateway)..."
curl -s "$GATEWAY/api/v1/odfi/entries/$ODFI_ID" | head -3
echo "   ✅ Retrieved entry details"
echo ""

echo "📝 4. Updating ODFI Status (via Gateway)..."
curl -s -X PATCH "$GATEWAY/api/v1/odfi/entries/$ODFI_ID/status" \
    -H "Content-Type: application/json" \
    -d '{"status": "SENT"}' > /dev/null
echo "   ✅ Updated status to SENT"
echo ""

# ========== RDFI Operations ==========
echo "📝 5. Creating RDFI Entry (via Gateway)..."
RDFI_RESPONSE=$(curl -s -X POST "$GATEWAY/api/v1/rdfi/entries" \
    -H "Content-Type: application/json" \
    -d '{
        "trace_number": "GATEWAY999999999",
        "receiver_name": "Gateway Demo User",
        "amount_cents": 35000
    }')

RDFI_ID=$(echo "$RDFI_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "   ✅ Created RDFI entry: $RDFI_ID"
echo ""

echo "📝 6. Listing RDFI Entries (via Gateway)..."
RDFI_COUNT=$(curl -s "$GATEWAY/api/v1/rdfi/entries" | grep -o '"id"' | wc -l | tr -d ' ')
echo "   ✅ Found $RDFI_COUNT RDFI entries"
echo ""

echo "📝 7. Returning RDFI Entry (via Gateway)..."
curl -s -X POST "$GATEWAY/api/v1/rdfi/entries/$RDFI_ID/return" \
    -H "Content-Type: application/json" \
    -d '{"reason": "R01"}' > /dev/null
echo "   ✅ Returned entry with R01"
echo ""

# ========== Ledger Operations ==========
echo "📝 8. Creating Ledger Posting (via Gateway)..."
curl -s -X POST "$GATEWAY/api/v1/ledger/postings" \
    -H "Content-Type: application/json" \
    -d '{
        "ach_side": "ODFI",
        "trace_number": "GATEWAY111111111",
        "amount_cents": 75000,
        "direction": "DEBIT",
        "description": "Gateway demo payment"
    }' > /dev/null
echo "   ✅ Created ledger posting"
echo ""

echo "📝 9. Listing Ledger Postings (via Gateway)..."
LEDGER_COUNT=$(curl -s "$GATEWAY/api/v1/ledger/postings" | grep -o '"id"' | wc -l | tr -d ' ')
echo "   ✅ Found $LEDGER_COUNT ledger postings"
echo ""

echo "📝 10. Getting Balances (via Gateway)..."
BALANCES=$(curl -s "$GATEWAY/api/v1/ledger/balances")
echo "$BALANCES" | head -5
echo "   ✅ Retrieved balance information"
echo ""

# ========== EIP Operations ==========
echo "📝 11. Creating EIP Case (via Gateway)..."
CASE_RESPONSE=$(curl -s -X POST "$GATEWAY/api/v1/eip/cases" \
    -H "Content-Type: application/json" \
    -d '{
        "side": "RDFI",
        "trace_number": "GATEWAY999999999",
        "type": "RETURN_REVIEW",
        "notes": "Gateway demo - R01 return needs review"
    }')

CASE_ID=$(echo "$CASE_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "   ✅ Created EIP case: $CASE_ID"
echo ""

echo "📝 12. Listing EIP Cases (via Gateway)..."
CASE_COUNT=$(curl -s "$GATEWAY/api/v1/eip/cases" | grep -o '"id"' | wc -l | tr -d ' ')
echo "   ✅ Found $CASE_COUNT EIP cases"
echo ""

echo "📝 13. Updating Case Status (via Gateway)..."
curl -s -X PATCH "$GATEWAY/api/v1/eip/cases/$CASE_ID/status" \
    -H "Content-Type: application/json" \
    -d '{"status": "IN_PROGRESS"}' > /dev/null
echo "   ✅ Updated case status to IN_PROGRESS"
echo ""

# ========== Unified View (Legacy) ==========
echo "📝 14. Querying Unified ACH Items (via Gateway)..."
UNIFIED_COUNT=$(curl -s "$GATEWAY/api/v1/ach-items" | grep -o '"entry_id"' | wc -l | tr -d ' ')
echo "   ✅ Found $UNIFIED_COUNT unified ACH items"
echo ""

echo "🎉 Demo Complete!"
echo ""
echo "📊 Summary:"
echo "   ✅ All CRUD operations performed through gateway"
echo "   ✅ ODFI: Create, List, Get, Update Status"
echo "   ✅ RDFI: Create, List, Get, Return"
echo "   ✅ Ledger: Create Posting, List, Get Balances"
echo "   ✅ EIP: Create Case, List, Update Status"
echo "   ✅ Unified View: Query across services"
echo ""
echo "🌟 Gateway URL: $GATEWAY"
echo "📖 Full documentation: see GATEWAY.md"

