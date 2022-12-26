package settings

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Keys_path string `yaml:"keys_path"`
}

func Read() *Settings {
	yamlf, err := os.ReadFile("./setting/settings_default.yml")
	if err != nil {
		panic("Cannot read settings file")
	}

	var settings Settings
	yaml.Unmarshal(yamlf, &settings)

	return &settings
}
