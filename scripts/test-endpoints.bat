@echo off
setlocal enabledelayedexpansion

set BASE_URL=http://localhost:8080

echo Testing DC_Handler API Endpoints
echo =================================
echo.

echo Testing /health...
curl -s "%BASE_URL%/health" | jq .
echo.

echo Testing /api/wp-posts...
curl -s "%BASE_URL%/api/wp-posts?status=draft&per_page=5" | jq .
echo.

echo Testing /api/optimize-content...
curl -s -X POST "%BASE_URL%/api/optimize-content" ^
  -H "Content-Type: application/json" ^
  -d "{\"title\": \"Test Title\", \"content\": \"Test content\", \"excerpt\": \"Test\"}" | jq .
echo.

echo Testing complete
