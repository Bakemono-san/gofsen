package middlewares

import (
	"encoding/json"
	"net/http"

	"gofsen/internal/types"
)

func JSONEncoder(next types.HandlerFunc) types.HandlerFunc {
	return func(ctx *types.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		next(ctx)
	}
}

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ParseJSONRequest(r *http.Request, dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}
