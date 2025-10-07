package gofsen

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Error("New() should return a router instance")
	}
}

func TestGET(t *testing.T) {
	app := New()

	app.GET("/test", func(c *Context) {
		c.JSON(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "test" {
		t.Errorf("Expected message 'test', got '%s'", response["message"])
	}
}

func TestPOST(t *testing.T) {
	app := New()

	app.POST("/users", func(c *Context) {
		var user map[string]string
		c.BindJSON(&user)
		c.Status(201).JSON(user)
	})

	body := bytes.NewBufferString(`{"name": "Alice"}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestParams(t *testing.T) {
	app := New()

	app.GET("/users/:id", func(c *Context) {
		id := c.Param("id")
		c.JSON(map[string]string{"id": id})
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["id"] != "123" {
		t.Errorf("Expected id '123', got '%s'", response["id"])
	}
}

func TestQueryParams(t *testing.T) {
	app := New()

	app.GET("/search", func(c *Context) {
		query := c.QueryParam("q")
		c.JSON(map[string]string{"query": query})
	})

	req := httptest.NewRequest("GET", "/search?q=gofsen", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["query"] != "gofsen" {
		t.Errorf("Expected query 'gofsen', got '%s'", response["query"])
	}
}

func TestMiddleware(t *testing.T) {
	app := New()

	middlewareCalled := false
	app.Use(func(c *Context) {
		middlewareCalled = true
		c.Next()
	})

	app.GET("/test", func(c *Context) {
		c.JSON(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if !middlewareCalled {
		t.Error("Middleware should have been called")
	}
}

func TestGroup(t *testing.T) {
	app := New()

	api := app.Group("/api/v1")
	api.GET("/users", func(c *Context) {
		c.JSON(map[string]string{"message": "users"})
	})

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotFound(t *testing.T) {
	app := New()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestCORS(t *testing.T) {
	app := New()
	app.Use(CORS())

	app.GET("/test", func(c *Context) {
		c.JSON(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	corsHeader := w.Header().Get("Access-Control-Allow-Origin")
	if corsHeader != "http://localhost:3000" {
		t.Errorf("Expected CORS header 'http://localhost:3000', got '%s'", corsHeader)
	}
}

func TestLogger(t *testing.T) {
	app := New()
	app.Use(Logger())

	app.GET("/test", func(c *Context) {
		c.JSON(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRecovery(t *testing.T) {
	app := New()
	app.Use(Recovery())

	app.GET("/panic", func(c *Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestCORSFromEnv(t *testing.T) {
	// Sauvegarder les variables d'environnement existantes
	originalOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	originalMethods := os.Getenv("CORS_ALLOWED_METHODS")
	originalHeaders := os.Getenv("CORS_ALLOWED_HEADERS")

	// Nettoyer après le test
	defer func() {
		os.Setenv("CORS_ALLOWED_ORIGINS", originalOrigins)
		os.Setenv("CORS_ALLOWED_METHODS", originalMethods)
		os.Setenv("CORS_ALLOWED_HEADERS", originalHeaders)
	}()

	// Test avec variables d'environnement personnalisées
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,https://myapp.com")
	os.Setenv("CORS_ALLOWED_METHODS", "GET,POST")
	os.Setenv("CORS_ALLOWED_HEADERS", "Content-Type,X-Custom-Header")

	app := New()
	app.Use(CORSFromEnv())

	app.GET("/test", func(c *Context) {
		c.JSON(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "http://localhost:3000" {
		t.Errorf("Expected CORS origin 'http://localhost:3000', got '%s'", corsOrigin)
	}

	corsMethods := w.Header().Get("Access-Control-Allow-Methods")
	if corsMethods != "GET, POST" {
		t.Errorf("Expected CORS methods 'GET, POST', got '%s'", corsMethods)
	}

	corsHeaders := w.Header().Get("Access-Control-Allow-Headers")
	if corsHeaders != "Content-Type, X-Custom-Header" {
		t.Errorf("Expected CORS headers 'Content-Type, X-Custom-Header', got '%s'", corsHeaders)
	}
}

func TestCORSFromEnvWithFallback(t *testing.T) {
	// Sauvegarder les variables d'environnement existantes
	originalCorsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	originalAllowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// Nettoyer après le test
	defer func() {
		os.Setenv("CORS_ALLOWED_ORIGINS", originalCorsOrigins)
		os.Setenv("ALLOWED_ORIGINS", originalAllowedOrigins)
	}()

	// Supprimer CORS_ALLOWED_ORIGINS et utiliser ALLOWED_ORIGINS
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com")

	app := New()
	app.Use(CORSFromEnv())

	app.GET("/test", func(c *Context) {
		c.JSON(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "https://example.com" {
		t.Errorf("Expected CORS origin 'https://example.com', got '%s'", corsOrigin)
	}
}

func TestCORSFromEnvWithDefaults(t *testing.T) {
	// Sauvegarder les variables d'environnement existantes
	originalCorsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	originalAllowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// Nettoyer après le test
	defer func() {
		os.Setenv("CORS_ALLOWED_ORIGINS", originalCorsOrigins)
		os.Setenv("ALLOWED_ORIGINS", originalAllowedOrigins)
	}()

	// Supprimer toutes les variables d'environnement CORS
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	os.Unsetenv("ALLOWED_ORIGINS")

	app := New()
	app.Use(CORSFromEnv())

	app.GET("/test", func(c *Context) {
		c.JSON(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	// Avec "*" par défaut, toutes les origines sont autorisées
	corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "http://localhost:3000" {
		t.Errorf("Expected CORS origin 'http://localhost:3000', got '%s'", corsOrigin)
	}
}
