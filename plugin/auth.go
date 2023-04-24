package plugin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/textproto"

	"github.com/google/uuid"
	authclientv1 "go.infratographer.com/permissions-api/pkg/client/v1"
	"go.infratographer.com/x/urnx"
)

const (
	// AuthorizationHeader is the name of the header that contains the authorization token
	AuthorizationHeader = "Authorization"
	// HTTPJSONEncoding is the JSON encoding
	HTTPJSONEncoding = "application/json"
)

var (
	// ErrNoValidToken is returned when no valid token is found in the request
	ErrNoValidToken = errors.New("no valid token found")
	// ErrCreatingAuthzClient is returned when an error occurs creating the authz client
	ErrCreatingAuthzClient = errors.New("error creating authz client")
	// ErrCheckingPermissions is returned when an error occurs checking permissions
	ErrCheckingPermissions = errors.New("error checking permissions")
	// ErrNoValidResourceID is returned when no valid resource ID is found in the request
	ErrNoValidResourceID = errors.New("no valid resource ID found")
	// ErrInvalidResourceUUID is returned when the resource ID is not a valid UUID
	ErrInvalidResourceUUID = errors.New("resource ID is not a valid UUID")
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
	resourceId := getResourceID(req, cfg.ResourceParam)
	if resourceId == "" {
		return false, ErrNoValidResourceID
	}

	resUUID, err := uuid.Parse(resourceId)
	if err != nil {
		return false, ErrInvalidResourceUUID
	}

	urn, err := urnx.Build("infratrographer", cfg.ResourceType, resUUID)
	if err != nil {
		logger.Error("error building urn from resource type and id", err)
		return false, ErrInvalidResourceUUID
	}

	btok := getAuthorizationHeader(req)
	if btok == "" {
		return false, ErrNoValidToken
	}

	httpcli := &http.Client{
		Transport: newTokenRoundTripper(btok, http.DefaultTransport),
	}
	authzcli, err := authclientv1.New(cfg.AuthorizationService.Endpoint.String(), httpcli)
	if err != nil {
		return false, ErrCreatingAuthzClient
	}

	allowed, err := authzcli.Allowed(ctx, cfg.Action, urn.String())
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrCheckingPermissions, err)
	}

	return allowed, nil
}

// getResourceID returns the resource ID from the request
func getResourceID(req RequestWrapper, paramName string) string {
	if req.Params() == nil {
		return ""
	}

	// this is what they do in the lura project :grimace:
	// https://github.com/luraproject/lura/blob/20f8788a5f61d85a73cc31247c7e3b0dcd86a6df/router/gin/endpoint.go#L117
	p := textproto.CanonicalMIMEHeaderKey(paramName[:1]) + paramName[1:]

	return req.Params()[p]
}
