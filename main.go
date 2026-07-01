package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"charm.land/huh/v2"
)

//go:embed template/*
var templateFS embed.FS

type ProjectConfig struct {
	ProjectName    string
	EnvPrefix      string
	EnvPrefixLower string
	ModulePath     string
	LogLevel       string
	PGURL          string
	Host           string
	Port           int
	AdminPwd       string
}

func main() {
	cfg, err := GetCfg()
	if err != nil {
		log.Fatal(err)
	}
	err = GenerateProject(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Project files generated succefully. Start initialization")
	err = AfterGenCmds(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Project initialized succefully. Good luck with development!")
}

func GetCfg() (*ProjectConfig, error) {
	var cfg ProjectConfig
	port := "3000"
	cfg.PGURL = "postgres://postgres:postgres@localhost:5432/gotham"
	cfg.Host = "127.0.0.1"
	cfg.AdminPwd = "strongpassword"
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Gotham").
				Description("Create Go + Templ + HTMX + AlpineJS web application"),
			huh.NewInput().Title("Application title").Placeholder("Gotham").Value(&cfg.ProjectName),
			huh.NewInput().Title("Module path").Placeholder("github.com/myGitHub/gotham").Value(&cfg.ModulePath),
			huh.NewSelect[string]().
				Title("Log level").
				Options(
					huh.NewOption("debug", "debug").Selected(true),
					huh.NewOption("info", "info"),
					huh.NewOption("warn", "warn"),
					huh.NewOption("error", "error"),
				).
				Value(&cfg.LogLevel),
		),
		huh.NewGroup(
			huh.NewInput().Title("Postgres URL").Value(&cfg.PGURL),
			huh.NewInput().Title("Application host").Inline(true).Value(&cfg.Host),
			huh.NewInput().
				Title("Application port").
				Inline(true).
				Value(&port).
				Validate(func(input string) error {
					_, err := strconv.Atoi(input)
					return err
				}),
			huh.NewInput().Title("Admin password for pprof").Inline(true).Value(&cfg.AdminPwd),
		),
	)
	if err := form.Run(); err != nil {
		return nil, err
	}
	cfg.Port, _ = strconv.Atoi(port)
	cfg.EnvPrefix = strings.ToUpper(strings.ReplaceAll(cfg.ProjectName, " ", "_"))
	cfg.EnvPrefixLower = strings.ToLower(cfg.EnvPrefix)
	return &cfg, nil
}

func GenerateProject(cfg *ProjectConfig) error {
	if err := os.MkdirAll(cfg.EnvPrefixLower, 0o755); err != nil {
		return err
	}

	return fs.WalkDir(templateFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(path, "template/")
		if relPath == "" || strings.HasPrefix(relPath, "template") {
			return nil
		}
		name := filepath.Base(path)
		if name == "env.template" {
			name = ".env"
		}
		destPath := filepath.Join(cfg.EnvPrefixLower, relPath)
		destDir := filepath.Dir(destPath)
		destPath = filepath.Join(destDir, name)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		content, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		isTextFile := ext == ".go" || ext == ".mod" || ext == ".sum" ||
			ext == ".yml" || ext == ".yaml" || ext == ".json" ||
			ext == ".templ" || ext == ".sql" || ext == ".env" || name == "justfile"

		var data []byte
		if isTextFile {
			if ext == ".templ" || name == "justfile" {
				contentStr := string(content)
				contentStr = strings.ReplaceAll(contentStr, "{{ .ProjectName }}", cfg.ProjectName)
				contentStr = strings.ReplaceAll(contentStr, "{{ .ModulePath }}", cfg.ModulePath)
				contentStr = strings.ReplaceAll(contentStr, "{{ .EnvPrefix }}", cfg.EnvPrefix)
				data = []byte(contentStr)
			} else {
				// Для остальных текстовых файлов используем шаблоны
				tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
				if err != nil {
					return err
				}
				var buf bytes.Buffer
				if err := tmpl.Execute(&buf, cfg); err != nil {
					return err
				}
				data = buf.Bytes()
			}
		} else {
			data = content
		}
		return os.WriteFile(destPath, data, 0o644)
	})
}

func AfterGenCmds(cfg *ProjectConfig) error {
	wd := cfg.EnvPrefixLower
	initCmd := exec.Command("go", "mod", "init", cfg.ModulePath)
	initCmd.Dir = wd
	if output, err := initCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod init failed with error: %w\nOutput: %s", err, string(output))
	}
	fmt.Printf("✅ go mod init %s completed\n", cfg.ModulePath)
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = wd
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed with error: %w\nOutput: %s", err, string(output))
	}
	fmt.Println("✅ go mod tidy completed")
	twCmd := exec.Command("npm", "install")
	twCmd.Dir = wd
	if output, err := twCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("npm install failed with error: %w\nOutput: %s", err, string(output))
	}
	templToolCmd := exec.Command("go", "get", "-tool", "github.com/a-h/templ/cmd/templ@latest")
	templToolCmd.Dir = wd
	if output, err := templToolCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go get -tool github.com/a-h/templ/cmd/templ@latest failed with error: %w\nOutput: %s", err, string(output))
	}
	genTemplCmd := exec.Command("go", "tool", "templ", "generate")
	genTemplCmd.Dir = wd
	if output, err := genTemplCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go tool templ generate failed with error: %w\nOutput: %s", err, string(output))
	}
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = wd
	if output, err := gitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init failed with error: %w\nOutput: %s", err, string(output))
	}
	return nil
}
