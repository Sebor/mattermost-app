package main

import (
	"mattermost-bot/internal/config"
	"mattermost-bot/internal/handlers"
	"mattermost-bot/internal/logger"
	"mattermost-bot/internal/meta"
	"net/http"
	"strconv"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mattermost/mattermost-plugin-apps/utils/httputils"
)

// main sets up the http server, with paths mapped for the static assets, the
// bindings' callback, and the other functions.
func main() {
	log := logger.GetLogger()
	cfg := config.GetConfig()

	handler := httputils.NewHandler()
	handler.HandleFunc("/manifest.json", httputils.DoHandleJSON(meta.Manifest))
	handler.HandleFunc("/bindings", handlers.Bindings)
	handler.HandleFunc("/create", handlers.Create)
	handler.HandleFunc("/create/withpost", handlers.CreateWithPost)
	handler.HandleFunc("/install", handlers.Install)
	handler.PathPrefix("/metrics").Handler(promhttp.Handler())
	handler.HandleFunc("/ping", httputils.DoHandleJSON(apps.NewDataResponse(nil)))
	handler.
		PathPrefix("/static/").
		Handler(http.FileServer(http.FS(meta.StaticFiles)))

	server := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.AppPort),
		Handler: handlers.LogRequest(handler),
	}
	log.Info().Msg("Listening on " + server.Addr)
	log.Info().Msg(
		"Use '/apps install http " +
			cfg.AppRootURL +
			"/manifest.json' to install the app") // matches manifest.json

	log.Error().Err(server.ListenAndServe()).Msg("")
}
