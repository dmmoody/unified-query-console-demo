#!/bin/bash

# ACH Concourse - System Verification Script

set -e

echo "🚀 ACH Concourse System Verification"
echo "===================================="
echo ""

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed"
    exit 1
fi

echo "✅ Docker Compose is installed"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "⚠️  Go is not installed (required for local development only)"
else
    echo "✅ Go is installed: $(go version)"
fi

echo ""
echo "📦 Building Docker images..."
docker-compose build

echo ""
echo "🏥 Starting services..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 10

echo ""
echo "🧪 Testing service health endpoints..."

# Function to test health endpoint
test_health() {
    local service=$1
    local port=$2
    
    if curl -sf "http://localhost:${port}/healthz" > /dev/null; then
        echo "  ✅ ${service} is healthy on port ${port}"
    else
        echo "  ❌ ${service} failed health check on port ${port}"
        return 1
    fi
}

test_health "Console" 8080
test_health "ODFI" 8081
test_health "RDFI" 8082
test_health "Ledger" 8083
test_health "EIP" 8084

echo ""
echo "🎉 All services are running!"
echo ""
echo "📋 Service Endpoints:"
echo "   Console:  http://localhost:8080/api/v1/ach-items"
echo "   ODFI:     http://localhost:8081/api/v1/entries"
echo "   RDFI:     http://localhost:8082/api/v1/entries"
echo "   Ledger:   http://localhost:8083/api/v1/postings"
echo "   EIP:      http://localhost:8084/api/v1/cases"
echo ""
echo "💡 To view logs: docker-compose logs -f [service-name]"
echo "💡 To stop: docker-compose down"

