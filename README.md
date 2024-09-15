# Mamlaka payment

Mamlaka payments
---

## Getting Started

Follow these instructions to set up a local version of the project for development and testing purposes. See the deployment section for notes on how to deploy the project to a live system.

---

## Makefile Commands

You can run the following commands using the `Makefile`:

### Run all make commands with clean tests
```bash
make all build
```

### Build the application
```bash
make build
```

### Run the application
```bash
make run
```

### Create DB container
```bash
make docker-run
```

### Shutdown DB container
```bash
make docker-down
```

### Live reload the application
```bash
make watch
```

### Run the test suite
```bash
make test
```

### Clean up binary from the last build
```bash
make clean
```

---

## Backend Services

### API Gateway
- Acts as the single entry point for all client requests, routing them to the appropriate services.
- Manages authentication, rate limiting, request validation and payments.

### Microservices
1. **User Service**: Manages user registration, authentication, profile management, and permissions.

---

## Database

- **Relational Databases**: PostgreSQL for structured data like user profiles, ads, and transactions.

---

## Integration Services

- **Payment API**: inbuilt payment processing system.

---

## Notification Service

- **Push Notification Service**: coming soon
- **Email Service**: coming soon

---

## Storage Solutions

- **Local Storage**: Local Storage for web apps (for caching and offline access to media and posts).

---

## Security and Compliance

- **Authentication and Authorization**: OAuth 2.0 and JWT for secure user authentication and role-based access control.
- **Data Encryption**: TLS (in transit) and AES-256 (at rest) to protect sensitive data.
- **Compliance**: Ensure adherence to GDPR, CCPA, and other relevant data protection regulations.

---

### Backend Technologies Overview

#### API Gateway
- **Languages**: Go
- **Frameworks**: Echo (Go)
- **Libraries**: JWT (authentication), Rate-limiter-flexible (rate limiting)

#### Microservices
- **Languages**: Go
- **Frameworks**: Echo (Go)
- **Libraries**:
    - GORM (ORM for Go)
    - Bcrypt (password hashing)
    - JWT (authentication)

#### Databases
- **Relational Databases**: PostgreSQL
- **ORM Tools**: GORM (Go),

#### API Collection
- ** Located in docs folder in the root of the project