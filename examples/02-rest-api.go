package main

import (
	"gofsen/internal/middlewares"
	"gofsen/internal/router"
	"gofsen/internal/types"
	"gofsen/internal/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Structures de données
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Created  string `json:"created"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Base de données simulée
var users = []User{
	{1, "Alice Dupont", "alice@example.com", "alice", time.Now().AddDate(0, -2, 0).Format(time.RFC3339)},
	{2, "Bob Martin", "bob@example.com", "bob", time.Now().AddDate(0, -1, -15).Format(time.RFC3339)},
	{3, "Charlie Durand", "charlie@example.com", "charlie", time.Now().AddDate(0, 0, -7).Format(time.RFC3339)},
}

func main() {
	r := router.NewRouter()

	// Middleware globaux
	r.Use(middlewares.LoggerMiddleware)
	r.Use(middlewares.RecoveryMiddleware)
	r.Use(middlewares.CorsMiddleware)

	// Routes publiques
	r.GET("/", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "API REST Users - Gofsen Example",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"GET /users":         "Liste tous les utilisateurs",
				"GET /users?limit=2": "Liste avec pagination",
				"POST /users":        "Créer un utilisateur (nécessite auth)",
				"GET /auth/profile":  "Profil utilisateur (nécessite auth)",
				"GET /auth/admin":    "Zone admin (nécessite auth)",
			},
			"auth": map[string]string{
				"header":  "Authorization: Bearer valid-token",
				"example": "curl -H 'Authorization: Bearer valid-token' http://localhost:3001/auth/profile",
			},
		})
	})

	// API publique
	api := r.Group("/api/v1")

	// Routes utilisateurs publiques
	api.GET("/users", func(ctx *types.Context) {
		// Pagination simple
		limitStr := ctx.QueryParam("limit")
		limit := len(users) // Par défaut, tous les users

		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		// Filtrage
		search := ctx.QueryParam("search")
		filteredUsers := users

		if search != "" {
			filteredUsers = []User{}
			for _, user := range users {
				if contains(user.Name, search) || contains(user.Username, search) || contains(user.Email, search) {
					filteredUsers = append(filteredUsers, user)
				}
			}
		}

		// Limitation
		if limit < len(filteredUsers) {
			filteredUsers = filteredUsers[:limit]
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"users":       filteredUsers,
			"total":       len(filteredUsers),
			"total_all":   len(users),
			"search":      search,
			"limit":       limit,
			"api_version": "v1",
		})
	})

	api.GET("/users/count", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"count":     len(users),
			"timestamp": time.Now(),
		})
	})

	// Health check
	api.GET("/health", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"database":  "connected (simulée)",
			"uptime":    "running",
			"users":     len(users),
			"timestamp": time.Now(),
		})
	})

	// Routes protégées par authentification
	tokenValidator := utils.NewTokenValidator()
	protected := r.Group("/auth")
	protected.Use(middlewares.AuthMiddleware(tokenValidator))

	protected.GET("/profile", func(ctx *types.Context) {
		// Simulation d'un utilisateur connecté
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"user": map[string]interface{}{
				"id":          999,
				"username":    "authenticated_user",
				"name":        "Utilisateur Connecté",
				"email":       "user@authenticated.com",
				"role":        "user",
				"permissions": []string{"read", "write"},
			},
			"session": map[string]interface{}{
				"authenticated": true,
				"login_time":    time.Now().Add(-2 * time.Hour),
				"expires":       time.Now().Add(6 * time.Hour),
			},
			"message": "Accès autorisé au profil utilisateur",
		})
	})

	protected.POST("/users", func(ctx *types.Context) {
		var newUser CreateUserRequest
		if err := ctx.BindJSON(&newUser); err != nil {
			ctx.Error(http.StatusBadRequest, "JSON invalide")
			return
		}

		// Validation simple
		if newUser.Name == "" || newUser.Email == "" || newUser.Username == "" {
			ctx.Error(http.StatusBadRequest, "Tous les champs sont requis (name, email, username)")
			return
		}

		// Vérifier si l'username existe déjà
		for _, user := range users {
			if user.Username == newUser.Username || user.Email == newUser.Email {
				ctx.Error(http.StatusConflict, "Username ou email déjà utilisé")
				return
			}
		}

		// Créer le nouvel utilisateur
		user := User{
			ID:       len(users) + 1,
			Name:     newUser.Name,
			Email:    newUser.Email,
			Username: newUser.Username,
			Created:  time.Now().Format(time.RFC3339),
		}

		users = append(users, user)

		ctx.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Utilisateur créé avec succès",
			"user":    user,
			"total":   len(users),
		})
	})

	protected.GET("/admin", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Zone administrateur",
			"access":  "granted",
			"admin_data": map[string]interface{}{
				"total_users":   len(users),
				"server_uptime": "running",
				"system_status": "healthy",
				"last_backup":   time.Now().Add(-6 * time.Hour),
				"pending_tasks": []string{"backup", "cleanup", "monitoring"},
			},
			"permissions": []string{"read", "write", "delete", "admin"},
		})
	})

	// Routes de démonstration d'erreurs
	r.GET("/demo/errors", func(ctx *types.Context) {
		errorType := ctx.QueryParam("type")

		switch errorType {
		case "400":
			ctx.Error(http.StatusBadRequest, "Exemple d'erreur de requête malformée")
		case "401":
			ctx.Error(http.StatusUnauthorized, "Exemple d'erreur d'authentification requise")
		case "403":
			ctx.Error(http.StatusForbidden, "Exemple d'erreur d'accès interdit")
		case "404":
			ctx.Error(http.StatusNotFound, "Exemple de resource non trouvée")
		case "500":
			ctx.Error(http.StatusInternalServerError, "Exemple d'erreur serveur interne")
		default:
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"message": "Démo des erreurs HTTP",
				"examples": map[string]string{
					"400": "/demo/errors?type=400",
					"401": "/demo/errors?type=401",
					"403": "/demo/errors?type=403",
					"404": "/demo/errors?type=404",
					"500": "/demo/errors?type=500",
				},
			})
		}
	})

	log.Println("🚀 API REST Users démarrée sur http://localhost:3001")
	log.Println("")
	log.Println("🔗 Endpoints publics:")
	log.Println("   GET  http://localhost:3001/")
	log.Println("   GET  http://localhost:3001/api/v1/users")
	log.Println("   GET  http://localhost:3001/api/v1/users?limit=2&search=alice")
	log.Println("   GET  http://localhost:3001/api/v1/health")
	log.Println("")
	log.Println("🔐 Endpoints protégés (nécessitent: Authorization: Bearer valid-token):")
	log.Println("   GET  http://localhost:3001/auth/profile")
	log.Println("   POST http://localhost:3001/auth/users")
	log.Println("   GET  http://localhost:3001/auth/admin")
	log.Println("")
	log.Println("📋 Commandes de test:")
	log.Println("   curl http://localhost:3001/api/v1/users")
	log.Println("   curl -H 'Authorization: Bearer valid-token' http://localhost:3001/auth/profile")
	log.Println("   curl -X POST -H 'Authorization: Bearer valid-token' -H 'Content-Type: application/json' \\")
	log.Println("        -d '{\"name\":\"Test User\",\"email\":\"test@example.com\",\"username\":\"testuser\"}' \\")
	log.Println("        http://localhost:3001/auth/users")

	if err := http.ListenAndServe(":3001", r); err != nil {
		log.Fatalf("Erreur serveur: %v\n", err)
	}
}

// Fonction helper pour la recherche
func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr ||
			(len(str) > len(substr) &&
				(str[:len(substr)] == substr ||
					str[len(str)-len(substr):] == substr ||
					containsMiddle(str, substr))))
}

func containsMiddle(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
