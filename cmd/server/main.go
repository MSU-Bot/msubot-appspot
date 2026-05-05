package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/MSU-Bot/msubot-appspot/internal/api"
	"github.com/MSU-Bot/msubot-appspot/internal/handlers"
	"github.com/MSU-Bot/msubot-appspot/internal/logging"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	r := chi.NewRouter()
	r.Use(logging.Recover)
	r.Use(chimw.RequestID)
	r.Use(logging.Chi)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	var h handlers.Server
	srv := api.HandlerFromMux(api.NewStrictHandler(&h, nil), r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	slog.Info("server listening", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, srv); err != nil {
		slog.Error("server exited", slog.Any("error", err))
		os.Exit(1)
	}
}
