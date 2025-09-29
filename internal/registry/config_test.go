package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestConfigFile(t *testing.T, dir, content string) {
	configPath := filepath.Join(dir, "config.yml")
	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create test config file: %v", err)
	}
}

func TestLoadConfig(t *testing.T) {
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	t.Run("it loads local config by default", func(t *testing.T) {
		tmpDir := t.TempDir()
		localDir := filepath.Join(tmpDir, "config", "local")
		os.MkdirAll(localDir, 0755)
		createTestConfigFile(t, localDir, "server:\n  port: \"8080\"")

		os.Chdir(tmpDir)
		os.Setenv("APP_ENV", "development")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}
		if cfg.Server.Port != "8080" {
			t.Errorf("expected port '8080', but got '%s'", cfg.Server.Port)
		}
	})
}
