#!/bin/bash
# ZeroTrace Development Setup with Fast Package Managers

echo "🚀 Setting up ZeroTrace with bun (Node.js) and uv (Python)"

# Add bun to PATH if not already there
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# Verify installations
echo "📦 Checking package managers..."
which bun && bun --version
which uv && uv --version

echo ""
echo "🔧 Installing dependencies..."

# Install Node.js dependencies with bun
echo "📦 Installing web-react dependencies with bun..."
cd web-react
bun install
cd ..

# Install Python dependencies with uv
echo "🐍 Installing enrichment-python dependencies with uv..."
cd enrichment-python
uv pip install -r requirements.txt
cd ..

echo ""
echo "✅ Setup complete! Use these commands:"
echo "  • bun run dev     (instead of npm run dev)"
echo "  • bun install     (instead of npm install)"
echo "  • uv pip install  (instead of pip install)"
echo "  • uv pip sync     (instead of pip install -r requirements.txt)"