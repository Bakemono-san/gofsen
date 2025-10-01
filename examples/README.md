# 📚 Exemples Gofsen Framework

Ce dossier contient des exemples pratiques d'utilisation du framework Gofsen.

## 🚀 Exemples disponibles

### 1. Hello World (`01-hello-world.go`)
**Exemple de base** - Premier contact avec Gofsen
- Routes GET/POST/PUT/DELETE simples
- Paramètres de requête
- Réponses JSON
- Gestion d'erreurs basique

**Pour tester :**
```bash
go run examples/01-hello-world.go
curl http://localhost:3000/
curl "http://localhost:3000/hello?name=Alice&lang=fr"
```

### 2. API REST (`02-rest-api.go`)
**API REST complète** avec utilisateurs
- CRUD operations
- Authentification avec middleware
- Pagination et recherche
- Validation des données
- Gestion d'erreurs avancée

**Pour tester :**
```bash
go run examples/02-rest-api.go
curl http://localhost:3001/api/v1/users
curl -H "Authorization: Bearer valid-token" http://localhost:3001/auth/profile
```

### 3. Démonstration Middlewares (`03-middleware-demo.go`)
**Showcase des middlewares** Gofsen
- Middlewares globaux et locaux
- Middleware custom (headers, timing)
- Chaînage de middlewares
- CORS, Recovery, Validation

**Pour tester :**
```bash
go run examples/03-middleware-demo.go
curl -I http://localhost:3002/ # Voir les headers
curl http://localhost:3002/panic # Test recovery
```

---

## 🎯 Comment utiliser ces exemples

### Prérequis
- Go 1.21+ installé
- Gofsen framework dans votre projet

### Étapes
1. **Cloner/télécharger** les exemples
2. **Naviguer** vers le dossier du projet
3. **Lancer** un exemple :
   ```bash
   go run examples/01-hello-world.go
   ```
4. **Tester** avec curl ou un navigateur

### Structure type d'un exemple
```go
package main

import (
    "gofsen/internal/router"
    "gofsen/internal/middlewares"
    "gofsen/internal/types"
    // ...
)

func main() {
    // 1. Créer le router
    r := router.NewRouter()
    
    // 2. Ajouter middlewares
    r.Use(middlewares.LoggerMiddleware)
    
    // 3. Définir routes
    r.GET("/", handlerFunction)
    
    // 4. Démarrer serveur
    http.ListenAndServe(":8080", r)
}
```

---

## 📖 Fonctionnalités démontrées

### ✅ Routing
- [x] Méthodes HTTP (GET, POST, PUT, DELETE, PATCH)
- [x] Groupes de routes
- [x] Paramètres de requête (`ctx.QueryParam()`)

### ✅ Middleware
- [x] Middleware globaux
- [x] Middleware locaux (par groupe)
- [x] Logger, Recovery, CORS, Auth
- [x] Middleware custom

### ✅ I/O & Data
- [x] Réponses JSON (`ctx.JSON()`)
- [x] Parsing JSON (`ctx.BindJSON()`)
- [x] Gestion d'erreurs (`ctx.Error()`)

### ✅ Sécurité
- [x] Authentification par token
- [x] Validation des données
- [x] CORS headers
- [x] Recovery de panics

---

## 🧪 Tests rapides

### Test complet automatique
```bash
# Démarrer le serveur principal
go run cmd/main.go

# Dans un autre terminal
./test-suite.sh
```

### Tests manuels par exemple

#### Exemple 1 - Hello World
```bash
go run examples/01-hello-world.go &
curl http://localhost:3000/
curl "http://localhost:3000/hello?name=World&lang=en"
curl -X POST -H "Content-Type: application/json" -d '{"test":"data"}' http://localhost:3000/echo
```

#### Exemple 2 - API REST  
```bash
go run examples/02-rest-api.go &
curl http://localhost:3001/api/v1/users
curl -H "Authorization: Bearer valid-token" http://localhost:3001/auth/profile
curl -X POST -H "Authorization: Bearer valid-token" -H "Content-Type: application/json" \
     -d '{"name":"Test","email":"test@example.com","username":"test"}' \
     http://localhost:3001/auth/users
```

#### Exemple 3 - Middlewares
```bash
go run examples/03-middleware-demo.go &
curl -I http://localhost:3002/ # Headers custom
curl http://localhost:3002/panic # Recovery middleware
curl -H "X-API-Key: demo-key-123" http://localhost:3002/secure/data
```

---

## 💡 Points clés à retenir

### 🎯 Patterns recommandés
1. **Router first** - Toujours créer le router en premier
2. **Middleware global** - Logger et Recovery en premier
3. **Groupes logiques** - Séparer API/Auth/Admin
4. **Validation early** - Valider les données dès la réception
5. **Erreurs consistantes** - Utiliser `ctx.Error()` partout

### 🚀 Bonnes pratiques
- **Structure claire** - Un handler par fonctionnalité
- **Logging approprié** - Logger les actions importantes
- **Gestion d'erreurs** - Toujours gérer les cas d'erreur
- **Headers de sécurité** - CORS, Content-Type, etc.
- **Codes de statut corrects** - 200, 201, 400, 401, 404, 500

### ⚡ Performance
- **Middleware minimal** - Seulement ce qui est nécessaire
- **Validation rapide** - Échouer vite sur les erreurs
- **JSON efficient** - Utiliser les structures Go
- **Logging async** - Ne pas bloquer les requêtes

---

## 🔗 Ressources supplémentaires

- [README principal](../README.md) - Documentation complète
- [Tests complets](../TEST_ROUTES.md) - Suite de tests
- [API Reference](../API.md) - Documentation API détaillée

---

**🎉 Amusez-vous bien avec Gofsen !**