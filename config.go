package gopherstatus

import (
	"os"

	"gopherstatus/status"

	"github.com/BurntSushi/toml"
)

type GSConfig struct {
	IP           string                          `toml:"ip"`
	Port         int                             `toml:"port"`
	TemplatePath string                          `toml:"template_path"`
	Status       map[string]interface{}          `toml:"status"`
	Managers     map[string]status.StatusManager `toml:"-"`
}

func ParseConfig(path string) (*GSConfig, error) {
	var config GSConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(data)
	if _, err := toml.Decode(text, &config); err != nil {
		return nil, err
	}
	config.Managers = make(map[string]status.StatusManager)
	return &config, nil
}
