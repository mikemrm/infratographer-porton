package plugin

import (
	"context"
	"net/http"
)

// RegisterHandlers is the method that will be called by the plugin loader to register the plugin.
func (r portonRegisterer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.serverPluginHandle)
}

func (r portonRegisterer) serverPluginHandle(ctx context.Context, _ map[string]interface{}, h http.Handler) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Debug("serverPluginHandle")
		h.ServeHTTP(w, req)
	}), nil
}
