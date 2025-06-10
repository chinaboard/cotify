# Cotify

Cotify is a modern Go-based web service that provides a robust and scalable backend infrastructure. Built with performance and reliability in mind, it leverages the power of Gin web framework and GORM for database operations.

## Features

- 🚀 High-performance HTTP server powered by Gin
- 💾 MySQL database integration with GORM
- 🔒 Secure and reliable data handling
- 📝 Structured logging with Logrus
- 🐳 Docker support for easy deployment
- 🧪 Comprehensive testing setup

## Prerequisites

- Go 1.24.0 or higher
- MySQL 5.7 or higher
- Docker and Docker Compose (for containerized deployment)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/chinaboard/cotify.git
cd cotify
```

2. Install dependencies:
```bash
go mod download
```

## Configuration

The application requires the following environment variables:

- `DB_NAME`: Database name (default: cotify)
- `DB_USER`: Database user (default: cotify)
- `DB_PASSWORD`: Database password
- `DB_PORT`: Database port (default: 3306)
- `SERVER_PORT`: Server port (default: 80)

## Running the Application

### Local Development

```bash
go run cmd/main.go
```

### Using Docker

The project includes a complete Docker setup with MariaDB and Adminer:

```bash
docker-compose up -d
```

This will start:
- MariaDB database on port 3306
- Adminer (database management tool) on port 8080
- Cotify application on port 80

## Project Structure

```
.
├── app/          # Application-specific code
├── cmd/          # Main applications
├── internal/     # Private application code
├── pkg/          # Public library code
├── sdk/          # Software development kit
└── examples/     # Example applications
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Logrus](https://github.com/sirupsen/logrus)