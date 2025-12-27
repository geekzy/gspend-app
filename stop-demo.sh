#!/bin/bash

echo "🛑 Stopping GSpend Demo Environment..."

# Stop and remove containers
docker-compose -f docker-compose.demo.yml down

echo "✅ Demo environment stopped."
echo ""
echo "💡 To completely remove demo data:"
echo "   docker-compose -f docker-compose.demo.yml down -v"
echo ""
echo "🚀 To start again:"
echo "   ./start-demo.sh"