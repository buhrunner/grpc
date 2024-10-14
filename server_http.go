package grpc

import (
	"github.com/buhrunner/grpc/v5/common"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

func (p *Plugin) createHTTPServer(middlewares map[string]common.Middleware) *http.Server {
	mux := http.NewServeMux()

	handler := p.middleware()

	for _, middleware := range middlewares {
		handler = middleware.Middleware(handler)
	}

	mux.Handle("/", handler)

	server := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{
			MaxConcurrentStreams:         uint32(p.config.MaxConcurrentStreams),
			PermitProhibitedCipherSuites: false,
		}),
	}

	return server
}

func (p *Plugin) middleware() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		p.server.ServeHTTP(writer, req)
	})
}
