package port_test

import (
	"context"
	"testing"

	"github.com/albertogviana/port-service/internal/repository/memory"

	"github.com/albertogviana/port-service/internal/port"
	"github.com/stretchr/testify/suite"
)

type ServiceUnitTestSuite struct {
	suite.Suite
}

func TestServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceUnitTestSuite))
}

func (s *ServiceUnitTestSuite) TestSavePort() {
	type tableTest struct {
		name  string
		data  *port.Port
		input *port.Port
	}

	tt := []tableTest{
		{
			name: "creating a new port",
			data: nil,
			input: &port.Port{
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
			},
		},
		{
			name: "updating a port",
			data: &port.Port{
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
			},
			input: &port.Port{
				Name:        "Ajman 2",
				City:        "Ajman 2",
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
			},
		},
	}

	s.T().Parallel()

	for _, tc := range tt {
		tc := tc

		s.Run(tc.name, func() {
			s.T().Parallel()

			svc := port.NewService(
				memory.NewInMemRepository(),
			)

			if tc.data != nil {
				s.NoError(svc.SavePort(context.Background(), tc.data))
			}

			s.NoError(svc.SavePort(context.Background(), tc.input))
		})
	}
}
