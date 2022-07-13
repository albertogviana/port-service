package handlers_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/albertogviana/port-service/internal/repository"

	"github.com/albertogviana/port-service/internal/web"

	"github.com/albertogviana/port-service/internal/handlers"
	"github.com/albertogviana/port-service/internal/port"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const routerPath = "/v1/port"

type PortHandlerUnitTestSuite struct {
	suite.Suite
	log *log.Logger
}

func TestPortHandlerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(PortHandlerUnitTestSuite))
}

func (ph *PortHandlerUnitTestSuite) SetupSuite() {
	ph.log = log.New()
	ph.log.SetFormatter(&log.JSONFormatter{})
}

func (ph *PortHandlerUnitTestSuite) TestCreateOrUpdatePort() {
	ph.Run("check create or update port router", func() {
		svc := port.NewService(repository.NewInMemRepository())

		router := mux.NewRouter()
		handlers.MakePortHandlers(
			router,
			svc,
			ph.log,
		)
		createUser := router.GetRoute("createOrUpdatePort")
		path, err := createUser.GetPathTemplate()
		ph.NoError(err)
		ph.Equal(routerPath, path)

		method, err := createUser.GetMethods()
		ph.NoError(err)
		ph.Equal([]string{http.MethodPost}, method)
	})

	type tableTest struct {
		name               string
		port               *port.Port
		payload            string
		expectedStatusCode int
		expectedResponse   *web.ErrorResponse
	}

	tt := []tableTest{
		{
			name: "create a port",
			port: nil,
			payload: `{
				"name": "Ajman",
				"city": "Ajman",
				"country": "United Arab Emirates",
				"alias": [],
				"regions": [],
				"coordinates": [
				  55.5136433,
				  25.4052165
				],
				"province": "Ajman",
				"timezone": "Asia/Dubai",
				"unlocs": [
				  "AEAJM"
				],
				"code": "52000"
			}`,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   nil,
		},
		{
			name: "update a port",
			port: &port.Port{
				Name:        "Abu Dhabi 2",
				City:        "Abu Dhabi",
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
			payload: `{
				"name": "Abu Dhabi",
				"coordinates": [
				  54.37,
				  24.47
				],
				"city": "Abu Dhabi",
				"province": "Abu Z¸aby [Abu Dhabi]",
				"country": "United Arab Emirates",
				"alias": [],
				"regions": [],
				"timezone": "Asia/Dubai",
				"unlocs": [
				  "AEAUH"
				],
				"code": "52001"
			}`,
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   nil,
		},
		{
			name: "invalid json",
			port: nil,
			payload: `{
				"name": "Abu Dhabi",
				"coordinates": [
				  54.37,
				  24.47
				],
				"city": "Abu Dhabi",
				"province": "Abu Z¸aby [Abu Dhabi]",
				"country": "United Arab Emirates",
				"alias": [],
				"regions": [],
				"timezone": "Asia/Dubai",
				"unlocs": [
				  "AEAUH"
				]
				"code": "52001"
			}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &web.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "failed to decode json",
			},
		},
	}

	for _, tc := range tt {
		tc := tc

		ph.Run(tc.name, func() {
			svc := port.NewService(repository.NewInMemRepository())

			if tc.port != nil {
				err := svc.SavePort(context.Background(), tc.port)
				ph.NoError(err)
			}

			router := mux.NewRouter()
			handlers.MakePortHandlers(
				router,
				svc,
				ph.log,
			)

			ts := httptest.NewServer(router)
			defer ts.Close()

			client := &http.Client{}
			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodPost,
				ts.URL+routerPath,
				strings.NewReader(tc.payload),
			)
			ph.NoError(err)

			resp, err := client.Do(req)
			ph.NoError(err)
			defer resp.Body.Close()
			ph.Equal(tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedResponse != nil {
				expectedResult, err := json.Marshal(tc.expectedResponse)
				ph.NoError(err)

				respBody, err := ioutil.ReadAll(resp.Body)
				ph.NoError(err)

				ph.Equal(expectedResult, respBody)
			}
		})
	}
}

func (ph *PortHandlerUnitTestSuite) TestGetPortByUnloc() {
	ph.Run("check get port by unloc router", func() {
		svc := port.NewService(repository.NewInMemRepository())

		router := mux.NewRouter()
		handlers.MakePortHandlers(
			router,
			svc,
			ph.log,
		)
		createUser := router.GetRoute("getPortByUnloc")
		path, err := createUser.GetPathTemplate()
		ph.NoError(err)
		ph.Equal("/v1/port/{unloc:[0-9A-Z]+}", path)

		method, err := createUser.GetMethods()
		ph.NoError(err)
		ph.Equal([]string{http.MethodGet}, method)
	})

	type tableTest struct {
		name               string
		port               *port.Port
		unloc              string
		expectedStatusCode int
		expectedResponse   *web.ErrorResponse
	}

	tt := []tableTest{
		{
			name: "get a port",
			port: &port.Port{
				Name:        "Abu Dhabi 2",
				City:        "Abu Dhabi",
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
			unloc:              "AEAJM",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   nil,
		},
		{
			name:               "port was not found",
			port:               nil,
			unloc:              "AEAJM",
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: &web.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "port with unloc AEAJM was not found",
			},
		},
	}

	for _, tc := range tt {
		tc := tc

		ph.Run(tc.name, func() {
			svc := port.NewService(repository.NewInMemRepository())

			if tc.port != nil {
				err := svc.SavePort(context.Background(), tc.port)
				ph.NoError(err)
			}

			router := mux.NewRouter()
			handlers.MakePortHandlers(
				router,
				svc,
				ph.log,
			)

			ts := httptest.NewServer(router)
			defer ts.Close()

			client := &http.Client{}
			req, err := http.NewRequestWithContext(
				context.Background(),
				http.MethodGet,
				ts.URL+routerPath+"/"+tc.unloc,
				nil,
			)
			ph.NoError(err)

			resp, err := client.Do(req)
			ph.NoError(err)
			defer resp.Body.Close()
			ph.Equal(tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedResponse != nil {
				expectedResult, err := json.Marshal(tc.expectedResponse)
				ph.NoError(err)

				respBody, err := ioutil.ReadAll(resp.Body)
				ph.NoError(err)

				ph.Equal(expectedResult, respBody)
			}
		})
	}
}
