package router

import (
	"gofsen/internal/types"
	"gofsen/internal/utils"
	"net/http"
)

type Router struct {
	routes      map[string]map[string]types.HandlerFunc
	middlewares []types.Middleware
}

type RouteGroup struct {
	prefix      string
	parent      *Router
	middlewares []types.Middleware
}

func NewRouter() *Router {
	r := &Router{
		routes:      make(map[string]map[string]types.HandlerFunc),
		middlewares: []types.Middleware{},
	}
	r.RegisterHealthRoutes()
	r.RegisterTestRoutes()      // Add test routes
	r.RegisterErrorDemoRoutes() // Add error demo routes
	return r
}

func (r *Router) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix:      prefix,
		parent:      r,
		middlewares: []types.Middleware{},
	}
}

func (g *RouteGroup) Use(mws ...types.Middleware) {
	g.middlewares = append(g.middlewares, mws...)
}

func (g *RouteGroup) Handle(method, path string, handler types.HandlerFunc) {
	fullPath := g.prefix + path

	finalHandler := handler

	for i := len(g.middlewares) - 1; i >= 0; i-- {
		finalHandler = g.middlewares[i](finalHandler)
	}

	g.parent.Handle(method, fullPath, finalHandler)
}

func (r *Router) Use(mw ...types.Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

func (r *RouteGroup) GET(path string, handler types.HandlerFunc) {
	r.Handle("GET", path, handler)
}

func (r *RouteGroup) POST(path string, handler types.HandlerFunc) {
	r.Handle("POST", path, handler)
}

func (r *RouteGroup) PUT(path string, handler types.HandlerFunc) {
	r.Handle("PUT", path, handler)
}

func (r *RouteGroup) DELETE(path string, handler types.HandlerFunc) {
	r.Handle("DELETE", path, handler)
}

func (r *RouteGroup) PATCH(path string, handler types.HandlerFunc) {
	r.Handle("PATCH", path, handler)
}

func (r *Router) Handle(method, path string, handler types.HandlerFunc, key ...string) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]types.HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) GET(path string, handler types.HandlerFunc) {
	r.Handle("GET", path, handler)
}

func (r *Router) POST(path string, handler types.HandlerFunc) {
	r.Handle("POST", path, handler)
}

func (r *Router) PUT(path string, handler types.HandlerFunc) {
	r.Handle("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler types.HandlerFunc) {
	r.Handle("DELETE", path, handler)
}

func (r *Router) PATCH(path string, handler types.HandlerFunc) {
	r.Handle("PATCH", path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &types.Context{
		Request: req,
		Writer:  w,
	}

	logger := utils.GetLogger()

	// Vérifier si la méthode existe pour cette route
	if methodRoutes, ok := r.routes[req.Method]; ok {
		if handler, ok := methodRoutes[req.URL.Path]; ok {
			// Route trouvée - exécuter le handler avec middlewares
			finalHandler := handler
			for i := len(r.middlewares) - 1; i >= 0; i-- {
				finalHandler = r.middlewares[i](finalHandler)
			}

			finalHandler(ctx)
			return
		}
	}

	// Route non trouvée - vérifier si le chemin existe avec une autre méthode
	pathExists := false
	allowedMethods := []string{}

	for method, routes := range r.routes {
		if _, exists := routes[req.URL.Path]; exists {
			pathExists = true
			allowedMethods = append(allowedMethods, method)
		}
	}

	if pathExists {
		// Le chemin existe mais pas pour cette méthode HTTP
		logger.LogMethodNotAllowed(ctx, allowedMethods)
		logger.SendDetailedError(ctx, http.StatusMethodNotAllowed,
			"Méthode HTTP non autorisée pour cette route",
			map[string]interface{}{
				"allowed_methods": allowedMethods,
				"suggestion":      "Essayez avec: " + allowedMethods[0],
			})
		return
	}

	// Route complètement inexistante
	logger.LogRouteNotFound(ctx)

	// Obtenir la liste des routes disponibles pour suggestions
	availableRoutes := []string{}
	for _, routes := range r.routes {
		for path := range routes {
			availableRoutes = append(availableRoutes, path)
		}
	}

	suggestions := utils.SuggestSimilarRoutes(req.URL.Path, availableRoutes)

	logger.SendDetailedError(ctx, http.StatusNotFound,
		"Route non trouvée",
		map[string]interface{}{
			"suggestions":      suggestions,
			"available_routes": availableRoutes,
			"tip":              "Vérifiez l'URL et la méthode HTTP",
		})
}
