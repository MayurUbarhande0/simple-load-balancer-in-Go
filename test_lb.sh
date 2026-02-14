#!/bin/bash

# Script to test the load balancer
set -e  # Exit on error
set -u  # Exit on undefined variable

# Check for required tools
command -v curl >/dev/null 2>&1 || { echo "Error: curl is required but not installed." >&2; exit 1; }

# Check if jq is available (optional, for pretty printing)
HAS_JQ=false
if command -v jq >/dev/null 2>&1; then
    HAS_JQ=true
fi

echo "Testing Load Balancer..."
echo "======================="
echo ""

LB_URL="http://localhost:8080"

# Check if load balancer is running
echo "1. Checking if load balancer is running..."
if ! curl -s -f "${LB_URL}/health" > /dev/null; then
    echo "❌ Load balancer is not responding. Please start it first."
    echo "   Run: ./lb"
    exit 1
fi
echo "✅ Load balancer is running"
echo ""

# Test health endpoint
echo "2. Testing /health endpoint..."
echo "Response:"
if [ "$HAS_JQ" = true ]; then
    curl -s "${LB_URL}/health" | jq .
else
    curl -s "${LB_URL}/health"
    echo ""
fi
echo ""

# Test metrics endpoint
echo "3. Testing /metrics endpoint..."
echo "Response:"
if [ "$HAS_JQ" = true ]; then
    curl -s "${LB_URL}/metrics" | jq .
else
    curl -s "${LB_URL}/metrics"
    echo ""
fi
echo ""

# Test load balancing
echo "4. Testing load balancing (10 requests)..."
for i in {1..10}; do
    echo "Request $i:"
    curl -s "${LB_URL}/" | head -n 1
done
echo ""

# Test concurrent requests
echo "5. Testing concurrent requests..."
echo "Sending 50 concurrent requests..."
for i in {1..50}; do
    curl -s "${LB_URL}/" > /dev/null &
done
wait
echo "✅ Concurrent requests completed"
echo ""

# Check metrics after load
echo "6. Final metrics..."
if [ "$HAS_JQ" = true ]; then
    curl -s "${LB_URL}/metrics" | jq .
else
    curl -s "${LB_URL}/metrics"
    echo ""
fi
echo ""

echo "======================="
echo "Testing completed!"
