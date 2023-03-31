package plugin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	authclientv1 "go.infratographer.com/permissions-api/pkg/client/v1"
)

const (
	// AuthorizationHeader is the name of the header that contains the authorization token
	AuthorizationHeader = "Authorization"
	// TenantHeader is the name of the header that contains the tenant ID
	TenantHeader = "X-Tenant-Id"
	// HTTPJSONEncoding is the JSON encoding
	HTTPJSONEncoding = "application/json"
)

var (
	// ErrNoValidToken is returned when no valid token is found in the request
	ErrNoValidToken = errors.New("no valid token found")
	// ErrNoValidTenant is returned when no valid tenant is found in the request
	ErrNoValidTenant = errors.New("no valid tenant found")
	// ErrCreatingAuthzClient is returned when an error occurs creating the authz client
	ErrCreatingAuthzClient = errors.New("error creating authz client")
	// ErrCheckingPermissions is returned when an error occurs checking permissions
	ErrCheckingPermissions = errors.New("error checking permissions")
)

type HTTPResponseError struct {
	Code         int    `json:"http_status_code"`
	Msg          string `json:"http_body,omitempty"`
	HTTPEncoding string `json:"http_encoding"`
}

// Error returns the error message
func (r HTTPResponseError) Error() string {
	return r.Msg
}

// StatusCode returns the status code returned by the backend
func (r HTTPResponseError) StatusCode() int {
	return r.Code
}

// Encoding returns the HTTP output encoding
func (r HTTPResponseError) Encoding() string {
	return r.HTTPEncoding
}

// getAuthorizationHeader returns the value of the Authorization header from the given request
func getAuthorizationHeader(req RequestWrapper) string {
	return getHeader(req, AuthorizationHeader)
}

// getTenantHeader returns the value of the X-Tenant-Id header from the given request
func getTenantHeader(req RequestWrapper) string {
	return getHeader(req, TenantHeader)
}

// getHeader returns the value of the given header from the given request
func getHeader(req RequestWrapper, header string) string {
	if req.Headers() == nil {
		return ""
	}
	if req.Headers()[header] == nil {
		return ""
	}
	if len(req.Headers()[header]) == 0 {
		return ""
	}

	return req.Headers()[header][0]
}

// getTenantFromPath returns the tenant ID from the path
// It will find the `tenants` path segment and fetch the next one
// This assumes that the path is in the format `/v1/tenants/{tenant_id}/...`
// and that the tenant ID is the next path segment
// If the path is not in this format, the tenant ID will be empty
// and the request will be rejected
func getTenantFromPath(req RequestWrapper) string {
	p := req.Path()

	segments := strings.Split(p, "/")
	for i, s := range segments {
		if s == "tenants" {
			if i+1 < len(segments) {
				return segments[i+1]
			}
		}
	}

	return ""
}

// tokenRoundTripper is a round tripper that adds the authorization token to the request
// It is used to create the authz client.
// Note that token validation is not performed here, it is performed by an earlier
// plugin in the API Gateway, as well as the authorization service.
// Note that this simple round tripper was created to avoid conflicting
// dependencies with https://pkg.go.dev/golang.org/x/oauth2 which is used by the
// krakend plugin builder.
type tokenRoundTripper struct {
	token string
	trans http.RoundTripper
}

func newTokenRoundTripper(token string, trans http.RoundTripper) *tokenRoundTripper {
	return &tokenRoundTripper{
		token: token,
		trans: trans,
	}
}

// RoundTrip adds the authorization token to the request and calls the default transport
// to perform the request.
func (t *tokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(AuthorizationHeader, t.token)
	return t.trans.RoundTrip(req)
}

// handleAuthorizationRequest handles the authorization request
// It returns a boolean indicating whether the request is authorized and an error
func handleAuthorizationRequest(ctx context.Context, req RequestWrapper, cfg *Config) (bool, error) {
	btok := getAuthorizationHeader(req)
	if btok == "" {
		return false, ErrNoValidToken
	}

	var tenantID string
	switch cfg.TenantSource {
	case HeaderTenantSource:
		tenantID = getTenantHeader(req)
	case PathTenantSource:
		tenantID = getTenantFromPath(req)
	}

	if tenantID == "" {
		return false, ErrNoValidTenant
	}

	httpcli := &http.Client{
		Transport: newTokenRoundTripper(btok, http.DefaultTransport),
	}
	authzcli, err := authclientv1.New(cfg.AuthorizationService.Endpoint.String(), httpcli)
	if err != nil {
		return false, ErrCreatingAuthzClient
	}

	allowed, err := authzcli.Allowed(ctx, cfg.Action, tenantID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrCheckingPermissions, err)
	}

	return allowed, nil
}
