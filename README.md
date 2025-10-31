# Crypsis - Enterprise File Encryption & Storage Service

A comprehensive, production-ready file encryption and secure storage solution with a modern React admin dashboard, built with Go backend, featuring OAuth2 authentication, KMS integration, and distributed storage capabilities.

## 🌟 Overview

Crypsis is an enterprise-grade platform that provides:
- **Secure File Storage**: End-to-end encrypted file management
- **Modern Admin Dashboard**: React + TypeScript frontend with beautiful UI
- **OAuth2 Authentication**: Industry-standard authentication flow
- **Key Management**: File-based or external KMS integration
- **Distributed Storage**: MinIO-backed object storage
- **Comprehensive Monitoring**: Built-in observability with Grafana, Prometheus, and Jaeger
- **Production-Ready**: Docker-compose deployment, horizontal scalability

## 🎯 Key Features

### Backend (Go + Gin)
- ✅ RESTful API with comprehensive endpoints
- ✅ AES-256-GCM encryption for all files
- ✅ OAuth2 + Hydra integration for secure authentication
- ✅ PostgreSQL for metadata storage
- ✅ MinIO for distributed file storage
- ✅ OpenTelemetry instrumentation for observability
- ✅ File integrity verification (SHA-256/512 hashing)
- ✅ Audit logging for all operations
- ✅ Key rotation and re-encryption support

### Frontend (React + TypeScript)
- ✅ Modern, responsive UI built with Tailwind CSS
- ✅ State management with Zustand
- ✅ Type-safe with TypeScript
- ✅ File upload with progress tracking
- ✅ Admin dashboard with system statistics
- ✅ Application (OAuth2 client) management
- ✅ User management interface
- ✅ Audit log viewer
- ✅ Security settings and re-keying
- ✅ Production-ready with Docker deployment

## 🏗️ Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                      Frontend (React)                         │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────────┐ │
│  │  Dashboard  │  │  File Mgmt   │  │  Admin/Apps/Logs    │ │
│  └─────────────┘  └──────────────┘  └─────────────────────┘ │
└───────────────────────────┬──────────────────────────────────┘
                            │ HTTPS/REST API
┌───────────────────────────┴──────────────────────────────────┐
│                   Backend (Go + Gin)                          │
│  ┌────────────┐  ┌─────────────┐  ┌─────────────────────┐   │
│  │  Handlers  │──│  Services   │──│   Repositories      │   │
│  └────────────┘  └─────────────┘  └─────────────────────┘   │
│  │ Auth │ File │  │ Encryption  │  │  Database Access    │   │
│  │ Admin│ Logs │  │ KMS │ OAuth2│  │  Storage I/O        │   │
└───────────────────────────┬──────────────────────────────────┘
                            │
       ┌────────────────────┼────────────────────┐
       │                    │                    │
┌──────▼─────┐    ┌─────────▼────────┐   ┌──────▼──────┐
│PostgreSQL  │    │      MinIO       │   │    Hydra    │
│ (Metadata) │    │   (File Store)   │   │  (OAuth2)   │
└────────────┘    └──────────────────┘   └─────────────┘
       │                    │                    │
       └────────────────────┴────────────────────┘
                            │
              ┌─────────────┴────────────────┐
              │   Observability Stack        │
              │ Prometheus│Jaeger│Grafana    │
              └──────────────────────────────┘
```

## 📋 Complete API Reference

### Public Endpoints
| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/admin/login` | Admin login | `{username, password}` |

### Client File Operations (OAuth2 Token Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/files` | Upload encrypted file |
| GET | `/api/files/{id}/download` | Download file |
| PUT | `/api/files/{id}/update` | Update file |
| DELETE | `/api/files/{id}/delete` | Delete file |
| GET | `/api/files/list` | List all files |
| GET | `/api/files/{id}/metadata` | Get file metadata |
| POST | `/api/files/encrypt` | Encrypt file (client-side) |
| POST | `/api/files/decrypt` | Decrypt file (client-side) |

### Admin Management (Admin Token Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/logout` | Logout |
| GET | `/api/admin/refresh-token` | Refresh access token |
| GET | `/api/admin/list` | List admins |
| PATCH | `/api/admin/username` | Update username |
| PATCH | `/api/admin/password` | Update password |
| DELETE | `/api/admin?id={id}` | Delete admin |
| POST | `/api/admin/add` | Add new admin |

