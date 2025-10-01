# 🧪 Gofsen Framework Test Routes

## Quick Start

1. **Start the server:**
   ```bash
   go run cmd/main.go
   ```

2. **Test all functionalities:**
   ```bash
   ./test-suite.sh
   ```

3. **Or visit test overview:**
   ```
   http://localhost:8081/test/all
   ```

---

## 📋 Available Test Routes

### 🧱 Basic Routing Methods
- `GET /test/get` - Test GET method
- `POST /test/post` - Test POST method  
- `PUT /test/put` - Test PUT method
- `DELETE /test/delete` - Test DELETE method
- `PATCH /test/patch` - Test PATCH method

### 🧭 Route Groups & Middleware
- `GET /test/group/basic` - Test route groups with local middleware

### 🔐 Middleware & Security
- `GET /test/logger` - Test logger middleware (check console)
- `GET /test/auth/protected` - Test auth middleware (requires token)
- `GET /test/auth/public` - Test public route (no auth)
- `GET /test/recovery` - Test panic recovery middleware
- `GET /test/cors/check` - Test CORS middleware

### ⚙️ Helpers & I/O
- `GET /test/json` - Test JSON response
- `GET /test/query?name=John&age=25` - Test query parameters
- `POST /test/bind` - Test JSON body parsing
- `GET /test/error?type=400` - Test error responses

### 🧪 Advanced Combined Tests
- `POST /test/multi/combined` - Test multiple features together

---

## 🔑 Authentication Testing

For protected routes, use this header:
```
Authorization: Bearer valid-token
```

**Example:**
```bash
curl -H "Authorization: Bearer valid-token" http://localhost:8081/test/auth/protected
```

---

## 🌐 CORS Testing

For CORS testing, add origin header:
```
Origin: https://example.com
```

**Example:**
```bash
curl -H "Origin: https://example.com" http://localhost:8081/test/cors/check
```

---

## 📤 JSON Body Testing

For POST/PUT requests with JSON:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","age":30}' \
  http://localhost:8081/test/bind
```

---

## 🎯 Complete Test Examples

### Test all basic methods:
```bash
curl http://localhost:8081/test/get
curl -X POST -H "Content-Type: application/json" -d '{"test":"data"}' http://localhost:8081/test/post
curl -X PUT http://localhost:8081/test/put
curl -X DELETE http://localhost:8081/test/delete
curl -X PATCH http://localhost:8081/test/patch
```

### Test middleware:
```bash
# Auth (should fail)
curl http://localhost:8081/test/auth/protected

# Auth (should succeed)
curl -H "Authorization: Bearer valid-token" http://localhost:8081/test/auth/protected

# CORS
curl -H "Origin: https://example.com" http://localhost:8081/test/cors/check

# Recovery (should return error instead of crashing)
curl http://localhost:8081/test/recovery
```

### Test I/O helpers:
```bash
# Query params
curl "http://localhost:8081/test/query?name=John&age=25&city=Paris"

# JSON binding
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","email":"bob@test.com","age":25}' \
  http://localhost:8081/test/bind

# Error responses
curl "http://localhost:8081/test/error?type=400"
curl "http://localhost:8081/test/error?type=404"
curl "http://localhost:8081/test/error"
```

---

## 🚀 Automated Testing

Run the complete test suite:
```bash
./test-suite.sh
```

This will test all functionalities and show you a complete report of what's working.

---

## ✅ Expected Results

All routes should return `"test_passed": true` in their JSON response when working correctly. The test suite script will show ✅ for passing tests and ❌ for failing ones.