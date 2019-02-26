package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	Port int `env:"port" envDefault:"3000"`

	BannerColor string `env:"BANNER_COLOR" envDefault:"green"`
}

func main() {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339Nano,
		},
	)

	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse configuration values")
	}

	srv := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: RegisterRouter(),
	}
	go func() {
		log.Info().Int("port", cfg.Port).Msg("Listening HTTP")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener
			log.Fatal().Err(err).Msg("Error when running HTTP server")
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
	// We received an interrupt signal, shut down.
	log.Info().Msg("Shutting down...")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}
	log.Info().Msg("Server has been stopped")
}
