package plugin

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	// AuthzServiceKey is the key used to retrieve the authorization server URL from the configuration
	AuthzServiceKey = "authz_service"
	// AuthnServiceEndpointKey is the key used to retrieve the authorization server endpoint from the configuration
	AuthnServiceEndpointKey = "endpoint"
	// AuthnServiceTimeoutKey is the key used to retrieve the authorization server timeout from the configuration
	AuthnServiceTimeoutKey = "timeout"
	// ActionKey is the key used to retrieve the action from the configuration
	ActionKey = "action"
	// ResourceTypeKey is the key used to retrieve the resource type from the configuration
	ResourceTypeKey = "resource_type"
	// ResourceParamKey is the key used to retrieve the resource param from the configuration
	ResourceParamKey = "resource_param"
)

var (
	// ErrInvalidConfig is returned when the configuration is not valid
	ErrInvalidConfig = errors.New("invalid config")
	// ErrConfigurationNotFound is returned when the configuration is not found
	ErrConfigurationNotFound = errors.New("configuration not found")
)

type AuthzService struct {
	// Endpoint is the URL of the authorization server
	Endpoint *url.URL `json:"endpoint"`
	// Timeout is the timeout for the authorization server in milliseconds
	// defaults to 1000
	Timeout int `json:"timeout"`
}

type Config struct {
	// AuthorizationService is the URL of the authorization server
	AuthorizationService *AuthzService `json:"authz_service"`
	// Action is the action to be performed
	Action string `json:"action"`
	// ResourceType is the name of resource type
	ResourceType string `json:"resource_type"`
	// ResourceParam is the name of the resource parameter
	ResourceParam string `json:"resource_param"`
}

// ParseConfig parses the configuration and returns a Config object
// The configuration is the expected krakend format.
func ParseConfig(cfg map[string]interface{}) (*Config, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	// Verify plugin configuration is present
	pconf, ok := cfg[PluginName].(map[string]interface{})
	if !ok {
		return nil, ErrConfigurationNotFound
	}

	if pconf == nil {
		return nil, ErrConfigurationNotFound
	}

	// Verify authorization service configuration
	if pconf[AuthzServiceKey] == nil {
		return nil, fmt.Errorf("%w: %s is missing", ErrInvalidConfig, AuthzServiceKey)
	}

	if _, ok := pconf[AuthzServiceKey].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("%w: %s should be a map", ErrInvalidConfig, AuthzServiceKey)
	}

	authzSvc := pconf[AuthzServiceKey].(map[string]interface{})
	if authzSvc == nil {
		return nil, fmt.Errorf("%w: %s is missing", ErrInvalidConfig, AuthzServiceKey)
	}

	// Verify authorization service endpoint
	authzURL, authzURLVerifyErr := stringRequired(authzSvc, AuthnServiceEndpointKey)
	if authzURLVerifyErr != nil {
		return nil, authzURLVerifyErr
	}

	parsedURL, err := url.Parse(authzURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s is not a valid URL", ErrInvalidConfig, AuthzServiceKey)
	}

	// Get and verify timeout
	tmout, tmoutErr := getOrDefault(authzSvc, AuthnServiceTimeoutKey, 1000)
	if tmoutErr != nil {
		return nil, fmt.Errorf("%w: %s is not a valid timeout", ErrInvalidConfig, AuthnServiceTimeoutKey)
	}

	// Verify action
	action, actionVerifyErr := stringRequired(pconf, ActionKey)
	if actionVerifyErr != nil {
		return nil, actionVerifyErr
	}

	// Verify resource type
	resourceType, resourceTypeVerifyErr := stringRequired(pconf, ResourceTypeKey)
	if resourceTypeVerifyErr != nil {
		return nil, resourceTypeVerifyErr
	}

	// Verify resource path param
	resourceParam, resourceParamVerifyErr := stringRequired(pconf, ResourceParamKey)
	if resourceParamVerifyErr != nil {
		return nil, resourceParamVerifyErr
	}

	return &Config{
		AuthorizationService: &AuthzService{
			Endpoint: parsedURL,
			Timeout:  tmout,
		},
		Action:        action,
		ResourceType:  resourceType,
		ResourceParam: resourceParam,
	}, nil
}

func typeRequired[T any](conf map[string]interface{}, key string) (T, error) {
	var empty T
	if conf == nil {
		return empty, ErrInvalidConfig
	}

	if conf[key] == nil {
		return empty, fmt.Errorf("%w: %s is missing", ErrInvalidConfig, key)
	}

	if _, ok := conf[key].(T); !ok {
		return empty, fmt.Errorf("%w: %s is not a string", ErrInvalidConfig, key)
	}

	return conf[key].(T), nil
}

func stringRequired(conf map[string]interface{}, key string) (string, error) {
	val, err := typeRequired[string](conf, key)
	if err != nil {
		return "", err
	}

	if val == "" {
		return "", fmt.Errorf("%w: %s is empty", ErrInvalidConfig, key)
	}

	return val, nil
}

func getOrDefault[T comparable](conf map[string]interface{}, key string, def T) (T, error) {
	// This is an empty variable for us to compare against
	var empty T

	if conf == nil {
		return def, nil
	}

	if conf[key] == nil {
		return def, nil
	}

	if _, ok := conf[key].(T); !ok {
		return def, fmt.Errorf("%w: %s is of the wrong type", ErrInvalidConfig, key)
	}

	if conf[key].(T) == empty {
		return def, nil
	}

	return conf[key].(T), nil
}
