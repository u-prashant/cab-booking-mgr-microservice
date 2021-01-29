package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_load(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		expectedConfig *Config
	}{
		{
			name: "valid path",
			args: args{"./config.yaml"},
			expectedConfig: &Config{
				DSN:     "root:root@tcp(127.0.0.1:3306)/cab",
				Address: "0.0.0.0:8080",
			},
			wantErr: false,
		},
		{
			name:           "invalid path",
			args:           args{"invalid.yaml"},
			expectedConfig: &Config{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultErr := load(tt.args.filename)
			require.Equal(t, tt.wantErr, resultErr != nil, "err: %s", resultErr)
			assert.Equal(t, tt.expectedConfig, App, tt.name)
		})
	}
}