### Application Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/apps` | Create OAuth2 application |
| GET | `/api/admin/apps` | List all applications |
| GET | `/api/admin/apps/{id}` | Get application details |
| DELETE | `/api/admin/apps/{id}` | Delete application |
| POST | `/api/admin/apps/{id}/recover` | Recover deleted app |
| PUT | `/api/admin/apps/{id}/rotate-secret` | Rotate client secret |

### File & Log Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/files` | List all files (admin view) |
| GET | `/api/admin/apps/{id}/files` | List files by application |
| GET | `/api/admin/logs` | View audit logs |
| POST | `/api/admin/files/re-key` | Re-encrypt with new key |

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd Crypsis

# Start all services
docker-compose up -d

# Initialize the database
docker-compose exec app ./scripts/init-db.sh

# Access the applications
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# MinIO Console: http://localhost:9001
# Grafana: http://localhost:3000 (monitoring)
```

### Option 2: Manual Setup

#### Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
# (Ensure PostgreSQL is running)

# Start the backend server
go run cmd/main.go
```

#### Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Set up environment variables
cp .env.example .env
# Edit .env: VITE_API_URL=http://localhost:8080

# Start development server
npm run dev

# Or build for production
npm run build
npm run preview
```

## ⚙️ Configuration

### Backend Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=crypsis
DB_PASSWORD=secure_password
DB_NAME=crypsis_db
DB_SSLMODE=disable

# Storage Configuration (MinIO)
STORAGE_ENDPOINT=localhost:9000
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_SSL=false
BUCKET_NAME=crypsis-files

# Security Configuration
HASH_METHOD=SHA256
ENC_METHOD=AES-256-GCM
HASH_ENCRYPTED_FILE=true
MKEY_PATH=./resources/master.key

# OAuth2 Configuration (Hydra)
HYDRA_PUBLIC_URL=http://localhost:4444
HYDRA_ADMIN_URL=http://localhost:4445

# Optional KMS Configuration
KMS_ENABLE=false
KMS_KEY_UID=your-kms-key-id
KMS_URL=https://kms.example.com
KEY_PATH=./cosmian/kms.key
CERT_PATH=./cosmian/kms.crt
CA_PATH=./cosmian/kms.server.p12

# Server Configuration
PORT=8080
GIN_MODE=release
```

### Frontend Environment Variables

```bash
# API Configuration
VITE_API_URL=http://localhost:8080

# Application Configuration
VITE_APP_NAME=Crypsis
VITE_APP_VERSION=1.0.0

# Feature Flags
VITE_ENABLE_FILE_ENCRYPTION=true
VITE_MAX_FILE_SIZE=104857600
```

## 🎨 Frontend Features

### Dashboard
- System statistics overview
- Recent activity feed
- Quick actions
- Storage usage charts

### File Management
- Drag-and-drop file upload
- Upload progress tracking
- File list with sorting and filtering
- Bulk operations
- File metadata viewer
- Download/update/delete operations

### Admin Management
- Create/edit/delete admin users
- Password management
- Role-based access control

### Application Management
- Register OAuth2 applications
- View client credentials
- Rotate secrets
- Manage redirect URIs
- Monitor application usage

### Audit Logs
- Comprehensive activity tracking
- Filtering by date, user, action
- Export logs
- Real-time updates

### Security Settings
- Key rotation (re-keying)
- Encryption method configuration
- Security alerts
- Access control policies

## 🔒 Security Features

### Encryption
- **AES-256-GCM**: Industry-standard authenticated encryption
- **Client-side encryption**: Optional pre-upload encryption
- **At-rest encryption**: All files encrypted before storage
- **Key derivation**: PBKDF2 or KMS-based key management

### Authentication & Authorization
- **OAuth2**: Standard protocol implementation
- **JWT tokens**: Secure, stateless authentication
- **Token introspection**: Real-time validation via Hydra
- **Role-based access**: Admin vs. client permissions

### Audit & Compliance
- **Complete audit trail**: All operations logged
- **Tamper-proof logs**: Immutable audit records
- **IP tracking**: Source IP for all operations
- **User agent logging**: Device/browser information

## 📊 Monitoring & Observability

### Quick Start Monitoring

```bash
# Start observability stack
./start-observability.sh

# Access dashboards
# Grafana: http://localhost:3000 (admin/admin)
# Prometheus: http://localhost:9090
# Jaeger: http://localhost:16686
```

### Available Metrics
- HTTP request count and duration
- File upload/download rates
- Encryption/decryption performance
- Database query performance
- Storage usage
- Active connections
- Error rates

### Distributed Tracing
- End-to-end request tracing
- Service dependency mapping
- Performance bottleneck identification
- Error correlation

## 🧪 Testing

### Backend Tests

