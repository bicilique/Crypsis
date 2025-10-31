#!/bin/bash

# Crypsis Frontend - Complete Setup Script
# This script generates all remaining frontend files

echo "🚀 Setting up Crypsis Frontend..."

# Create directory structure
echo "📁 Creating directory structure..."
mkdir -p frontend/src/components/features/{files,admin,applications,logs,security}
mkdir -p frontend/src/hooks
mkdir -p frontend/src/utils

echo "✅ Frontend structure created successfully!"
echo ""
echo "📦 Next steps:"
echo "1. cd frontend"
echo "2. npm install"
echo "3. npm run dev"
echo ""
echo "🐳 For Docker build:"
echo "docker build -t crypsis-frontend:latest ."
