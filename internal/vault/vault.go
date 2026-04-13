package vault

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/zarar/vaultfs/defaults"
)

// Config represents the vault configuration.
type Config struct {
	Vault   VaultConfig          `yaml:"vault"`
	Presets map[string]Preset    `yaml:"presets"`
}

type VaultConfig struct {
	Path string `yaml:"path"`
}

type Preset struct {
	Directories []string `yaml:"directories"`
}

// Discover resolves the vault path using the priority chain:
// 1. Explicit flag value
// 2. VAULTFS_PATH environment variable
// 3. Walk up from cwd looking for .vaultfs/
// 4. Fall back to ~/vault-fs
func Discover(flagValue string, cwd string) (string, error) {
	if flagValue != "" {
		return filepath.Clean(flagValue), nil
	}

	if envPath := os.Getenv("VAULTFS_PATH"); envPath != "" {
		return filepath.Clean(envPath), nil
	}

	if cwd != "" {
		if found := walkUp(cwd); found != "" {
			return found, nil
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "vault-fs"), nil
}

// walkUp searches for .vaultfs/ directory from dir upward.
func walkUp(dir string) string {
	dir = filepath.Clean(dir)
	for {
		candidate := filepath.Join(dir, ".vaultfs")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// LoadConfig loads the vault config from .vaultfs/config.yaml, falling back
// to embedded defaults. User presets are merged on top of defaults.
func LoadConfig(vaultPath string) (*Config, error) {
	defaultCfg, err := loadDefaultConfig()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(vaultPath, ".vaultfs", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultCfg, nil
		}
		return nil, err
	}

	var userCfg Config
	if err := yaml.Unmarshal(data, &userCfg); err != nil {
		return nil, err
	}

	return mergeConfigs(defaultCfg, &userCfg), nil
}

// LoadDefaultConfig loads the embedded default configuration.
func LoadDefaultConfig() (*Config, error) {
	return loadDefaultConfig()
}

// MarshalConfig serializes a config to YAML bytes.
func MarshalConfig(cfg *Config) ([]byte, error) {
	return yaml.Marshal(cfg)
}

func loadDefaultConfig() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(defaults.ConfigYAML, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// mergeConfigs merges user config on top of defaults.
// User presets are added; default presets are preserved if not overridden.
func mergeConfigs(base, user *Config) *Config {
	result := &Config{
		Vault:   base.Vault,
		Presets: make(map[string]Preset),
	}

	// Copy default presets
	for k, v := range base.Presets {
		result.Presets[k] = v
	}

	// Override with user values
	if user.Vault.Path != "" {
		result.Vault = user.Vault
	}
	for k, v := range user.Presets {
		result.Presets[k] = v
	}

	return result
}
