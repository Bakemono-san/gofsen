package middlewares

import (
	"gofsen/internal/types"
	"gofsen/internal/utils"
	"net/http"
)

func CorsMiddleware(next types.HandlerFunc) types.HandlerFunc {
	return func(ctx *types.Context) {

		if ctx.Request.Header.Get("Origin") != "" {
			origin := ctx.Request.Header.Get("Origin")

			if utils.IsAllowedOrigin(origin) {
				utils.SetCORSHeaders(ctx, origin)
			}
		}

		if ctx.Request.Method == "OPTIONS" {
			ctx.Writer.WriteHeader(http.StatusOK)
			return
		}

		next(ctx)
	}
}
