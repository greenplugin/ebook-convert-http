package http_server_fx

import (
	"context"
	"ebook-convert-http-wrapper/internal/infra/http_server_fx/handlers"
	"flag"
	"fmt"
	"go.uber.org/fx"
	"log"
	"net"
	"net/http"
)

type HTTPServerParams struct {
	fx.In

	Lc       fx.Lifecycle
	Handlers []handlers.Handler `group:"http_handlers"`
}

func NewHTTPServer(p HTTPServerParams) *http.Server {
	port := flag.String("port", "12600", "port to start the server (default 12600)")
	flag.StringVar(port, "p", "12600", "port to start the server (short parameter)")

	srv := &http.Server{Addr: ":" + *port}
	mux := http.NewServeMux()
	srv.Handler = mux
	log.Printf("Starting HTTP server at %s", srv.Addr)
	for _, h := range p.Handlers {
		log.Printf("Registering handler %T", h)
		h.Register(mux)
	}

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
