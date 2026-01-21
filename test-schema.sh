#!/bin/bash

# Test script to explore the Unraid GraphQL API schema

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: ./test-schema.sh <UNRAID_URL> <API_KEY>"
    echo "Example: ./test-schema.sh http://192.168.1.100 your-api-key"
    exit 1
fi

UNRAID_URL="$1"
API_KEY="$2"

echo "Testing connection to: ${UNRAID_URL}/graphql"
echo ""

# Test 1: Basic connectivity
echo "=== Test 1: Basic Connectivity ==="
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" "${UNRAID_URL}/graphql"
echo ""

# Test 2: Introspection query to get schema
echo "=== Test 2: Schema Introspection ==="
curl -s -X POST "${UNRAID_URL}/graphql" \
  -H "Content-Type: application/json" \
  -H "x-api-key: ${API_KEY}" \
  -d '{
    "query": "{ __schema { queryType { name fields { name description } } } }"
  }' | jq '.' 2>/dev/null || echo "Failed to parse JSON"
echo ""

# Test 3: Get available types
echo "=== Test 3: Available Types ==="
curl -s -X POST "${UNRAID_URL}/graphql" \
  -H "Content-Type: application/json" \
  -H "x-api-key: ${API_KEY}" \
  -d '{
    "query": "{ __schema { types { name kind } } }"
  }' | jq '.data.__schema.types[] | select(.kind == "OBJECT") | .name' 2>/dev/null | head -20
echo ""

# Test 4: Try some common query names
echo "=== Test 4: Testing Common Queries ==="
for query in "systemInfo" "server" "status" "info" "version"; do
    echo -n "Trying query: $query ... "
    result=$(curl -s -X POST "${UNRAID_URL}/graphql" \
      -H "Content-Type: application/json" \
      -H "x-api-key: ${API_KEY}" \
      -d "{\"query\": \"{ $query }\" }" 2>/dev/null)

    if echo "$result" | grep -q "Cannot query field"; then
        echo "❌ Not available"
    elif echo "$result" | grep -q "error"; then
        echo "⚠️  Error: $(echo $result | jq -r '.errors[0].message' 2>/dev/null)"
    else
        echo "✓ Available!"
        echo "$result" | jq '.' 2>/dev/null
    fi
done
