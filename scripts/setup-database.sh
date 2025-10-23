#!/bin/bash

# ZeroTrace Database Setup Script
# This script sets up the PostgreSQL database for ZeroTrace development

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Database configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-zerotrace}
DB_SSLMODE=${DB_SSLMODE:-disable}

echo -e "${BLUE}🚀 ZeroTrace Database Setup${NC}"
echo -e "${BLUE}===========================${NC}"

# Function to print status
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if PostgreSQL is running
echo -e "${BLUE}Checking PostgreSQL connection...${NC}"
if ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER > /dev/null 2>&1; then
    print_error "PostgreSQL is not running or not accessible"
    print_warning "Please start PostgreSQL and ensure it's accessible at $DB_HOST:$DB_PORT"
    exit 1
fi
print_status "PostgreSQL is running"

# Check if database exists
echo -e "${BLUE}Checking if database exists...${NC}"
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    print_warning "Database '$DB_NAME' already exists"
    read -p "Do you want to drop and recreate it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}Dropping existing database...${NC}"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;"
        print_status "Database dropped"
    else
        print_warning "Using existing database"
    fi
fi

# Create database if it doesn't exist
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo -e "${BLUE}Creating database '$DB_NAME'...${NC}"
    createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    print_status "Database created"
fi

# Set environment variables for migration
export DB_HOST=$DB_HOST
export DB_PORT=$DB_PORT
export DB_USER=$DB_USER
export DB_PASSWORD=$DB_PASSWORD
export DB_NAME=$DB_NAME
export DB_SSLMODE=$DB_SSLMODE

# Run migrations
echo -e "${BLUE}Running database migrations...${NC}"
cd api-go/migrations

# Check if Go is available
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

# Build and run migration tool
echo -e "${BLUE}Building migration tool...${NC}"
go mod init migrations 2>/dev/null || true
go mod tidy
go build -o migrate migrate.go

echo -e "${BLUE}Running migrations...${NC}"
./migrate

# Clean up
rm -f migrate
rm -f go.mod go.sum

print_status "Database setup completed successfully!"

echo -e "${BLUE}Database Configuration:${NC}"
echo -e "  Host: $DB_HOST"
echo -e "  Port: $DB_PORT"
echo -e "  Database: $DB_NAME"
echo -e "  User: $DB_USER"
echo -e "  SSL Mode: $DB_SSLMODE"

echo -e "${BLUE}Connection String:${NC}"
echo -e "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE"

echo -e "${GREEN}🎉 ZeroTrace database is ready!${NC}"
