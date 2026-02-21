package gopherstatus

import (
	"os"
	"testing"
)

// Ensure that we can parse a good config as expected.
func TestParseConfig(t *testing.T) {
	sample_config := `
ip = "0.0.0.0"
port = 3000
template_path = "templates"

[status.restreamer]
type = "docker"
container_name = "restreamer"
image = "datarhei/restreamer:latest"
arguments = [
    "-e",
    "RS_USERNAME=admin",
    "-e",
    "RS_PASSWORD=datarhei",
    "-p",
    "8080:8080",
    "-v",
    "/mnt/restreamer/db:/restreamer/db"
]
`
	f, err := os.CreateTemp("", "temp")
	if err != nil {
		t.Errorf("Could not make temp file.")
	}
	_, err = f.Write([]byte(sample_config))
	if err != nil {
		t.Errorf("Could not write to temp file.")
	}
	config, err := ParseConfig(f.Name())
	if err != nil {
		t.Errorf("Could not parse config: %s", err)
	}

	if config.IP != "0.0.0.0" {
		t.Errorf("Expected `0.0.0.0` but got `%s` for IP.", config.IP)
	}

	if config.Port != 3000 {
		t.Errorf("Expected `3000` but got `%d` for Port.", config.Port)
	}

	if config.TemplatePath != "templates" {
		t.Errorf("Expected `templates` but got `%s` for TemplatePath.", config.TemplatePath)
	}

	os.Remove(f.Name())
}

// Ensure that when the config file does not exist that an error is returned.
func TestParseConfigFileNotExist(t *testing.T) {
	config, err := ParseConfig("/tmp/path/not/exist.toml")
	if err == nil {
		t.Errorf("Expected error when file does not exist.")
	}
	if config != nil {
		t.Errorf("Expected config to be nil when an error occurs.")
	}
}

// Ensure that an error is returned when we cannot parse the toml file.
func TestParseConfigNotTOMLFile(t *testing.T) {
	f, err := os.CreateTemp("", "temp")
	if err != nil {
		t.Errorf("Could not make temp file.")
	}
	_, err = f.Write([]byte("bad toml file."))
	if err != nil {
		t.Errorf("Could not write to temp file.")
	}
	config, err := ParseConfig(f.Name())

	if err == nil {
		t.Errorf("Expected an error but did not receive one.")
	}

	if config != nil {
		t.Errorf("Expected config to be nil when an error is returned but it was not.")
	}

	os.Remove(f.Name())
}
