package settings

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Keys_path           string `yaml:"keys_path"`
	Active_connections  int    `yaml:"active_connections"`
	Default_public_node string `yaml:"default_public_node"`
}

func Read() *Settings {
	yamlf, err := os.ReadFile("./settings/settings_default.yml")
	if err != nil {
		panic("Cannot read settings file")
	}

	var settings Settings
	yaml.Unmarshal(yamlf, &settings)

	return &settings
}
