package middlewares

import (
	"gofsen/internal/types"
	"gofsen/internal/utils"
	"net/http"
)

func AuthMiddleware(validator types.TokenValidator) types.Middleware {
	return func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			logger := utils.GetLogger()
			token := ctx.Request.Header.Get("Authorization")

			if token == "" {
				logger.LogAuthFailure(ctx, "Missing Authorization header")
				logger.SendDetailedError(ctx, http.StatusUnauthorized,
					"En-tête Authorization manquant",
					map[string]interface{}{
						"required_format": "Authorization: Bearer <token>",
						"example":         "Authorization: Bearer valid-token",
					})
				return
			}

			if !validator.ValidateToken(token) {
				logger.LogAuthFailure(ctx, "Invalid token")
				logger.SendDetailedError(ctx, http.StatusUnauthorized,
					"Token d'authentification invalide",
					map[string]interface{}{
						"token_format": "Bearer <token>",
						"note":         "Vérifiez que votre token est valide et non expiré",
					})
				return
			}

			next(ctx)
		}
	}
}
