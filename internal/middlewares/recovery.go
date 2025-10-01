package middlewares

import (
	"fmt"
	"gofsen/internal/types"
	"gofsen/internal/utils"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware(next types.HandlerFunc) types.HandlerFunc {
	return func(ctx *types.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger := utils.GetLogger()

				err := fmt.Errorf("panic recovered: %v", r)
				logger.LogServerError(ctx, err)

				stackTrace := string(debug.Stack())

				logger.SendDetailedError(ctx, http.StatusInternalServerError,
					"Erreur interne du serveur",
					map[string]interface{}{
						"panic_message": fmt.Sprintf("%v", r),
						"stack_trace":   stackTrace,
						"recovery_note": "L'application a récupéré d'une erreur critique",
					})
			}
		}()
		next(ctx)
	}
}
