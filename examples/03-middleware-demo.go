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
	r := router.NewRouter()

	// Middleware globaux
	r.Use(middlewares.LoggerMiddleware)
	r.Use(middlewares.RecoveryMiddleware)

	// Middleware custom pour ajouter des headers
	customHeaderMiddleware := func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			ctx.Writer.Header().Set("X-API-Version", "1.0")
			ctx.Writer.Header().Set("X-Powered-By", "Gofsen")
			next(ctx)
		}
	}

	// Middleware de timing
	timingMiddleware := func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			start := time.Now()
			next(ctx)
			duration := time.Since(start)
			ctx.Writer.Header().Set("X-Response-Time", duration.String())
		}
	}

	// Appliquer les middlewares custom globalement
	r.Use(customHeaderMiddleware)
	r.Use(timingMiddleware)

	// Route principale
	r.GET("/", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Démo des Middlewares Gofsen",
			"middlewares": []string{
				"LoggerMiddleware (global)",
				"RecoveryMiddleware (global)",
				"CustomHeaderMiddleware (global)",
				"TimingMiddleware (global)",
			},
			"check_headers": "Regardez les headers X-API-Version, X-Powered-By, X-Response-Time",
		})
	})

	// Groupe avec middleware spécifique
	apiGroup := r.Group("/api")

	// Middleware spécifique au groupe API
	apiMiddleware := func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			ctx.Writer.Header().Set("X-API-Group", "activated")
			log.Println("🔧 Middleware API group activé pour:", ctx.Request.URL.Path)
			next(ctx)
		}
	}

	apiGroup.Use(apiMiddleware)

	apiGroup.GET("/status", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"status": "API Group avec middleware local",
			"path":   ctx.Request.URL.Path,
			"middlewares": []string{
				"Global: Logger, Recovery, CustomHeader, Timing",
				"Local: ApiMiddleware",
			},
		})
	})

	// Groupe avec CORS
	corsGroup := r.Group("/cors")
	corsGroup.Use(middlewares.CorsMiddleware)

	corsGroup.GET("/test", func(ctx *types.Context) {
		origin := ctx.Request.Header.Get("Origin")
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":      "Test CORS",
			"origin":       origin,
			"cors_enabled": true,
			"note":         "Testez avec: curl -H 'Origin: https://example.com' http://localhost:3002/cors/test",
		})
	})

	// Route pour tester la récupération de panic
	r.GET("/panic", func(ctx *types.Context) {
		panicType := ctx.QueryParam("type")

		switch panicType {
		case "nil":
			var ptr *string
			_ = *ptr // Provoque un panic
		case "index":
			arr := []string{"a", "b"}
			_ = arr[10] // Index out of bounds
		default:
			panic("Panic de démonstration - le middleware Recovery va gérer ça!")
		}
	})

	// Routes avec middleware en chaîne
	secureGroup := r.Group("/secure")

	// Middleware de validation
	validateMiddleware := func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			apiKey := ctx.Request.Header.Get("X-API-Key")
			if apiKey == "" {
				ctx.Error(http.StatusBadRequest, "X-API-Key header requis")
				return
			}
			if apiKey != "demo-key-123" {
				ctx.Error(http.StatusUnauthorized, "X-API-Key invalide")
				return
			}
			next(ctx)
		}
	}

	// Middleware de limitation
	requestCountMiddleware := func(next types.HandlerFunc) types.HandlerFunc {
		count := 0
		return func(ctx *types.Context) {
			count++
			ctx.Writer.Header().Set("X-Request-Count", string(rune(count+48))) // Convert to ASCII
			log.Printf("🔢 Request #%d to %s", count, ctx.Request.URL.Path)
			next(ctx)
		}
	}

	secureGroup.Use(validateMiddleware)
	secureGroup.Use(requestCountMiddleware)

	secureGroup.GET("/data", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Données sécurisées",
			"data":    "Informations confidentielles",
			"middlewares": []string{
				"Global: Logger, Recovery, CustomHeader, Timing",
				"Local: Validate, RequestCount",
			},
			"note": "Cette route nécessite X-API-Key: demo-key-123",
		})
	})

	// Route pour tester la chaîne complète
	r.GET("/demo", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Démo complète des middlewares",
			"routes": map[string]interface{}{
				"GET /":            "Route de base avec middlewares globaux",
				"GET /api/status":  "Route avec middleware de groupe",
				"GET /cors/test":   "Route avec middleware CORS",
				"GET /panic":       "Test du middleware Recovery",
				"GET /secure/data": "Route avec validation et compteur",
			},
			"tests": map[string]string{
				"headers": "curl -I http://localhost:3002/",
				"cors":    "curl -H 'Origin: https://example.com' http://localhost:3002/cors/test",
				"panic":   "curl http://localhost:3002/panic",
				"secure":  "curl -H 'X-API-Key: demo-key-123' http://localhost:3002/secure/data",
				"invalid": "curl -H 'X-API-Key: wrong' http://localhost:3002/secure/data",
			},
		})
	})

	log.Println("🚀 Serveur Middleware Demo démarré sur http://localhost:3002")
	log.Println("")
	log.Println("🔧 Middlewares actifs:")
	log.Println("   Global: Logger, Recovery, CustomHeader, Timing")
	log.Println("   /api/*: ApiMiddleware")
	log.Println("   /cors/*: CorsMiddleware")
	log.Println("   /secure/*: Validate + RequestCount")
	log.Println("")
	log.Println("📋 Tests recommandés:")
	log.Println("   curl -I http://localhost:3002/ (voir headers)")
	log.Println("   curl http://localhost:3002/panic (test recovery)")
	log.Println("   curl -H 'X-API-Key: demo-key-123' http://localhost:3002/secure/data")
	log.Println("   curl -H 'Origin: https://example.com' http://localhost:3002/cors/test")

	if err := http.ListenAndServe(":3002", r); err != nil {
		log.Fatalf("Erreur serveur: %v\n", err)
	}
}
