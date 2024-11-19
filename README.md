# Wind Surf Go

A modern Go web application featuring user authentication and RESTful API design.

## Features

- User Authentication System
  - Registration with username and password
  - Login with JWT token generation
  - Secure password hashing using bcrypt
- RESTful API Design
  - Versioned API endpoints (v1)
  - Standardized response format
  - Clean project structure
- Database Integration
  - MySQL database
  - GORM ORM integration
  - Automatic schema migration

## Project Structure

```
wind-surf-go/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── config/           # Configuration files
│   ├── handler/          # HTTP request handlers
│   ├── model/           # Database models
│   ├── router/          # Route definitions
│   └── utils/           # Utility functions
└── README.md
```

## API Endpoints

### User Management

#### Register User
```
POST /v1/api/users/register
Content-Type: application/json

Request:
{
    "username": "string",    // min: 3, max: 50 characters
    "password": "string"     // min: 6 characters
}

Response:
{
    "code": 0,              // 0: success, 1: error
    "message": "string",
    "data": {
        "token": "string",
        "username": "string",
        "user_id": number
    }
}
```

#### User Login
```
POST /v1/api/users/login
Content-Type: application/json

Request:
{
    "username": "string",
    "password": "string"
}

Response:
{
    "code": 0,              // 0: success, 1: error
    "message": "string",
    "data": {
        "token": "string",
        "username": "string",
        "user_id": number
    }
}
```

## Response Format

All API responses follow a standardized format:

### Success Response
```json
{
    "code": 0,
    "message": "success",
    "data": {
        // Response data
    }
}
```

### Error Response
```json
{
    "code": 1,
    "message": "error description"
}
```

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/yourusername/wind-surf-go.git
cd wind-surf-go
```

2. Set up MySQL database
```bash
mysql -u root -p
```

```sql
CREATE DATABASE wind_surf CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

3. Configure the application
Update the MySQL configuration in `internal/config/config.yaml` with your database credentials.

4. Install dependencies
```bash
go mod tidy
```

5. Run the application
```bash
go run cmd/main.go
```

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [GORM](https://gorm.io) - ORM library with MySQL driver
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Viper](https://github.com/spf13/viper) - Configuration management

## Configuration

The application uses a YAML configuration file located at `internal/config/config.yaml`. Configure the following settings:

```yaml
server:
  port: "8080"  # Server port

mysql:
  host: "localhost"      # MySQL host
  port: 3306            # MySQL port
  username: "root"      # MySQL username
  password: "123456"    # MySQL password
  database: "wind_surf" # Database name
  charset: "utf8mb4"    # Character set
  parseTime: true       # Parse time values
  loc: "Local"         # Time zone
```

## Security Considerations

- JWT tokens expire after 24 hours
- Passwords are hashed using bcrypt
- Username uniqueness is enforced
- Input validation for username and password

## License

This project is licensed under the MIT License - see the LICENSE file for details.
