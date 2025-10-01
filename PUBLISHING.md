# 🚀 Guide de Publication - Gofsen Package

## 📦 Étapes pour publier Gofsen comme package Go

### 1. Préparer le Repository GitHub

```bash
# 1. Créer un nouveau repository sur GitHub : github.com/bakemono/gofsen
# 2. Cloner et pousser le code

git init
git add .
git commit -m "Initial release: Gofsen v1.0.0"
git branch -M main
git remote add origin https://github.com/bakemono/gofsen.git
git push -u origin main
```

### 2. Créer un Tag de Version

```bash
# Créer un tag pour la version
git tag v1.0.0
git push origin v1.0.0
```

### 3. Publier sur Go Modules

Une fois le repository public sur GitHub, le package sera automatiquement disponible via :

```bash
go get github.com/bakemono/gofsen
```

### 4. Structure des Fichiers pour Publication

```
github.com/bakemono/gofsen/
├── gofsen.go              # Package principal
├── gofsen_test.go         # Tests
├── go.mod                 # Module definition
├── README_PACKAGE.md      # Documentation package
├── LICENSE                # License MIT
├── example/
│   └── main.go           # Exemple d'utilisation
└── examples-standalone/
    └── standalone-example.go  # Version standalone
```

### 5. Vérification avant Publication

```bash
# Tester le package
go test ./...

# Vérifier le module
go mod tidy
go mod verify

# Tester l'import externe
mkdir /tmp/test-gofsen
cd /tmp/test-gofsen
go mod init test
go get github.com/bakemono/gofsen
```

### 6. Documentation pkg.go.dev

Une fois publié, la documentation sera automatiquement générée sur :
- https://pkg.go.dev/github.com/bakemono/gofsen

### 7. Utilisation par les Développeurs

Les développeurs pourront ensuite utiliser Gofsen comme ceci :

```go
// main.go
package main

import "github.com/bakemono/gofsen"

func main() {
    app := gofsen.New()
    
    app.GET("/", func(c *gofsen.Context) {
        c.JSON(map[string]string{
            "message": "Hello from Gofsen!",
        })
    })
    
    app.Listen("8080")
}
```

```bash
# Installation
go mod init my-project
go get github.com/bakemono/gofsen
go run main.go
```

### 8. Maintenance et Mises à Jour

Pour publier de nouvelles versions :

```bash
# Apporter des modifications
git add .
git commit -m "feat: add new feature"
git push

# Créer une nouvelle version
git tag v1.1.0
git push origin v1.1.0
```

### 9. Promotion

- **GitHub** : README.md attrayant avec badges
- **Reddit** : r/golang pour l'annonce
- **Twitter** : #golang #webframework
- **Dev.to** : Article de blog sur le framework
- **Discord/Slack** : Communautés Go

### 10. Checklist Publication

- [x] Package principal (`gofsen.go`)
- [x] Tests complets (`gofsen_test.go`)
- [x] Documentation (`README_PACKAGE.md`)
- [x] Licence (`LICENSE`)
- [x] Exemple d'utilisation (`example/main.go`)
- [x] Go module configuré (`go.mod`)
- [ ] Repository GitHub public
- [ ] Tag de version `v1.0.0`
- [ ] Tests passent (`go test ./...`)

## 🌟 Commandes Finales

```bash
# Dernière vérification
go test ./...
go mod tidy

# Publier
git add .
git commit -m "Release v1.0.0: Gofsen HTTP Framework"
git tag v1.0.0
git push origin main --tags

# Le package sera disponible immédiatement via :
# go get github.com/bakemono/gofsen
```

---

**Note**: Remplacez `bakemono` par votre nom d'utilisateur GitHub réel.