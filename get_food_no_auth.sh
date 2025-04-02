#!/bin/bash

# Configuration
API_URL="http://localhost:3000"

curl -s "${API_URL}/item/week"

# Cleanup
rm -f "$COOKIE_JAR"

echo "Setup complete!"
