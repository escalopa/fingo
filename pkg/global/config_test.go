package global

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	type testConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	defer close(done)

	deleteFile := func(name string) {
		wg.Add(1)
		defer wg.Done()
		<-done
		os.Remove(name)
	}

	tests := []struct {
		name            string
		cfg             interface{}
		createConigFile func(t *testing.T) (string, string, string)
		wantErr         bool
	}{
		{
			name: "success yaml",
			cfg:  testConfig{},
			createConigFile: func(t *testing.T) (string, string, string) {
				file, err := os.Create("/tmp/fingo_config.yaml")
				go deleteFile("fingo_config.yaml")
				require.NoError(t, err)
				file.WriteString("host: localhost\nport: 8080")
				return filepath.Base(file.Name()), "/tmp", filepath.Ext(file.Name())[1:]
			},
			wantErr: false,
		},
		{
			name: "success json",
			cfg:  testConfig{},
			createConigFile: func(t *testing.T) (string, string, string) {
				file, err := os.Create("/tmp/fingo_config.json")
				go deleteFile("fingo_config.json")
				require.NoError(t, err)
				file.WriteString(`{"host": "localhost", "port": 8080}`)
				return filepath.Base(file.Name()), "/tmp", filepath.Ext(file.Name())[1:]
			},
			wantErr: false,
		},
		{
			name: "success toml",
			cfg:  testConfig{},
			createConigFile: func(t *testing.T) (string, string, string) {
				file, err := os.Create("/tmp/fingo_config.toml")
				go deleteFile("fingo_config.toml")
				require.NoError(t, err)
				file.WriteString("host = \"localhost\"\nport = 8080")
				return filepath.Base(file.Name()), "/tmp", filepath.Ext(file.Name())[1:]
			},
			wantErr: false,
		},
		{
			name: "failed to read config",
			createConigFile: func(t *testing.T) (string, string, string) {
				file, err := os.Create("/tmp/fingo_config.yaml")
				go deleteFile("fingo_config.yaml")
				require.NoError(t, err)
				file.WriteString("host: localhost\nport: 8080")
				return filepath.Base(file.Name()), "/tmp", "json"
			},
			wantErr: true,
		},
		{
			name: "failed to unmarshall config",
			cfg:  func() {}, // upsupported type for unmarshalling config
			createConigFile: func(t *testing.T) (string, string, string) {
				file, err := os.Create("/tmp/fingo_config.yaml")
				go deleteFile("fingo_config.yaml")
				require.NoError(t, err)
				file.WriteString("host: localhost\nport: 8080")
				return filepath.Base(file.Name()), "/tmp", filepath.Ext(file.Name())[1:]
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, path, form := tt.createConigFile(t)
			defer func() { done <- struct{}{} }() // remove file after test
			err := LoadConfig(&tt.cfg, name, path, form)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				x, ok := tt.cfg.(testConfig)
				require.True(t, ok)
				require.NotEmpty(t, x)
				require.Equal(t, "localhost", x.Host)
				require.Equal(t, 8080, x.Port)
			}
		})
	}

	wg.Wait()
}
