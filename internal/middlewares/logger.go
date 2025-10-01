package middlewares

import (
	"gofsen/internal/types"
	"log"
	"time"
)

func LoggerMiddleware(next types.HandlerFunc) types.HandlerFunc {
	return func(ctx *types.Context) {
		start := time.Now()
		log.Printf("Started %s %s", ctx.Request.Method, ctx.Request.URL.Path)

		next(ctx)

		duration := time.Since(start)
		log.Printf("Completed in %v", duration)
	}
}
