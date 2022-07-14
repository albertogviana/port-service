package commands

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/albertogviana/port-service/internal/repository"

	logger "github.com/albertogviana/port-service/internal/log"
	"github.com/albertogviana/port-service/internal/port"
	"github.com/urfave/cli/v2"

	"github.com/go-sql-driver/mysql"
)

// Import imports the port data from a file to a database.
func Import(c *cli.Context) error {
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

	svc := port.NewService(repository.NewMySQLRepository(dbConn))

	filename := c.String("file")
	if filename == "" {
		return fmt.Errorf("the parameter file is required")
	}

	// logging
	log := logger.NewLogger(
		"development",
		os.Stdout,
		true,
	)

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error to read [file=%v]: %w", filename, err)
	}

	_, err = f.Stat()
	if err != nil {
		return fmt.Errorf("could not obtain stat, handle error: %w", err)
	}

	r := bufio.NewReader(f)
	d := json.NewDecoder(r)

	ctx := context.Background()

	count := 0
	for d.More() {
		m := make(map[string]*port.Port, 1)

		if err := d.Decode(&m); err != nil {
			return fmt.Errorf("error to decode: %w", err)
		}

		for _, v := range m {
			count++
			err := svc.SavePort(ctx, v)
			if err != nil {
				return fmt.Errorf("error to save port: %w", err)
			}
		}
	}

	log.Info(fmt.Sprintf("%v ports imported", count))

	return nil
}
