package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/handlename/awsc/internal/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_determineConfigPath(t *testing.T) {
	homedir := os.Getenv("HOME")
	defer os.Setenv("HOME", homedir)

	resetEnvs := func() error {
		envs := []string{
			"HOME",
			env.EnvConfigPath,
			env.EnvDefaultConfigDir,
		}

		for _, e := range envs {
			if err := os.Unsetenv(e); err != nil {
				return err
			}
		}

		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	testdataRoot := filepath.Join(filepath.Dir(filename), "testdata", "cli_determine_config_path")

	tests := []struct {
		name    string
		env     map[string]string
		want    string
		wantErr bool
		errBody string
	}{
		{
			name: "no envs",
			env: map[string]string{
				"HOME": filepath.Join(testdataRoot, "home0"),
			},
			want: "",
		},
		{
			name: fmt.Sprintf("set %s", env.EnvConfigPath),
			env: map[string]string{
				env.EnvConfigPath: filepath.Join(testdataRoot, "specific", "config.toml"),
			},
			want: filepath.Join(testdataRoot, "specific", "config.toml"),
		},
		{
			name: "exists ~/.config/awsc/config.yaml",
			env: map[string]string{
				"HOME": filepath.Join(testdataRoot, "home1"),
			},
			want: filepath.Join(testdataRoot, "home1", ".config", "awsc", "config.yaml"),
		},
		{
			name: "exists ~/.awsc/config.yaml",
			env: map[string]string{
				"HOME": filepath.Join(testdataRoot, "home2"),
			},
			want: filepath.Join(testdataRoot, "home2", ".awsc", "config.yaml"),
		},
		{
			name: "a file is placed in path than expects a directory",
			env: map[string]string{
				"HOME": filepath.Join(testdataRoot, "home3"),
			},
			wantErr: true,
			errBody: "expects a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, resetEnvs())

			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			path, err := determineConigPath()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errBody)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, path)
		})
	}
}
