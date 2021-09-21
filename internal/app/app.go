package app

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/mebr0/squirrel-bot/internal/config"
	"github.com/mebr0/squirrel-bot/internal/repo"
	"github.com/mebr0/squirrel-bot/internal/service"
	"github.com/mebr0/squirrel-bot/internal/telegram"
	"github.com/mebr0/squirrel-bot/pkg/database"
	"github.com/mebr0/squirrel-bot/pkg/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg := config.LoadConfig(configPath)

	// Logging
	log, err := logging.NewLogger(cfg.Log.Level)

	if err != nil {
		fmt.Println("error initializing logger instance - " + err.Error())
		return
	}

	//nolint:errcheck
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()

	// Database
	db, err := database.NewPostgres(cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.User, cfg.DB.Password, cfg.DB.SSLMode)

	if err != nil {
		log.Fatal("error connecting to postgres - " + err.Error())
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		log.Fatal("error creating postgres driver - " + err.Error())
	}

	// Migrations
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", cfg.DB.Name, driver)

	if err != nil {
		log.Fatal("error reading migrations - " + err.Error())
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error migrating database schema - " + err.Error())
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)

	if err != nil {
		log.Fatal("error initializing telegram client - " + err.Error())
	}

	repos := repo.NewRepos(db)
	services := service.NewServices(repos)

	// Telegram bot
	bot := telegram.NewBot(botApi, services, log, cfg.Game)
	go func() {
		if err = bot.Start(); err != nil {
			log.Fatal("error while bot polling - " + err.Error())
		}
	}()

	log.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = bot.Stop(ctx); err != nil {
		log.Error("failed to stop bot - " + err.Error())
	}

	if err = db.Close(); err != nil {
		log.Error("failed to close database connection - " + err.Error())
	}
}
