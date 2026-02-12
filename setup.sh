#!/usr/bin/env bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== DC_Handler Phase 1 Setup ===${NC}\n"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed: $(go version)${NC}\n"

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}Creating .env from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ .env created - Please update with your credentials${NC}\n"
else
    echo -e "${GREEN}✓ .env already exists${NC}\n"
fi

# Download dependencies
echo -e "${YELLOW}Downloading Go dependencies...${NC}"
go mod download
go mod tidy

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Dependencies installed${NC}\n"
else
    echo -e "${RED}Failed to install dependencies${NC}"
    exit 1
fi

# Build binary
echo -e "${YELLOW}Building application...${NC}"
go build -o dc-handler .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}\n"
else
    echo -e "${RED}Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}=== Setup Complete ===${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo -e "1. Update .env with your WordPress and API credentials"
echo -e "2. Start Ollama: ${YELLOW}ollama serve${NC}"
echo -e "3. Run the application: ${YELLOW}go run main.go${NC}"
echo -e "4. Or use Docker: ${YELLOW}docker-compose up${NC}\n"
