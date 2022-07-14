package port_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/albertogviana/port-service/internal/port"
	"github.com/albertogviana/port-service/internal/repository"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/suite"
)

type ServiceIntegrationTestSuite struct {
	suite.Suite
	db *sql.DB
}

func TestServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceIntegrationTestSuite))
}

func (s *ServiceIntegrationTestSuite) SetupSuite() {
	// For simplicity the credentials are hardcoded.
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 "localhost",
		DBName:               "ports-test",
		User:                 "root",
		Passwd:               "root",
		AllowNativePasswords: true,
		ParseTime:            true,
		Timeout:              10 * time.Second,
	}

	dbConn, err := sql.Open("mysql", cfg.FormatDSN())
	s.NoError(err)

	dbConn.SetMaxOpenConns(10)
	dbConn.SetMaxIdleConns(10)

	s.NoError(dbConn.Ping())

	s.db = dbConn
}

func (s *ServiceIntegrationTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *ServiceIntegrationTestSuite) TestStoreState() {
	s.Run("create, update and get port", func() {
		s.T().Parallel()

		svc := port.NewService(repository.NewMySQLRepository(s.db))

		p, err := svc.GetPortByUnLoc(context.Background(), "AEAJM")
		s.ErrorIs(err, port.ErrPortNotFound)
		s.Nil(p)

		err = svc.SavePort(context.Background(), &port.Port{
			Name:        "Ajman",
			City:        "Ajman",
			Country:     "United Arab Emirates",
			Alias:       []string{},
			Regions:     []string{},
			Coordinates: []float64{55.5136433, 25.4052165},
			Province:    "Ajman",
			Timezone:    "Asia/Dubai",
			Unlocs: []string{
				"AEAJM",
			},
			Code: "52000",
		})

		s.NoError(err)

		updatePort := &port.Port{
			Name:        "Ajman 2",
			City:        "Ajman",
			Country:     "United Arab Emirates",
			Alias:       []string{},
			Regions:     []string{},
			Coordinates: []float64{55.5136433, 25.4052165},
			Province:    "Ajman",
			Timezone:    "Asia/Dubai",
			Unlocs: []string{
				"AEAJM",
			},
			Code: "52000",
		}

		err = svc.SavePort(context.Background(), updatePort)

		s.NoError(err)

		p, err = svc.GetPortByUnLoc(context.Background(), "AEAJM")
		s.NoError(err)

		s.EqualValues(updatePort, p)
	})
}
