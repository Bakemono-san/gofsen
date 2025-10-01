package main

// 🚀 Gofsen Framework - Exemple Standalone
// Cet exemple peut être copié et utilisé directement dans un nouveau projet Go

import (
	"log"
	"net/http"
	"time"
)

// ============================================================================
// TYPES & STRUCTURES CORE GOFSEN
// ============================================================================

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Params  map[string]string
}

type HandlerFunc func(*Context)
type Middleware func(HandlerFunc) HandlerFunc

type Router struct {
	routes      map[string]map[string]HandlerFunc
	middlewares []Middleware
}

type RouteGroup struct {
	prefix      string
	parent      *Router
	middlewares []Middleware
}

// ============================================================================
// CONTEXT HELPERS
// ============================================================================

func (c *Context) JSON(status int, data interface{}) error {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	
	// Simple JSON encoding pour l'exemple
	jsonStr := `{"message":"` + data.(map[string]string)["message"] + `"}`
	_, err := c.Writer.Write([]byte(jsonStr))
	return err
}

func (c *Context) QueryParam(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Error(status int, message string) {
	c.JSON(status, map[string]string{"error": message})
}

// ============================================================================
// ROUTER CORE
// ============================================================================

func NewRouter() *Router {
	return &Router{
		routes:      make(map[string]map[string]HandlerFunc),
		middlewares: []Middleware{},
	}
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

func (r *Router) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix:      prefix,
		parent:      r,
		middlewares: []Middleware{},
	}
}

func (g *RouteGroup) Use(mws ...Middleware) {
	g.middlewares = append(g.middlewares, mws...)
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (g *RouteGroup) Handle(method, path string, handler HandlerFunc) {
	fullPath := g.prefix + path
	
	finalHandler := handler
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		finalHandler = g.middlewares[i](finalHandler)
	}
	
	g.parent.Handle(method, fullPath, finalHandler)
}

// HTTP Methods
func (r *Router) GET(path string, handler HandlerFunc) {
	r.Handle("GET", path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.Handle("POST", path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) {
	r.Handle("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.Handle("DELETE", path, handler)
}

func (g *RouteGroup) GET(path string, handler HandlerFunc) {
	g.Handle("GET", path, handler)
}

func (g *RouteGroup) POST(path string, handler HandlerFunc) {
	g.Handle("POST", path, handler)
}

// ServeHTTP implémente http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request: req,
		Writer:  w,
	}

	if methodRoutes, ok := r.routes[req.Method]; ok {
		if handler, ok := methodRoutes[req.URL.Path]; ok {
			finalHandler := handler
			for i := len(r.middlewares) - 1; i >= 0; i-- {
				finalHandler = r.middlewares[i](finalHandler)
			}

			finalHandler(ctx)
			return
		}
	}
	
	// 404 Not Found
	ctx.Error(http.StatusNotFound, "Route not found: " + req.Method + " " + req.URL.Path)
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

func LoggerMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		log.Printf("Started %s %s", ctx.Request.Method, ctx.Request.URL.Path)
		
		next(ctx)
		
		duration := time.Since(start)
		log.Printf("Completed in %v", duration)
	}
}

func RecoveryMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				ctx.Error(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next(ctx)
	}
}

// ============================================================================
// EXEMPLE D'APPLICATION
// ============================================================================

func main() {
	// Créer le router
	r := NewRouter()

	// Middleware globaux
	r.Use(RecoveryMiddleware)
	r.Use(LoggerMiddleware)

	// Routes de base
	r.GET("/", func(ctx *Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "🚀 Gofsen Framework - Standalone Example",
		})
	})

	r.GET("/hello", func(ctx *Context) {
		name := ctx.QueryParam("name")
		if name == "" {
			name = "World"
		}
		
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Hello " + name + "!",
		})
	})

	// Route de test d'erreur
	r.GET("/error", func(ctx *Context) {
		ctx.Error(http.StatusBadRequest, "Ceci est un exemple d'erreur")
	})

	// Route de test panic
	r.GET("/panic", func(ctx *Context) {
		panic("Test du middleware recovery!")
	})

	// Groupe de routes API
	api := r.Group("/api")
	
	// Middleware local pour le groupe API
	api.Use(func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			ctx.Writer.Header().Set("X-API-Version", "1.0")
			next(ctx)
		}
	})

	api.GET("/status", func(ctx *Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "API is running",
		})
	})

	api.GET("/time", func(ctx *Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "Current time: " + time.Now().Format(time.RFC3339),
		})
	})

	// Démarrer le serveur
	log.Println("🚀 Gofsen Standalone Example running at http://localhost:3000")
	log.Println("")
	log.Println("🔗 Try these endpoints:")
	log.Println("   GET  http://localhost:3000/")
	log.Println("   GET  http://localhost:3000/hello?name=Alice")
	log.Println("   GET  http://localhost:3000/api/status")
	log.Println("   GET  http://localhost:3000/api/time")
	log.Println("   GET  http://localhost:3000/error")
	log.Println("   GET  http://localhost:3000/panic")
	log.Println("")
	log.Println("📋 Test commands:")
	log.Println("   curl http://localhost:3000/")
	log.Println("   curl http://localhost:3000/hello?name=Go")
	log.Println("   curl http://localhost:3000/api/status")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v\n", err)
	}
}

// ============================================================================
// UTILISATION:
// 
// 1. Copiez ce fichier dans un nouveau projet Go
// 2. Initialisez le module: go mod init mon-projet
// 3. Lancez: go run main.go
// 4. Testez: curl http://localhost:3000/
//
// Ce fichier contient toutes les fonctionnalités core de Gofsen
// dans un seul fichier pour faciliter l'adoption et les tests.
// ============================================================================