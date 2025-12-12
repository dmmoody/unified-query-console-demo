#!/bin/bash

# ACH Concourse - Database Seeding Script (Shell Version)
# This script seeds all databases with realistic test data
#
# ⚠️  For high-volume seeding (1000+ records), use the Go seeder instead:
#     go run ./cmd/seed
#     It uses concurrent workers and is ~10x faster.
#
# This shell script creates 620 records total.

set -e

echo "🌱 ACH Concourse - Database Seeding (Shell - 620 records)"
echo "========================================================="
echo "💡 For 2000+ records, use: go run ./cmd/seed"
echo ""

# Check if services are running
if ! curl -sf http://localhost:8081/healthz > /dev/null 2>&1; then
    echo "❌ Services are not running. Please run 'make up' first."
    exit 1
fi

echo "✅ Services are running"
echo ""

# Seed ODFI entries
echo "📝 Seeding ODFI entries (150 records)..."
for i in $(seq 1 150); do
    TRACE_NUM=$(printf "%015d" $((1000000000000 + i)))
    COMPANY_NAMES=("ACME Corp" "TechStart Inc" "Global Traders" "MegaCorp LLC" "SmallBiz Co" "Enterprise Solutions" "Digital Payments" "FinTech Group" "Payment Solutions" "Commerce Partners")
    COMPANY_NAME=${COMPANY_NAMES[$((i % 10))]}
    SEC_CODES=("PPD" "CCD" "WEB" "TEL")
    SEC_CODE=${SEC_CODES[$((i % 4))]}
    AMOUNT=$((RANDOM % 100000 + 1000))
    STATUSES=("PENDING" "PENDING" "SENT" "SENT" "SENT" "CANCELLED")
    STATUS=${STATUSES[$((i % 6))]}
    
    curl -s -X POST http://localhost:8081/api/v1/entries \
        -H "Content-Type: application/json" \
        -d "{
            \"trace_number\": \"$TRACE_NUM\",
            \"company_name\": \"$COMPANY_NAME\",
            \"sec_code\": \"$SEC_CODE\",
            \"amount_cents\": $AMOUNT
        }" > /dev/null
    
    # Update status for non-PENDING entries
    if [ "$STATUS" != "PENDING" ]; then
        ENTRY_ID=$(curl -s "http://localhost:8081/api/v1/entries?trace_number=$TRACE_NUM" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        if [ ! -z "$ENTRY_ID" ]; then
            curl -s -X PATCH "http://localhost:8081/api/v1/entries/$ENTRY_ID/status" \
                -H "Content-Type: application/json" \
                -d "{\"status\": \"$STATUS\"}" > /dev/null
        fi
    fi
    
    if [ $((i % 30)) -eq 0 ]; then
        echo "  Created $i ODFI entries..."
    fi
done
echo "✅ ODFI entries created: 150"
echo ""

# Seed RDFI entries
echo "📝 Seeding RDFI entries (150 records)..."
for i in $(seq 1 150); do
    TRACE_NUM=$(printf "%015d" $((2000000000000 + i)))
    RECEIVER_NAMES=("John Smith" "Jane Doe" "Robert Johnson" "Mary Williams" "James Brown" "Patricia Davis" "Michael Miller" "Linda Wilson" "David Moore" "Barbara Taylor")
    RECEIVER_NAME=${RECEIVER_NAMES[$((i % 10))]}
    AMOUNT=$((RANDOM % 80000 + 500))
    STATUSES=("RECEIVED" "RECEIVED" "RECEIVED" "POSTED" "POSTED" "RETURNED")
    STATUS=${STATUSES[$((i % 6))]}
    
    curl -s -X POST http://localhost:8082/api/v1/entries \
        -H "Content-Type: application/json" \
        -d "{
            \"trace_number\": \"$TRACE_NUM\",
            \"receiver_name\": \"$RECEIVER_NAME\",
            \"amount_cents\": $AMOUNT
        }" > /dev/null
    
    # Return some entries
    if [ "$STATUS" == "RETURNED" ]; then
        ENTRY_ID=$(curl -s "http://localhost:8082/api/v1/entries?trace_number=$TRACE_NUM" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        RETURN_CODES=("R01" "R02" "R03" "R04" "R10")
        RETURN_CODE=${RETURN_CODES[$((i % 5))]}
        if [ ! -z "$ENTRY_ID" ]; then
            curl -s -X POST "http://localhost:8082/api/v1/entries/$ENTRY_ID/return" \
                -H "Content-Type: application/json" \
                -d "{\"reason\": \"$RETURN_CODE\"}" > /dev/null
        fi
    fi
    
    if [ $((i % 30)) -eq 0 ]; then
        echo "  Created $i RDFI entries..."
    fi
done
echo "✅ RDFI entries created: 150"
echo ""

# Seed Ledger postings
echo "📝 Seeding Ledger postings (200 records)..."
for i in $(seq 1 200); do
    if [ $((i % 2)) -eq 0 ]; then
        ACH_SIDE="ODFI"
        TRACE_NUM=$(printf "%015d" $((1000000000000 + (i / 2))))
    else
        ACH_SIDE="RDFI"
        TRACE_NUM=$(printf "%015d" $((2000000000000 + ((i + 1) / 2))))
    fi
    
    DIRECTIONS=("DEBIT" "CREDIT")
    DIRECTION=${DIRECTIONS[$((i % 2))]}
    AMOUNT=$((RANDOM % 90000 + 1000))
    DESCRIPTIONS=("Payment processing" "Vendor payment" "Payroll deposit" "Invoice payment" "Refund processing" "Settlement transfer" "Account funding" "Bill payment")
    DESCRIPTION=${DESCRIPTIONS[$((i % 8))]}
    
    curl -s -X POST http://localhost:8083/api/v1/postings \
        -H "Content-Type: application/json" \
        -d "{
            \"ach_side\": \"$ACH_SIDE\",
            \"trace_number\": \"$TRACE_NUM\",
            \"amount_cents\": $AMOUNT,
            \"direction\": \"$DIRECTION\",
            \"description\": \"$DESCRIPTION\"
        }" > /dev/null
    
    if [ $((i % 40)) -eq 0 ]; then
        echo "  Created $i ledger postings..."
    fi