```bash
cd backend

# Run unit tests
go test ./internal/...

# Run integration tests
go test -tags=integration ./test/...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/services -run TestFileService
```

### Frontend Tests

```bash
cd frontend

# Run type checking
npm run type-check

# Run linter
npm run lint

# Build to verify
npm run build
```

### Performance Testing

```bash
cd performance_test

# Run smoke test
k6 run scripts/k6_smoke_test.js

# Run load test
k6 run scripts/k6_load_test.js

# Run stress test
k6 run scripts/k6_stress_test.js

# Run spike test
k6 run scripts/k6_spike_test.js
```

## 📦 Deployment

### Docker Production Deployment

```bash
# Build images
docker-compose build

# Deploy to production
docker-compose -f docker-compose.prod.yml up -d

# Scale services
docker-compose up -d --scale app=3
```

### Frontend Production Build

```bash
cd frontend

# Build optimized production bundle
npm run build

# Preview production build
npm run preview

# Deploy dist/ folder to your web server
# or use the Docker image
```

### Backend Binary Deployment

```bash
cd backend

# Build binary
go build -o crypsis-server cmd/main.go

# Run binary
./crypsis-server
```

## 🔧 Development

### Project Structure

```
Crypsis/
├── backend/
│   ├── cmd/                    # Application entry point
│   ├── internal/
│   │   ├── config/            # Configuration
│   │   ├── delivery/http/     # HTTP handlers & routes
│   │   ├── entity/            # Database models
│   │   ├── helper/            # Utility functions
│   │   ├── model/             # Request/Response models
│   │   ├── repository/        # Data access layer
│   │   ├── services/          # Business logic
│   │   └── delivery/middlewere/ # HTTP middleware
│   ├── scripts/               # Helper scripts
│   └── test/                  # Test files
├── frontend/
│   ├── public/                # Static assets
│   ├── src/
│   │   ├── components/        # React components
│   │   │   ├── features/     # Feature-specific components
│   │   │   ├── layout/       # Layout components
│   │   │   └── ui/           # Reusable UI components
│   │   ├── pages/            # Page components
│   │   ├── services/         # API services
│   │   ├── stores/           # Zustand state stores
│   │   ├── types/            # TypeScript type definitions
│   │   ├── utils/            # Utility functions
│   │   └── constants/        # Application constants
│   ├── Dockerfile            # Frontend Docker image
│   ├── nginx.conf            # Nginx configuration
│   └── package.json          # NPM dependencies
├── config/                    # Configuration files
├── data/                      # Persistent data
├── docker-compose.yaml        # Service orchestration
└── README.md                  # This file
```

### Adding New Features

#### Backend
1. Define interface in `internal/interfaces/`
2. Implement service in `internal/services/`
3. Add handler in `internal/delivery/http/`
4. Define routes in `routes.go`
5. Add tests in `test/`

#### Frontend
1. Define types in `src/types/`
2. Create API service in `src/services/`
3. Create Zustand store in `src/stores/`
4. Build UI components in `src/components/`
5. Create page in `src/pages/`

## 📖 API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data
  }
}
```

### Success with Count
```json
{
  "success": true,
  "message": "Data retrieved",
  "count": 42,
  "data": [...]
}
```

### Error Response
```json
{
  "success": false,
  "message": "Operation failed",
  "error": "Detailed error message"
}
```

## 🎓 Usage Examples

### Upload File (cURL)

```bash
# Get access token first (OAuth2)
TOKEN="your_access_token"

# Upload file
curl -X POST http://localhost:8080/api/files \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/your/file.pdf"
```

### Admin Login

```bash
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

### List Files

```bash
curl -X GET "http://localhost:8080/api/files/list?offset=0&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

### Create OAuth2 Application

```bash
curl -X POST http://localhost:8080/api/admin/apps \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My App",
    "uri": "https://myapp.com",
    "redirectUri": "https://myapp.com/callback"
  }'
```

## 🆘 Troubleshooting

### Backend Issues

**Database connection failed**
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check connection
psql -h localhost -U crypsis -d crypsis_db
```

**MinIO connection failed**
```bash
# Check MinIO is running
docker-compose ps minio

# Access MinIO console
# http://localhost:9001
```

### Frontend Issues

**Cannot reach API**
- Verify `VITE_API_URL` in `.env`
- Check backend is running: `curl http://localhost:8080/health`
- Check CORS configuration in backend

