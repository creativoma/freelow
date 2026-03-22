package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Config struct {
	Active  string   `json:"active"`
	Clients []Client `json:"clients"`
}

func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".freelow", "clients.json"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func Save(config *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *Config) GetActive() (*Client, error) {
	if c.Active == "" {
		return nil, fmt.Errorf("no active client set")
	}
	for i, client := range c.Clients {
		if client.ID == c.Active {
			return &c.Clients[i], nil
		}
	}
	return nil, fmt.Errorf("active client %q not found", c.Active)
}

func (c *Config) FindByID(id string) (*Client, error) {
	for i, client := range c.Clients {
		if client.ID == id {
			return &c.Clients[i], nil
		}
	}
	return nil, fmt.Errorf("client %q not found", id)
}

// ValidColor returns true if s is an ANSI number (0-255) or a hex color (#RGB / #RRGGBB).
func ValidColor(s string) bool {
	if len(s) == 4 || len(s) == 7 {
		if s[0] == '#' {
			for _, c := range s[1:] {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
					return false
				}
			}
			return true
		}
	}
	// ANSI 0-255
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
		n = n*10 + int(c-'0')
	}
	return len(s) > 0 && n >= 0 && n <= 255
}

func ToSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return strings.Trim(b.String(), "-")
}
