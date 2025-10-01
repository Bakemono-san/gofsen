package utils

import "gofsen/internal/types"

func SetCORSHeaders(ctx *types.Context, origin string) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
}

func IsAllowedOrigin(origin string) bool {
	allowedOrigins := []string{"https://example.com", "https://another-example.com"}
	for _, o := range allowedOrigins {
		if o == origin {
			return true
		}
	}
	return false
}
