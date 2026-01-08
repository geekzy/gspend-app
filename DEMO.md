# GSpend Demo Environment

This demo environment provides a complete GSpend application with realistic dummy data for dashboard preview and testing.

## Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Ports 80, 8081, 8082, 27017, and 6379 available

### Start Demo
```bash
make demo-start
```

### Stop Demo
```bash
make demo-stop
```

### Other Demo Commands
```bash
make demo-restart     # Restart demo environment
make demo-clean       # Clean demo environment and data
make demo-logs        # View demo logs
make demo-status      # Check demo status
make help             # See all available commands
```

## Demo Data

The demo environment includes:

### 👤 Demo User
- **Email:** demo@gspend.com
- **Password:** password

### 💰 Financial Data
- **3 months** of sample transactions (income and expenses)
- **8 expense categories:** Food & Dining, Transportation, Shopping, Entertainment, Bills & Utilities, Healthcare, Education, Travel
- **4 income categories:** Salary, Freelance, Investment, Other Income
- **Monthly budget** with realistic spending vs. planned amounts
- **Income records:** Regular salary, freelance projects, investment returns

### 📊 Dashboard Features
- Expense breakdown by category (pie chart)
- Monthly spending trends (line chart)
- Budget vs. actual spending comparison
- Recent transactions list
- Income vs. expense summary

## Access Points

- **Frontend Dashboard:** http://localhost
- **Auth Service API:** http://localhost:8081
- **Financial Service API:** http://localhost:8082
- **MongoDB:** localhost:27017
- **Redis:** localhost:6379

## Sample Data Details

### Transactions
- **~200 transactions** across 3 months
- Realistic amounts and descriptions
- Various payment methods (Credit Card, Debit Card, Cash, etc.)
- Proper categorization

### Categories
- **System categories** with icons and colors
- **Expense categories:** Food (🍽️), Transport (🚗), Shopping (🛍️), etc.
- **Income categories:** Salary (💰), Freelance (💻), Investment (📈), etc.

### Budget
- **Monthly budget** for current month
- **5 main expense categories** with planned vs. spent amounts
- Realistic budget allocations

### Income
- **Monthly salary:** $5,000
- **Freelance project:** $1,200 (one-time)
- **Investment returns:** $300/month

## Development

### View Logs
```bash
docker-compose -f docker-compose.demo.yml logs -f
```

### Restart Services
```bash
docker-compose -f docker-compose.demo.yml restart
```

### Clean Reset
```bash
make demo-clean
make demo-start
```

### Individual Service Logs
```bash
# Auth service
docker-compose -f docker-compose.demo.yml logs -f auth-service

# Financial service
docker-compose -f docker-compose.demo.yml logs -f financial-service

# Frontend
docker-compose -f docker-compose.demo.yml logs -f frontend
```

## API Testing

### Login to get JWT token
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@gspend.com","password":"password"}'
```

### Get dashboard data
```bash
# Replace YOUR_JWT_TOKEN with the token from login
curl -X GET http://localhost:8082/api/v1/dashboard \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get transactions
```bash
curl -X GET http://localhost:8082/api/v1/transactions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Troubleshooting

### Services not starting
1. Check if ports are available: `lsof -i :80,8081,8082,27017,6379`
2. Ensure Docker has enough resources (4GB+ RAM recommended)
3. Check logs: `docker-compose -f docker-compose.demo.yml logs`

### Frontend not loading
1. Wait for all services to be healthy (can take 2-3 minutes)
2. Check nginx logs: `docker-compose -f docker-compose.demo.yml logs nginx`
3. Verify backend services: `curl http://localhost:8081/api/v1/auth/health`

### Database connection issues
1. Check MongoDB logs: `docker-compose -f docker-compose.demo.yml logs mongodb`
2. Verify data seeding completed: `docker-compose -f docker-compose.demo.yml logs data-seeder`

### Reset everything
```bash
make demo-clean
make demo-start
```

## Architecture

The demo environment runs:
1. **MongoDB** - Database with seeded dummy data
2. **Redis** - Caching layer
3. **Data Seeder** - One-time container that populates the database
4. **Auth Service** - User authentication and JWT management
5. **Financial Service** - Core financial data management
6. **Frontend** - Vue.js dashboard application
7. **Nginx** - Reverse proxy and static file serving

All services run in Docker containers with proper health checks and dependency management.