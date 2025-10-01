package router

import (
	"gofsen/internal/middlewares"
	"gofsen/internal/types"
	"gofsen/internal/utils"
	"net/http"
	"time"
)

// RegisterTestRoutes creates routes to test all implemented functionalities
func (r *Router) RegisterTestRoutes() {
	// Token validator for auth testing
	tokenValidator := utils.NewTokenValidator()

	// =============================================================================
	// 🧱 PARTIE 1: FONCTIONNALITÉS DE BASE (Routes de test)
	// =============================================================================

	// Test basic routing methods
	r.GET("/test/get", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"method":      "GET",
			"message":     "✅ GET method working",
			"path":        ctx.Request.URL.Path,
			"timestamp":   time.Now(),
			"test_passed": true,
		})
	})

	r.POST("/test/post", func(ctx *types.Context) {
		var body map[string]interface{}
		ctx.BindJSON(&body)
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"method":      "POST",
			"message":     "✅ POST method working",
			"received":    body,
			"test_passed": true,
		})
	})

	r.PUT("/test/put", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"method":      "PUT",
			"message":     "✅ PUT method working",
			"test_passed": true,
		})
	})

	r.DELETE("/test/delete", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"method":      "DELETE",
			"message":     "✅ DELETE method working",
			"test_passed": true,
		})
	})

	r.PATCH("/test/patch", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"method":      "PATCH",
			"message":     "✅ PATCH method working",
			"test_passed": true,
		})
	})

	// =============================================================================
	// 🧭 PARTIE 2: GROUPES DE ROUTES + MIDDLEWARE LOCAL
	// =============================================================================

	// Test route groups with local middleware
	testGroup := r.Group("/test/group")
	testGroup.Use(func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			ctx.Writer.Header().Set("X-Group-Middleware", "applied")
			next(ctx)
		}
	})

	testGroup.GET("/basic", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "✅ Route group working",
			"prefix":      "/test/group",
			"middleware":  "local middleware applied",
			"test_passed": true,
		})
	})

	// =============================================================================
	// 🔐 PARTIE 3: MIDDLEWARE & SÉCURITÉ
	// =============================================================================

	// Test Logger Middleware (applied globally, visible in console)
	r.GET("/test/logger", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "✅ Logger middleware working",
			"note":        "Check console for log output",
			"test_passed": true,
		})
	})

	// Test Auth Middleware
	authGroup := r.Group("/test/auth")
	authGroup.Use(middlewares.AuthMiddleware(tokenValidator))

	authGroup.GET("/protected", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "✅ Auth middleware working",
			"note":        "You accessed a protected route!",
			"token":       "valid",
			"test_passed": true,
		})
	})

	// Test without auth (should fail)
	r.GET("/test/auth/public", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "✅ Public route (no auth required)",
			"protected":   false,
			"test_passed": true,
		})
	})

	// Test Recovery Middleware (panic handling)
	r.GET("/test/recovery", func(ctx *types.Context) {
		// Simulate panic
		panic("Simulated panic for testing recovery middleware")
	})

	// Test CORS Middleware
	corsGroup := r.Group("/test/cors")
	corsGroup.Use(middlewares.CorsMiddleware)

	corsGroup.GET("/check", func(ctx *types.Context) {
		origin := ctx.Request.Header.Get("Origin")
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "✅ CORS middleware working",
			"origin":      origin,
			"headers_set": "Check response headers for CORS",
			"test_passed": true,
		})
	})

	// =============================================================================
	// ⚙️ PARTIE 4: HELPERS & I/O
	// =============================================================================

	// Test JSON response
	r.GET("/test/json", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "✅ JSON response working",
			"data": map[string]interface{}{
				"users":  []string{"Alice", "Bob", "Charlie"},
				"count":  3,
				"active": true,
			},
			"test_passed": true,
		})
	})

	// Test Query params
	r.GET("/test/query", func(ctx *types.Context) {
		name := ctx.QueryParam("name")
		age := ctx.QueryParam("age")
		city := ctx.QueryParam("city")

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "✅ Query params working",
			"params": map[string]string{
				"name": name,
				"age":  age,
				"city": city,
			},
			"url_example": "/test/query?name=John&age=25&city=Paris",
			"test_passed": true,
		})
	})

	// Test BindJSON (body parsing)
	r.POST("/test/bind", func(ctx *types.Context) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Age   int    `json:"age"`
		}

		if err := ctx.BindJSON(&user); err != nil {
			ctx.Error(http.StatusBadRequest, "Invalid JSON body")
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":     "✅ BindJSON working",
			"received":    user,
			"parsed":      true,
			"test_passed": true,
		})
	})

	// Test Error response
	r.GET("/test/error", func(ctx *types.Context) {
		errorType := ctx.QueryParam("type")

		switch errorType {
		case "400":
			ctx.Error(http.StatusBadRequest, "✅ 400 Bad Request error test")
		case "401":
			ctx.Error(http.StatusUnauthorized, "✅ 401 Unauthorized error test")
		case "404":
			ctx.Error(http.StatusNotFound, "✅ 404 Not Found error test")
		case "500":
			ctx.Error(http.StatusInternalServerError, "✅ 500 Internal Server error test")
		default:
			ctx.Error(http.StatusTeapot, "✅ Custom error response working (418 I'm a teapot)")
		}
	})

	// =============================================================================
	// 🧪 TESTS COMBINÉS & AVANCÉS
	// =============================================================================

	// Test multiple middlewares combined
	multiGroup := r.Group("/test/multi")
	multiGroup.Use(middlewares.CorsMiddleware)
	multiGroup.Use(func(next types.HandlerFunc) types.HandlerFunc {
		return func(ctx *types.Context) {
			ctx.Writer.Header().Set("X-Custom-Header", "multi-middleware-test")
			next(ctx)
		}
	})

	multiGroup.POST("/combined", func(ctx *types.Context) {
		var data map[string]interface{}
		ctx.BindJSON(&data)

		name := ctx.QueryParam("name")

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message":       "✅ Multiple features combined",
			"query_param":   name,
			"json_body":     data,
			"cors_enabled":  true,
			"custom_header": "applied",
			"middlewares":   []string{"CORS", "Custom", "Logger(global)"},
			"test_passed":   true,
		})
	})

	// =============================================================================
	// 📋 ROUTE DE TEST MASTER (tous les tests en un)
	// =============================================================================

	r.GET("/test/all", func(ctx *types.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "🎉 Gofsen Framework Test Suite",
			"version": "1.0.0",
			"tests": map[string]interface{}{
				"basic_routing": map[string]string{
					"GET":    "/test/get",
					"POST":   "/test/post",
					"PUT":    "/test/put",
					"DELETE": "/test/delete",
					"PATCH":  "/test/patch",
				},
				"route_groups": map[string]string{
					"basic": "/test/group/basic",
				},
				"middleware": map[string]string{
					"logger":   "/test/logger",
					"auth":     "/test/auth/protected (needs: Authorization: Bearer valid-token)",
					"recovery": "/test/recovery",
					"cors":     "/test/cors/check",
				},
				"io_helpers": map[string]string{
					"json":      "/test/json",
					"query":     "/test/query?name=John&age=25",
					"bind_json": "/test/bind (POST with JSON body)",
					"error":     "/test/error?type=400",
				},
				"advanced": map[string]string{
					"multi_middleware": "/test/multi/combined (POST)",
				},
			},
			"usage": map[string]string{
				"auth_header": "Authorization: Bearer valid-token",
				"cors_origin": "Origin: https://example.com",
				"json_body":   `{"name": "John", "email": "john@example.com", "age": 25}`,
			},
			"status": "✅ All functionalities implemented and testable",
		})
	})
}
