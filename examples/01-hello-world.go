package main

import (
	"gofsen/internal/middlewares"
	"gofsen/internal/router"
	"gofsen/internal/types"
	"log"
	"net/http"
	"time"
)

func main() {
	// Créer le router
	r := router.NewRouter()

	// Ajouter un middleware global pour logger toutes les requêtes
	r.Use(middlewares.LoggerMiddleware)

	// Route simple GET
	r.GET("/", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":   "Bienvenue sur Gofsen! 🚀",
			"framework": "Gofsen",
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Route avec paramètres de requête
	r.GET("/hello", func(ctx *types.Context) {
		name := ctx.QueryParam("name")
		if name == "" {
			name = "World"
		}

		language := ctx.QueryParam("lang")
		var greeting string
		switch language {
		case "fr":
			greeting = "Bonjour"
		case "es":
			greeting = "Hola"
		case "de":
			greeting = "Hallo"
		default:
			greeting = "Hello"
		}

		ctx.JSON(http.StatusOK, map[string]string{
			"greeting": greeting + " " + name + "!",
			"language": language,
		})
	})

	// Route POST avec JSON
	r.POST("/echo", func(ctx *types.Context) {
		var body map[string]interface{}
		if err := ctx.BindJSON(&body); err != nil {
			ctx.Error(http.StatusBadRequest, "Invalid JSON")
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":      "Echo de votre message",
			"received":     body,
			"content_type": ctx.Request.Header.Get("Content-Type"),
			"timestamp":    time.Now(),
		})
	})

	// Routes pour différentes méthodes HTTP
	r.PUT("/data", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"method":  "PUT",
			"message": "Données mises à jour",
		})
	})

	r.DELETE("/data", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"method":  "DELETE",
			"message": "Données supprimées",
		})
	})

	// Route de test d'erreur
	r.GET("/error", func(ctx *types.Context) {
		errorType := ctx.QueryParam("type")

		switch errorType {
		case "400":
			ctx.Error(http.StatusBadRequest, "Erreur de requête")
		case "404":
			ctx.Error(http.StatusNotFound, "Resource non trouvée")
		case "500":
			ctx.Error(http.StatusInternalServerError, "Erreur serveur")
		default:
			ctx.Error(http.StatusTeapot, "Je suis une théière! ☕")
		}
	})

	// Route de status/health check
	r.GET("/status", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"uptime":    "running",
			"server":    "Gofsen",
			"timestamp": time.Now(),
			"endpoints": []string{
				"GET /",
				"GET /hello?name=John&lang=fr",
				"POST /echo",
				"PUT /data",
				"DELETE /data",
				"GET /error?type=404",
				"GET /status",
			},
		})
	})

	log.Println("🚀 Serveur Hello World démarré sur http://localhost:3000")
	log.Println("")
	log.Println("🔗 Endpoints disponibles:")
	log.Println("   GET  http://localhost:3000/")
	log.Println("   GET  http://localhost:3000/hello?name=Alice&lang=fr")
	log.Println("   POST http://localhost:3000/echo")
	log.Println("   GET  http://localhost:3000/status")
	log.Println("")
	log.Println("📋 Testez avec:")
	log.Println("   curl http://localhost:3000/")
	log.Println("   curl http://localhost:3000/hello?name=Alice&lang=fr")
	log.Println("   curl -X POST -H 'Content-Type: application/json' -d '{\"message\":\"test\"}' http://localhost:3000/echo")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Erreur serveur: %v\n", err)
	}
}
