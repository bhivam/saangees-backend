#!/bin/bash
API_URL="http://localhost:3000"
ADMIN_PHONE="1234567890"
ADMIN_PASSWORD="admin_secure_password123"

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

add_menu_item() {
    local name=$1
    local description=$2
    local base_price=$3
    local date=$4
    
    RESPONSE=$(curl -s -b "$COOKIE_JAR" -X POST "${API_URL}/item/create" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"${name}\",
            \"description\": \"${description}\",
            \"base_price\": ${base_price},
            \"date\": \"${date}\",
            \"spice_options\": [
                {
                    \"name\": \"Mild\",
                    \"description\": \"Subtle flavors, minimal heat\"
                },
                {
                    \"name\": \"Medium\",
                    \"description\": \"Balanced spice level\"
                },
                {
                    \"name\": \"Hot\",
                    \"description\": \"Intense heat for spice lovers\"
                }
            ]
        }")
    
    echo "Added '${name}' for ${date}"
    echo "Response: ${RESPONSE}"
}

echo "Adding menu items for next week..."

# Thursday (April 3, 2025)
add_menu_item "Butter Chicken" "Tender chicken in a rich, creamy tomato sauce" 14.99 "2025-04-03T05:00:00Z"
add_menu_item "Dal Makhani" "Creamy black lentils simmered overnight" 12.99 "2025-04-03T05:00:00Z"
add_menu_item "Naan" "Fresh baked flatbread" 3.99 "2025-04-03T05:00:00Z"
# Friday (April 4, 2025)
add_menu_item "Chicken Biryani" "Fragrant rice cooked with spiced chicken" 15.99 "2025-04-04T05:00:00Z"
add_menu_item "Raita" "Cooling yogurt with mild spices" 4.99 "2025-04-04T05:00:00Z"
add_menu_item "Mixed Vegetable Curry" "Assorted vegetables in aromatic gravy" 11.99 "2025-04-04T05:00:00Z"
# Monday (April 7, 2025)
add_menu_item "Fish Curry" "Fresh fish cooked in coconut curry" 16.99 "2025-04-07T05:00:00Z"
add_menu_item "Palak Paneer" "Cottage cheese in creamy spinach sauce" 13.99 "2025-04-07T05:00:00Z"
add_menu_item "Jeera Rice" "Basmati rice with cumin" 4.99 "2025-04-07T05:00:00Z"
# Tuesday (April 8, 2025)
add_menu_item "Lamb Rogan Josh" "Tender lamb in rich Kashmiri gravy" 17.99 "2025-04-08T05:00:00Z"
add_menu_item "Aloo Gobi" "Potato and cauliflower curry" 11.99 "2025-04-08T05:00:00Z"
add_menu_item "Garlic Naan" "Flatbread with roasted garlic" 4.99 "2025-04-08T05:00:00Z"
# Wednesday (April 9, 2025)
add_menu_item "Chicken Tikka Masala" "Grilled chicken in spiced tomato cream sauce" 15.99 "2025-04-09T05:00:00Z"
add_menu_item "Chana Masala" "Spiced chickpeas in onion-tomato gravy" 11.99 "2025-04-09T05:00:00Z"
add_menu_item "Pulao Rice" "Fragrant basmati rice with vegetables" 5.99 "2025-04-09T05:00:00Z"


echo "Verifying menu items..."

curl -s -b "$COOKIE_JAR" "${API_URL}/item/week" \
  -H "Content-Type: application/json"
rm -f "$COOKIE_JAR"
echo "Setup complete!"
