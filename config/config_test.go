package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/peienxie/go-bank/config"
	"github.com/stretchr/testify/assert"
)

func createDefaultEnvFile(t *testing.T) {
	// create default app.env in current folder
	envs := make(map[string]string)
	envs["DB_SOURCE"] = "default_source"
	envs["DB_DRIVER"] = "default_driver"
	envs["SERVER_ADDRESS"] = "default_address"

	var envString string
	for k, v := range envs {
		envString += fmt.Sprintf("%s=%s\n", k, v)
	}

	err := os.WriteFile("app.env", []byte(envString), 0644)
	assert.NoError(t, err)
}

func cleanupEnvFile(t *testing.T) {
	err := os.Remove("app.env")
	assert.NoError(t, err)
}

func TestLoadDefaultConfig(t *testing.T) {
	createDefaultEnvFile(t)

	config, err := config.LoadConfig(".")
	assert.NoError(t, err)
	assert.NotEmpty(t, config)

	assert.Equal(t, "default_driver", config.DBDriver)
	assert.Equal(t, "default_source", config.DBSource)
	assert.Equal(t, "default_address", config.ServerAddress)

	cleanupEnvFile(t)
}

func TestOverrideConfigByEnvironmentVariables(t *testing.T) {
	os.Setenv("GOBANK_DB_SOURCE", "mysource")
	os.Setenv("GOBANK_DB_DRIVER", "mydriver")
	os.Setenv("GOBANK_SERVER_ADDRESS", "localhost:8080")

	config, err := config.LoadConfig(".")
	assert.NoError(t, err)
	assert.NotEmpty(t, config)

	assert.Equal(t, "mydriver", config.DBDriver)
	assert.Equal(t, "mysource", config.DBSource)
	assert.Equal(t, "localhost:8080", config.ServerAddress)
}
