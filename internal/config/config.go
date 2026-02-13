package config

import (
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

const configFileName = ".clo.yaml"

type Config struct {
	Commands []Command `yaml:"commands"`
}

type Command struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Command     string `yaml:"command"`
	Category    string `yaml:"category"`
}

// Load reads ~/.clo.yaml. Creates a default file if it doesn't exist.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, configFileName)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefault(path)
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the command list to ~/.clo.yaml, replacing the file contents.
func Save(commands []Command) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, configFileName)

	cfg := &Config{Commands: commands}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Categories returns sorted unique category names from the command list.
func (c *Config) Categories() []string {
	seen := make(map[string]struct{})
	for _, cmd := range c.Commands {
		if cmd.Category != "" {
			seen[cmd.Category] = struct{}{}
		}
	}
	cats := make([]string, 0, len(seen))
	for cat := range seen {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	return cats
}

func createDefault(path string) (*Config, error) {
	cfg := &Config{
		Commands: []Command{
			{Name: "flush-dns", Description: "Flush macOS DNS cache", Command: "sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder", Category: "network"},
			{Name: "docker-prune", Description: "Remove all stopped containers, dangling images, and unused networks", Command: "docker system prune -af", Category: "docker"},
			{Name: "git-undo", Description: "Undo last commit keeping changes staged", Command: "git reset --soft HEAD~1", Category: "git"},
		},
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, err
	}
	return cfg, nil
}
