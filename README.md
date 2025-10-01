# 🚀 Gofsen - Framework HTTP léger pour Go

> **Un framework web moderne, simple et performant inspiré d'Express.js**

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](#-tests)
[![Progress](https://img.shields.io/badge/progress-70%25-orange.svg)](#-checklist)

---

## ✨ Pourquoi Gofsen ?

**Gofsen** combine la **simplicité d'Express.js** avec la **performance de Go**. Créé pour les développeurs qui veulent :

- 🚀 **Démarrage rapide** - API intuitive, zéro configuration
- 🔧 **Middleware flexible** - Système de middleware puissant et modulaire  
- 🛡️ **Sécurité intégrée** - Auth, CORS, Recovery, messages d'erreur détaillés
- ⚡ **Performance native** - Toute la vitesse de Go, sans compromis
- 📦 **Zéro dépendance** - Seulement la standard library Go
- 🧪 **Suite de tests complète** - Validation de toutes les fonctionnalités
- 💡 **Developer Experience** - Messages d'erreur clairs, exemples pratiques

---

## 🚀 Installation & Démarrage rapide

### Installation

```bash
go mod init mon-projet
go get github.com/username/gofsen
```

### Hello World en 30 secondes

```go
package main

import (
    "gofsen/internal/router"
    "gofsen/internal/middlewares"
    "net/http"
    "log"
)

func main() {
    // Créer le router
    r := router.NewRouter()
    
    // Ajouter un middleware global
    r.Use(middlewares.LoggerMiddleware)
    
    // Définir une route simple
    r.GET("/hello", func(ctx *types.Context) {
        ctx.JSON(http.StatusOK, map[string]string{
            "message": "Hello, Gofsen! 🎉",
        })
    })
    
    // Démarrer le serveur
    log.Println("🚀 Serveur démarré sur http://localhost:8080")
    http.ListenAndServe(":8080", r)
}
```

**Test :**
```bash
curl http://localhost:8080/hello
# {"message":"Hello, Gofsen! 🎉"}
```

---

## 📚 Guide d'utilisation

### 🧭 Routing de base

```go
r := router.NewRouter()

// Méthodes HTTP supportées
r.GET("/users", getUsersHandler)
r.POST("/users", createUserHandler)  
r.PUT("/users/:id", updateUserHandler)
r.DELETE("/users/:id", deleteUserHandler)
r.PATCH("/users/:id", patchUserHandler)
```

### 🗂️ Groupes de routes

```go
// Créer un groupe avec préfixe
api := r.Group("/api/v1")

// Ajouter des routes au groupe
api.GET("/users", getUsersHandler)
api.POST("/users", createUserHandler)

// Middleware spécifique au groupe
api.Use(middlewares.AuthMiddleware(tokenValidator))
api.GET("/profile", getProfileHandler) // Protégé par auth
```

### 📥📤 Gestion des données

#### Réponses JSON
```go
r.GET("/data", func(ctx *types.Context) {
    data := map[string]interface{}{
        "users": []string{"Alice", "Bob"},
        "count": 2,
        "timestamp": time.Now(),
    }
    ctx.JSON(http.StatusOK, data)
})
```

#### Lecture des paramètres de requête
```go
r.GET("/search", func(ctx *types.Context) {
    query := ctx.QueryParam("q")
    limit := ctx.QueryParam("limit")
    
    ctx.JSON(http.StatusOK, map[string]string{
        "query": query,
        "limit": limit,
    })
})
// GET /search?q=golang&limit=10
```

#### Parsing du JSON body
```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

r.POST("/users", func(ctx *types.Context) {
    var user User
    if err := ctx.BindJSON(&user); err != nil {
        ctx.Error(http.StatusBadRequest, "Invalid JSON")
        return
    }
    
    // Traitement...
    ctx.JSON(http.StatusCreated, user)
})
```

#### Gestion des erreurs
```go
r.GET("/error-example", func(ctx *types.Context) {
    if someCondition {
        ctx.Error(http.StatusNotFound, "Resource not found")
        return
    }
    
    ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
})
```

---

## 🔧 Middleware

### Middleware globaux

```go
r := router.NewRouter()

// Appliqué à toutes les routes
r.Use(middlewares.LoggerMiddleware)
r.Use(middlewares.RecoveryMiddleware)
r.Use(middlewares.CorsMiddleware)
```

### Middleware locaux (par groupe)

```go
adminGroup := r.Group("/admin")

// Middleware spécifique aux routes admin
adminGroup.Use(middlewares.AuthMiddleware(tokenValidator))
adminGroup.Use(func(next types.HandlerFunc) types.HandlerFunc {
    return func(ctx *types.Context) {
        // Logique custom
        next(ctx)
    }
})
```

### 🔐 Middleware d'authentification

```go
// Créer un validateur de token
tokenValidator := utils.NewTokenValidator()

// Appliquer sur les routes protégées
protected := r.Group("/api")
protected.Use(middlewares.AuthMiddleware(tokenValidator))

protected.GET("/profile", func(ctx *types.Context) {
    ctx.JSON(http.StatusOK, map[string]string{
        "message": "Accès autorisé !",
    })
})
```

**Test avec authentification :**
```bash
# Sans token (échec)
curl http://localhost:8080/api/profile
# {"error":"Missing Authorization Header"}

# Avec token (succès)
curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/profile
# {"message":"Accès autorisé !"}
```

### 🌐 Middleware CORS

```go
r.Use(middlewares.CorsMiddleware)

// Configure automatiquement :
// - Access-Control-Allow-Origin
// - Access-Control-Allow-Methods  
// - Access-Control-Allow-Headers
// - Gestion des requêtes OPTIONS
```

### 🛡️ Middleware Recovery

```go
r.Use(middlewares.RecoveryMiddleware)

r.GET("/panic", func(ctx *types.Context) {
    panic("Oops!") // Le serveur ne crash pas
})
// Retourne : {"error":"Internal Server Error"}
```

---

## 🧪 Exemples concrets

### API REST complète

```go
package main

import (
    "gofsen/internal/router"
    "gofsen/internal/middlewares"
    "gofsen/internal/types"
    "gofsen/internal/utils"
    "net/http"
    "log"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {1, "Alice", "alice@example.com"},
    {2, "Bob", "bob@example.com"},
}

func main() {
    r := router.NewRouter()
    
    // Middleware globaux
    r.Use(middlewares.LoggerMiddleware)
    r.Use(middlewares.RecoveryMiddleware)
    r.Use(middlewares.CorsMiddleware)
    
    // Routes publiques
    r.GET("/health", func(ctx *types.Context) {
        ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
    })
    
    // API v1
    api := r.Group("/api/v1")
    
    // Routes publiques API
    api.GET("/users", func(ctx *types.Context) {
        ctx.JSON(http.StatusOK, map[string]interface{}{
            "users": users,
            "count": len(users),
        })
    })
    
    // Routes protégées
    protected := api.Group("")
    tokenValidator := utils.NewTokenValidator()
    protected.Use(middlewares.AuthMiddleware(tokenValidator))
    
    protected.POST("/users", func(ctx *types.Context) {
        var newUser User
        if err := ctx.BindJSON(&newUser); err != nil {
            ctx.Error(http.StatusBadRequest, "Invalid JSON")
            return
        }
        
        newUser.ID = len(users) + 1
        users = append(users, newUser)
        
        ctx.JSON(http.StatusCreated, newUser)
    })
    
    protected.GET("/profile", func(ctx *types.Context) {
        ctx.JSON(http.StatusOK, map[string]string{
            "message": "Profile privé",
            "user":    "authenticated",
        })
    })
    
    log.Println("🚀 API REST démarrée sur http://localhost:8080")
    log.Println("📋 Essayez : curl http://localhost:8080/api/v1/users")
    http.ListenAndServe(":8080", r)
}
```

### Microservice avec validation

```go
func main() {
    r := router.NewRouter()
    r.Use(middlewares.LoggerMiddleware)
    
    // Middleware de validation custom
    validateJSON := func(next types.HandlerFunc) types.HandlerFunc {
        return func(ctx *types.Context) {
            if ctx.Request.Header.Get("Content-Type") != "application/json" {
                ctx.Error(http.StatusBadRequest, "Content-Type must be application/json")
                return
            }
            next(ctx)
        }
    }
    
    api := r.Group("/api")
    api.Use(validateJSON)
    
    api.POST("/webhook", func(ctx *types.Context) {
        var payload map[string]interface{}
        if err := ctx.BindJSON(&payload); err != nil {
            ctx.Error(http.StatusBadRequest, "Invalid JSON payload")
            return
        }
        
        // Traitement webhook...
        log.Printf("Webhook reçu: %+v", payload)
        
        ctx.JSON(http.StatusOK, map[string]string{
            "status": "processed",
            "id":     "webhook-123",
        })
    })
    
    http.ListenAndServe(":8080", r)
}
```

---

## 🧪 Tests & Validation

### Suite de tests complète

Gofsen inclut une suite de tests complète pour valider toutes les fonctionnalités :

```bash
# Démarrer le serveur
go run cmd/main.go

# Lancer les tests automatisés
./test-suite.sh

# Ou tester manuellement
curl http://localhost:8081/test/all
```

### 🚨 Messages d'erreur intelligents

Gofsen fournit des **messages d'erreur détaillés** pour améliorer l'expérience de développement :

```bash
# Test des erreurs détaillées
curl http://localhost:8081/demo/errors

# 404 avec suggestions de routes similaires
curl http://localhost:8081/route-inexistante

# 405 avec méthodes autorisées
curl -X POST http://localhost:8081/health

# 401 avec aide pour l'authentification
curl http://localhost:8081/auth/profile

# 500 avec stack trace en mode debug
curl http://localhost:8081/demo/errors/panic
```

**Fonctionnalités des erreurs :**
- 🔍 **Suggestions intelligentes** - Routes similaires pour les 404
- 📋 **Méthodes autorisées** - Liste des méthodes HTTP valides pour les 405
- 🔐 **Aide d'authentification** - Exemples de headers requis pour les 401
- 💥 **Recovery avancé** - Stack traces détaillées pour le debugging
- 📊 **Logging contextuel** - Informations complètes dans les logs

### Tests manuels

```bash
# Test de base
curl http://localhost:8081/test/all

# Test d'authentification  
curl -H "Authorization: Bearer valid-token" http://localhost:8081/test/auth/protected

# Test CORS
curl -H "Origin: https://example.com" http://localhost:8081/test/cors/check

# Test JSON avec body
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@test.com"}' \
  http://localhost:8081/test/bind

# Test des erreurs détaillées
curl http://localhost:8081/demo/errors/validate?name=
curl http://localhost:8081/demo/errors/codes?code=400
```

---

## 📊 Performance

Gofsen est conçu pour la performance :

- ⚡ **Zéro allocation** pour le routing de base
- 🎯 **Lookup O(1)** pour les routes statiques  
- 🏗️ **Architecture middleware** optimisée
- 📦 **Empreinte mémoire minimale**

### Benchmarks

```bash
# À venir : suite de benchmarks
go test -bench=. ./...
```

---

## 🤝 Comparaison avec d'autres frameworks

| Feature              | Gofsen | Gin  | Echo | Gorilla |
|---------------------|--------|------|------|---------|
| Routing             | ✅     | ✅   | ✅   | ✅      |
| Middleware          | ✅     | ✅   | ✅   | ✅      |
| Zero dependencies   | ✅     | ❌   | ❌   | ❌      |
| Built-in CORS       | ✅     | ❌   | ✅   | ❌      |
| Built-in Auth       | ✅     | ❌   | ❌   | ❌      |
| Learning curve      | 🟢 Easy | 🟢 Easy | 🟡 Medium | 🔴 Hard |

---

## 🗺️ Roadmap

### ✅ Fonctionnalités actuelles (v1.0)

- ✅ **Routing HTTP complet** (GET, POST, PUT, DELETE, PATCH)
- ✅ **Groupes de routes** avec middleware locaux
- ✅ **Middleware système** : Logger, Auth, Recovery, CORS
- ✅ **Helpers I/O** : JSON, Query params, Error handling
- ✅ **Messages d'erreur intelligents** avec suggestions et contexte
- ✅ **Suite de tests complète** avec 20+ endpoints de validation
- ✅ **Documentation & exemples** pratiques d'utilisation
- ✅ **Developer Experience** optimisée avec debugging avancé

### 🔄 Prochaines versions

#### v1.1 - Routing avancé
- 🚧 Routes dynamiques (`/users/:id`)
- 🚧 Wildcard routes (`/files/*path`)
- 🚧 Route parameter injection

#### v1.2 - Performance & Qualité  
- 🔲 Rate limiting middleware
- 🔲 Suite de benchmarks
- 🔲 Tests unitaires étendus
- 🔲 Optimisations performance

#### v1.3 - Ecosystem
- 🔲 Documentation Swagger/OpenAPI
- 🔲 Intégrations ORM (GORM)
- 🔲 Plugin system
- 🔲 WebSocket support

---

## 🆘 Support & Communauté

### 📚 Documentation
- [Guide API complet](./API.md)
- [Exemples d'utilisation](./examples/)
- [Tests de fonctionnalités](./TEST_ROUTES.md)

### 🐛 Bugs & Feature Requests
- [Issues GitHub](https://github.com/username/gofsen/issues)
- [Discussions](https://github.com/username/gofsen/discussions)

### 🤝 Contribution
Contributions bienvenues ! Voir [CONTRIBUTING.md](./CONTRIBUTING.md)

---

## 📄 License

MIT License - voir [LICENSE](./LICENSE) pour les détails.

---

## 🙏 Remerciements

Inspiré par les excellents frameworks :
- [Express.js](https://expressjs.com/) pour l'API intuitive
- [Gin](https://gin-gonic.com/) pour l'approche middleware  
- [Echo](https://echo.labstack.com/) pour les performances

---

**Fait avec ❤️ en Go**

> *"Simple, rapide, efficace - comme Go devrait l'être"*