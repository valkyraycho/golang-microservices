# Microservices E-Commerce Platform

A modern e-commerce platform built with Go microservices architecture, showcasing best practices in distributed systems design and cloud-native development.

## Technology Stack

### Backend Services

-   **Go** (v1.23.4) - Modern, concurrent programming language for scalable services
-   **gRPC** - High-performance RPC framework for service-to-service communication
-   **Protocol Buffers** - Efficient, language-agnostic data serialization
-   **GraphQL** - Flexible API gateway using gqlgen
-   **PostgreSQL** - Reliable relational database for persistent storage

### API Layer

-   **GraphQL Gateway** - Single entry point consolidating multiple microservices
-   **gRPC Services**:
    -   Account Service - User account management
    -   Catalog Service - Product catalog operations
    -   Order Service - Order processing and management

### Infrastructure

-   **Docker** - Containerization of microservices
-   **Docker Compose** - Local development and service orchestration

## Architecture Overview

The application follows a microservices architecture with:

-   **Service Isolation** - Each service (Account, Catalog, Order) runs independently
-   **API Gateway Pattern** - GraphQL layer aggregating underlying gRPC services
-   **Domain-Driven Design** - Services organized around business capabilities
-   **Protocol Buffer Contracts** - Strong typing and versioning of service interfaces

## Key Features

### Modular Design

-   Independent scaling of services
-   Language-agnostic service interfaces
-   Clear separation of concerns

### Modern API Design

-   GraphQL for flexible data querying
-   gRPC for efficient inter-service communication
-   Protocol Buffers for type-safe contracts

### Cloud-Native Architecture

-   Containerized services
-   Environment-based configuration
-   Ready for container orchestration

### Data Persistence

-   PostgreSQL for reliable data storage
-   Repository pattern implementation
-   Efficient database connection management
