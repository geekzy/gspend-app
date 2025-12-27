#!/bin/bash

echo "🚀 Starting GSpend Demo Environment..."
echo "This will create a complete demo environment with dummy data."
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Stop any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose -f docker-compose.demo.yml down -v

# Build and start services
echo "🔨 Building and starting services..."
docker-compose -f docker-compose.demo.yml up --build -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
echo "This may take a few minutes on first run..."

# Function to check service health
check_service() {
    local service_name=$1
    local url=$2
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "$url" > /dev/null 2>&1; then
            echo "✅ $service_name is ready"
            return 0
        fi
        echo "⏳ Waiting for $service_name... (attempt $attempt/$max_attempts)"
        sleep 10
        attempt=$((attempt + 1))
    done

    echo "❌ $service_name failed to start"
    return 1
}

# Wait for auth service
check_service "Auth Service" "http://localhost:8081/api/v1/auth/health"

# Wait for financial service
check_service "Financial Service" "http://localhost:8082/api/v1/health"

# Wait for frontend (through nginx)
check_service "Frontend" "http://localhost/"

echo ""
echo "🎉 Demo environment is ready!"
echo ""
echo "📊 Dashboard Preview:"
echo "   URL: http://localhost"
echo ""
echo "🔐 Demo Login Credentials:"
echo "   Email: demo@gspend.com"
echo "   Password: password"
echo ""
echo "📈 What's included:"
echo "   • 3 months of sample transactions"
echo "   • Multiple expense categories (Food, Transport, Shopping, etc.)"
echo "   • Income records (Salary, Freelance, Investments)"
echo "   • Monthly budget with spending tracking"
echo "   • Dashboard with charts and analytics"
echo ""
echo "🛠️  Useful commands:"
echo "   • View logs: docker-compose -f docker-compose.demo.yml logs -f"
echo "   • Stop demo: docker-compose -f docker-compose.demo.yml down"
echo "   • Restart: docker-compose -f docker-compose.demo.yml restart"
echo ""
echo "🌐 Service URLs (for development):"
echo "   • Frontend: http://localhost"
echo "   • Auth API: http://localhost:8081"
echo "   • Financial API: http://localhost:8082"
echo "   • MongoDB: localhost:27017"
echo "   • Redis: localhost:6379"
echo ""