**Build errors**
```bash
# Clear cache and rebuild
rm -rf node_modules package-lock.json
npm install
npm run build
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## � Acknowledgments

- Go Gin framework
- React and the React ecosystem
- MinIO for object storage
- Ory Hydra for OAuth2
- OpenTelemetry for observability
- Tailwind CSS for styling

## 📞 Support

For questions, issues, or contributions:
- Create an issue on GitHub
- Contact the development team
- Check the documentation in `/docs`

---

**Built with ❤️ for enterprise-grade security**

### Core Functionality
- **Secure File Storage**: Upload, download, update, and delete files with enterprise-level security
- **End-to-End Encryption**: Client-side file encryption/decryption with AES-256
- **Key Management Service (KMS)**: Optional integration with external KMS for enhanced security
- **OAuth2 Authentication**: Industry-standard authentication and authorization
- **Admin Management**: Complete administrative interface for user and application management
- **Audit Logging**: Comprehensive activity tracking and monitoring

### Security Features
- **Multiple Encryption Methods**: Configurable encryption algorithms
- **Hash Verification**: File integrity verification using configurable hash methods
- **Token-based Authentication**: JWT tokens with introspection
- **Access Control**: Role-based access control for files and applications
- **Secure Key Storage**: File-based or KMS-based key management

### Storage & Infrastructure
- **MinIO Integration**: Distributed object storage backend
- **PostgreSQL Database**: Reliable metadata and configuration storage
- **Docker Support**: Complete containerization with Docker Compose
- **Nginx Proxy**: Load balancing and SSL termination
- **Hydra OAuth2**: Scalable OAuth2 server integration

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Admin Panel   │    │  Client Apps    │    │      KMS        │
│   (Frontend)    │    │  (OAuth2 Apps)  │    │   (Optional)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Crypsis API   │
                    │  (Go + Gin)     │
                    └─────────────────┘
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  PostgreSQL     │    │     MinIO       │    │    Hydra        │
│  (Metadata)     │    │   (Storage)     │    │   (OAuth2)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📋 API Endpoints

### Public Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/login` | Admin authentication |

### Client File Operations (Requires OAuth2 Token)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/files` | Upload a file |
| GET | `/api/files/{id}/download` | Download a file |
| PUT | `/api/files/{id}/update` | Update a file |
| DELETE | `/api/files/{id}/delete` | Delete a file |
| GET | `/api/files/list` | List all files |
| GET | `/api/files/{id}/metadata` | Get file metadata |
| POST | `/api/files/encrypt` | Encrypt a file |
| POST | `/api/files/decrypt` | Decrypt a file |

### Admin Management (Requires Admin Token)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/logout` | Admin logout |
| GET | `/api/admin/refresh-token` | Refresh access token |
| GET | `/api/admin/list` | List all admins |
| PATCH | `/api/admin/username` | Update admin username |
| PATCH | `/api/admin/password` | Update admin password |
| DELETE | `/api/admin` | Delete admin |
| POST | `/api/admin/add` | Add new admin |

### Application Management (Requires Admin Token)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/apps` | Create new application |
| GET | `/api/admin/apps` | List all applications |
| GET | `/api/admin/apps/{id}` | Get application details |
| DELETE | `/api/admin/apps/{id}` | Delete application |
| POST | `/api/admin/apps/{id}/recover` | Recover deleted application |
| PUT | `/api/admin/apps/{id}/rotate-secret` | Rotate application secret |

### File Management (Requires Admin Token)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/files` | List all files |
| GET | `/api/admin/apps/{id}/files` | List files by application |
| GET | `/api/admin/logs` | View audit logs |
| POST | `/api/admin/files/re-key` | Re-encrypt files with new key |

## 🛠️ Installation & Setup

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for development)
- PostgreSQL 13+
- MinIO Server

### Quick Start with Docker

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Crypsis
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

3. **Initialize the database**
   ```bash
   docker-compose exec app ./scripts/init-db.sh
   ```

4. **Access the service**
   - API: `http://localhost:8080`
   - MinIO Console: `http://localhost:9001`

### Development Setup

