package app_fx

import (
	"ebook-convert-http-wrapper/internal/contract"
	"ebook-convert-http-wrapper/internal/infra/http_server_fx"
	"ebook-convert-http-wrapper/internal/infra/http_server_fx/handlers"
	"github.com/spf13/pflag"
	"go.uber.org/fx"
	"net/http"
)

func Start() {
	var port contract.Port = "12600"

	pflag.VarP(&port, "port", "p", "port to start the server (default 12600)")
	pflag.Parse()

	fx.New(
		fx.Supply(port),
		fx.Provide(
			http_server_fx.NewHTTPServer,
			asHandler(handlers.NewConvert),
			asHandler(handlers.NewHealth),
			asHandler(handlers.NewRecipes),
		),
		fx.Invoke(func(*http.Server) {}),
	).
		Run()
}

func asHandler(h any) any {
	return fx.Annotate(
		h,
		fx.As(new(handlers.Handler)),
		fx.ResultTags(`group:"http_handlers"`),
	)
}
