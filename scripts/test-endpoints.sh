#!/usr/bin/env bash
set -e

# Test all endpoints
BASE_URL="http://localhost:8080"

echo "🧪 Testing DC_Handler API Endpoints"
echo "=================================="
echo ""

# Health check
echo "1️⃣  Testing /health..."
curl -s "$BASE_URL/health" | jq . || echo "Failed"
echo ""

# Get posts
echo "2️⃣  Testing /api/wp-posts..."
curl -s "$BASE_URL/api/wp-posts?status=draft&per_page=5" | jq . || echo "Failed"
echo ""

# Optimize content
echo "3️⃣  Testing /api/optimize-content..."
curl -s -X POST "$BASE_URL/api/optimize-content" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Title",
    "content": "This is test content",
    "excerpt": "Test excerpt"
  }' | jq . || echo "Failed"
echo ""

echo "✓ Testing complete"
