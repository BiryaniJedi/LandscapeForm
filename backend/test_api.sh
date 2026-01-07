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

echo "================================"
echo "USER ENDPOINT TESTS"
echo "================================"
echo ""

# Test 14: Create User (Auth)
echo "14. Creating User (Auth)..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Alice","last_name":"Johnson","date_of_birth":"1990-05-15T00:00:00Z","username":"alice.johnson","password":"password123"}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
USER_RESPONSE=$(echo "$RESPONSE" | sed '$d')
USER_ID=$(echo "$USER_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "$HTTP_CODE"

if check_status 200 "$HTTP_CODE" "Create user (registration)"; then
    echo "   Created user ID: $USER_ID"
fi
echo ""

# Test 15: Create Second User
echo "15. Creating Second User..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Bob","last_name":"Smith","date_of_birth":"1985-08-20T00:00:00Z","username":"bob.smith","password":"password456"}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
USER2_RESPONSE=$(echo "$RESPONSE" | sed '$d')
USER2_ID=$(echo "$USER2_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if check_status 200 "$HTTP_CODE" "Create second user"; then
    echo "   Created user ID: $USER2_ID"
fi
echo ""

# echo "16. Logging in to first User..."
# RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/auth/login" \
#   -H "Content-Type: application/json" \
#   -d '{"user_name":"alice.johnson","password":"password123"}')
# HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
# USER2_RESPONSE=$(echo "$RESPONSE" | sed '$d')
# USER2_ID=$(echo "$USER2_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

# if check_status 201 "$HTTP_CODE" "Create second user"; then
#     echo "   Created user ID: $USER2_ID"
# fi
# echo ""

# Test 16: Get User by ID
echo "16. Getting User by ID..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/users/$USER_ID")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
GET_USER_RESPONSE=$(echo "$RESPONSE" | sed '$d')
USERNAME=$(echo "$GET_USER_RESPONSE" | grep -o '"user_name":"[^"]*"' | cut -d'"' -f4)
PENDING=$(echo "$GET_USER_RESPONSE" | grep -o '"pending":[^,}]*' | cut -d':' -f2)

if check_status 200 "$HTTP_CODE" "Get user by ID"; then
    echo "   Username: $USERNAME"
    echo "   Pending: $PENDING"
fi
echo ""

# Test 17: List All Users (Admin endpoint - no auth yet)
echo "17. Listing All Users..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/users")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
LIST_USERS_RESPONSE=$(echo "$RESPONSE" | sed '$d')
USER_COUNT=$(echo "$LIST_USERS_RESPONSE" | grep -o '"count":[0-9]*' | cut -d':' -f2)

if check_status 200 "$HTTP_CODE" "List all users"; then
    echo "   Found $USER_COUNT users"
fi
echo ""

# Test 18: List Users Sorted by First Name (ASC)
echo "18. Listing Users Sorted by First Name (ASC)..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/users?sort_by=first_name&order=ASC")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
check_status 200 "$HTTP_CODE" "List users with sorting"
echo ""

# Test 19: Update User
echo "19. Updating User..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/api/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Alice","last_name":"Johnson-Williams","date_of_birth":"1990-05-15T00:00:00Z","user_name":"alice.williams","password":""}')
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
UPDATED_USER_RESPONSE=$(echo "$RESPONSE" | sed '$d')

if check_status 200 "$HTTP_CODE" "Update user"; then
    echo "   User updated successfully"
fi
echo ""

# Test 20: Get Updated User
echo "20. Getting Updated User..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/users/$USER_ID")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
GET_UPDATED_USER=$(echo "$RESPONSE" | sed '$d')
UPDATED_LASTNAME=$(echo "$GET_UPDATED_USER" | grep -o '"last_name":"[^"]*"' | cut -d'"' -f4)
UPDATED_USERNAME=$(echo "$GET_UPDATED_USER" | grep -o '"user_name":"[^"]*"' | cut -d'"' -f4)

if check_status 200 "$HTTP_CODE" "Get updated user"; then
    echo "   Updated last name: $UPDATED_LASTNAME"
    echo "   Updated username: $UPDATED_USERNAME"
fi
echo ""

# Test 21: Approve User (Admin endpoint)
echo "21. Approving User Registration (Admin)..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/users/$USER_ID/approve")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
APPROVED_USER=$(echo "$RESPONSE" | sed '$d')

if check_status 200 "$HTTP_CODE" "Approve user registration"; then
    echo "   User approved successfully"
fi
echo ""

# Test 22: Verify User is Approved
echo "22. Verifying User is Approved..."
RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/api/users/$USER_ID")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
VERIFIED_USER=$(echo "$RESPONSE" | sed '$d')
PENDING_STATUS=$(echo "$VERIFIED_USER" | grep -o '"pending":[^,}]*' | cut -d':' -f2)

if check_status 200 "$HTTP_CODE" "Verify user approval"; then
    echo "   Pending status: $PENDING_STATUS (should be false)"
fi
echo ""

# Test 23: Delete User (Admin endpoint)
echo "23. Deleting Second User (Admin)..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/api/users/$USER2_ID")
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
DELETE_RESPONSE=$(echo "$RESPONSE" | sed '$d')

if check_status 200 "$HTTP_CODE" "Delete user"; then
    echo "   User deleted successfully"
fi
echo ""

# Test 24: Verify Deleted User Returns 404
echo "24. Trying to Get Deleted User (Should Return 404)..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/users/$USER2_ID")
check_status 404 "$HTTP_CODE" "Get deleted user (expect 404)"
echo ""

# Test 25: Delete First User (Cleanup)
echo "25. Deleting First User (Cleanup)..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/api/users/$USER_ID")
check_status 200 "$HTTP_CODE" "Delete first user (cleanup)"
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
