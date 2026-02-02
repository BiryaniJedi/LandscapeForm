# LandscapeForm

> Enterprise-grade form management system for pesticide application tracking and compliance

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://go.dev/)
[![Node Version](https://img.shields.io/badge/Node-18+-339933?logo=node.js)](https://nodejs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-336791?logo=postgresql)](https://www.postgresql.org/)

---

## Overview

LandscapeForm is a comprehensive digital form management system designed for landscape maintenance companies to track pesticide applications, manage chemical inventories, and maintain compliance with regulatory requirements. The application provides a secure, user-friendly interface for employees to document their work while giving administrators powerful oversight and reporting capabilities.

### Key Features

- **Digital Form Management**: Create, edit, and manage shrub and lawn pesticide application forms
- **Chemical Database**: Centralized repository of EPA-registered chemicals with brand names, active ingredients, and application rates
- **User Management**: Role-based access control with employee approval workflows
- **PDF Generation**: Professional PDF exports for customer records and regulatory compliance
- **Audit Trail**: Comprehensive tracking of all form submissions and modifications
- **Responsive Design**: Accessible from desktop, tablet, and mobile devices
- **Secure Authentication**: JWT-based authentication with encrypted password storage

---

## Technology Stack

### Backend
- **Language**: Go 1.25+
- **Web Framework**: Chi v5
- **Database**: PostgreSQL 16+
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt

### Frontend
- **Framework**: Next.js 16 (React 19)
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS 4
- **PDF Generation**: React-PDF
- **State Management**: React Context API

### Infrastructure
- **Web Server**: Nginx
- **Process Manager**: systemd (backend), PM2 (frontend)
- **SSL/TLS**: Let's Encrypt
- **CI/CD**: GitHub Actions

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Client Browser                     │
│            (Desktop / Tablet / Mobile)               │
└────────────────────┬────────────────────────────────┘
                     │ HTTPS
                     ↓
┌─────────────────────────────────────────────────────┐
│                  Nginx Reverse Proxy                 │
│              (SSL Termination, Routing)              │
└────────┬───────────────────────────────────┬────────┘
         │                                   │
         ↓                                   ↓
┌─────────────────┐              ┌──────────────────┐
│  Next.js Server │              │   Go API Server  │
│   (Port 3000)   │──────────────│   (Port 8000)    │
│                 │   API Calls  │                  │
└─────────────────┘              └────────┬─────────┘
                                          │
                                          ↓
                               ┌────────────────────┐
                               │   PostgreSQL DB    │
                               │    (Port 5432)     │
                               └────────────────────┘
```

---

## Getting Started

### Prerequisites

- **Go** 1.25 or higher ([download](https://go.dev/dl/))
- **Node.js** 18 or higher ([download](https://nodejs.org/))
- **PostgreSQL** 16 or higher ([download](https://www.postgresql.org/download/))
- **Git** ([download](https://git-scm.com/downloads))

### Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/landscapeform.git
cd landscapeform
```

#### 2. Database Setup

**Start PostgreSQL** (using Docker):
```bash
cd db
docker-compose up -d
```

**Or configure local PostgreSQL**:
```sql
CREATE DATABASE landscapeform;
CREATE USER landscapeform_user WITH ENCRYPTED PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE landscapeform TO landscapeform_user;
```

**Run migrations**:
```bash
psql -U landscapeform_user -d landscapeform -h localhost < db/migrations/schema.sql
```

#### 3. Backend Setup

```bash
cd backend

# Create environment file
cp .env.example .env

# Edit .env with your configuration:
# PORT=8000
# DATABASE_URL=postgresql://landscapeform_user:password@localhost:5432/landscapeform?sslmode=disable
# JWT_SECRET=your-secret-key-here

# Download dependencies
go mod download

# Run the server
go run ./cmd/api
```

The API server will start on `http://localhost:8000`

#### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Create environment file
cp .env.local.example .env.local

# Edit .env.local:
# NEXT_PUBLIC_API_URL=http://localhost:8000/api

# Run development server
npm run dev
```

The application will be available at `http://localhost:3000`

---

## Development

### Project Structure

```
landscapeform/
├── backend/                    # Go backend API
│   ├── cmd/
│   │   └── api/
│   │       └── main.go        # Application entry point
│   ├── internal/
│   │   ├── auth/              # JWT authentication logic
│   │   ├── db/                # Database connection pooling
│   │   ├── dto/               # Data transfer objects
│   │   ├── forms/             # Form business logic
│   │   ├── chemicals/         # Chemical database logic
│   │   ├── users/             # User management logic
│   │   ├── handlers/          # HTTP request handlers
│   │   └── middleware/        # HTTP middleware (auth, CORS, logging)
│   ├── go.mod                 # Go module dependencies
│   └── .env                   # Environment variables (gitignored)
│
├── frontend/                  # Next.js frontend application
│   ├── app/                   # Next.js App Router
│   │   ├── (auth)/           # Authentication pages
│   │   ├── dashboard/        # User dashboard
│   │   ├── forms/            # Form management pages
│   │   ├── admin/            # Admin pages
│   │   ├── settings/         # User settings
│   │   ├── layout.tsx        # Root layout
│   │   └── page.tsx          # Landing page
│   ├── lib/
│   │   ├── api/              # API client functions
│   │   ├── components/       # Shared React components
│   │   ├── common/           # Utility functions
│   │   └── pdf/              # PDF generation components
│   ├── public/               # Static assets
│   ├── package.json          # NPM dependencies
│   └── .env.local            # Environment variables (gitignored)
│
├── db/
│   ├── migrations/           # Database schema migrations
│   └── docker-compose.yml    # PostgreSQL development setup
│
├── .github/
│   └── workflows/
│       └── deploy.yml        # CI/CD pipeline
│
├── PROJECT.md                # Database schema documentation
├── DEPLOYMENT_PROPOSAL.md    # Deployment options and costs
├── DEPLOYMENT_SUMMARY.md     # One-page deployment summary
└── README.md                 # This file
```

### API Endpoints

#### Authentication
```
POST   /api/auth/register      Create new user account
POST   /api/auth/login         Authenticate user
POST   /api/auth/logout        End user session
GET    /api/auth/me            Get current user details
```

#### Forms
```
GET    /api/forms              List user's forms (paginated)
POST   /api/forms/shrub        Create shrub application form
POST   /api/forms/lawn         Create lawn application form
GET    /api/forms/{id}         Get form by ID
PUT    /api/forms/shrub/{id}   Update shrub form
PUT    /api/forms/lawn/{id}    Update lawn form
DELETE /api/forms/{id}         Delete form
GET    /api/forms/{id}/print   Get form for PDF export
```

#### Chemicals
```
GET    /api/chemicals                      List all chemicals
GET    /api/chemicals/category/{category}  List chemicals by category
POST   /api/admin/chemicals                Create chemical (admin only)
PUT    /api/admin/chemicals/{id}           Update chemical (admin only)
DELETE /api/admin/chemicals/{id}           Delete chemical (admin only)
```

#### Users (Admin Only)
```
GET    /api/users              List all users
GET    /api/users/{id}         Get user by ID
PUT    /api/users/{id}         Update user
DELETE /api/users/{id}         Delete user
POST   /api/users/{id}/approve Approve pending user
```

### Database Schema

See [PROJECT.md](PROJECT.md) for detailed database schema documentation.

Key tables:
- `users` - User accounts and authentication
- `forms` - Form metadata (type, status, ownership)
- `shrubs` - Shrub application details
- `lawns` - Lawn application details
- `pesticide_applications` - Chemical applications per form
- `chemicals` - Chemical database (EPA-registered products)

---

## Testing

### Backend Tests

```bash
cd backend
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/auth
```

### Frontend Tests

```bash
cd frontend
npm test

# Run with coverage
npm test -- --coverage

# Run in watch mode
npm test -- --watch
```

### Manual API Testing

```bash
cd backend
./test_api.sh
```

---

## Building for Production

### Backend

```bash
cd backend
go build -o api ./cmd/api

# Cross-compile for Linux (from macOS/Windows)
GOOS=linux GOARCH=amd64 go build -o api ./cmd/api
```

### Frontend

```bash
cd frontend
npm run build

# Test production build locally
npm start
```

---

## Deployment

For detailed deployment instructions, see the deployment documentation:
- **Cost Analysis**: See [DEPLOYMENT_PROPOSAL.md](DEPLOYMENT_PROPOSAL.md)
- **Quick Summary**: See [DEPLOYMENT_SUMMARY.md](DEPLOYMENT_SUMMARY.md)
- **Technical Guide**: Contact the development team for detailed deployment procedures

### Production Requirements

- **Server**: 4GB RAM, 2 vCPUs, 80GB SSD (recommended)
- **OS**: Ubuntu 22.04 LTS or similar
- **Domain**: Custom domain with SSL certificate
- **Database**: PostgreSQL 16+ (can be on same server or managed service)

### Environment Variables

**Backend** (`.env`):
```bash
PORT=8000
DATABASE_URL=postgresql://user:password@host:port/database?sslmode=disable
JWT_SECRET=your-secure-secret-key
```

**Frontend** (`.env.local`):
```bash
NEXT_PUBLIC_API_URL=https://yourdomain.com/api
```

---

## Security

### Authentication & Authorization
- JWT token-based authentication with HTTP-only cookies
- Passwords hashed using bcrypt (cost factor: 10)
- Role-based access control (employee, admin)
- User approval workflow for new accounts

### Data Protection
- HTTPS/TLS encryption for all data in transit
- Database encryption at rest
- Parameterized SQL queries (protection against SQL injection)
- Input validation and sanitization
- CORS configuration for API security

### Security Best Practices
- Regular dependency updates
- Security headers (X-Frame-Options, X-Content-Type-Options, etc.)
- Rate limiting on API endpoints
- Session timeout and token expiration
- Audit logging of critical operations

---

## Monitoring & Maintenance

### Application Logs

**Backend**:
```bash
# Production (systemd)
sudo journalctl -u landscapeform-api -f

# Development
go run ./cmd/api
```

**Frontend**:
```bash
# Production (PM2)
pm2 logs landscapeform-frontend

# Development
npm run dev
```

### Health Checks

- **API Health**: `GET /api/health`
- **Database Connection**: Automatic connection pooling with health checks
- **Uptime Monitoring**: Configure with external service (UptimeRobot, Pingdom, etc.)

### Backup Strategy

- **Database**: Automated daily backups with 7-day retention
- **Application Files**: Version controlled via Git
- **User Uploads**: Included in database backups (if file storage added)

---

## Performance

### Optimization Features

- **Database**: Connection pooling, indexed queries, optimized schema
- **Frontend**: Code splitting, lazy loading, optimized images
- **Caching**: Static asset caching, API response caching (where appropriate)
- **CDN**: Cloudflare or similar for static assets (production)

### Capacity

- **Concurrent Users**: 50+ simultaneous users on recommended hardware
- **Response Time**: < 200ms API response time (average)
- **Database**: Supports 100,000+ records without performance degradation
- **Scalability**: Horizontal scaling supported via load balancer

---

## Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Mobile browsers (iOS Safari, Chrome Mobile)

---

## Contributing

### Development Workflow

1. Create a feature branch from `main`
2. Make your changes with descriptive commit messages
3. Test thoroughly (unit tests, integration tests, manual testing)
4. Submit a pull request with detailed description
5. Code review and approval required before merge

### Commit Message Format

```
type(scope): brief description

Detailed explanation if needed

Fixes #issue-number
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Code Style

- **Go**: Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- **TypeScript/React**: Follow [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- **Formatting**: `gofmt` for Go, `prettier` for TypeScript/React

---

## Roadmap

### Current Version (v1.0)
- ✅ User authentication and authorization
- ✅ Form creation and management (shrub, lawn)
- ✅ Chemical database management
- ✅ PDF export functionality
- ✅ Admin dashboard
- ✅ Responsive design

### Future Enhancements (v1.1+)
- 📋 Form templates and customization
- 📊 Reporting and analytics dashboard
- 📱 Progressive Web App (PWA) for offline capability
- 🔔 Email notifications for form approvals
- 📈 Data export to Excel/CSV
- 🔍 Advanced search and filtering
- 🗓️ Calendar view for scheduled applications
- 📷 Photo upload for application sites
- 🌐 Multi-language support

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Support

### Documentation
- [Database Schema](PROJECT.md)
- [Deployment Guide](DEPLOYMENT_PROPOSAL.md)
- [API Documentation](docs/API.md) *(coming soon)*

### Contact

For technical support, feature requests, or bug reports:
- **Email**: [Your support email]
- **Issue Tracker**: [GitHub Issues](https://github.com/YOUR_USERNAME/landscapeform/issues)
- **Documentation**: [Wiki](https://github.com/YOUR_USERNAME/landscapeform/wiki)

### Professional Services

For enterprise deployments, custom development, or training:
- Contact: [Your professional email]
- Website: [Your website]

---

## Acknowledgments

Built with modern, enterprise-grade technologies:
- [Go](https://go.dev/) - Backend programming language
- [Chi](https://github.com/go-chi/chi) - Lightweight HTTP router
- [Next.js](https://nextjs.org/) - React framework
- [PostgreSQL](https://www.postgresql.org/) - Relational database
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- [React PDF](https://react-pdf.org/) - PDF generation library

---

<p align="center">
  <strong>LandscapeForm</strong> - Professional pesticide application tracking made simple
</p>

<p align="center">
  Made with ❤️ for landscape professionals
</p>
