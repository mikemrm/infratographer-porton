package plugin

import (
	"context"
	"io"
	"net/http"
)

// RegisterClients is the method that will be called by the plugin loader to register the plugin.
func (r portonRegisterer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.clientPluginHandle)
}

func (r portonRegisterer) clientPluginHandle(ctx context.Context, _ map[string]interface{}) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger.Debug("clientPluginHandle")

		// We simply forward the request to the backend
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy headers, status codes, and body from the backend to the response writer
		for k, hs := range resp.Header {
			for _, h := range hs {
				w.Header().Add(k, h)
			}
		}
		w.WriteHeader(resp.StatusCode)
		if resp.Body == nil {
			return
		}
		io.Copy(w, resp.Body)
		resp.Body.Close()
	}), nil
}
