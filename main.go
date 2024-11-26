package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/jcbwlkr/deck-stats/internal/database"
	"github.com/jcbwlkr/deck-stats/internal/domains/magic"
	"github.com/jcbwlkr/deck-stats/internal/handlers"
	"github.com/jcbwlkr/deck-stats/internal/moxfield"
	"github.com/jcbwlkr/deck-stats/internal/services/users"
)

func main() {
	if err := run(); err != nil {
		slog.Error("main application failure", "error", err)
		os.Exit(1)
	}
}

func run() error {

	var config struct {
		AddressServer string `envconfig:"address_server" default:"localhost:9040"`

		DBName string `envconfig:"db_name"`
		DBHost string `envconfig:"db_host"`
		DBPort int    `envconfig:"db_port"`
		DBUser string `envconfig:"db_user"`
		DBPass string `envconfig:"db_pass"`
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

	mc := moxfield.NewClient(1 * time.Second)
	userService := users.NewService(db)
	magicService := magic.NewService(db, userService, mc)
	defer magicService.Wait()

	app := handlers.App(magicService)

	// TODO(jlw) graceful shutdown

	slog.Info("deck-stats api running", "address", config.AddressServer)
	return http.ListenAndServe(config.AddressServer, app)
}
