# HRMS Setup Guide

## Quick Start

### Option 1: Docker (Recommended)

1. **Create environment file** (optional - Docker uses defaults):
   ```bash
   cp .env.example .env
   # Edit .env if needed (for local dev, change DB_HOST=localhost)
   ```

2. **Start all services**:
   ```bash
   docker-compose up -d
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health

### Option 2: Local Development

#### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+ (or use Docker for DB only)

#### Setup Steps

1. **Start PostgreSQL** (if not using Docker):
   ```bash
   # Option A: Use Docker for PostgreSQL only
   docker-compose up -d postgres
   
   # Option B: Use local PostgreSQL
   # Make sure PostgreSQL is running and create database:
   createdb hrms
   ```

2. **Backend Setup**:
   ```bash
   cd backend
   
   # Create .env file
   cp ../.env.example .env
   
   # For local dev, update DB_HOST in .env:
   # DB_HOST=localhost
   
   # Download dependencies
   go mod download
   
   # Run backend
   go run main.go
   ```

3. **Frontend Setup** (new terminal):
   ```bash
   cd frontend
   
   # Create .env file
   cp .env.example .env
   
   # Install dependencies
   npm install
   
   # Start dev server
   npm start
   ```

## Environment Variables

### Backend (.env in project root or backend directory)

Key variables to configure:

- **Database**: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- **JWT Secrets**: `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET` (CHANGE IN PRODUCTION!)
- **OpenAI** (optional): `OPENAI_API_KEY` for AI assistant
- **SMTP** (optional): Email configuration for notifications

### Frontend (.env in frontend directory)

- **API URL**: `REACT_APP_API_URL` (default: http://localhost:8080/api)

## First User Setup

After starting the application, register your first admin user:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@hrms.com",
    "password": "admin123",
    "role": "admin"
  }'
```

Then login at http://localhost:3000

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL status
docker-compose ps postgres

# View PostgreSQL logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Backend Issues

```bash
# Check backend logs
docker-compose logs backend

# Verify Go dependencies
cd backend
go mod tidy
```

### Frontend Issues

```bash
# Clear and reinstall dependencies
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Port Conflicts

If ports 3000, 8080, or 5432 are already in use:

- **Frontend**: Change port in `frontend/package.json` or use `PORT=3001 npm start`
- **Backend**: Change `SERVER_PORT` in `.env`
- **PostgreSQL**: Change port mapping in `docker-compose.yml`

## Production Deployment

For production:

1. **Change JWT secrets** to strong random strings
2. **Set ENVIRONMENT=production**
3. **Configure proper database credentials**
4. **Enable SMTP** for email notifications
5. **Set up SSL/TLS** for database connections
6. **Configure New Relic** if using observability
7. **Set up proper file storage** (S3 for production)

## Support

For issues or questions, check:
- README.md for API documentation
- PRD.md for feature requirements
- Logs: `docker-compose logs -f`
