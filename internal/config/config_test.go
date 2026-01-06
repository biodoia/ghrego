package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Setup env
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DATABASE_URL")

	cfg := Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DatabaseURL)
	assert.Equal(t, 30*time.Second, cfg.ServerTimeout) // default
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				Port:           "8080",
				DatabaseURL:    "postgres://...",
				ServerTimeout:  10 * time.Second,
				MaxRequestSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "missing port",
			cfg: &Config{
				Port:        "",
				DatabaseURL: "postgres://...",
			},
			wantErr: true,
		},
		{
			name: "missing db url",
			cfg: &Config{
				Port:        "8080",
				DatabaseURL: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
