package commands

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/albertogviana/port-service/internal/handlers"
	"github.com/gorilla/mux"

	logger "github.com/albertogviana/port-service/internal/log"
	"github.com/albertogviana/port-service/internal/port"
	"github.com/albertogviana/port-service/internal/repository"
	"github.com/go-sql-driver/mysql"
	muxHandlers "github.com/gorilla/handlers"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

// Start starts the rest api.
func Start(c *cli.Context) error {
	dbHost := "127.0.0.1"
	if os.Getenv("PORT_DB_HOST") != "" {
		dbHost = os.Getenv("PORT_DB_HOST")
	}

	// For simplicity the credentials are hardcoded.
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 dbHost,
		DBName:               "ports",
		User:                 "root",
		Passwd:               "root",
		AllowNativePasswords: true,
		ParseTime:            true,
		Timeout:              10 * time.Second,
	}

	dbConn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("error during the database connection: %w", err)
	}

	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(10)

	if err := dbConn.Ping(); err != nil {
		return fmt.Errorf("error during the database connection: %w", err)
	}

	defer dbConn.Close()

	// logging
	log := logger.NewLogger(
		"development",
		os.Stdout,
		true,
	)

	router := mux.NewRouter()

	var svc *port.Service
	{
		svc = port.NewService(repository.NewMySQLRepository(dbConn))

		handlers.MakePortHandlers(
			router,
			svc,
			log,
		)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: muxHandlers.RecoveryHandler()(router),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Errorf("listen and Serve problem: %w", err))
		}
	}()

	log.Infof("Started HTTP on %s", srv.Addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Block until we receive our signal.
	<-signalCh

	log.Info("Shutdown port service down")

	cancel()

	var g errgroup.Group

	g.Go(func() error {
		srv.Shutdown(ctx)

		return nil
	})

	g.Wait()

	return nil
}