done
echo "✅ Ledger postings created: 200"
echo ""

# Seed EIP cases
echo "📝 Seeding EIP cases (120 records)..."
for i in $(seq 1 120); do
    if [ $((i % 2)) -eq 0 ]; then
        SIDE="ODFI"
        TRACE_NUM=$(printf "%015d" $((1000000000000 + (i / 2))))
    else
        SIDE="RDFI"
        TRACE_NUM=$(printf "%015d" $((2000000000000 + ((i + 1) / 2))))
    fi
    
    TYPES=("RETURN_REVIEW" "NOC_REVIEW" "CUSTOMER_DISPUTE")
    TYPE=${TYPES[$((i % 3))]}
    
    NOTES_OPTIONS=(
        "Customer called to dispute charge"
        "Return received from bank, needs review"
        "NOC received, account number correction needed"
        "Duplicate transaction reported"
        "Unauthorized transaction claim"
        "Amount discrepancy reported"
        "Timing issue with settlement"
        "Customer requested investigation"
        "Bank returned entry with R03"
        "Notification of change received"
    )
    NOTES=${NOTES_OPTIONS[$((i % 10))]}
    
    curl -s -X POST http://localhost:8084/api/v1/cases \
        -H "Content-Type: application/json" \
        -d "{
            \"side\": \"$SIDE\",
            \"trace_number\": \"$TRACE_NUM\",
            \"type\": \"$TYPE\",
            \"notes\": \"$NOTES\"
        }" > /dev/null
    
    # Update some case statuses
    STATUSES=("OPEN" "OPEN" "OPEN" "IN_PROGRESS" "IN_PROGRESS" "RESOLVED")
    STATUS=${STATUSES[$((i % 6))]}
    if [ "$STATUS" != "OPEN" ]; then
        CASE_ID=$(curl -s "http://localhost:8084/api/v1/cases?trace_number=$TRACE_NUM" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
        if [ ! -z "$CASE_ID" ]; then
            curl -s -X PATCH "http://localhost:8084/api/v1/cases/$CASE_ID/status" \
                -H "Content-Type: application/json" \
                -d "{\"status\": \"$STATUS\"}" > /dev/null
        fi
    fi
    
    if [ $((i % 30)) -eq 0 ]; then
        echo "  Created $i EIP cases..."
    fi
done
echo "✅ EIP cases created: 120"
echo ""

echo "🎉 Database seeding completed!"
echo ""
echo "📊 Summary:"
echo "  ODFI entries:      150 records"
echo "  RDFI entries:      150 records"
echo "  Ledger postings:   200 records"
echo "  EIP cases:         120 records"
echo "  ─────────────────────────────"
echo "  Total:             620 records"
echo ""
echo "🔍 Verify data:"
echo "  Console: curl http://localhost:8080/api/v1/ach-items"
echo "  ODFI:    curl http://localhost:8081/api/v1/entries"
echo "  RDFI:    curl http://localhost:8082/api/v1/entries"
echo "  Ledger:  curl http://localhost:8083/api/v1/postings"
echo "  EIP:     curl http://localhost:8084/api/v1/cases"
echo ""
echo "  Balances: curl http://localhost:8083/api/v1/balances"

