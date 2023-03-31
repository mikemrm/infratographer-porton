// package plugin implements the porton krakend.io plugin
package plugin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

/*
The portonRegisterer object implements the Registerer interface
for Lura's server and client plugins.

We must be careful not to change the `RegisterClients` and `RegisterHandlers`
methods' signatures, as they are called by the plugin loader.

It is not recommended to import the packages that define these interfaces.
That is:

* "github.com/luraproject/lura/v2/transport/http/client/plugin"
* "github.com/luraproject/lura/v2/transport/http/server/plugin"

This is to avoid any dependency issues that may arise from the plugin loader
importing the same packages as the main application. We instead stick to just
relying on the interfaces' signatures.

Also note that the plugin loader interface does not work if the `*Registerer`
variable instances are pointers. We need to use values instead. We decided to
build on top of `string` since that's what's working and documented in the
official plugin loader example.
*/
type portonRegisterer string

var unkownTypeErr = errors.New("unknow request type")

// PluginName is the name of the plugin
const PluginName = "porton"

// RequestWrapper is an interface for passing proxy request between the krakend pipe
// and the loaded plugins
type RequestWrapper interface {
	Params() map[string]string
	Headers() map[string][]string
	Body() io.ReadCloser
	Method() string
	URL() *url.URL
	Query() url.Values
	Path() string
}

// NewPortonRegisterer returns a new portonRegisterer object
func NewPortonRegisterer(name string) portonRegisterer {
	return portonRegisterer(name)
}

// RegisterClients is the method that will be called by the plugin loader to register the plugin.
func (r portonRegisterer) RegisterModifiers(f func(
	name string,
	factoryFunc func(map[string]interface{}) func(interface{}) (interface{}, error),
	appliesToRequest bool,
	appliesToResponse bool,
)) {
	f(string(r), r.requestModPluginHandle, true, false)
}

// requestModPluginHandle is the function that will be called by the plugin loader to register the plugin.
// It returns a function that will be called by the krakend pipe to handle the request.
func (r portonRegisterer) requestModPluginHandle(conf map[string]interface{}) func(interface{}) (interface{}, error) {
	cfg, err := ParseConfig(conf)
	if err != nil {
		logger.Error(err)
		return func(interface{}) (interface{}, error) {
			return nil, fmt.Errorf("error parsing config: %w", err)
		}
	}

	return func(input interface{}) (interface{}, error) {
		req, ok := input.(RequestWrapper)
		if !ok {
			return nil, unkownTypeErr
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.AuthorizationService.Timeout)*time.Millisecond)
		defer cancel()

		allowed, err := handleAuthorizationRequest(ctx, req, cfg)
		if err != nil {
			logger.Error(err)
			return nil, HTTPResponseError{
				Code:         http.StatusInternalServerError,
				Msg:          "error handling request",
				HTTPEncoding: HTTPJSONEncoding,
			}
		}

		if !allowed {
			logger.Info("not allowed")
			return nil, HTTPResponseError{
				Code:         http.StatusForbidden,
				Msg:          "not allowed",
				HTTPEncoding: HTTPJSONEncoding,
			}
		}

		logger.Info("allowed")
		return input, nil
	}
}
