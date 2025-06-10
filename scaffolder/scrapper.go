package scaffolder

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/gostart/config"
	"github.com/SwanHtetAungPhyo/gostart/spinner"
	"github.com/SwanHtetAungPhyo/gostart/templates"
	"time"

	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

type Scaffolder struct {
	config  *config.Config
	spinner *spinner.Spinner
}

func NewScaffolder(config *config.Config) *Scaffolder {
	return &Scaffolder{
		config:  config,
		spinner: spinner.NewSpinner(),
	}
}
func (s *Scaffolder) createDirectoryStructure() error {
	return os.MkdirAll(s.config.ProjectDir, 0755)
}

func (s *Scaffolder) initializeGoModule() error {
	return s.runCommand("go", "mod", "init", s.config.ModuleName)
}

func (s *Scaffolder) generateMainFile() error {
	cmdDir := filepath.Join(s.config.ProjectDir, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		return err
	}

	generator := templates.TemplateGenerator{}
	mainContent := generator.GetMainTemplate(s.config.AppType, s.config.Framework)
	mainPath := filepath.Join(cmdDir, "main.go")
	return os.WriteFile(mainPath, []byte(mainContent), 0644)
}

func (s *Scaffolder) createInternalStructure() error {
	dirs := []string{
		"internal/model",
		"internal/repository",
		"internal/service",
		"internal/handler",
		"internal/config",
		"pkg/utils",
		"api",
		"docs",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(s.config.ProjectDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
		readmePath := filepath.Join(fullPath, ".gitkeep")
		err := os.WriteFile(readmePath, []byte(""), 0644)
		if err != nil {
			color.Red(err.Error())
			return err
		}
	}

	return nil
}

func (s *Scaffolder) generateDockerfile() error {
	dockerPath := filepath.Join(s.config.ProjectDir, "Dockerfile")
	generator := templates.TemplateGenerator{}
	content := generator.GetDockerTemplate()
	return os.WriteFile(dockerPath, []byte(content), 0644)
}

func (s *Scaffolder) setupAir() error {
	if _, err := exec.LookPath("air"); err != nil {
		color.Yellow("⚠️  Air is not installed.")
		if s.askToInstallAir() {
			if err := s.installAir(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("air installation skipped")
		}
	}
	return s.generateAirConfig()
}

func (s *Scaffolder) askToInstallAir() bool {
	sel := promptui.Select{
		Label: "Install air for hot reload? (requires Go 1.16+)",
		Items: []string{"yes", "no"},
	}
	_, result, _ := sel.Run()
	return result == "yes"
}

func (s *Scaffolder) installAir() error {
	color.Cyan("📦 Installing air...")
	cmd := exec.Command("go", "install", "github.com/air-verse/air@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install air: %w", err)
	}

	color.Green("✅ Air installed successfully")
	color.Yellow("⚠️  Make sure $GOPATH/bin is in your PATH")
	return nil
}

func (s *Scaffolder) generateAirConfig() error {
	cmd := exec.Command("air", "init")
	cmd.Dir = s.config.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate air config: %w", err)
	}

	color.Green("✅ Air configuration generated")
	return nil
}

func (s *Scaffolder) generateMakefile() error {
	makefilePath := filepath.Join(s.config.ProjectDir, "Makefile")
	generator := templates.TemplateGenerator{}
	content := generator.GetMakefileTemplate()
	return os.WriteFile(makefilePath, []byte(content), 0644)
}

func (s *Scaffolder) generateGitignore() error {
	gitignorePath := filepath.Join(s.config.ProjectDir, ".gitignore")
	generator := templates.TemplateGenerator{}
	content := generator.GetGitignoreTemplate()
	return os.WriteFile(gitignorePath, []byte(content), 0644)
}

func (s *Scaffolder) tidyGoMod() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = s.config.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to tidy go.mod: %w", err)
	}

	color.Green("✅ go.mod tidied successfully")
	return nil
}

