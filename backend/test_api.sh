#!/bin/bash

# API Testing Script for Landscaping Forms
# This script tests all the API endpoints

BASE_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
PASSED=0
FAILED=0

# Function to check HTTP status code
check_status() {
    local expected=$1
    local actual=$2
    local test_name=$3

    if [ "$actual" -eq "$expected" ]; then
        echo -e "${GREEN}✓ PASS${NC} - $test_name (HTTP $actual)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} - $test_name (Expected HTTP $expected, got $actual)"
        ((FAILED++))
        return 1
    fi
}

echo "================================"
echo "API Testing Script"
echo "================================"
echo ""

# Test 1: Health Check
echo "1. Testing Health Check..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health")
check_status 200 "$HTTP_CODE" "Health check endpoint"
echo ""

# Test 2: Create Shrub Form
echo "2. Creating Shrub Form..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/forms/shrub" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe","home_phone":"555-1234","num_shrubs":5}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
SHRUB_RESPONSE=$(echo "$RESPONSE" | sed '$d')
SHRUB_ID=$(echo "$SHRUB_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if check_status 201 "$HTTP_CODE" "Create shrub form"; then
    echo "   Created form ID: $SHRUB_ID"
fi
echo ""

# Test 3: Create Pesticide Form
echo "3. Creating Pesticide Form..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/forms/pesticide" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Jane","last_name":"Smith","home_phone":"555-5678","pesticide_name":"RoundUp"}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
PESTICIDE_RESPONSE=$(echo "$RESPONSE" | sed '$d')
PESTICIDE_ID=$(echo "$PESTICIDE_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if check_status 201 "$HTTP_CODE" "Create pesticide form"; then
    echo "   Created form ID: $PESTICIDE_ID"
fi
echo ""

# Test 4: List All Forms
echo "4. Listing All Forms..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/forms")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
LIST_RESPONSE=$(echo "$RESPONSE" | sed '$d')
COUNT=$(echo "$LIST_RESPONSE" | grep -o '"count":[0-9]*' | cut -d':' -f2)

if check_status 200 "$HTTP_CODE" "List all forms"; then
    echo "   Found $COUNT forms"
fi
echo ""

# Test 5: List Forms Sorted by First Name
echo "5. Listing Forms Sorted by First Name (ASC)..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/forms?sort_by=first_name&order=ASC")
check_status 200 "$HTTP_CODE" "List forms with sorting"
echo ""

# Test 6: Get Shrub Form by ID
echo "6. Getting Shrub Form by ID..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/forms/$SHRUB_ID")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
FORM_RESPONSE=$(echo "$RESPONSE" | sed '$d')
FORM_TYPE=$(echo "$FORM_RESPONSE" | grep -o '"form_type":"[^"]*"' | cut -d'"' -f4)

if check_status 200 "$HTTP_CODE" "Get form by ID"; then
    echo "   Form type: $FORM_TYPE"
fi
echo ""

# Test 7: Update Shrub Form
echo "7. Updating Shrub Form..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/forms/$SHRUB_ID" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe Jr.","home_phone":"555-9999","num_shrubs":10}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
UPDATED_RESPONSE=$(echo "$RESPONSE" | sed '$d')
UPDATED_NAME=$(echo "$UPDATED_RESPONSE" | grep -o '"last_name":"[^"]*"' | cut -d'"' -f4)

if check_status 200 "$HTTP_CODE" "Update shrub form"; then
    echo "   Updated last name: $UPDATED_NAME"
fi
echo ""

# Test 8: Update Pesticide Form
echo "8. Updating Pesticide Form..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/forms/$PESTICIDE_ID" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Jane","last_name":"Smith-Johnson","home_phone":"555-0000","pesticide_name":"Organic Spray"}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
UPDATED_RESPONSE=$(echo "$RESPONSE" | sed '$d')
UPDATED_PESTICIDE=$(echo "$UPDATED_RESPONSE" | grep -o '"pesticide_name":"[^"]*"' | cut -d'"' -f4)

if check_status 200 "$HTTP_CODE" "Update pesticide form"; then
    echo "   Updated pesticide: $UPDATED_PESTICIDE"
fi
echo ""

# Test 9: List Forms After Updates
echo "9. Listing All Forms After Updates..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/forms")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
LIST_RESPONSE=$(echo "$RESPONSE" | sed '$d')
COUNT=$(echo "$LIST_RESPONSE" | grep -o '"count":[0-9]*' | cut -d':' -f2)

if check_status 200 "$HTTP_CODE" "List forms after updates"; then
    echo "   Total forms: $COUNT"
fi
echo ""

# Test 10: Delete Pesticide Form
echo "10. Deleting Pesticide Form..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/api/forms/$PESTICIDE_ID")
check_status 200 "$HTTP_CODE" "Delete pesticide form"
echo ""

# Test 11: List Forms After Deletion
echo "11. Listing Forms After Deletion..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/forms")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
LIST_RESPONSE=$(echo "$RESPONSE" | sed '$d')
COUNT=$(echo "$LIST_RESPONSE" | grep -o '"count":[0-9]*' | cut -d':' -f2)

if check_status 200 "$HTTP_CODE" "List forms after deletion"; then
    echo "   Remaining forms: $COUNT"
fi
echo ""

# Test 12: Try to Get Deleted Form (Should Return 404)
echo "12. Trying to Get Deleted Form (Should Return 404)..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/forms/$PESTICIDE_ID")
check_status 404 "$HTTP_CODE" "Get deleted form (expect 404)"
echo ""

# Test 13: Delete Shrub Form (Cleanup)
echo "13. Deleting Shrub Form (Cleanup)..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/api/forms/$SHRUB_ID")
check_status 200 "$HTTP_CODE" "Delete shrub form (cleanup)"
echo ""

# Final summary
echo "================================"
echo "Test Results Summary"
echo "================================"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total:  $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
