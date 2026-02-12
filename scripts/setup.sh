#!/usr/bin/env bash
set -e

# Check requirements
check_requirement() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed"
        return 1
    fi
    echo "✓ $1 found"
}

echo "🔍 Checking requirements..."
check_requirement "go" || exit 1

echo ""
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

echo ""
echo "✓ Environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Copy .env.example to .env and update credentials"
echo "2. Run 'docker-compose up' or 'go run main.go'"
echo ""
