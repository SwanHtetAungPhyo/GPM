package wizzard

import (
	"errors"
	"fmt"
	"github.com/SwanHtetAungPhyo/gostart/config"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Dependency struct {
	Name       string
	URL        string
	ImportPath string
}

func showColorfulBanner() {
	lines := []string{
		"  ██████  ██████  ███    ███ ",
		" ██       ██   ██ ████  ████ ",
		" ██   ███ ██████  ██ ████ ██ ",
		" ██    ██ ██      ██  ██  ██ ",
		"  ██████  ██      ██      ██ ",
	}

	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgYellow),
		color.New(color.FgGreen),
		color.New(color.FgCyan),
		color.New(color.FgMagenta),
	}

	println()
	for i, line := range lines {
		colors[i%len(colors)].Println(line)
	}

	color.New(color.FgBlue, color.Bold).Println("\n  GPM - Go Package Manager & Scaffolder")
	color.New(color.FgMagenta).Println("   Built by Swan Htet Aung Phyo")
	color.New(color.FgGreen).Println("   Making Go development workflow faster ⚡")
	color.Yellow("─────────────────────────────────────────────")
	color.Red("Reminder: Please make sure u do not have the same named folder as the module path you declared \n . I let the generator to make the folder with same name in the go mod path name. ")
}

type Wizard struct{}

func NewWizard() *Wizard {
	return &Wizard{}
}

func (w *Wizard) Run() (*config.Config, error) {
	configuartion := &config.Config{}

	if err := w.getModuleName(configuartion); err != nil {
		return nil, err
	}

	if err := w.getAppType(configuartion); err != nil {
		return nil, err
	}

	if configuartion.AppType == "web" {
		if err := w.getFramework(configuartion); err != nil {
			return nil, err
		}
	}

	configuartion.UseDocker = w.yesNo("Will you use Docker?")
	configuartion.UseAir = w.yesNo("Include air.toml (hot reload)?")
	configuartion.UseMakefile = w.yesNo("Include Makefile?")

	configuartion.ProjectDir = filepath.Base(configuartion.ModuleName)
	if err := w.getDependencies(configuartion); err != nil {
		return nil, err
	}

	return configuartion, nil
}

func (w *Wizard) getModuleName(config *config.Config) error {
	showColorfulBanner()
	prompt := promptui.Prompt{
		Label:   "What is the name of the Go module?",
		Default: "github.com/yourname/project",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("module name cannot be empty")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}
	config.ModuleName = strings.TrimSpace(result)
	return nil
}

func (w *Wizard) getAppType(config *config.Config) error {
	appTypes := []string{"cli", "cobra", "web"}
	sel := promptui.Select{
		Label: "What type of application?",
		Items: appTypes,
	}

	_, result, err := sel.Run()
	if err != nil {
		return err
	}
	config.AppType = result
	return nil
}

func (w *Wizard) getFramework(config *config.Config) error {
	frameworks := []string{"fiber", "gin", "echo"}
	sel := promptui.Select{
		Label: "Choose web framework",
		Items: frameworks,
	}

	_, result, err := sel.Run()
	if err != nil {
		return err
	}
	config.Framework = result
	return nil
}

func (w *Wizard) yesNo(label string) bool {
	sel := promptui.Select{
		Label: label,
		Items: []string{"yes", "no"},
	}
	_, result, _ := sel.Run()
	return result == "yes"
}

func (w *Wizard) getDependencies(config *config.Config) error {
	if !w.yesNo("Would you like to add popular Go dependencies?") {
		return nil
	}

	categories := w.getCuratedCategories()

	for {
		selectedCategory, err := w.selectCategory(categories)
		if err != nil {
			return err
		}

		dependencies := w.getDependenciesForCategory(selectedCategory)
		selectedDeps, err := w.selectDependencies(dependencies)
		if err != nil {
			return err
		}

		config.SelectedDependencies = append(config.SelectedDependencies, selectedDeps...)

		if !w.yesNo("Browse another category?") {
			break
		}
	}

	return nil
}

