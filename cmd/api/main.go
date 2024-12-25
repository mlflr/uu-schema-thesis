package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"thesis.lefler.eu/internal/data"
	e "thesis.lefler.eu/internal/error"
	"thesis.lefler.eu/internal/handler"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn struct {
			views string
		}
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type dbConns struct {
	views *sql.DB
}

type application struct {
	config   config
	handlers handler.Handlers
	logger   *slog.Logger
	errors   e.Errors
	models   data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn.views, "db-dsn", os.Getenv("VIEWS_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbConns := dbConns{}

	err := openDB(cfg, cfg.db.dsn.views, dbConns.views, "views")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer dbConns.views.Close()

	logger.Info("database connection pool established")

	models := data.NewModels(dbConns.views)
	errors := e.NewErrors(logger)

	app := &application{
		config:   cfg,
		logger:   logger,
		models:   models,
		errors:   errors,
		handlers: handler.NewHandlers(&errors, &models),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "localhost", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg config, dsn string, dbConn *sql.DB, method string) error {

	if dsn == "" {
		return fmt.Errorf("missing %s DSN", method)
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	dbConn.SetMaxOpenConns(cfg.db.maxOpenConns)
	dbConn.SetMaxIdleConns(cfg.db.maxIdleConns)
	dbConn.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbConn.PingContext(ctx)
	if err != nil {
		dbConn.Close()
		return err
	}

	return nil
}
