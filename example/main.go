// Example usage of Gofsen framework as an external package
package main

import (
	"github.com/Bakemono-san/gofsen"
)

// User représente un utilisateur
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Base de données simulée
var users = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com"},
	{ID: 2, Name: "Bob", Email: "bob@example.com"},
}

func main() {
	// Créer une nouvelle instance Gofsen
	app := gofsen.New()

	// Middlewares globaux
	app.Use(gofsen.Logger())
	app.Use(gofsen.Recovery())
	app.Use(gofsen.CORS())

	// Routes de base
	app.GET("/", func(c *gofsen.Context) {
		c.JSON(map[string]string{
			"message": "Welcome to Gofsen Framework!",
			"version": gofsen.Version,
		})
	})

	app.GET("/health", func(c *gofsen.Context) {
		c.JSON(map[string]string{
			"status":    "OK",
			"framework": "Gofsen",
		})
	})

	// API Users
	api := app.Group("/api/v1")

	// GET /api/v1/users
	api.GET("/users", func(c *gofsen.Context) {
		c.JSON(map[string]interface{}{
			"users": users,
			"count": len(users),
		})
	})

	// GET /api/v1/users/:id
	api.GET("/users/:id", func(c *gofsen.Context) {
		id := c.Param("id")

		for _, user := range users {
			if user.ID == parseID(id) {
				c.JSON(user)
				return
			}
		}

		c.Error(404, "User not found")
	})

	// POST /api/v1/users
	api.POST("/users", func(c *gofsen.Context) {
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			c.Error(400, "Invalid JSON")
			return
		}

		newUser.ID = len(users) + 1
		users = append(users, newUser)

		c.Status(201).JSON(newUser)
	})

	// Route avec query parameters
	app.GET("/search", func(c *gofsen.Context) {
		query := c.QueryParam("q")
		limit := c.QueryParam("limit")

		c.JSON(map[string]string{
			"query": query,
			"limit": limit,
		})
	})

	// Middleware d'authentification simple pour les routes protégées
	authMiddleware := func(c *gofsen.Context) {
		token := c.Request.Header.Get("Authorization")
		if token != "Bearer secret" {
			c.Error(401, "Unauthorized")
			return
		}
		c.Next()
	}

	// Groupe protégé
	protected := app.Group("/protected")
	protected.Use(authMiddleware)

	protected.GET("/profile", func(c *gofsen.Context) {
		c.JSON(map[string]string{
			"message": "This is a protected route",
			"user":    "authenticated",
		})
	})

	// Afficher toutes les routes
	app.PrintRoutes()

	// Démarrer le serveur
	app.Listen("3000")
}

// Helper function
func parseID(id string) int {
	// Simple conversion pour l'exemple
	if id == "1" {
		return 1
	}
	if id == "2" {
		return 2
	}
	return 0
}