1. **Install dependencies**
   ```bash
   cd src
   go mod download
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run the application**
   ```bash
   go run cmd/main.go
   ```

## ⚙️ Configuration

### Environment Variables

#### Database Configuration
- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

#### Storage Configuration
- `STORAGE_ENDPOINT`: MinIO endpoint
- `STORAGE_ACCESS_KEY`: MinIO access key
- `STORAGE_SECRET_KEY`: MinIO secret key
- `STORAGE_SSL`: Enable SSL for storage (true/false)
- `BUCKET_NAME`: Default bucket name

#### Security Configuration
- `HASH_METHOD`: File hashing method (SHA256, SHA512)
- `ENC_METHOD`: Encryption method (AES-256-GCM)
- `HASH_ENCRYPTED_FILE`: Hash encrypted files (true/false)
- `MKEY_PATH`: Master key file path

#### OAuth2 Configuration
- `HYDRA_PUBLIC_URL`: Hydra public endpoint
- `HYDRA_ADMIN_URL`: Hydra admin endpoint

#### KMS Configuration (Optional)
- `KMS_ENABLE`: Enable KMS integration (true/false)
- `KMS_KEY_UID`: KMS key identifier
- `KMS_URL`: KMS service URL
- `KEY_PATH`: Client key file path
- `CERT_PATH`: Client certificate path
- `CA_PATH`: CA certificate path

## 🔐 Authentication & Authorization

### OAuth2 Flow

1. **Application Registration**: Admin creates OAuth2 applications
2. **Token Generation**: Applications authenticate to receive access tokens
3. **API Access**: Include Bearer token in Authorization header
4. **Token Introspection**: Server validates tokens with Hydra

### Authentication Headers
```
Authorization: Bearer <your-access-token>
Content-Type: application/json
```

## 📊 Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error information"
}
```

### List Response with Count
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "count": 10,
  "data": [
    // Array of items
  ]
}
```

## 🚦 Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## 🧪 Testing

### Unit Tests
```bash
cd src
go test ./...
```

### Integration Tests
```bash
cd src
go test -tags=integration ./test/...
```

## 📊 Observability & Performance Monitoring

Crypsis includes a complete observability stack for monitoring, tracing, and analyzing application performance.

### Quick Start

```bash
./start-observability.sh
```

Then open **Grafana** at http://localhost:3000 (admin/admin)

### What's Included

- **📈 Grafana**: Unified dashboards with real-time metrics visualization
- **📊 Prometheus**: Time-series metrics database with custom metric collection
- **🔍 Jaeger**: Distributed tracing for end-to-end request tracking
- **📡 OpenTelemetry**: Automatic instrumentation with traces and metrics

### Features

✅ **Pre-configured Dashboards**: Ready-to-use performance dashboards  
✅ **Real-time Monitoring**: Live metrics as your application runs  
✅ **Distributed Tracing**: Track every request through the entire stack  
✅ **Performance Analysis**: Identify bottlenecks and optimize  
✅ **Resource Monitoring**: CPU, memory, goroutines tracking  
✅ **Custom Metrics**: Encryption time, file size, operation duration  

### Access URLs

| Service | URL | Purpose |
|---------|-----|---------|
| **Grafana** | http://localhost:3000 | Main visualization dashboard |
| **Prometheus** | http://localhost:9090 | Metrics query engine |
| **Jaeger** | http://localhost:16686 | Trace visualization |

### Documentation

- **[Observability Summary](OBSERVABILITY_SUMMARY.md)**: Quick overview and getting started
- **[Complete Guide](OBSERVABILITY_COMPLETE_GUIDE.md)**: Detailed architecture and usage
- **[Quick Reference](OBSERVABILITY_QUICK_REFERENCE.md)**: Commands and troubleshooting

### Performance Testing

Run comprehensive load tests with k6:

```bash
cd performance_test
k6 run scripts/k6_load_test.js
```

Watch real-time metrics in Grafana as tests run, then analyze results to identify optimization opportunities.

## 📝 Development

### Project Structure
```
src/
├── cmd/                    # Application entrypoint
├── internal/
│   ├── config/            # Configuration management
│   ├── delivery/http/     # HTTP handlers and routes
│   ├── entity/            # Database entities
│   ├── interfaces/        # Service interfaces
│   ├── middleware/        # HTTP middleware
│   ├── model/             # Request/response models
│   ├── repository/        # Data access layer
│   ├── services/          # Business logic
│   ├── storage/           # Storage implementations
│   └── utils/             # Utility functions
└── test/                  # Test files
```

### Adding New Features

1. **Define Interface**: Add service interface in `interfaces/`
2. **Implement Service**: Create service implementation in `services/`
3. **Add Handler**: Create HTTP handler in `delivery/http/`
4. **Define Routes**: Add routes in `routes.go`
5. **Add Tests**: Create tests in `test/`

## 🛡️ Security Considerations

- **Encryption**: All files are encrypted using AES-256-GCM
- **Key Management**: Support for both file-based and KMS-based key storage
- **Access Control**: OAuth2-based authentication with token introspection
- **Audit Logging**: Complete activity tracking for compliance
- **Input Validation**: Comprehensive request validation
- **HTTPS**: SSL/TLS encryption for all communications

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation
