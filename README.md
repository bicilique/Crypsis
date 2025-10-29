# Crypsis - Enterprise File Encryption & Storage Service

A comprehensive, enterprise-grade file encryption and secure storage solution built with Go, featuring OAuth2 authentication, KMS integration, and distributed storage capabilities.

## 🚀 Features

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
