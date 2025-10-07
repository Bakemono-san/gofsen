// Package gofsen provides a lightweight, Express.js-inspired HTTP framework for Go
package gofsen

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Version de Gofsen
const Version = "v1.0.0"

// Route structure pour définir une route
type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
	Pattern *regexp.Regexp
	Params  []string
}

// Context encapsule les informations de la requête et réponse
type Context struct {
	Request         *http.Request
	ResponseWriter  http.ResponseWriter
	Params          map[string]string
	Query           map[string]string
	middleware      []MiddlewareFunc
	middlewareIndex int
}

// HandlerFunc définit le type de fonction pour les handlers
type HandlerFunc func(*Context)

// MiddlewareFunc définit le type de fonction pour les middlewares
type MiddlewareFunc func(*Context)

// Router structure principale du framework
type Router struct {
	routes      []Route
	middlewares []MiddlewareFunc
	groups      map[string]*RouteGroup
}

// RouteGroup pour organiser les routes
type RouteGroup struct {
	prefix      string
	middlewares []MiddlewareFunc
	router      *Router
}

// New crée une nouvelle instance du router Gofsen
func New() *Router {
	return &Router{
		routes: make([]Route, 0),
		groups: make(map[string]*RouteGroup),
	}
}

// Use ajoute un middleware global
func (r *Router) Use(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

// Group crée un groupe de routes avec un préfixe
func (r *Router) Group(prefix string) *RouteGroup {
	group := &RouteGroup{
		prefix: prefix,
		router: r,
	}
	r.groups[prefix] = group
	return group
}

// addRoute ajoute une route au router
func (r *Router) addRoute(method, path string, handler HandlerFunc) {
	// Convertir les paramètres dynamiques en regex
	pattern, params := convertPathToRegex(path)

	route := Route{
		Method:  method,
		Path:    path,
		Handler: handler,
		Pattern: pattern,
		Params:  params,
	}
	r.routes = append(r.routes, route)
}

// Méthodes HTTP
func (r *Router) GET(path string, handler HandlerFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.addRoute("POST", path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc) {
	r.addRoute("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) PATCH(path string, handler HandlerFunc) {
	r.addRoute("PATCH", path, handler)
}

// RouteGroup methods
func (g *RouteGroup) Use(middleware MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middleware)
}

func (g *RouteGroup) GET(path string, handler HandlerFunc) {
	g.router.addRoute("GET", g.prefix+path, handler)
}

func (g *RouteGroup) POST(path string, handler HandlerFunc) {
	g.router.addRoute("POST", g.prefix+path, handler)
}

func (g *RouteGroup) PUT(path string, handler HandlerFunc) {
	g.router.addRoute("PUT", g.prefix+path, handler)
}

func (g *RouteGroup) DELETE(path string, handler HandlerFunc) {
	g.router.addRoute("DELETE", g.prefix+path, handler)
}

func (g *RouteGroup) PATCH(path string, handler HandlerFunc) {
	g.router.addRoute("PATCH", g.prefix+path, handler)
}

// ServeHTTP implémente l'interface http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:         req,
		ResponseWriter:  w,
		Params:          make(map[string]string),
		Query:           parseQuery(req.URL.RawQuery),
		middleware:      r.middlewares,
		middlewareIndex: -1,
	}

	// Trouver la route correspondante
	route, params := r.findRoute(req.Method, req.URL.Path)
	if route == nil {
		ctx.Status(404).JSON(map[string]string{"error": "Route not found"})
		return
	}

	// Ajouter les paramètres de route au contexte
	ctx.Params = params

	// Exécuter les middlewares puis le handler
	finalHandler := func(c *Context) {
		route.Handler(c)
	}
	ctx.middleware = append(ctx.middleware, finalHandler)
	ctx.Next()
}

