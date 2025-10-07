package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// These variables are populated at build time via -ldflags (e.g., by GoReleaser)
var (
	version = "dev" // semantic version, e.g. v1.2.3
	commit  = ""    // git commit SHA
	date    = ""    // build date (RFC3339 or yyyy-mm-dd)
)

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

type ProjectConfig struct {
	Name        string
	Module      string
	Port        string
	UseCORS     bool
	UseAuth     bool
	UseDatabase bool
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "new", "create":
		createProject()
	case "generate", "gen":
		if len(os.Args) < 3 {
			fmt.Println("❌ Error: Please specify what to generate (route, middleware, handler)")
			showGenerateHelp()
			return
		}
		generateCode(os.Args[2])
	case "version", "-v", "--version":
		printVersion()
	case "help", "-h", "--help":
		showHelp()
	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		showHelp()
	}
}

func showHelp() {
	helpText := `🚀 Gofsen CLI 🇸🇳 - Framework HTTP pour Go

UTILISATION:
  gofsen-cli <command> [arguments]

COMMANDES:
  new, create              Créer un nouveau projet Gofsen
  generate, gen <type>     Générer du code (route, middleware, handler)
  version, -v, --version   Afficher la version
  help, -h, --help         Afficher cette aide

EXEMPLES:
  gofsen-cli new                    # Créer un nouveau projet interactif
  gofsen-cli gen route users        # Générer des routes CRUD pour 'users'
  gofsen-cli gen middleware auth    # Générer un middleware d'authentification
  gofsen-cli gen handler products   # Générer un handler pour 'products'

Pour plus d'informations: https://github.com/Bakemono-san/gofsen`

	fmt.Println(helpText)
}

func printVersion() {
	if commit != "" || date != "" {
		fmt.Printf("Gofsen CLI 🇸🇳 %s (commit: %s, date: %s)\n", version, commit, date)
		return
	}
	fmt.Printf("Gofsen CLI 🇸🇳 %s\n", version)
}

func showGenerateHelp() {
	helpText := `GÉNÉRATION DE CODE:
  gofsen-cli gen route <name>       # Routes CRUD complètes
  gofsen-cli gen middleware <name>  # Middleware personnalisé
  gofsen-cli gen handler <name>     # Handler/Controller

EXEMPLES:
  gofsen-cli gen route users        # Génère les routes GET, POST, PUT, DELETE pour users
  gofsen-cli gen middleware cors    # Génère un middleware CORS personnalisé
  gofsen-cli gen handler auth       # Génère un handler d'authentification`

	fmt.Println(helpText)
}

