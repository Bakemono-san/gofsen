# 🚀 Gofsen - Go HTTP Framework

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Gofsen** est un framework HTTP léger et expressif pour Go, inspiré d'Express.js. Il offre une API simple et intuitive pour créer des serveurs web performants.

## 📦 Installation

```bash
go get github.com/bakemono/gofsen
```

## 🚀 Démarrage Rapide

```go
package main

import "github.com/bakemono/gofsen"

func main() {
    app := gofsen.New()
    
    app.GET("/", func(c *gofsen.Context) {
        c.JSON(map[string]string{
            "message": "Hello Gofsen!",
        })
    })
    
    app.Listen("3000")
}
```

## 📚 Fonctionnalités

### ✅ Routing HTTP
- **Méthodes HTTP** : GET, POST, PUT, DELETE, PATCH
- **Routes paramétrées** : `/users/:id`
- **Groupes de routes** : `/api/v1`
- **Query parameters** : `?name=value`

### ✅ Middleware System
- **Logger** : Journalisation automatique des requêtes
- **Recovery** : Récupération des paniques
- **CORS** : Support complet avec configuration
- **Middleware personnalisés** : Créez vos propres middlewares

### ✅ Request/Response Helpers
- **JSON** : Parsing et envoi automatique
- **Query Params** : Accès facile aux paramètres
- **Route Params** : Support des paramètres dynamiques
- **Error Handling** : Gestion d'erreurs intégrée

## 💡 Exemples d'utilisation

### Basic Server
```go
package main

import "github.com/bakemono/gofsen"

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

import "github.com/bakemono/gofsen"

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

### Middleware Personnalisé
```go
package main

import (
    "log"
    "github.com/bakemono/gofsen"
)

func main() {
    app := gofsen.New()
    
    // Middleware d'authentification
    authMiddleware := func(c *gofsen.Context) {
        token := c.Request.Header.Get("Authorization")
        if token == "" {
            c.Error(401, "Missing token")
            return
        }
        log.Printf("User authenticated with token: %s", token)
        c.Next() // Important: continuer vers le handler suivant
    }
    
    // Routes protégées
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
app := gofsen.New()                    // Créer une instance
app.Use(middleware)                    // Ajouter un middleware global
app.GET(path, handler)                 // Route GET
app.POST(path, handler)                // Route POST  
app.PUT(path, handler)                 // Route PUT
app.DELETE(path, handler)              // Route DELETE
app.PATCH(path, handler)               // Route PATCH
app.Group(prefix)                      // Créer un groupe de routes
app.Listen(port)                       // Démarrer le serveur
app.PrintRoutes()                      // Afficher les routes
```

### Context Methods
```go
// Request
c.Param("id")                          // Paramètre de route
c.QueryParam("name")                   // Paramètre de query
c.BindJSON(&struct{})                  // Parser le JSON

// Response
c.JSON(data)                           // Réponse JSON
c.Text("Hello")                        // Réponse texte
c.HTML("<h1>Hello</h1>")              // Réponse HTML
c.Status(200)                          // Code de statut
c.Error(404, "Not found")             // Erreur avec code

// Middleware
c.Next()                               // Middleware suivant
```

### Middlewares Built-in
```go
gofsen.Logger()                        // Logger des requêtes
gofsen.Recovery()                      // Récupération des paniques
gofsen.CORS()                          // CORS avec config par défaut
gofsen.CORSWithConfig(config)          // CORS avec config personnalisée
```

## 🔧 Configuration CORS
```go
corsConfig := gofsen.CORSConfig{
    AllowOrigins: []string{"http://localhost:3000", "https://myapp.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
}

app.Use(gofsen.CORSWithConfig(corsConfig))
```

## 📊 Comparaison avec d'autres frameworks

| Framework | Import | Philosophie |
|-----------|--------|-------------|
| **Gofsen** | `github.com/bakemono/gofsen` | Simple, Express.js-like |
| Gin | `github.com/gin-gonic/gin` | Performance, minimaliste |
| Fiber | `github.com/gofiber/fiber/v2` | Express.js pour Go |
| Echo | `github.com/labstack/echo/v4` | Haute performance |

### Exemple de migration depuis Gin
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

## 🧪 Tests
```bash
go test ./...
```

## 📈 Performance
Gofsen est optimisé pour offrir des performances excellentes avec une API simple :
- Routing rapide avec regex optimisées
- Middleware chain efficace
- Zero allocation dans les cas courants

## 🤝 Contribution
Les contributions sont les bienvenues ! 

1. Fork le projet
2. Créez votre branche (`git checkout -b feature/amazing-feature`)
3. Commit vos changements (`git commit -m 'Add amazing feature'`)
4. Push vers la branche (`git push origin feature/amazing-feature`)
5. Ouvrez une Pull Request

## 📄 License
Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

## 🌟 Pourquoi Gofsen ?
- **Simple** : API intuitive inspirée d'Express.js
- **Léger** : Aucune dépendance externe
- **Rapide** : Performance optimisée
- **Flexible** : Système de middleware extensible
- **Prêt pour la production** : Gestion d'erreurs robuste

---

Made with ❤️ by [Bakemono](https://github.com/bakemono)