// findRoute trouve la route correspondante à la méthode et au chemin
func (r *Router) findRoute(method, path string) (*Route, map[string]string) {
	for _, route := range r.routes {
		if route.Method == method {
			if route.Pattern != nil {
				// Route avec paramètres dynamiques
				if matches := route.Pattern.FindStringSubmatch(path); matches != nil {
					params := make(map[string]string)
					for i, param := range route.Params {
						if i+1 < len(matches) {
							params[param] = matches[i+1]
						}
					}
					return &route, params
				}
			} else if route.Path == path {
				// Route exacte
				return &route, make(map[string]string)
			}
		}
	}
	return nil, nil
}

// convertPathToRegex convertit un chemin avec paramètres en regex
func convertPathToRegex(path string) (*regexp.Regexp, []string) {
	if !strings.Contains(path, ":") {
		return nil, nil
	}

	var params []string
	regexPath := path

	// Remplacer :param par ([^/]+)
	paramRegex := regexp.MustCompile(`:([^/]+)`)
	matches := paramRegex.FindAllStringSubmatch(path, -1)

	for _, match := range matches {
		params = append(params, match[1])
		regexPath = strings.Replace(regexPath, match[0], "([^/]+)", 1)
	}

	pattern, err := regexp.Compile("^" + regexPath + "$")
	if err != nil {
		return nil, nil
	}

	return pattern, params
}

// Listen démarre le serveur sur le port spécifié
func (r *Router) Listen(port string) error {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("🚀 Gofsen server listening on http://localhost%s", port)
	return http.ListenAndServe(port, r)
}

// Context methods

// Next exécute le middleware suivant dans la chaîne
func (c *Context) Next() {
	c.middlewareIndex++
	if c.middlewareIndex < len(c.middleware) {
		c.middleware[c.middlewareIndex](c)
	}
}

// Status définit le code de statut HTTP
func (c *Context) Status(code int) *Context {
	c.ResponseWriter.WriteHeader(code)
	return c
}

// JSON envoie une réponse JSON
func (c *Context) JSON(data interface{}) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.ResponseWriter).Encode(data)
}

// Text envoie une réponse texte
func (c *Context) Text(text string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/plain")
	c.ResponseWriter.Write([]byte(text))
}

// HTML envoie une réponse HTML
func (c *Context) HTML(html string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/html")
	c.ResponseWriter.Write([]byte(html))
}

// Param récupère un paramètre de route
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// QueryParam récupère un paramètre de query string
func (c *Context) QueryParam(key string) string {
	return c.Query[key]
}

// BindJSON parse le body JSON dans une structure
func (c *Context) BindJSON(v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// Error envoie une réponse d'erreur
func (c *Context) Error(code int, message string) {
	c.Status(code).JSON(map[string]interface{}{
		"error":  message,
		"status": code,
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Middlewares prédéfinis

// Logger middleware pour logger les requêtes
func Logger() MiddlewareFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("%s %s - %v",
			c.Request.Method,
			c.Request.URL.Path,
			duration,
		)
	}
}

// Recovery middleware pour récupérer les panics
func Recovery() MiddlewareFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC: %v", r)
				c.Error(500, "Internal Server Error")
			}
		}()
		c.Next()
	}
}

// CORS middleware pour gérer les CORS avec support des variables d'environnement
func CORS() MiddlewareFunc {
	return CORSWithConfig(CORSConfig{
		AllowOrigins: getCORSOriginsFromEnv(),
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	})
}

// getCORSOriginsFromEnv récupère les origines CORS depuis les variables d'environnement
func getCORSOriginsFromEnv() []string {
	// Essayer CORS_ALLOWED_ORIGINS d'abord
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins == "" {
		// Fallback vers ALLOWED_ORIGINS
		corsOrigins = os.Getenv("ALLOWED_ORIGINS")
	}

	if corsOrigins == "" {
		// Valeur par défaut si aucune variable d'environnement n'est définie
		return []string{"*"}
	}

	// Séparer les origines par des virgules et nettoyer les espaces
	origins := strings.Split(corsOrigins, ",")
	var cleanOrigins []string
	for _, origin := range origins {
		cleaned := strings.TrimSpace(origin)
		if cleaned != "" {
			cleanOrigins = append(cleanOrigins, cleaned)
		}
	}

	return cleanOrigins
}

