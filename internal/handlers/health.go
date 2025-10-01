package handlers

import (
	"gofsen/internal/types"

	"net/http"
)

func HealthCheckHandler(ctx *types.Context) {
	ctx.BindJSON(map[string]string{})
	ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