func (s *Scaffolder) CreateProject() error {
	color.Cyan("\n🚀 Scaffolding %s...", s.config.ModuleName)
	s.spinner.Start()
	defer s.spinner.Stop()

	steps := []struct {
		name string
		fn   func() error
	}{
		{"creating directory structure", s.createDirectoryStructure},
		{"initializing Go module", s.initializeGoModule},
		{"generating main file", s.generateMainFile},
		{"creating internal structure", s.createInternalStructure},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("failed %s: %w", step.name, err)
		}
	}

	optionalSteps := []struct {
		enabled bool
		name    string
		fn      func() error
	}{
		{s.config.UseDocker, "generating Dockerfile", s.generateDockerfile},
		{s.config.UseAir, "setting up Air", s.setupAir},
		{s.config.UseMakefile, "generating Makefile", s.generateMakefile},
	}

	for _, step := range optionalSteps {
		if step.enabled {
			if err := step.fn(); err != nil {
				color.Yellow("⚠️  %s failed: %v", step.name, err)
			}
		}
	}

	fileCreators := []struct {
		name string
		fn   func() error
	}{
		{"generating .gitignore", s.generateGitignore},
		{"creating .env file", s.envFileCreation},
		{"creating README.md", s.readmeFileCreation},
		{"tidying go.mod", s.tidyGoMod},
	}

	for _, creator := range fileCreators {
		if err := creator.fn(); err != nil {
			color.Yellow("⚠️  %s failed: %v", creator.name, err)
		}
	}
	if err := s.installDependencies(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

func (s *Scaffolder) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = s.config.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *Scaffolder) envFileCreation() error {
	envPath := filepath.Join(s.config.ProjectDir, ".env")
	content := `# Environment variables
# Add your environment variables here
`
	return os.WriteFile(envPath, []byte(content), 0644)
}

func (s *Scaffolder) readmeFileCreation() error {
	readmePath := filepath.Join(s.config.ProjectDir, "README.md")
	content := fmt.Sprintf(`# %s
This is a Go project scaffolded by the Go Scaffolder (GPM) 


## Project Structure
- **cmd/**: Contains the main application entry point.
- **internal/**: Contains internal packages.
- **pkg/**: Contains reusable packages.
- **api/**: Contains API definitions.
- **docs/**: Contains documentation scaffolder.
- **Makefile**: Contains build and run commands.
- **Dockerfile**: Contains Docker build instructions.
- **.gitignore**: Specifies scaffolder to ignore in Git.
`, s.config.ModuleName)
	return os.WriteFile(readmePath, []byte(content), 0644)
}

func (s *Scaffolder) installDependencies() error {
	if len(s.config.SelectedDependencies) == 0 {
		return nil
	}

	color.Cyan("\n📦 Installing selected dependencies...")
	s.spinner.Start()
	defer s.spinner.Stop()
	successCount := 0
	failedDeps := make([]string, 0)

	for _, dep := range s.config.SelectedDependencies {
		if err := s.installSingleDependency(dep); err != nil {
			color.Yellow("⚠️  Failed to install %s: %v", dep, err)
			failedDeps = append(failedDeps, dep)
		} else {
			color.Green("✓ Installed %s", dep)
			successCount++
		}
	}

	totalDeps := len(s.config.SelectedDependencies)
	if successCount == totalDeps {
		color.Green("\n✅ Successfully installed all %d dependencies", successCount)
	} else {
		color.Yellow("\n✅ Successfully installed %d/%d dependencies", successCount, totalDeps)
		if len(failedDeps) > 0 {
			color.Red("Failed dependencies: %v", failedDeps)
			color.Yellow("💡 You can manually install failed dependencies later with: go get <package>")
		}
	}

	return nil
}
func (s *Scaffolder) installSingleDependency(dep string) error {
	time.Sleep(100 * time.Millisecond)

	cmd := exec.Command("go", "get", dep)
	cmd.Dir = s.config.ProjectDir
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		time.Sleep(500 * time.Millisecond)
		retryCmd := exec.Command("go", "get", dep)
		retryCmd.Dir = s.config.ProjectDir
		retryCmd.Stdout = nil
		retryCmd.Stderr = nil

		if retryErr := retryCmd.Run(); retryErr != nil {
			return fmt.Errorf("failed after retry: %w", retryErr)
		}
	}

	return nil
}
