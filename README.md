[![Go Tests](https://github.com/afteracademy/goserve/actions/workflows/go-test.yml/badge.svg)](https://github.com/afteracademy/goserve/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/afteracademy/goserve/v2)](https://goreportcard.com/report/github.com/afteracademy/goserve/v2)
[![GoDoc](https://pkg.go.dev/badge/github.com/afteracademy/goserve/v2)](https://pkg.go.dev/github.com/afteracademy/goserve/v2)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

<div align="center">

# GoServe

### Production-Ready Go Backend Architecture Framework

![Banner](.docs/goserve-banner.png)

**A comprehensive, production-ready framework for building scalable Go backend services with PostgreSQL, MongoDB, Redis, and NATS microservices support.**

[![Documentation](https://img.shields.io/badge/📚_Read_Documentation-goserve.afteracademy.com-blue?style=for-the-badge)](http://goserve.afteracademy.com)

</div>

## Getting Started
- [Quick Start Guide](https://goserve.afteracademy.com/getting-started)
- [Framework Architecture](https://goserve.afteracademy.com/architecture)
- [Core Concepts](https://goserve.afteracademy.com/core-concepts)

## Features

- **Clean Architecture** - Well-structured, maintainable codebase following Go best practices
- **HTTP Server** - Built on Gin framework with middleware support
- **Authentication** - JWT-based authentication and authorization
- **Multiple Databases** - Support for PostgreSQL, MongoDB, and Redis
- **Microservices** - NATS-based microservice communication patterns
- **Validation** - Comprehensive request/response validation using validator v10
- **Error Handling** - Structured error handling and API responses
- **Testing** - Extensive test coverage with mocking support
- **DTOs** - Type-safe data transfer objects for common types (UUID, ObjectID, Slug, Pagination)
- **Performance** - Optimized for high-throughput production workloads

## Core Packages

| Package | Description |
|---------|-------------|
| **network** | HTTP routing, middleware, request/response handling, validation |
| **mongo** | MongoDB connection, query builder, validation utilities |
| **postgres** | PostgreSQL database connectivity and operations |
| **redis** | Redis caching and key-value store operations |
| **micro** | NATS microservice framework for message-based communication |
| **dto** | Common DTOs (MongoID, UUID, Slug, Pagination) |
| **utility** | Helper functions for formatting, mapping, random generation |
| **middleware** | HTTP middleware (error catcher, 404 handler) |

## Example Projects

Real-world applications built with GoServe:

1. **[PostgreSQL API Server](https://github.com/afteracademy/goserve-example-api-server-postgres)**  
   Complete REST API with PostgreSQL, JWT authentication, and clean architecture

2. **[MongoDB API Server](https://github.com/afteracademy/goserve-example-api-server-mongo)**  
   MongoDB-based backend with flexible schema design

3. **[Microservices Example](https://github.com/afteracademy/gomicro)**  
   NATS-based microservices communication patterns

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Authentication**: JWT tokens
- **Databases**: 
  - PostgreSQL ([pgx](https://github.com/jackc/pgx))
  - MongoDB ([mongo-driver](https://github.com/mongodb/mongo-go-driver))
  - Redis ([go-redis](https://github.com/redis/go-redis))
- **Validation**: [validator](https://github.com/go-playground/validator)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Messaging**: [NATS](https://github.com/nats-io/nats.go)
- **Testing**: [Testify](https://github.com/stretchr/testify)

## Documentation

<div align="center">

### [**Read the Full Documentation on pkg.go.dev**](https://goserve.afteracademy.com)

Comprehensive guides, API references, and examples for all packages

</div>

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md).

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/goserve.git`
3. Create a feature branch: `git checkout -b feature/your-feature`
4. Make changes and add tests
5. Run tests: `go test ./...`
6. Commit with clear messages: `git commit -m "Add feature: description"`
7. Push to your fork: `git push origin feature/your-feature`
8. Open a Pull Request

## Learn More

Subscribe to **AfterAcademy** on YouTube for in-depth tutorials and concept explanations:

[![YouTube](https://img.shields.io/badge/YouTube-Subscribe-red?style=for-the-badge&logo=youtube&logoColor=white)](https://www.youtube.com/@afteracad)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support This Project

If you find GoServe useful, please consider:

- **Starring** this repository
- **Reporting** bugs and issues
- **Suggesting** new features
- **Contributing** code improvements
- **Sharing** with the community

## Security

For security concerns, please review our [Security Policy](SECURITY.md).

---

<div align="center">

**Built with love by [AfterAcademy](https://github.com/afteracademy)**

[Documentation](https://goserve.afteracademy.com) • [Quick Start Guide](https://goserve.afteracademy.com/getting-started) • [Contributing](CONTRIBUTING.md) • [License](LICENSE)

---
</div>