func (w *Wizard) getCuratedCategories() map[string][]Dependency {
	return map[string][]Dependency{
		"🌐 Web Frameworks": {
			{Name: "Fiber", ImportPath: "github.com/gofiber/fiber/v2", URL: "https://github.com/gofiber/fiber"},
			{Name: "Gin", ImportPath: "github.com/gin-gonic/gin", URL: "https://github.com/gin-gonic/gin"},
			{Name: "Echo", ImportPath: "github.com/labstack/echo/v4", URL: "https://github.com/labstack/echo"},
			{Name: "Chi Router", ImportPath: "github.com/go-chi/chi/v5", URL: "https://github.com/go-chi/chi"},
			{Name: "Gorilla Mux", ImportPath: "github.com/gorilla/mux", URL: "https://github.com/gorilla/mux"},
		},
		"🗄️ Database & ORM": {
			{Name: "GORM", ImportPath: "gorm.io/gorm", URL: "https://github.com/go-gorm/gorm"},
			{Name: "PostgreSQL Driver", ImportPath: "gorm.io/driver/postgres", URL: "https://github.com/go-gorm/postgres"},
			{Name: "MySQL Driver", ImportPath: "gorm.io/driver/mysql", URL: "https://github.com/go-gorm/mysql"},
			{Name: "SQLite Driver", ImportPath: "gorm.io/driver/sqlite", URL: "https://github.com/go-gorm/sqlite"},
			{Name: "MongoDB Driver", ImportPath: "go.mongodb.org/mongo-driver/mongo", URL: "https://github.com/mongodb/mongo-go-driver"},
			{Name: "Redis", ImportPath: "github.com/redis/go-redis/v9", URL: "https://github.com/redis/go-redis"},
		},
		"⚙️ CLI Tools": {
			{Name: "Cobra", ImportPath: "github.com/spf13/cobra", URL: "https://github.com/spf13/cobra"},
			{Name: "Viper", ImportPath: "github.com/spf13/viper", URL: "https://github.com/spf13/viper"},
			{Name: "PromptUI", ImportPath: "github.com/manifoldco/promptui", URL: "https://github.com/manifoldco/promptui"},
			{Name: "Survey", ImportPath: "github.com/AlecAivazis/survey/v2", URL: "https://github.com/AlecAivazis/survey"},
			{Name: "Color", ImportPath: "github.com/fatih/color", URL: "https://github.com/fatih/color"},
		},
		"📊 Logging & Monitoring": {
			{Name: "Logrus", ImportPath: "github.com/sirupsen/logrus", URL: "https://github.com/sirupsen/logrus"},
			{Name: "Zap", ImportPath: "go.uber.org/zap", URL: "https://github.com/uber-go/zap"},
			{Name: "Zerolog", ImportPath: "github.com/rs/zerolog", URL: "https://github.com/rs/zerolog"},
			{Name: "Prometheus Client", ImportPath: "github.com/prometheus/client_golang", URL: "https://github.com/prometheus/client_golang"},
		},
		"🔐 Authentication & Security": {
			{Name: "JWT Go", ImportPath: "github.com/golang-jwt/jwt/v5", URL: "https://github.com/golang-jwt/jwt"},
			{Name: "OAuth2", ImportPath: "golang.org/x/oauth2", URL: "https://github.com/golang/oauth2"},
			{Name: "Bcrypt", ImportPath: "golang.org/x/crypto/bcrypt", URL: "https://golang.org/x/crypto"},
			{Name: "CORS", ImportPath: "github.com/rs/cors", URL: "https://github.com/rs/cors"},
			{Name: "Rate Limiter", ImportPath: "github.com/ulule/limiter/v3", URL: "https://github.com/ulule/limiter"},
		},
		"🧪 Testing": {
			{Name: "Testify", ImportPath: "github.com/stretchr/testify", URL: "https://github.com/stretchr/testify"},
			{Name: "GoMock", ImportPath: "go.uber.org/mock", URL: "https://github.com/uber-go/mock"},
			{Name: "Ginkgo", ImportPath: "github.com/onsi/ginkgo/v2", URL: "https://github.com/onsi/ginkgo"},
			{Name: "Gomega", ImportPath: "github.com/onsi/gomega", URL: "https://github.com/onsi/gomega"},
		},
		"🌐 HTTP & API": {
			{Name: "Resty", ImportPath: "github.com/go-resty/resty/v2", URL: "https://github.com/go-resty/resty"},
			{Name: "GraphQL", ImportPath: "github.com/99designs/gqlgen", URL: "https://github.com/99designs/gqlgen"},
			{Name: "WebSocket", ImportPath: "github.com/gorilla/websocket", URL: "https://github.com/gorilla/websocket"},
			{Name: "gRPC", ImportPath: "google.golang.org/grpc", URL: "https://github.com/grpc/grpc-go"},
		},
		"📦 Utilities": {
			{Name: "UUID", ImportPath: "github.com/google/uuid", URL: "https://github.com/google/uuid"},
			{Name: "ULID", ImportPath: "github.com/oklog/ulid/v2", URL: "https://github.com/oklog/ulid"},
			{Name: "Validator", ImportPath: "github.com/go-playground/validator/v10", URL: "https://github.com/go-playground/validator"},
			{Name: "GoDotEnv", ImportPath: "github.com/joho/godotenv", URL: "https://github.com/joho/godotenv"},
			{Name: "Copier", ImportPath: "github.com/jinzhu/copier", URL: "https://github.com/jinzhu/copier"},
		},
		"📁 File & Data Processing": {
			{Name: "Excel", ImportPath: "github.com/xuri/excelize/v2", URL: "https://github.com/xuri/excelize"},
			{Name: "CSV", ImportPath: "encoding/csv", URL: "https://pkg.go.dev/encoding/csv"},
			{Name: "YAML", ImportPath: "gopkg.in/yaml.v3", URL: "https://github.com/go-yaml/yaml"},
			{Name: "TOML", ImportPath: "github.com/BurntSushi/toml", URL: "https://github.com/BurntSushi/toml"},
			{Name: "Afero (File System)", ImportPath: "github.com/spf13/afero", URL: "https://github.com/spf13/afero"},
		},
	}
}