// CORSConfig configuration pour CORS
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

// CORSFromEnv crée un middleware CORS configuré depuis les variables d'environnement
// Variables supportées:
// - CORS_ALLOWED_ORIGINS ou ALLOWED_ORIGINS: origines autorisées (séparées par des virgules)
// - CORS_ALLOWED_METHODS: méthodes autorisées (séparées par des virgules)
// - CORS_ALLOWED_HEADERS: headers autorisés (séparés par des virgules)
func CORSFromEnv() MiddlewareFunc {
	config := CORSConfig{
		AllowOrigins: getCORSOriginsFromEnv(),
		AllowMethods: getCORSMethodsFromEnv(),
		AllowHeaders: getCORSHeadersFromEnv(),
	}
	return CORSWithConfig(config)
}

// getCORSMethodsFromEnv récupère les méthodes CORS depuis les variables d'environnement
func getCORSMethodsFromEnv() []string {
	corsMethods := os.Getenv("CORS_ALLOWED_METHODS")
	if corsMethods == "" {
		// Valeurs par défaut
		return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	}

	methods := strings.Split(corsMethods, ",")
	var cleanMethods []string
	for _, method := range methods {
		cleaned := strings.TrimSpace(strings.ToUpper(method))
		if cleaned != "" {
			cleanMethods = append(cleanMethods, cleaned)
		}
	}

	return cleanMethods
}

// getCORSHeadersFromEnv récupère les headers CORS depuis les variables d'environnement
func getCORSHeadersFromEnv() []string {
	corsHeaders := os.Getenv("CORS_ALLOWED_HEADERS")
	if corsHeaders == "" {
		// Valeurs par défaut
		return []string{"Content-Type", "Authorization"}
	}

	headers := strings.Split(corsHeaders, ",")
	var cleanHeaders []string
	for _, header := range headers {
		cleaned := strings.TrimSpace(header)
		if cleaned != "" {
			cleanHeaders = append(cleanHeaders, cleaned)
		}
	}

	return cleanHeaders
}

// CORSWithConfig CORS avec configuration personnalisée
func CORSWithConfig(config CORSConfig) MiddlewareFunc {
	return func(c *Context) {
		origin := c.Request.Header.Get("Origin")

		// Vérifier si l'origine est autorisée
		allowed := false
		for _, allowedOrigin := range config.AllowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				c.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
			}
		}

		c.ResponseWriter.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		c.ResponseWriter.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		c.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}

		c.Next()
	}
}

// parseQuery parse la query string
func parseQuery(rawQuery string) map[string]string {
	query := make(map[string]string)
	if rawQuery == "" {
		return query
	}

	pairs := strings.Split(rawQuery, "&")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			query[parts[0]] = parts[1]
		}
	}
	return query
}

// Helpers pour debug et informations

// Routes retourne toutes les routes enregistrées
func (r *Router) Routes() []Route {
	// Trier les routes par méthode puis par chemin
	routes := make([]Route, len(r.routes))
	copy(routes, r.routes)

	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Method != routes[j].Method {
			return routes[i].Method < routes[j].Method
		}
		return routes[i].Path < routes[j].Path
	})

	return routes
}

// PrintRoutes affiche toutes les routes enregistrées
func (r *Router) PrintRoutes() {
	routes := r.Routes()
	fmt.Println("\n🗺️  Routes enregistrées:")
	fmt.Println("========================")

	for _, route := range routes {
		fmt.Printf("%-7s %s\n", route.Method, route.Path)
	}
	fmt.Println()
}
