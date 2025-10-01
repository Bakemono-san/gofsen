package router

import (
	"gofsen/internal/types"
	"gofsen/internal/utils"
	"net/http"
)

// Helper function to mask auth tokens
func maskAuthHeader(auth string) string {
	if auth == "" {
		return ""
	}
	if len(auth) <= 10 {
		return "***"
	}
	return auth[:10] + "***"
}

// RegisterErrorDemoRoutes ajoute des routes pour démontrer les messages d'erreur améliorés
func (r *Router) RegisterErrorDemoRoutes() {
	// Route principale de démonstration des erreurs
	r.GET("/demo/errors", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "🚨 Démonstration des messages d'erreur Gofsen",
			"features": []string{
				"Messages d'erreur détaillés",
				"Logging contextuel",
				"Suggestions de routes",
				"Gestion 404/405 intelligente",
			},
			"tests": map[string]interface{}{
				"route_404": map[string]string{
					"url":         "/demo/nonexistent",
					"description": "Route inexistante avec suggestions",
				},
				"method_405": map[string]string{
					"url":         "POST /demo/errors (cette route n'accepte que GET)",
					"description": "Méthode non autorisée",
				},
				"auth_401": map[string]string{
					"url":         "/demo/errors/protected",
					"description": "Erreur d'authentification détaillée",
				},
				"validation_400": map[string]string{
					"url":         "/demo/errors/validate?name=",
					"description": "Erreur de validation avec aide",
				},
				"panic_500": map[string]string{
					"url":         "/demo/errors/panic",
					"description": "Erreur serveur avec stack trace",
				},
			},
			"tips": []string{
				"Les messages d'erreur incluent suggestions et contexte",
				"Les logs montrent des détails pour le debug",
				"Mode debug vs production configurable",
			},
		})
	})

	// Route pour tester la validation avec erreurs détaillées
	r.GET("/demo/errors/validate", func(ctx *types.Context) {
		logger := utils.GetLogger()
		name := ctx.QueryParam("name")
		email := ctx.QueryParam("email")

		// Validation avec messages détaillés
		if name == "" {
			logger.SendDetailedError(ctx, http.StatusBadRequest,
				"Paramètre 'name' requis",
				map[string]interface{}{
					"missing_parameter": "name",
					"example_url":       "/demo/errors/validate?name=John&email=john@example.com",
					"validation_rules": map[string]string{
						"name":  "Obligatoire, minimum 2 caractères",
						"email": "Optionnel, format email valide",
					},
				})
			return
		}

		if len(name) < 2 {
			logger.SendDetailedError(ctx, http.StatusBadRequest,
				"Paramètre 'name' trop court",
				map[string]interface{}{
					"provided_value":  name,
					"provided_length": len(name),
					"minimum_length":  2,
					"suggestion":      "Utilisez un nom d'au moins 2 caractères",
				})
			return
		}

		response := map[string]interface{}{
			"message": "✅ Validation réussie",
			"data": map[string]string{
				"name":  name,
				"email": email,
			},
		}

		if email == "" {
			response["note"] = "Email optionnel non fourni"
		}

		ctx.JSON(http.StatusOK, response)
	})

	// Route protégée pour tester les erreurs d'auth détaillées
	tokenValidator := utils.NewTokenValidator()
	errorGroup := r.Group("/demo/errors")
	errorGroup.Use(func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			logger := utils.GetLogger()
			token := ctx.Request.Header.Get("Authorization")

			if token == "" {
				logger.LogAuthFailure(ctx, "Demo route - missing auth")
				logger.SendDetailedError(ctx, http.StatusUnauthorized,
					"Démonstration d'erreur d'authentification",
					map[string]interface{}{
						"demo_note":       "Cette route est protégée pour la démonstration",
						"required_header": "Authorization: Bearer valid-token",
						"example_curl":    "curl -H 'Authorization: Bearer valid-token' " + ctx.Request.URL.String(),
						"valid_tokens":    []string{"Bearer valid-token"},
					})
				return
			}

			if !tokenValidator.ValidateToken(token) {
				logger.LogAuthFailure(ctx, "Demo route - invalid token")
				logger.SendDetailedError(ctx, http.StatusUnauthorized,
					"Token d'authentification invalide",
					map[string]interface{}{
						"provided_token": maskAuthHeader(token),
						"valid_format":   "Bearer <token>",
						"demo_tokens":    []string{"Bearer valid-token"},
						"note":           "Pour la démonstration, utilisez 'Bearer valid-token'",
					})
				return
			}

			next(ctx)
		}
	})

	errorGroup.GET("/protected", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":    "✅ Accès autorisé à la route protégée",
			"demo":       "Cette route nécessite une authentification valide",
			"your_token": "valide",
		})
	})

	// Route pour tester les panics avec recovery détaillé
	r.GET("/demo/errors/panic", func(ctx *types.Context) {
		panicType := ctx.QueryParam("type")

		switch panicType {
		case "nil":
			var ptr *string
			_ = *ptr // Provoque un panic nil pointer
		case "slice":
			arr := []string{"a", "b"}
			_ = arr[10] // Index out of bounds
		case "map":
			m := make(map[string]interface{})
			m["key"] = "value"
			_ = m["key"].(int) // Type assertion panic
		default:
			panic("Démonstration d'une panic - le middleware Recovery va gérer cette erreur avec des détails complets!")
		}
	})

	// Route pour tester différents codes d'erreur
	r.GET("/demo/errors/codes", func(ctx *types.Context) {
		logger := utils.GetLogger()
		code := ctx.QueryParam("code")

		switch code {
		case "400":
			logger.SendDetailedError(ctx, http.StatusBadRequest,
				"Exemple d'erreur 400 - Requête malformée",
				map[string]interface{}{
					"common_causes": []string{
						"JSON malformé",
						"Paramètres manquants",
						"Format de données incorrect",
					},
					"fix_suggestions": []string{
						"Vérifiez le format JSON",
						"Ajoutez les paramètres requis",
						"Consultez la documentation API",
					},
				})
		case "403":
			logger.SendDetailedError(ctx, http.StatusForbidden,
				"Exemple d'erreur 403 - Accès interdit",
				map[string]interface{}{
					"reason":               "Permissions insuffisantes",
					"required_permissions": []string{"read", "write"},
					"your_permissions":     []string{"read"},
					"contact":              "Contactez l'administrateur pour plus de permissions",
				})
		case "409":
			logger.SendDetailedError(ctx, http.StatusConflict,
				"Exemple d'erreur 409 - Conflit de données",
				map[string]interface{}{
					"conflict_reason":   "Ressource déjà existante",
					"existing_resource": "user@example.com",
					"suggestion":        "Utilisez PUT pour modifier ou choisissez un autre identifiant",
				})
		case "429":
			logger.SendDetailedError(ctx, http.StatusTooManyRequests,
				"Exemple d'erreur 429 - Trop de requêtes",
				map[string]interface{}{
					"rate_limit":    "100 requêtes/minute",
					"current_usage": "150 requêtes/minute",
					"retry_after":   "60 seconds",
					"upgrade_info":  "Contactez-nous pour augmenter votre limite",
				})
		default:
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"message": "🎯 Test des codes d'erreur HTTP",
				"available_codes": map[string]string{
					"400": "/demo/errors/codes?code=400",
					"403": "/demo/errors/codes?code=403",
					"409": "/demo/errors/codes?code=409",
					"429": "/demo/errors/codes?code=429",
				},
				"note": "Chaque code d'erreur inclut des détails contextuels utiles",
			})
		}
	})
}
