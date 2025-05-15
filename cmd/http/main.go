package main

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/lulzshadowwalker/personal/internal/config"
	"github.com/lulzshadowwalker/personal/internal/template"
)

func main() {
	component := template.Hello("zaya")
	
	http.Handle("/", templ.Handler(component))

	slog.Info("starting server", "port", config.Port())
	http.ListenAndServe(":" + config.Port(), nil)
}
