package utils

import (
	"encoding/json"
	"gofsen/internal/types"
	"log"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message,omitempty"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Timestamp string `json:"timestamp"`
	Code      int    `json:"code"`
	RequestID string `json:"request_id,omitempty"`
}

type DetailedErrorResponse struct {
	ErrorResponse
	Details interface{}       `json:"details,omitempty"`
	Trace   string            `json:"trace,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

type GofsenLogger struct {
	Level                LogLevel
	EnableDetailedErrors bool
}

var defaultLogger = &GofsenLogger{
	Level:                LogInfo,
	EnableDetailedErrors: true,
}

func GetLogger() *GofsenLogger {
	return defaultLogger
}

func (gl *GofsenLogger) LogRouteNotFound(ctx *types.Context) {
	log.Printf("🔍 [404] Route not found: %s %s", ctx.Request.Method, ctx.Request.URL.Path)

	if gl.Level <= LogDebug {
		log.Printf("📝 Available routes debug info:")
		log.Printf("   - Method: %s", ctx.Request.Method)
		log.Printf("   - Path: %s", ctx.Request.URL.Path)
		log.Printf("   - Query: %s", ctx.Request.URL.RawQuery)
		log.Printf("   - User-Agent: %s", ctx.Request.Header.Get("User-Agent"))
		log.Printf("   - Origin: %s", ctx.Request.Header.Get("Origin"))
	}
}

func (gl *GofsenLogger) LogMethodNotAllowed(ctx *types.Context, allowedMethods []string) {
	log.Printf("❌ [405] Method not allowed: %s %s (allowed: %v)",
		ctx.Request.Method, ctx.Request.URL.Path, allowedMethods)
}

func (gl *GofsenLogger) LogServerError(ctx *types.Context, err error) {
	log.Printf("💥 [500] Server error on %s %s: %v",
		ctx.Request.Method, ctx.Request.URL.Path, err)

	if gl.Level <= LogDebug {
		log.Printf("📝 Error context:")
		log.Printf("   - User-Agent: %s", ctx.Request.Header.Get("User-Agent"))
		log.Printf("   - Content-Type: %s", ctx.Request.Header.Get("Content-Type"))
		log.Printf("   - Content-Length: %s", ctx.Request.Header.Get("Content-Length"))
	}
}

func (gl *GofsenLogger) LogAuthFailure(ctx *types.Context, reason string) {
	log.Printf("🔐 [401] Auth failure on %s %s: %s",
		ctx.Request.Method, ctx.Request.URL.Path, reason)

	if gl.Level <= LogDebug {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader != "" {
			log.Printf("   - Auth header present: %s...", authHeader[:min(len(authHeader), 20)])
		} else {
			log.Printf("   - No Authorization header")
		}
	}
}

func (gl *GofsenLogger) SendDetailedError(ctx *types.Context, code int, message string, details interface{}) {
	errorResp := ErrorResponse{
		Error:     http.StatusText(code),
		Message:   message,
		Path:      ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Timestamp: time.Now().Format(time.RFC3339),
		Code:      code,
	}

	if gl.EnableDetailedErrors && gl.Level <= LogDebug {
		detailedResp := DetailedErrorResponse{
			ErrorResponse: errorResp,
			Details:       details,
			Headers: map[string]string{
				"User-Agent":    ctx.Request.Header.Get("User-Agent"),
				"Content-Type":  ctx.Request.Header.Get("Content-Type"),
				"Authorization": maskAuthHeader(ctx.Request.Header.Get("Authorization")),
				"Origin":        ctx.Request.Header.Get("Origin"),
			},
		}

		ctx.Writer.Header().Set("Content-Type", "application/json")
		ctx.Writer.WriteHeader(code)
		json.NewEncoder(ctx.Writer).Encode(detailedResp)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(code)
	json.NewEncoder(ctx.Writer).Encode(errorResp)
}

func maskAuthHeader(auth string) string {
	if auth == "" {
		return ""
	}
	if len(auth) <= 10 {
		return "***"
	}
	return auth[:10] + "***"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SuggestSimilarRoutes(requestedPath string, availableRoutes []string) []string {
	suggestions := []string{}

	for _, route := range availableRoutes {
		if len(route) > 0 && len(requestedPath) > 0 {
			if route[:1] == requestedPath[:1] ||
				(len(route) > 3 && len(requestedPath) > 3 && route[len(route)-3:] == requestedPath[len(requestedPath)-3:]) {
				suggestions = append(suggestions, route)
			}
		}
	}

	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return suggestions
}

var ErrorMessages = map[int]string{
	400: "Requête malformée",
	401: "Authentification requise",
	403: "Accès interdit",
	404: "Ressource non trouvée",
	405: "Méthode non autorisée",
	409: "Conflit de données",
	422: "Données invalides",
	429: "Trop de requêtes",
	500: "Erreur interne du serveur",
	502: "Service indisponible",
	503: "Service temporairement indisponible",
}

func GetFriendlyErrorMessage(code int) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "Erreur inconnue"
}
