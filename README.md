# 🚀 Gofsen - HTTP Framework for Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![GitHub release](https://img.shields.io/github/v/release/Bakemono-san/gofsen.svg)](https://github.com/Bakemono-san/gofsen/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bakemono-san/gofsen)](https://goreportcard.com/report/github.com/Bakemono-san/gofsen)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/Bakemono-san/gofsen.svg)](https://pkg.go.dev/github.com/Bakemono-san/gofsen)

**Gofsen** is a lightweight, Express.js-inspired HTTP framework for Go. Simple, fast, and powerful.

---

## 📦 Installation

```bash
go get github.com/Bakemono-san/gofsen
```

### CLI (gofsen-cli) Installation

You can install the CLI in multiple ways:

- Using go install (latest from main):

```bash
go install github.com/Bakemono-san/gofsen/cmd/gofsen-cli@latest
```

- Using Homebrew (macOS/Linux):

```bash
brew tap Bakemono-san/homebrew-tap
brew install gofsen-cli
```

- Download prebuilt binaries (Windows/macOS/Linux) from the
[Releases](https://github.com/Bakemono-san/gofsen/releases) page.

## 🚀 Quick Start

```go
package main

import "github.com/Bakemono-san/gofsen"

func main() {
    app := gofsen.New()
    
    app.GET("/", func(c *gofsen.Context) {
        c.JSON(map[string]string{
            "message": "Hello Gofsen!",
            "version": gofsen.Version,
        })
    })
    
    app.Listen("8080")
}
```

## ✨ Features

### ✅ HTTP Routing

- **HTTP Methods**: GET, POST, PUT, DELETE, PATCH
- **Route Parameters**: `/users/:id`
- **Route Groups**: `/api/v1`
- **Query Parameters**: `?name=value`

### ✅ Middleware System

- **Logger**: Automatic request logging
- **Recovery**: Panic recovery
- **CORS**: Complete CORS support with configuration
- **Custom Middleware**: Create your own middlewares

### ✅ Request/Response Helpers

- **JSON**: Automatic parsing and sending
- **Query Params**: Easy access to parameters
- **Route Params**: Dynamic parameter support
- **Error Handling**: Built-in error management

## 💡 Examples

### Basic Server

```go
package main

import "github.com/Bakemono-san/gofsen"

func main() {
    app := gofsen.New()
    
    // Middlewares
    app.Use(gofsen.Logger())
    app.Use(gofsen.Recovery())
    app.Use(gofsen.CORS())
    
    // Routes
    app.GET("/health", func(c *gofsen.Context) {
        c.JSON(map[string]string{"status": "OK"})
    })
    
    app.Listen("8080")
}
```

### REST API

```go
package main

import "github.com/Bakemono-san/gofsen"

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    app := gofsen.New()
    app.Use(gofsen.Logger())
    
    api := app.Group("/api/v1")
    
    // GET /api/v1/users
    api.GET("/users", getUsers)
    
    // GET /api/v1/users/:id
    api.GET("/users/:id", getUser)
    
    // POST /api/v1/users
    api.POST("/users", createUser)
    
    app.Listen("3000")
}

func getUsers(c *gofsen.Context) {
    users := []User{{ID: 1, Name: "Alice"}}
    c.JSON(users)
}

func getUser(c *gofsen.Context) {
    id := c.Param("id")
    c.JSON(map[string]string{"id": id})
}

func createUser(c *gofsen.Context) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        c.Error(400, "Invalid JSON")
        return
    }
    c.Status(201).JSON(user)
}
```

### Custom Middleware

```go
package main

import (
    "log"
    "github.com/Bakemono-san/gofsen"
)

func main() {
    app := gofsen.New()
    
    // Authentication middleware
    authMiddleware := func(c *gofsen.Context) {
        token := c.Request.Header.Get("Authorization")
        if token == "" {
            c.Error(401, "Missing token")
            return
        }
        log.Printf("User authenticated with token: %s", token)
        c.Next() // Important: continue to next handler
    }
    
    // Protected routes
    protected := app.Group("/protected")
    protected.Use(authMiddleware)
    
    protected.GET("/profile", func(c *gofsen.Context) {
        c.JSON(map[string]string{"message": "Protected route"})
    })
    
    app.Listen("3000")
}
```

## 🛠️ API Reference

### Router Methods

```go
app := gofsen.New()                    // Create instance
app.Use(middleware)                    // Add global middleware
app.GET(path, handler)                 // GET route
app.POST(path, handler)                // POST route  
app.PUT(path, handler)                 // PUT route
app.DELETE(path, handler)              // DELETE route
app.PATCH(path, handler)               // PATCH route
app.Group(prefix)                      // Create route group
app.Listen(port)                       // Start server
app.PrintRoutes()                      // Print routes
```

### Context Methods

```go
// Request
c.Param("id")                          // Route parameter
c.QueryParam("name")                   // Query parameter
c.BindJSON(&struct{})                  // Parse JSON

// Response
c.JSON(data)                           // JSON response
c.Text("Hello")                        // Text response
c.HTML("<h1>Hello</h1>")              // HTML response
c.Status(200)                          // Status code
c.Error(404, "Not found")             // Error with code

// Middleware
c.Next()                               // Next middleware
```

### Built-in Middlewares

```go
gofsen.Logger()                        // Request logger
gofsen.Recovery()                      // Panic recovery
gofsen.CORS()                          // CORS with defaults
gofsen.CORSFromEnv()                   // CORS from environment variables
gofsen.CORSWithConfig(config)          // CORS with custom config
```

## 🔧 CORS Configuration

### Manual Configuration

```go
corsConfig := gofsen.CORSConfig{
    AllowOrigins: []string{"http://localhost:3000", "https://myapp.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
}

app.Use(gofsen.CORSWithConfig(corsConfig))
```

### Environment Variables Configuration

Gofsen supports configuring CORS through environment variables for easier deployment and configuration management:

```go
// Use CORS configured from environment variables
app.Use(gofsen.CORSFromEnv())
```

**Supported Environment Variables:**

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `CORS_ALLOWED_ORIGINS` | Allowed origins (comma-separated) | `*` | `http://localhost:3000,https://myapp.com` |
| `ALLOWED_ORIGINS` | Fallback for allowed origins | `*` | `https://example.com` |
| `CORS_ALLOWED_METHODS` | Allowed HTTP methods | `GET,POST,PUT,DELETE,PATCH,OPTIONS` | `GET,POST,PUT` |
| `CORS_ALLOWED_HEADERS` | Allowed headers | `Content-Type,Authorization` | `Content-Type,X-API-Key` |

**Example .env file:**

```bash
# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://myapp.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With
```

**Usage in code:**

```go
package main

import "github.com/Bakemono-san/gofsen"

func main() {
    app := gofsen.New()
    
    // CORS will automatically read from environment variables
    app.Use(gofsen.CORSFromEnv())
    
    app.GET("/", func(c *gofsen.Context) {
        c.JSON(map[string]string{"message": "CORS configured from env!"})
    })
    
    app.Listen("8080")
}
```

## 📊 Framework Comparison

| Framework | Import | Philosophy |
|-----------|--------|------------|
| **Gofsen** | `github.com/Bakemono-san/gofsen` | Simple, Express.js-like |
| Gin | `github.com/gin-gonic/gin` | Performance, minimalist |
| Fiber | `github.com/gofiber/fiber/v2` | Express.js for Go |
| Echo | `github.com/labstack/echo/v4` | High performance |

### Migration from Gin

```go
// Gin
r := gin.Default()
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"id": id})
})

// Gofsen
app := gofsen.New()
app.GET("/users/:id", func(c *gofsen.Context) {
    id := c.Param("id")
    c.JSON(map[string]string{"id": id})
})
```

## 🧪 Testing

```bash
go test ./...
```

## 📈 Performance

Gofsen is optimized for excellent performance with a simple API:

- Fast routing with optimized regex
- Efficient middleware chain
- Zero allocation in common cases

## 🤝 Contributing

Contributions are welcome!

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🌟 Why Choose Gofsen?

- **Simple**: Intuitive API inspired by Express.js
- **Lightweight**: No external dependencies
- **Fast**: Optimized performance
- **Flexible**: Extensible middleware system
- **Production-ready**: Robust error handling

---

Made with ❤️ by [Bakemono](https://github.com/Bakemono-san)

## 📚 Links

- **Documentation**: [pkg.go.dev/github.com/Bakemono-san/gofsen](https://pkg.go.dev/github.com/Bakemono-san/gofsen)
- **Repository**: [github.com/Bakemono-san/gofsen](https://github.com/Bakemono-san/gofsen)
- **Issues**: [github.com/Bakemono-san/gofsen/issues](https://github.com/Bakemono-san/gofsen/issues)
- **Releases**: [github.com/Bakemono-san/gofsen/releases](https://github.com/Bakemono-san/gofsen/releases)
