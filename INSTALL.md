# 🚀 Installation & Utilisation de Gofsen

## 📦 Installation

### Option 1: Cloner le projet complet
```bash
git clone <repository-url> gofsen
cd gofsen

# Tester que tout fonctionne
go run cmd/main.go
```

### Option 2: Utiliser Gofsen dans votre projet
```bash
# Créer votre nouveau projet
mkdir mon-projet
cd mon-projet
go mod init mon-projet

# Copier les fichiers Gofsen nécessaires
# Vous pouvez copier le dossier 'internal' ou adapter les imports
```

## 🔧 Structure pour utilisation externe

Pour utiliser Gofsen dans votre propre projet, voici la structure recommandée :

```
votre-projet/
├── go.mod
├── main.go
├── gofsen/          # Copiez le dossier internal ici
│   ├── router/
│   ├── middlewares/
│   ├── types/
│   └── utils/
└── examples/        # Optionnel
```

## 🚀 Quick Start pour un nouveau projet

### 1. Créer votre main.go
```go
package main

import (
    "votre-projet/gofsen/router"
    "votre-projet/gofsen/middlewares"
    "votre-projet/gofsen/types"
    "log"
    "net/http"
)

func main() {
    // Créer le router
    r := router.NewRouter()
    
    // Ajouter middlewares
    r.Use(middlewares.LoggerMiddleware)
    r.Use(middlewares.RecoveryMiddleware)
    
    // Vos routes
    r.GET("/", func(ctx *types.Context) {
        ctx.JSON(http.StatusOK, map[string]string{
            "message": "Hello Gofsen!",
        })
    })
    
    // Démarrer le serveur
    log.Println("🚀 Serveur démarré sur http://localhost:8080")
    http.ListenAndServe(":8080", r)
}
```

### 2. Adapter les imports
Si vous copiez les fichiers dans votre projet, modifiez les imports dans les fichiers Go :

**Avant (dans le projet original) :**
```go
import "gofsen/internal/types"
```

**Après (dans votre projet) :**
```go
import "votre-projet/gofsen/types"
```

## 🎯 Exemples prêts à l'emploi

### Hello World minimal
```go
package main

import (
    "votre-projet/gofsen/router"
    "votre-projet/gofsen/types"
    "net/http"
    "log"
)

func main() {
    r := router.NewRouter()
    
    r.GET("/hello", func(ctx *types.Context) {
        name := ctx.QueryParam("name")
        if name == "" {
            name = "World"
        }
        
        ctx.JSON(http.StatusOK, map[string]string{
            "greeting": "Hello " + name + "!",
        })
    })
    
    log.Println("🚀 http://localhost:8080/hello?name=Alice")
    http.ListenAndServe(":8080", r)
}
```

## 🔄 Migration depuis d'autres frameworks

### Depuis Gin
```go
// Gin
r.GET("/users/:id", func(c *gin.Context) {
    c.JSON(200, gin.H{"id": c.Param("id")})
})

// Gofsen (similaire)
r.GET("/users/:id", func(ctx *types.Context) {
    // Note: Routes dynamiques en développement
    ctx.JSON(http.StatusOK, map[string]string{"id": "123"})
})
```

### Depuis Echo
```go
// Echo
e.GET("/users", func(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "users"})
})

// Gofsen (très similaire)
r.GET("/users", func(ctx *types.Context) {
    ctx.JSON(http.StatusOK, map[string]string{"message": "users"})
})
```

## 📋 Checklist pour un nouveau projet

- [ ] **Copier les fichiers** Gofsen dans votre projet
- [ ] **Adapter les imports** selon votre structure
- [ ] **Tester** que le serveur démarre : `go run main.go`
- [ ] **Ajouter vos routes** et middlewares
- [ ] **Configurer** l'authentification si nécessaire
- [ ] **Tester** avec les endpoints de démonstration

## 🚨 Points d'attention

### Imports
⚠️ **Les imports doivent être adaptés** selon votre structure de projet.

### Go Version
✅ **Go 1.21+** recommandé (testé avec Go 1.23.5)

### Dépendances
✅ **Zéro dépendance externe** - utilise seulement la standard library Go

## 🆘 Problèmes courants

### "Package not found"
```bash
# Solution: Vérifiez vos imports et la structure des dossiers
go mod tidy
```

### "Cannot find module"
```bash
# Solution: Initialisez votre module Go
go mod init votre-projet
```

### "Permission denied" pour test-suite.sh
```bash
# Solution: Rendez le script exécutable
chmod +x test-suite.sh
```

## 🎉 Vous êtes prêt !

Une fois l'installation terminée, vous pouvez :

1. **Démarrer** avec les exemples fournis
2. **Consulter** la documentation complète (README.md)
3. **Tester** les fonctionnalités (/test/all)
4. **Développer** votre API

---

**Support :** Consultez le README.md principal pour la documentation complète.