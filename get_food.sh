#!/bin/bash

# Configuration
API_URL="http://localhost:3000"
ADMIN_PHONE="1234567890"
ADMIN_PASSWORD="admin_secure_password123"

# Step 2: Get authentication token
echo "Getting authentication token..."
COOKIE_JAR="cookies.txt"
rm -f "$COOKIE_JAR"

TOKEN_RESPONSE=$(curl -s -c "$COOKIE_JAR" -X POST "${API_URL}/token" \
  -H "Content-Type: application/json" \
  -d "{
    \"phone_number\": \"${ADMIN_PHONE}\",
    \"password\": \"${ADMIN_PASSWORD}\"
  }")

echo "Token response: ${TOKEN_RESPONSE}"

curl -s -b "$COOKIE_JAR" "${API_URL}/item/week" \
  -H "Content-Type: application/json"

# Cleanup
rm -f "$COOKIE_JAR"

echo "Setup complete!"
