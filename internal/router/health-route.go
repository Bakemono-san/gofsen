package router

import (
	"fmt"
	"gofsen/internal/handlers"
	"gofsen/internal/middlewares"
	"gofsen/internal/types"
	"gofsen/internal/utils"
)

func (r *Router) RegisterHealthRoutes() {

	tokenValidator := utils.NewTokenValidator()

	r.GET("/health", handlers.HealthCheckHandler)
	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware(tokenValidator))
	api.GET("/me", func(ctx *types.Context) {
		fmt.Fprintln(ctx.Writer, `{"user":"profile data"}`)
	})
}
