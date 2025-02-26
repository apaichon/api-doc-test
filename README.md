# API Documentation and Testing Project

This project contains a Go API server implementation along with Slidev presentations for API documentation and testing.

## Project Structure

```
api-doc-test/
├── slide/          # Slidev presentation files
└── src/
    └── api/        # Go API server implementation
```

## Prerequisites

- Node.js (v16 or higher)
- Go (v1.22 or higher)
- SQLite3

## Setting Up the Presentation (Slidev)

1. Navigate to the slide directory:
```bash
cd slide
```

2. Install dependencies:
```bash
npm install
```

3. Start the presentation in development mode:
```bash
npm run dev
```

4. View the presentation at http://localhost:3030

## Setting Up the API Server

1. Navigate to the API directory:
```bash
cd src/api
```

2. Install Go dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
# Create SQLite database and tables
sqlite3 payment.db < data/tables.sql
```

4. Configure the environment:
- Copy `.env.example` to `.env`
- Update the configuration values as needed

5. Start the API server:
```bash
# Run a single instance
go run cmd/server/main.go

# Run multiple instances with load balancer
go run cmd/loadbalancer/main.go -n 3
```

The API will be available at:
- Single instance: http://localhost:4000
- Load balanced instances: http://localhost:3999

## API Documentation

The API documentation is available at:
- Swagger UI: http://localhost:4000/swagger/index.html
- GraphQL Playground: http://localhost:4000/graphql

## Development

### Running Tests
```bash
go test ./...
```

### Generating API Documentation
```bash
swag init -g cmd/server/main.go
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