func createProject() {
	fmt.Println("🚀 Création d'un nouveau projet Gofsen")

	config := ProjectConfig{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("📝 Nom du projet: ")
	name, _ := reader.ReadString('\n')
	config.Name = strings.TrimSpace(name)
	if config.Name == "" {
		config.Name = "my-gofsen-app"
	}

	fmt.Printf("📦 Module Go (github.com/username/%s): ", config.Name)
	module, _ := reader.ReadString('\n')
	config.Module = strings.TrimSpace(module)
	if config.Module == "" {
		config.Module = fmt.Sprintf("github.com/username/%s", config.Name)
	}

	fmt.Print("🌐 Port (8080): ")
	port, _ := reader.ReadString('\n')
	config.Port = strings.TrimSpace(port)
	if config.Port == "" {
		config.Port = "8080"
	}

	config.UseCORS = askYesNo("🛡️ Inclure middleware CORS? (y/N): ")
	config.UseAuth = askYesNo("🔐 Inclure middleware Auth? (y/N): ")
	config.UseDatabase = askYesNo("🗄️ Inclure configuration database? (y/N): ")

	fmt.Printf("\n🎯 Création du projet '%s'...\n", config.Name)

	if err := generateProject(config); err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	successMessage := fmt.Sprintf(`✅ Projet créé avec succès!

📁 Structure créée:
  %s/
  ├── main.go
  ├── go.mod
  ├── .env.example
  ├── handlers/
  ├── middleware/
  └── README.md

🚀 Pour commencer:
  cd %s
  go mod tidy
  go run main.go

🌐 Votre serveur sera disponible sur: http://localhost:%s`,
		config.Name, config.Name, config.Port)

	fmt.Println(successMessage)
}

func askYesNo(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func generateProject(config ProjectConfig) error {
	if err := os.MkdirAll(config.Name, 0755); err != nil {
		return err
	}

	folders := []string{"handlers", "middleware"}
	for _, folder := range folders {
		if err := os.MkdirAll(filepath.Join(config.Name, folder), 0755); err != nil {
			return err
		}
	}

	if err := generateMainFile(config); err != nil {
		return err
	}

	if err := generateGoMod(config); err != nil {
		return err
	}

	if err := generateEnvFile(config); err != nil {
		return err
	}

	if err := generateReadme(config); err != nil {
		return err
	}

	if err := generateBaseHandlers(config); err != nil {
		return err
	}

	if config.UseAuth {
		if err := generateAuthMiddleware(config); err != nil {
			return err
		}
	}

	return nil
}

func generateMainFile(config ProjectConfig) error {
	mainTemplate := `package main

import (
	"{{.Module}}/handlers"
{{if .UseAuth}}	"{{.Module}}/middleware"{{end}}
	"github.com/Bakemono-san/gofsen"
	"log"
)

func main() {
	// Créer une nouvelle instance Gofsen
	app := gofsen.New()

	// Middlewares globaux
	app.Use(gofsen.Logger())
	app.Use(gofsen.Recovery())
{{if .UseCORS}}	app.Use(gofsen.CORSFromEnv()){{end}}
{{if .UseAuth}}	app.Use(middleware.AuthMiddleware()){{end}}

	// Routes de base
	app.GET("/", handlers.HomeHandler)
	app.GET("/health", handlers.HealthHandler)

	// Groupes d'API
	api := app.Group("/api/v1")
	api.GET("/status", handlers.StatusHandler)

	// Afficher les routes
	app.PrintRoutes()

	// Démarrer le serveur
	log.Printf("🚀 Serveur %s démarré sur http://localhost:{{.Port}}", "{{.Name}}")
	app.Listen("{{.Port}}")
}
`

	t, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(config.Name, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, config)
}

func generateGoMod(config ProjectConfig) error {
	content := fmt.Sprintf(`module %s

go 1.21

require github.com/Bakemono-san/gofsen v1.2.0
`, config.Module)

	return os.WriteFile(filepath.Join(config.Name, "go.mod"), []byte(content), 0644)
}

func generateEnvFile(config ProjectConfig) error {
	envContent := `# Configuration Gofsen

# Port du serveur
PORT=` + config.Port + `

# Configuration CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With

# Configuration Auth (optionnel)
JWT_SECRET=your-super-secret-key-change-this-in-production
AUTH_ENABLED=true

# Configuration Database (optionnel)
DATABASE_URL=postgres://user:password@localhost:5432/dbname
DATABASE_ENABLED=false
`

	return os.WriteFile(filepath.Join(config.Name, ".env.example"), []byte(envContent), 0644)
}

func generateReadme(config ProjectConfig) error {
	readmeContent := `# ` + config.Name + `

Projet Gofsen généré automatiquement.

## 🚀 Démarrage

` + "```bash" + `
# Installer les dépendances
go mod tidy

# Copier la configuration
cp .env.example .env

# Démarrer le serveur
go run main.go
` + "```" + `

## 📁 Structure

- ` + "`main.go`" + ` - Point d'entrée de l'application
- ` + "`handlers/`" + ` - Handlers/Controllers
- ` + "`middleware/`" + ` - Middlewares personnalisés
- ` + "`.env.example`" + ` - Configuration d'exemple

## 🌐 Endpoints

- ` + "`GET /`" + ` - Page d'accueil
- ` + "`GET /health`" + ` - Health check
- ` + "`GET /api/v1/status`" + ` - Status API

## 📚 Documentation

- Framework Gofsen: https://github.com/Bakemono-san/gofsen
- Documentation: https://pkg.go.dev/github.com/Bakemono-san/gofsen
`

	return os.WriteFile(filepath.Join(config.Name, "README.md"), []byte(readmeContent), 0644)
}

func generateBaseHandlers(config ProjectConfig) error {
	handlersContent := `package handlers

import (
	"github.com/Bakemono-san/gofsen"
)

// HomeHandler handler pour la page d'accueil
func HomeHandler(c *gofsen.Context) {
	c.JSON(map[string]interface{}{
		"message":   "Bienvenue sur ` + config.Name + `!",
		"framework": "Gofsen",
		"version":   "1.2.0",
	})
}

// HealthHandler handler pour le health check
func HealthHandler(c *gofsen.Context) {
	c.JSON(map[string]interface{}{
		"status":    "OK",
		"service":   "` + config.Name + `",
		"framework": "Gofsen",
	})
}

// StatusHandler handler pour le status de l'API
func StatusHandler(c *gofsen.Context) {
	c.JSON(map[string]interface{}{
		"api":     "v1",
		"status":  "running",
		"service": "` + config.Name + `",
	})
}
`

	return os.WriteFile(filepath.Join(config.Name, "handlers", "base.go"), []byte(handlersContent), 0644)
}

func generateAuthMiddleware(config ProjectConfig) error {
	authContent := `package middleware

import (
	"strings"
	"github.com/Bakemono-san/gofsen"
)

// AuthMiddleware middleware d'authentification basique
func AuthMiddleware() gofsen.MiddlewareFunc {
	return func(c *gofsen.Context) {
		// Vérifier le header Authorization
		authHeader := c.Request.Header.Get("Authorization")
		
		if authHeader == "" {
			c.Error(401, "Missing Authorization header")
			return
		}
		
		// Vérifier le format Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(401, "Invalid Authorization format")
			return
		}
		
		token := strings.TrimPrefix(authHeader, "Bearer ")
		
		// TODO: Implémenter la validation du token JWT ici
		if token == "" {
			c.Error(401, "Invalid token")
			return
		}
		
		// Continuer vers le handler suivant
		c.Next()
	}
}
`

	return os.WriteFile(filepath.Join(config.Name, "middleware", "auth.go"), []byte(authContent), 0644)
}

func generateCode(codeType string) {
	switch codeType {
	case "route", "routes":
		generateRoutes()
	case "middleware":
		generateMiddleware()
	case "handler", "handlers":
		generateHandler()
	default:
		fmt.Printf("❌ Type de génération inconnu: %s\n", codeType)
		showGenerateHelp()
	}
}

func generateRoutes() {
	if len(os.Args) < 4 {
		fmt.Println("❌ Usage: gofsen-cli gen route <name>")
		return
	}

	name := os.Args[3]
	titleName := capitalizeFirst(name)
	fmt.Printf("🛤️ Génération des routes CRUD pour '%s'...\n", name)

	routeContent := fmt.Sprintf(`package handlers

import (
	"github.com/Bakemono-san/gofsen"
)

// %sHandler contient les handlers pour %s
type %sHandler struct {
	// TODO: Ajouter les dépendances (DB, services, etc.)
}

// New%sHandler crée une nouvelle instance du handler
func New%sHandler() *%sHandler {
	return &%sHandler{}
}

// GetAll%s récupère tous les %s
func (h *%sHandler) GetAll%s(c *gofsen.Context) {
	// TODO: Implémenter la logique de récupération
	c.JSON(map[string]interface{}{
		"message": "Liste des %s",
		"data":    []interface{}{},
	})
}

// Get%s récupère un %s par ID
func (h *%sHandler) Get%s(c *gofsen.Context) {
	id := c.Param("id")
	
	// TODO: Implémenter la logique de récupération par ID
	c.JSON(map[string]interface{}{
		"message": "%s trouvé",
		"id":      id,
		"data":    map[string]interface{}{},
	})
}

// Create%s crée un nouveau %s
func (h *%sHandler) Create%s(c *gofsen.Context) {
	// TODO: Définir la structure de données
	var data map[string]interface{}
	
	if err := c.BindJSON(&data); err != nil {
		c.Error(400, "Données invalides")
		return
	}
	
	// TODO: Implémenter la logique de création
	c.Status(201).JSON(map[string]interface{}{
		"message": "%s créé avec succès",
		"data":    data,
	})
}

// Update%s met à jour un %s
func (h *%sHandler) Update%s(c *gofsen.Context) {
	id := c.Param("id")
	
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.Error(400, "Données invalides")
		return
	}
	
	// TODO: Implémenter la logique de mise à jour
	c.JSON(map[string]interface{}{
		"message": "%s mis à jour",
		"id":      id,
		"data":    data,
	})
}

// Delete%s supprime un %s
func (h *%sHandler) Delete%s(c *gofsen.Context) {
	id := c.Param("id")
	
	// TODO: Implémenter la logique de suppression
	c.JSON(map[string]interface{}{
		"message": "%s supprimé",
		"id":      id,
	})
}
`,
		titleName, name, titleName,
		titleName, titleName, titleName, titleName,
		titleName, name, titleName, titleName, name,
		titleName, name, titleName, titleName, titleName,
		titleName, name, titleName, titleName, titleName,
		titleName, name, titleName, titleName, titleName,
		titleName, name, titleName, titleName, titleName)

	filename := fmt.Sprintf("handlers/%s.go", name)
	if err := os.WriteFile(filename, []byte(routeContent), 0644); err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("✅ Routes CRUD générées: %s\n", filename)

	successMessage := fmt.Sprintf(`
📋 Routes créées:
  GET    /%s          # Liste tous les %s
  GET    /%s/:id      # Récupère un %s par ID
  POST   /%s          # Crée un nouveau %s
  PUT    /%s/:id      # Met à jour un %s
  DELETE /%s/:id      # Supprime un %s

🔧 N'oubliez pas d'enregistrer les routes dans main.go:
  %sHandler := handlers.New%sHandler()
  // Puis ajouter les routes individuellement ou créer un groupe`,
		name, name, name, name, name, name, name, name, name, name,
		name, titleName)

	fmt.Println(successMessage)
}

func generateMiddleware() {
	if len(os.Args) < 4 {
		fmt.Println("❌ Usage: gofsen-cli gen middleware <name>")
		return
	}

	name := os.Args[3]
	titleName := capitalizeFirst(name)
	fmt.Printf("🔧 Génération du middleware '%s'...\n", name)

	middlewareContent := fmt.Sprintf(`package middleware

import (
	"github.com/Bakemono-san/gofsen"
	"log"
)

// %sMiddleware middleware personnalisé pour %s
func %sMiddleware() gofsen.MiddlewareFunc {
	return func(c *gofsen.Context) {
		// TODO: Implémenter la logique du middleware
		log.Printf("Middleware %s exécuté pour: %%s %%s", c.Request.Method, c.Request.URL.Path)
		
		// Continuer vers le handler suivant
		c.Next()
	}
}

// %sWithConfig middleware %s avec configuration
func %sWithConfig(config %sConfig) gofsen.MiddlewareFunc {
	return func(c *gofsen.Context) {
		// TODO: Utiliser la configuration
		log.Printf("Middleware %s avec config exécuté")
		
		c.Next()
	}
}

// %sConfig configuration pour le middleware %s
type %sConfig struct {
	// TODO: Définir les options de configuration
	Enabled bool
	Debug   bool
}
`,
		titleName, name, titleName, name,
		titleName, name, titleName, titleName, name,
		titleName, name, titleName)

	filename := fmt.Sprintf("middleware/%s.go", name)
	if err := os.WriteFile(filename, []byte(middlewareContent), 0644); err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("✅ Middleware généré: %s\n", filename)

	successMessage := fmt.Sprintf(`
🔧 Pour utiliser le middleware:
  app.Use(middleware.%sMiddleware())
  
  # Ou avec configuration:
  config := middleware.%sConfig{Enabled: true, Debug: false}
  app.Use(middleware.%sWithConfig(config))`,
		titleName, titleName, titleName)

	fmt.Println(successMessage)
}

func generateHandler() {
	if len(os.Args) < 4 {
		fmt.Println("❌ Usage: gofsen-cli gen handler <name>")
		return
	}

	name := os.Args[3]
	titleName := capitalizeFirst(name)
	fmt.Printf("📝 Génération du handler '%s'...\n", name)

	handlerContent := fmt.Sprintf(`package handlers

import (
	"github.com/Bakemono-san/gofsen"
)

// %sHandler handler pour %s
func %sHandler(c *gofsen.Context) {
	// TODO: Implémenter la logique du handler
	c.JSON(map[string]interface{}{
		"message": "Handler %s",
		"path":    c.Request.URL.Path,
		"method":  c.Request.Method,
	})
}

// %sStatus handler pour le status de %s
func %sStatus(c *gofsen.Context) {
	c.JSON(map[string]interface{}{
		"service": "%s",
		"status":  "OK",
		"version": "1.0.0",
	})
}
`,
		titleName, name, titleName, name,
		titleName, name, titleName, name)

	filename := fmt.Sprintf("handlers/%s.go", name)
	if err := os.WriteFile(filename, []byte(handlerContent), 0644); err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("✅ Handler généré: %s\n", filename)

	successMessage := fmt.Sprintf(`
📝 Handlers créés:
  %sHandler  # Handler principal
  %sStatus   # Handler de status

🔧 Pour utiliser les handlers:
  app.GET("/%s", handlers.%sHandler)
  app.GET("/%s/status", handlers.%sStatus)`,
		titleName, titleName, name, titleName, name, titleName)

	fmt.Println(successMessage)
}
