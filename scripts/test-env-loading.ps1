# Test script to verify environment variable loading
# This temporarily disables YAML config to test env-only loading

Write-Host "===================================" -ForegroundColor Cyan
Write-Host "Environment Variable Loading Test" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Cyan
Write-Host ""

# Check if .env file exists
if (-not (Test-Path ".env")) {
    Write-Host "❌ .env file not found!" -ForegroundColor Red
    Write-Host "Please create .env file first. Use .env.example as a template." -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Found .env file" -ForegroundColor Green
Write-Host ""

# Backup config.yaml if it exists
if (Test-Path "configs\config.yaml") {
    Write-Host "📦 Backing up configs\config.yaml to configs\config.yaml.bak" -ForegroundColor Yellow
    Copy-Item "configs\config.yaml" "configs\config.yaml.bak" -Force
    Remove-Item "configs\config.yaml"
}

Write-Host "🧪 Testing application with environment variables only..." -ForegroundColor Cyan
Write-Host ""

# Load environment variables from .env file
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^#][^=]+)=(.*)$') {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim()
        Set-Item -Path "env:$name" -Value $value
        Write-Host "Loaded: $name" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "🚀 Starting application..." -ForegroundColor Green
Write-Host ""

# Build and run
go run cmd\server\main.go

# Cleanup - restore config.yaml
if (Test-Path "configs\config.yaml.bak") {
    Write-Host ""
    Write-Host "🔄 Restoring configs\config.yaml" -ForegroundColor Yellow
    Move-Item "configs\config.yaml.bak" "configs\config.yaml" -Force
}

Write-Host ""
Write-Host "✅ Test complete!" -ForegroundColor Green
