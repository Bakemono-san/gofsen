#!/bin/bash

# 🧪 Gofsen Framework Test Suite
# Test script for all implemented functionalities

echo "🎉 Testing Gofsen Framework Functionalities"
echo "============================================="

BASE_URL="http://localhost:8081"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local method=$1
    local url=$2
    local description=$3
    local extra_args=$4
    
    echo -e "\n${YELLOW}Testing: ${description}${NC}"
    echo "Command: curl -X ${method} ${extra_args} ${BASE_URL}${url}"
    
    response=$(curl -s -X ${method} ${extra_args} ${BASE_URL}${url})
    
    if [[ $response == *"test_passed\":true"* ]] || [[ $response == *"✅"* ]]; then
        echo -e "${GREEN}✅ PASSED${NC}"
    else
        echo -e "${RED}❌ FAILED${NC}"
    fi
    
    echo "Response: ${response}" | jq 2>/dev/null || echo "Response: ${response}"
}

echo -e "\n🔍 Checking server status..."
if ! curl -s ${BASE_URL}/health > /dev/null; then
    echo -e "${RED}❌ Server not running at ${BASE_URL}${NC}"
    echo "Please start the server with: go run cmd/main.go"
    exit 1
fi

echo -e "${GREEN}✅ Server is running${NC}"

# =============================================================================
# 🧱 PART 1: BASIC ROUTING
# =============================================================================
echo -e "\n\n🧱 PART 1: BASIC ROUTING METHODS"
echo "=================================="

test_endpoint "GET" "/test/get" "GET method"
test_endpoint "POST" "/test/post" "POST method" "-H 'Content-Type: application/json' -d '{\"message\":\"test\"}'"
test_endpoint "PUT" "/test/put" "PUT method"
test_endpoint "DELETE" "/test/delete" "DELETE method"
test_endpoint "PATCH" "/test/patch" "PATCH method"

# =============================================================================
# 🧭 PART 2: ROUTE GROUPS & LOCAL MIDDLEWARE
# =============================================================================
echo -e "\n\n🧭 PART 2: ROUTE GROUPS & MIDDLEWARE"
echo "====================================="

test_endpoint "GET" "/test/group/basic" "Route group with local middleware"

# =============================================================================
# 🔐 PART 3: MIDDLEWARE & SECURITY
# =============================================================================
echo -e "\n\n🔐 PART 3: MIDDLEWARE & SECURITY"
echo "================================="

test_endpoint "GET" "/test/logger" "Logger middleware (check console)"

echo -e "\n${YELLOW}Testing: Auth middleware (without token - should fail)${NC}"
response=$(curl -s ${BASE_URL}/test/auth/protected)
if [[ $response == *"Missing Authorization Header"* ]]; then
    echo -e "${GREEN}✅ PASSED (correctly rejected)${NC}"
else
    echo -e "${RED}❌ FAILED${NC}"
fi
echo "Response: ${response}" | jq 2>/dev/null || echo "Response: ${response}"

test_endpoint "GET" "/test/auth/protected" "Auth middleware (with token)" "-H 'Authorization: Bearer valid-token'"

test_endpoint "GET" "/test/auth/public" "Public route (no auth required)"

echo -e "\n${YELLOW}Testing: Recovery middleware (panic handling)${NC}"
response=$(curl -s ${BASE_URL}/test/recovery)
if [[ $response == *"Internal Server Error"* ]]; then
    echo -e "${GREEN}✅ PASSED (panic recovered)${NC}"
else
    echo -e "${RED}❌ FAILED${NC}"
fi
echo "Response: ${response}" | jq 2>/dev/null || echo "Response: ${response}"

test_endpoint "GET" "/test/cors/check" "CORS middleware" "-H 'Origin: https://example.com'"

# =============================================================================
# ⚙️ PART 4: HELPERS & I/O
# =============================================================================
echo -e "\n\n⚙️ PART 4: HELPERS & I/O"
echo "========================="

test_endpoint "GET" "/test/json" "JSON response"

test_endpoint "GET" "/test/query?name=John&age=25&city=Paris" "Query parameters"

test_endpoint "POST" "/test/bind" "BindJSON (body parsing)" "-H 'Content-Type: application/json' -d '{\"name\":\"Alice\",\"email\":\"alice@example.com\",\"age\":30}'"

test_endpoint "GET" "/test/error?type=400" "Error response (400)"
test_endpoint "GET" "/test/error?type=404" "Error response (404)"
test_endpoint "GET" "/test/error" "Error response (custom)"

# =============================================================================
# 🧪 ADVANCED TESTS
# =============================================================================
echo -e "\n\n🧪 ADVANCED: COMBINED FEATURES"
echo "==============================="

test_endpoint "POST" "/test/multi/combined?name=TestUser" "Multiple middleware + features" "-H 'Content-Type: application/json' -H 'Origin: https://example.com' -d '{\"action\":\"test\",\"data\":\"combined\"}'"

# =============================================================================
# 📋 FINAL SUMMARY
# =============================================================================
echo -e "\n\n📋 COMPLETE TEST SUMMARY"
echo "========================="

test_endpoint "GET" "/test/all" "Complete test suite overview"

echo -e "\n\n🎉 ${GREEN}Testing Complete!${NC}"
echo -e "📚 For detailed documentation, visit: ${BASE_URL}/test/all"
echo -e "🔗 Try the endpoints manually or with Postman/Insomnia"

echo -e "\n💡 ${YELLOW}Quick test commands:${NC}"
echo "• Health check: curl ${BASE_URL}/health"
echo "• Test overview: curl ${BASE_URL}/test/all | jq"
echo "• Auth test: curl -H 'Authorization: Bearer valid-token' ${BASE_URL}/test/auth/protected"
echo "• CORS test: curl -H 'Origin: https://example.com' ${BASE_URL}/test/cors/check"