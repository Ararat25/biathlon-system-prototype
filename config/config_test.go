package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_ValidConfig проверяет работу функции на корректных данных
func TestLoadConfig_ValidConfig(t *testing.T) {
	configPath := "test_config.json"
	content := `{
		"laps": 3,
		"lapLen": 1000,
		"penaltyLen": 150,
		"firingLines": 2,
		"start": "12:00:00.000",
		"startDelta": "00:01:00"
	}`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)
	defer os.Remove(configPath)

	cfg, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.Equal(t, 3, cfg.Laps)
	assert.Equal(t, "12:00:00.000", cfg.Start)
}
