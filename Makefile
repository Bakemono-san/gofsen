# 🚀 Makefile pour Gofsen Framework

.PHONY: help run test examples clean build install

# Variables
APP_NAME=gofsen
BUILD_DIR=build
EXAMPLES_DIR=examples

# Couleurs pour l'affichage
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Affiche cette aide
	@echo "$(GREEN)🚀 Gofsen Framework - Commandes disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(GREEN)📋 Exemples d'utilisation:$(NC)"
	@echo "  make run          # Démarre le serveur principal"
	@echo "  make test         # Lance tous les tests"
	@echo "  make examples     # Lance les exemples"
	@echo "  make standalone   # Lance l'exemple standalone"

run: ## Démarre le serveur principal Gofsen
	@echo "$(GREEN)🚀 Démarrage du serveur Gofsen...$(NC)"
	@go run cmd/main.go

test: ## Lance la suite de tests complète
	@echo "$(GREEN)🧪 Lancement des tests Gofsen...$(NC)"
	@chmod +x test-suite.sh
	@./test-suite.sh

build: ## Compile l'application
	@echo "$(GREEN)🔨 Compilation de Gofsen...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go
	@echo "$(GREEN)✅ Compilé dans $(BUILD_DIR)/$(APP_NAME)$(NC)"

standalone: ## Lance l'exemple standalone (un seul fichier)
	@echo "$(GREEN)🚀 Démarrage de l'exemple standalone...$(NC)"
	@go run standalone-example.go

example-hello: ## Lance l'exemple Hello World
	@echo "$(GREEN)👋 Démarrage de l'exemple Hello World...$(NC)"
	@go run $(EXAMPLES_DIR)/01-hello-world.go

example-api: ## Lance l'exemple API REST
	@echo "$(GREEN)🌐 Démarrage de l'exemple API REST...$(NC)"
	@go run $(EXAMPLES_DIR)/02-rest-api.go

example-middleware: ## Lance l'exemple Middleware Demo
	@echo "$(GREEN)🔧 Démarrage de l'exemple Middleware...$(NC)"
	@go run $(EXAMPLES_DIR)/03-middleware-demo.go

examples: ## Lance tous les exemples (ports différents)
	@echo "$(GREEN)📚 Tous les exemples Gofsen:$(NC)"
	@echo ""
	@echo "$(YELLOW)1. Hello World (port 3000)$(NC)"
	@go run $(EXAMPLES_DIR)/01-hello-world.go &
	@sleep 2
	@echo ""
	@echo "$(YELLOW)2. API REST (port 3001)$(NC)"
	@go run $(EXAMPLES_DIR)/02-rest-api.go &
	@sleep 2
	@echo ""
	@echo "$(YELLOW)3. Middleware Demo (port 3002)$(NC)"
	@go run $(EXAMPLES_DIR)/03-middleware-demo.go &
	@sleep 2
	@echo ""
	@echo "$(GREEN)✅ Tous les exemples sont démarrés!$(NC)"
	@echo "$(GREEN)📋 URLs disponibles:$(NC)"
	@echo "  http://localhost:3000 - Hello World"
	@echo "  http://localhost:3001 - API REST"
	@echo "  http://localhost:3002 - Middleware Demo"
	@echo ""
	@echo "$(YELLOW)Appuyez sur Ctrl+C pour arrêter tous les serveurs$(NC)"
	@wait

quick-test: ## Test rapide des endpoints principaux
	@echo "$(GREEN)⚡ Test rapide des endpoints...$(NC)"
	@curl -s http://localhost:8081/health || echo "$(RED)❌ Serveur non démarré$(NC)"
	@curl -s http://localhost:8081/test/all | head -5 || echo "$(RED)❌ Tests non disponibles$(NC)"

docs: ## Ouvre la documentation
	@echo "$(GREEN)📚 Ouverture de la documentation...$(NC)"
	@echo "$(GREEN)README.md:$(NC) Documentation principale"
	@echo "$(GREEN)INSTALL.md:$(NC) Guide d'installation"
	@echo "$(GREEN)TEST_ROUTES.md:$(NC) Documentation des tests"
	@echo "$(GREEN)CHECKLIST.md:$(NC) Progression du développement"

clean: ## Nettoie les fichiers de build
	@echo "$(GREEN)🧹 Nettoyage...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "$(GREEN)✅ Nettoyage terminé$(NC)"

install: ## Installe les dépendances
	@echo "$(GREEN)📦 Installation des dépendances...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dépendances installées$(NC)"

format: ## Formate le code Go
	@echo "$(GREEN)✨ Formatage du code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formaté$(NC)"

vet: ## Vérifie le code avec go vet
	@echo "$(GREEN)🔍 Vérification du code...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Code vérifié$(NC)"

check: format vet ## Formate et vérifie le code
	@echo "$(GREEN)✅ Code formaté et vérifié$(NC)"

dev: ## Mode développement (redémarre automatiquement)
	@echo "$(GREEN)🔄 Mode développement (Ctrl+C pour arrêter)$(NC)"
	@while true; do \
		go run cmd/main.go & \
		PID=$$!; \
		inotifywait -e modify,create,delete -r . --exclude='\.git|build'; \
		kill $$PID 2>/dev/null; \
		sleep 1; \
	done

serve-docs: ## Démarre un serveur pour la documentation
	@echo "$(GREEN)📖 Serveur de documentation sur http://localhost:8080$(NC)"
	@python3 -m http.server 8080 2>/dev/null || python -m SimpleHTTPServer 8080

package: build ## Crée un package de distribution
	@echo "$(GREEN)📦 Création du package...$(NC)"
	@tar -czf $(BUILD_DIR)/$(APP_NAME)-$(shell date +%Y%m%d).tar.gz \
		README.md INSTALL.md CHECKLIST.md TEST_ROUTES.md \
		examples/ internal/ cmd/ test-suite.sh standalone-example.go \
		go.mod Makefile
	@echo "$(GREEN)✅ Package créé: $(BUILD_DIR)/$(APP_NAME)-$(shell date +%Y%m%d).tar.gz$(NC)"

# Cibles par défaut
.DEFAULT_GOAL := help

# Pour éviter les conflits avec des fichiers du même nom
.SUFFIXES: