package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/albertogviana/port-service/internal/entity"
	logger "github.com/albertogviana/port-service/internal/log"
	"github.com/albertogviana/port-service/internal/port"
	"github.com/urfave/cli/v2"
	"os"
)

// Import displays the latest state running in the edge node
func Import(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	svc := port.NewService(port.NewRepositoryInMem())

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

	count := 0
	for d.More() {
		m := make(map[string]*entity.Port, 1)

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
