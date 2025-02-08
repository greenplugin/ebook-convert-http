package app_fx

import (
	"ebook-convert-http-wrapper/internal/infra/http_server_fx"
	"ebook-convert-http-wrapper/internal/infra/http_server_fx/handlers"
	"go.uber.org/fx"
	"net/http"
)

func Start() {
	fx.New(
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
