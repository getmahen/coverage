package handlers

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/credomobile/coverage/entity"
	"bitbucket.org/credomobile/coverage/validators"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCsaHappyPath(t *testing.T) {
	testCases := []struct {
		desc           string
		zipCode        string
		respondCovered bool
		causeTimeout   bool
		statusCode     int
		csaFound       bool
		expectedCsa    string
	}{
		{
			desc:           "Happy path with valid zipcode and csa found",
			zipCode:        "94105",
			respondCovered: true,
			causeTimeout:   false,
			statusCode:     http.StatusOK,
			csaFound:       true,
			expectedCsa:    "abc",
		},
		{
			desc:           "Happy path with valid zipcode with no csa",
			zipCode:        "94106",
			respondCovered: false,
			causeTimeout:   false,
			statusCode:     http.StatusOK,
			csaFound:       false,
			expectedCsa:    "",
		},
	}

	for _, tC := range testCases {
		csaValidator := validators.NewCsaValidator()
		csaService := MockCsa{}

		if tC.csaFound {
			csaService.On("GetCsa", mock.Anything, tC.zipCode).Return(entity.CsaResponse{CsaFound: tC.csaFound, Csa: tC.expectedCsa}, nil)
		} else {
			csaService.On("GetCsa", mock.Anything, tC.zipCode).Return(entity.CsaResponse{CsaFound: tC.csaFound, Csa: tC.expectedCsa}, nil)
		}

		t.Run(tC.desc, func(t *testing.T) {

			r := chi.NewRouter()
			r.Get("/v1/csa", GetCsa(csaValidator, &csaService))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/csa?zipcode=%s", ts.URL, tC.zipCode), nil)
			res, err := ts.Client().Do(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			body, _ := ioutil.ReadAll(res.Body)
			if tC.csaFound {
				assert.Contains(t, string(body), `{"Result":{"CsaFound":true,"Csa":"abc"}}`)
			} else {
				assert.Contains(t, string(body), `{"Result":{"CsaFound":false,"Csa":""}}`)
			}
		})
	}
}

func TestGetCsaSadPathValidationErrors(t *testing.T) {
	testCases := []struct {
		desc             string
		zipCode          string
		expectedResponse string
		statusCode       int
	}{
		{
			desc:             "Missing zipcode",
			zipCode:          "",
			expectedResponse: `{"Errors":[{"message":"Missing required property","path":"zipcode"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Invalid zipcode with exceeded length",
			zipCode:          "94105678907",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"zipcode"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Invalid zipcode with shortened length",
			zipCode:          "9410",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"zipcode"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Invalid zipcode with alphanumeric characters",
			zipCode:          "abc",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"zipcode"}]}`,
			statusCode:       http.StatusBadRequest,
		},
	}

	for _, tC := range testCases {
		csaValidator := validators.NewCsaValidator()
		csaService := MockCsa{}

		t.Run(tC.desc, func(t *testing.T) {

			r := chi.NewRouter()
			r.Get("/v1/csa", GetCsa(csaValidator, &csaService))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/csa?zipcode=%s", ts.URL, tC.zipCode), nil)
			res, err := ts.Client().Do(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)

			body, _ := ioutil.ReadAll(res.Body)
			assert.Contains(t, string(body), tC.expectedResponse)
		})
	}
}

func TestGetCsaSadPathInternalServerError(t *testing.T) {
	csaValidator := validators.NewCsaValidator()
	csaService := MockCsa{}
	zipCode := "94105"

	csaService.On("GetCsa", mock.Anything, zipCode).Return(entity.CsaResponse{}, errors.New("Fake error"))

	r := chi.NewRouter()
	r.Get("/v1/csa", GetCsa(csaValidator, &csaService))
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/csa?zipcode=%s", ts.URL, zipCode), nil)
	res, err := ts.Client().Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	body, _ := ioutil.ReadAll(res.Body)
	assert.Contains(t, string(body), `{"message":"There is a problem on the server. Please try again later"}`)
}

type MockCsa struct {
	mock.Mock
}

func (c *MockCsa) GetCsa(ctx context.Context, zipCode string) (entity.CsaResponse, error) {
	args := c.Called(ctx, zipCode)
	return args.Get(0).(entity.CsaResponse), errOrNil(args.Get(1))
}
