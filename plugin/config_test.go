package plugin

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseURL(t *testing.T, u string) *url.URL {
	t.Helper()

	parsed, err := url.Parse(u)
	require.NoError(t, err, "failed to parse URL: %s", u)

	return parsed
}

func TestParseConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     map[string]interface{}
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
						"timeout":  2000,
					},
					"action":         "read",
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want: &Config{
				AuthorizationService: &AuthzService{
					Endpoint: mustParseURL(t, "http://authz"),
					Timeout:  2000,
				},
				Action:        "read",
				ResourceType:  "test",
				ResourceParam: "test_id",
			},
			wantErr: false,
		},
		{
			name: "valid config with default timeout",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
					},
					"action":         "read",
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want: &Config{
				AuthorizationService: &AuthzService{
					Endpoint: mustParseURL(t, "http://authz"),
					Timeout:  1000,
				},
				Action:        "read",
				ResourceType:  "test",
				ResourceParam: "test_id",
			},
			wantErr: false,
		},
		{
			name: "invalid config - missing authz_service",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"action": "read",
				},
				"resource_type":  "test",
				"resource_param": "test_id",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - authz_service empty",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service":  map[string]interface{}{},
					"action":         "read",
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - missing authz_service.endpoint",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"timeout": 2000,
					},
					"action":         "read",
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - empty authz_service.endpoint",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "",
						"timeout":  2000,
					},
					"action":         "read",
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - invalid authz_service.endpoint",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "alksdjflkads://laksdjfadkls-  asjkldhf askldjl;a",
						"timeout":  2000,
					},
					"action":         "read",
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - invalid authz_service.timeout",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
						"timeout":  "invalid",
					},
					"action":         "read",
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - invalid action",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
						"timeout":  2000,
					},
					"action":         1234,
					"resource_type":  "test",
					"resource_param": "test_id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - missing resource_type",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
						"timeout":  2000,
					},
					"action":         "read",
					"resource_param": "test_id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - missing resource_param",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
						"timeout":  2000,
					},
					"action":        "read",
					"resource_type": "test",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - invalid resource_type type",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
						"timeout":  2000,
					},
					"action":        "read",
					"resource_type": 1234,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid config - invalid resource_param type",
			cfg: map[string]interface{}{
				PluginName: map[string]interface{}{
					"authz_service": map[string]interface{}{
						"endpoint": "http://authz",
						"timeout":  2000,
					},
					"action":         "read",
					"resource_param": 1234,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				require.NoError(t, err, "ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.want, got, "ParseConfig() = %v, want %v", got, tt.want)
		})
	}
}