func (w *Wizard) selectCategory(categories map[string][]Dependency) (string, error) {
	items := make([]string, 0, len(categories))
	for category := range categories {
		items = append(items, category)
	}

	sel := promptui.Select{
		Label: "Select a category",
		Items: items,
		Size:  15,
	}

	_, result, err := sel.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (w *Wizard) getDependenciesForCategory(category string) []Dependency {
	categories := w.getCuratedCategories()
	return categories[category]
}

func (w *Wizard) selectDependencies(deps []Dependency) ([]string, error) {
	if len(deps) == 0 {
		color.Yellow("No dependencies found in this category")
		return nil, nil
	}

	items := make([]string, len(deps))
	for i, dep := range deps {
		items[i] = fmt.Sprintf("%-20s - %s", dep.Name, dep.ImportPath)
	}

	searcher := func(input string, index int) bool {
		dep := deps[index]
		return strings.Contains(strings.ToLower(dep.Name), strings.ToLower(input)) ||
			strings.Contains(strings.ToLower(dep.ImportPath), strings.ToLower(input))
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "➤ {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "✓ {{ . | green }}",
		Details: `
--------- Dependency Details ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Import Path:" | faint }}	{{ .ImportPath }}
{{ "Repository:" | faint }}	{{ .URL }}`,
	}

	prompt := promptui.Select{
		Label:             "Select dependencies (use arrow keys, type to filter)",
		Items:             items,
		Templates:         templates,
		Size:              15,
		Searcher:          searcher,
		StartInSearchMode: true,
	}

	var selectedDeps []string
	for {
		idx, _, err := prompt.Run()
		if err != nil {
			if errors.Is(err, promptui.ErrInterrupt) {
				break
			}
			return nil, err
		}

		selectedDep := deps[idx]
		selectedDeps = append(selectedDeps, selectedDep.ImportPath)
		color.Green("✓ Added: %s", selectedDep.ImportPath)

		deps = append(deps[:idx], deps[idx+1:]...)
		items = append(items[:idx], items[idx+1:]...)
		prompt.Items = items

		if len(deps) == 0 {
			color.Yellow("No more dependencies in this category")
			break
		}

		if !w.yesNo("Add another dependency from this category?") {
			break
		}
	}

	return selectedDeps, nil
}

func (w *Wizard) installDependencies(config *config.Config, projectDir string) error {
	if len(config.SelectedDependencies) == 0 {
		return nil
	}

	color.Cyan("📦 Installing selected dependencies...")

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(projectDir); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			color.Red("failed to change to project directory: %w", err.Error())
		}
	}(originalDir)

	successCount := 0
	failedDeps := make([]string, 0)

	for _, dep := range config.SelectedDependencies {
		if err := w.installSingleDependency(dep); err != nil {
			color.Red("⚠️  Failed to install %s: %v", dep, err)
			failedDeps = append(failedDeps, dep)
		} else {
			color.Green("✓ Installed %s", dep)
			successCount++
		}
	}

	totalDeps := len(config.SelectedDependencies)
	if successCount == totalDeps {
		color.Green("✅ Successfully installed all %d dependencies", successCount)
	} else {
		color.Yellow("✅ Successfully installed %d/%d dependencies", successCount, totalDeps)
		if len(failedDeps) > 0 {
			color.Red("Failed dependencies: %v", failedDeps)
			color.Yellow("💡 You can manually install failed dependencies later with: go get <package>")
		}
	}

	return nil
}

func (w *Wizard) installSingleDependency(dep string) error {
	time.Sleep(100 * time.Millisecond)

	cmd := exec.Command("go", "get", dep)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		time.Sleep(500 * time.Millisecond)
		retryCmd := exec.Command("go", "get", dep)
		retryCmd.Stdout = nil
		retryCmd.Stderr = nil

		if retryErr := retryCmd.Run(); retryErr != nil {
			return fmt.Errorf("failed after retry: %w", retryErr)
		}
	}

	return nil
}
func (w *Wizard) installDependenciesBatch(config *config.Config, projectDir string) error {
	if len(config.SelectedDependencies) == 0 {
		return nil
	}

	color.Cyan("📦 Installing selected dependencies...")

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(projectDir); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			color.Red("failed to change to project directory: %w", err.Error())
		}
	}(originalDir)

	args := append([]string{"get"}, config.SelectedDependencies...)
	cmd := exec.Command("go", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red("⚠️  Failed to install some dependencies: %v", err)
		color.Yellow("Output: %s", string(output))
		color.Yellow("🔄 Falling back to individual installation...")
		return w.installDependencies(config, ".")
	}

	color.Green("✅ Successfully installed all %d dependencies", len(config.SelectedDependencies))
	return nil
}
