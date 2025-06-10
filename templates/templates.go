package templates

type TemplateGenerator struct{}

func (tg *TemplateGenerator) GetMainTemplate(appType, framework string) string {
	switch appType {
	case "cli":
		return tg.GetCliMainTemplate()
	case "cobra":
		return tg.GetCobraMainTemplate()
	case "web":
		return tg.GetWebMainTemplate(framework)
	}
	return ""
}

func (tg *TemplateGenerator) GetCliMainTemplate() string {
	return `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello CLI Application!")
	
	if len(os.Args) > 1 {
		fmt.Printf("Arguments: %v\n", os.Args[1:])
	}
}`
}

func (tg *TemplateGenerator) GetCobraMainTemplate() string {
	return `package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your application",
	Long:  "A longer description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello Cobra Application!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
}

func main() {
	Execute()
}`
}

func (tg *TemplateGenerator) GetWebMainTemplate(framework string) string {
	switch framework {
	case "fiber":
		return tg.GetFiberMainTemplate()
	case "gin":
		return tg.GetGinMainTemplate()
	default:
		return tg.GetEchoMainTemplate()
	}
}

func (tg *TemplateGenerator) GetFiberMainTemplate() string {
	return `package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Go Fiber App",
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello Fiber!",
			"status":  "success",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	log.Println("🚀 Server starting on :8080")
	log.Fatal(app.Listen(":8080"))
}`
}

func (tg *TemplateGenerator) GetGinMainTemplate() string {
	return `package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Gin!",
			"status":  "success",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.Run(":8080")
}`
}

func (tg *TemplateGenerator) GetEchoMainTemplate() string {
	return `package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello Echo!",
			"status":  "success",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	e.Logger.Info(" Server starting on :8080")
	e.Start(":8080")
}`
}

func (tg *TemplateGenerator) GetDockerTemplate() string {
	return `FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]`
}

func (tg *TemplateGenerator) GetMakefileTemplate() string {
	return `.PHONY: all build run test clean fmt vet lint tidy docker-build docker-run help

APP_NAME ?= server
DOCKER_IMAGE ?= $(APP_NAME):latest
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: fmt vet lint test build

build:
	@echo "🔨 Building application..."
	@go build -o bin/$(APP_NAME) ./cmd

run:
	@echo "🚀 Running application..."
	@go run ./cmd

test:
	@echo "🧪 Running tests..."
	@go test -v ./...

test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

fmt:
	@echo "📝 Formatting code..."
	@gofmt -s -w $(GO_FILES)

vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

tidy:
	@echo "📦 Tidying go modules..."
	@go mod tidy

docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE)

help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//'`
}

func (tg *TemplateGenerator) GetGitignoreTemplate() string {
	return `*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
*.test
*.out
coverage.html
go.work
vendor/
go.sum
.vscode/
.idea/
*.swp
*.swo
*~
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
*.log
air.log
.env
.env.local
.env.*.local
tmp/`
}
