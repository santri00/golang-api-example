#!/usr/bin/env bash
# seed_and_query.sh
# Seeds vehicle data via POST and queries each record back via GET.
#
# Usage:
#   ./scripts/seed_and_query.sh              # default: http://localhost:3001
#   BASE_URL=http://localhost:8080 ./scripts/seed_and_query.sh
#
# Requires: curl, jq (for JSON pretty-printing)
#   Install jq on Windows via winget: winget install jqlang.jq

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:3001}"
REQUESTS_DIR="$(dirname "$0")/Requests"
DIVIDER="======================================================"

# ---------- helpers ----------------------------------------------------------

check_deps() {
    if ! command -v curl &>/dev/null; then
        echo "ERROR: curl is required but not installed." >&2
        exit 1
    fi
    if ! command -v jq &>/dev/null; then
        echo "ERROR: jq is required for JSON pretty-printing." >&2
        echo "       Install via winget: winget install jqlang.jq" >&2
        exit 1
    fi
}

print_header() {
    echo ""
    echo "$DIVIDER"
    echo "  $1"
    echo "$DIVIDER"
}

# Extract a top-level string field from a flat JSON response using sed.
# Usage: extract_field <json> <field_name>
extract_field() {
    echo "$1" | sed -n "s/.*\"$2\":\"\([^\"]*\)\".*/\1/p"
}

# ---------- preflight --------------------------------------------------------

check_deps

print_header "Checking server health at $BASE_URL"
ping_response=$(curl -sf --max-time 5 "$BASE_URL/inventory/v1/ping" || true)
if [ -z "$ping_response" ]; then
    echo "ERROR: Server not reachable at $BASE_URL. Is it running?" >&2
    exit 1
fi
echo "Server is up. Ping response:"
echo "$ping_response" | jq .

# ---------- seed + query per file --------------------------------------------

stock_ids=()

for file in "$REQUESTS_DIR"/AddVehicle_Request_*.json; do
    filename=$(basename "$file")

    print_header "POST $filename"

    # Read the JSON file as-is and send it as the request body.
    post_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        --data @"$file" \
        "$BASE_URL/inventory/v1/vehicles")

    echo "Raw response: $post_response"

    stock_id=$(extract_field "$post_response" "stock_id")

    if [ -z "$stock_id" ]; then
        echo "WARNING: No stock_id returned for $filename — skipping GET calls." >&2
        continue
    fi

    echo "Stock ID: $stock_id"
    stock_ids+=("$stock_id")

    # --- GET by stock ID -----------------------------------------------------
    print_header "GET /inventory/v1/vehicles/$stock_id"
    curl -s "$BASE_URL/inventory/v1/vehicles/$stock_id" | jq .
done

# ---------- list all vehicles ------------------------------------------------

print_header "GET /inventory/v1/vehicles (list all)"
curl -s "$BASE_URL/inventory/v1/vehicles" | jq .

# ---------- summary ----------------------------------------------------------

print_header "Summary"
echo "Total vehicles seeded: ${#stock_ids[@]}"
echo ""
echo "Stock IDs created:"
for id in "${stock_ids[@]}"; do
    echo "  - $id"
done
echo ""
