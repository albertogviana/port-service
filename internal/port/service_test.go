package port_test

import (
	"context"
	"testing"

	"github.com/albertogviana/port-service/internal/entity"
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
		data  *entity.Port
		input *entity.Port
	}

	tt := []tableTest{
		{
			name: "creating a new port",
			data: nil,
			input: &entity.Port{
				Name:      "Ajman",
				City:      "Ajman",
				Country:   "United Arab Emirates",
				Alias:     []string{},
				Regions:   []string{},
				Latitude:  55.5136433,
				Longitude: 25.4052165,
				Province:  "Ajman",
				Timezone:  "Asia/Dubai",
				Unloc:     "AEAJM",
				Code:      "52000",
			},
		},
		{
			name: "updating a port",
			data: &entity.Port{
				Name:      "Ajman",
				City:      "Ajman",
				Country:   "United Arab Emirates",
				Alias:     []string{},
				Regions:   []string{},
				Latitude:  55.5136433,
				Longitude: 25.4052165,
				Province:  "Ajman",
				Timezone:  "Asia/Dubai",
				Unloc:     "AEAJM",
				Code:      "52000",
			},
			input: &entity.Port{
				Name:      "Ajman 2",
				City:      "Ajman 2",
				Country:   "United Arab Emirates",
				Alias:     []string{},
				Regions:   []string{},
				Latitude:  55.5136433,
				Longitude: 25.4052165,
				Province:  "Ajman",
				Timezone:  "Asia/Dubai",
				Unloc:     "AEAJM",
				Code:      "52000",
			},
		},
	}

	s.T().Parallel()

	for _, tc := range tt {
		tc := tc

		s.Run(tc.name, func() {
			s.T().Parallel()

			svc := port.NewService(
				port.NewRepositoryInMem(),
			)

			if tc.data != nil {
				s.NoError(svc.SavePort(context.Background(), tc.data))
			}

			s.NoError(svc.SavePort(context.Background(), tc.input))
		})
	}
}
