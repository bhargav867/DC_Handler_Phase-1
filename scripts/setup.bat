@echo off
echo Checking requirements...

where go >nul 2>nul
if errorlevel 1 (
    echo Error: Go is not installed
    exit /b 1
)

echo OK: Go found

echo.
echo Downloading dependencies...
go mod download
go mod tidy

echo.
echo Setup complete!
echo.
echo Next steps:
echo 1. Copy .env.example to .env and update credentials
echo 2. Run 'docker-compose up' or 'go run main.go'
echo.
