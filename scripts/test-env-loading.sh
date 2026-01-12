#!/bin/bash

# Test script to verify environment variable loading
# This temporarily disables YAML config to test env-only loading

echo "==================================="
echo "Environment Variable Loading Test"
echo "==================================="
echo ""

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "❌ .env file not found!"
    echo "Please create .env file first. Use .env.example as a template."
    exit 1
fi

echo "✅ Found .env file"
echo ""

# Backup config.yaml if it exists
if [ -f "configs/config.yaml" ]; then
    echo "📦 Backing up configs/config.yaml to configs/config.yaml.bak"
    cp configs/config.yaml configs/config.yaml.bak
    rm configs/config.yaml
fi

echo "🧪 Testing application with environment variables only..."
echo ""

# Load env vars and run the app
export $(cat .env | grep -v '^#' | xargs)

# Build and run
go run cmd/server/main.go

# Cleanup - restore config.yaml
if [ -f "configs/config.yaml.bak" ]; then
    echo ""
    echo "🔄 Restoring configs/config.yaml"
    mv configs/config.yaml.bak configs/config.yaml
fi

echo ""
echo "✅ Test complete!"
