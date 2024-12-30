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
			views           string
			expandDeprecate string
			branches        string
		}
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
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

	flag.StringVar(&cfg.db.dsn.views, "db-dsn-views", os.Getenv("VIEWS_DB_DSN"), "PostgreSQL DSN for Views method")
	flag.StringVar(&cfg.db.dsn.expandDeprecate, "db-dsn-expand-deprecate", os.Getenv("EXPAND_DEPRECATE_DB_DSN"), "PostgreSQL DSN for Expand & Deprecate method")
	flag.StringVar(&cfg.db.dsn.branches, "db-dsn-branches", os.Getenv("BRANCHES_DB_DSN"), "PostgreSQL DSN for Branches method")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	dbConns := data.DbConns{
		Views: nil,
	}

	var err error

	dbConns.Views, err = openDB(cfg, cfg.db.dsn.views, "views")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dbConns.ExpandDeprecate, err = openDB(cfg, cfg.db.dsn.expandDeprecate, "expandDeprecate")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dbConns.Branches, err = openDB(cfg, cfg.db.dsn.branches, "branches")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer dbConns.Views.Close()
	defer dbConns.ExpandDeprecate.Close()
	defer dbConns.Branches.Close()

	logger.Info("database connection pool established")

	models := data.NewModels(dbConns)
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

func openDB(cfg config, dsn string, method string) (*sql.DB, error) {

	if dsn == "" {
		return nil, fmt.Errorf("missing %s DSN", method)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
