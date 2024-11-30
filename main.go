package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/jcbwlkr/deck-stats/internal/auth"
	"github.com/jcbwlkr/deck-stats/internal/database"
	"github.com/jcbwlkr/deck-stats/internal/domains/magic"
	"github.com/jcbwlkr/deck-stats/internal/domains/users"
	"github.com/jcbwlkr/deck-stats/internal/handlers"
	"github.com/jcbwlkr/deck-stats/internal/services/moxfield"
)

func main() {
	if err := run(); err != nil {
		slog.Error("main application failure", "error", err)
		os.Exit(1)
	}
}

func run() error {

	var config struct {
		AddressServer string `envconfig:"address_server"`

		JWTSecret string `envconfig:"jwt_secret"`

		DBName string `envconfig:"db_name"`
		DBHost string `envconfig:"db_host"`
		DBPort int    `envconfig:"db_port"`
		DBUser string `envconfig:"db_user"`
		DBPass string `envconfig:"db_pass"`

		MoxfieldUserAgent string `envconfig:"moxfield_user_agent"`
	}
	envconfig.MustProcess("", &config)

	db, err := database.Open(
		config.DBUser,
		config.DBName,
		config.DBPass,
		config.DBHost,
		config.DBPort,
	)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	mc := moxfield.NewClient(config.MoxfieldUserAgent, 1*time.Second)
	userService := users.NewService(db)
	magicService := magic.NewService(db, userService, mc)
	defer magicService.Wait()
	authenticator := auth.NewAuthenticator(config.JWTSecret)

	app := handlers.App(magicService, userService, authenticator)

	// TODO(jlw) graceful shutdown

	slog.Info("deck-stats api running", "address", config.AddressServer)
	return http.ListenAndServe(config.AddressServer, app)
}
