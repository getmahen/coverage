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

func TestCoverageCheckHappyPath(t *testing.T) {
	testCases := []struct {
		desc            string
		zipCode         string
		carrierID       string
		respondCovered  bool
		causeTimeout    bool
		statusCode      int
		expectedCovered bool
	}{
		{
			desc:            "Happy path with valid zipcode and carriedID with coverage",
			zipCode:         "94105",
			carrierID:       "1",
			respondCovered:  true,
			causeTimeout:    false,
			statusCode:      http.StatusOK,
			expectedCovered: true,
		},
		{
			desc:            "Happy path with valid zipcode and carriedID with no coverage",
			zipCode:         "94106",
			carrierID:       "2",
			respondCovered:  false,
			causeTimeout:    false,
			statusCode:      http.StatusOK,
			expectedCovered: false,
		},
	}

	for _, tC := range testCases {
		coverageCheckValidator := validators.NewCoverageCheckValidator()
		coveragecheckService := MockCoverageCheck{}

		if tC.respondCovered {
			coveragecheckService.On("Verify", mock.Anything, tC.zipCode, tC.carrierID).Return(entity.CoverageCheckResponse{IsCovered: true}, nil)
		} else {
			coveragecheckService.On("Verify", mock.Anything, tC.zipCode, tC.carrierID).Return(entity.CoverageCheckResponse{IsCovered: false}, nil)
		}

		t.Run(tC.desc, func(t *testing.T) {

			r := chi.NewRouter()
			r.Get("/v1/coveragecheck", CheckCoverage(coverageCheckValidator, &coveragecheckService))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/coveragecheck?zipcode=%s&carrierid=%s", ts.URL, tC.zipCode, tC.carrierID), nil)
			res, err := ts.Client().Do(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)

			body, _ := ioutil.ReadAll(res.Body)
			if tC.expectedCovered {
				assert.Contains(t, string(body), `{"Result":{"IsCovered":true}}`)
			} else {
				assert.Contains(t, string(body), `{"Result":{"IsCovered":false}}`)
			}
			coveragecheckService.AssertExpectations(t)
		})
	}
}

func TestCoverageCheckSadPathValidationErrors(t *testing.T) {
	testCases := []struct {
		desc             string
		zipCode          string
		carrierID        string
		expectedResponse string
		statusCode       int
	}{
		{
			desc:             "Missing zipcode and carriedID",
			zipCode:          "",
			carrierID:        "",
			expectedResponse: `{"Errors":[{"message":"Missing required property","path":"zipcode"},{"message":"Missing required property","path":"carrierid"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Missing zipcode and a valid carriedID",
			zipCode:          "",
			carrierID:        "1",
			expectedResponse: `{"Errors":[{"message":"Missing required property","path":"zipcode"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Invalid zipcode with exceeded length and a valid carriedID",
			zipCode:          "941055",
			carrierID:        "1",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"zipcode"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Invalid zipcode with shortened length and a valid carriedID",
			zipCode:          "9410",
			carrierID:        "1",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"zipcode"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Invalid zipcode and an invalid carriedID",
			zipCode:          "abc",
			carrierID:        "3",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"zipcode"},{"message":"Illegal value for property","path":"carrierid"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Invalid zipcode and an invalid with non number carriedID",
			zipCode:          "abc",
			carrierID:        "a",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"zipcode"},{"message":"Illegal value for property","path":"carrierid"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Valid zipcode and an invalid carriedID",
			zipCode:          "94105",
			carrierID:        "3",
			expectedResponse: `{"Errors":[{"message":"Illegal value for property","path":"carrierid"}]}`,
			statusCode:       http.StatusBadRequest,
		},
		{
			desc:             "Valid zipcode and an missing carriedID",
			zipCode:          "94105",
			carrierID:        "",
			expectedResponse: `{"Errors":[{"message":"Missing required property","path":"carrierid"}]}`,
			statusCode:       http.StatusBadRequest,
		},
	}

	for _, tC := range testCases {
		coverageCheckValidator := validators.NewCoverageCheckValidator()
		coveragecheckService := MockCoverageCheck{}

		t.Run(tC.desc, func(t *testing.T) {

			r := chi.NewRouter()
			r.Get("/v1/coveragecheck", CheckCoverage(coverageCheckValidator, &coveragecheckService))
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/coveragecheck?zipcode=%s&carrierid=%s", ts.URL, tC.zipCode, tC.carrierID), nil)
			res, err := ts.Client().Do(req)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)

			body, _ := ioutil.ReadAll(res.Body)
			assert.Contains(t, string(body), tC.expectedResponse)
			coveragecheckService.AssertExpectations(t)
		})
	}
}

func TestCoverageCheckSadPathInternalServerError(t *testing.T) {
	coverageCheckValidator := validators.NewCoverageCheckValidator()
	coveragecheckService := MockCoverageCheck{}
	zipCode := "94105"
	carrierID := "1"

	coveragecheckService.On("Verify", mock.Anything, zipCode, carrierID).Return(entity.CoverageCheckResponse{}, errors.New("Fake error"))

	r := chi.NewRouter()
	r.Get("/v1/coveragecheck", CheckCoverage(coverageCheckValidator, &coveragecheckService))
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/coveragecheck?zipcode=%s&carrierid=%s", ts.URL, zipCode, carrierID), nil)
	res, err := ts.Client().Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	body, _ := ioutil.ReadAll(res.Body)
	assert.Contains(t, string(body), `{"message":"There is a problem on the server. Please try again later"}`)
	coveragecheckService.AssertExpectations(t)
}

type MockCoverageCheck struct {
	mock.Mock
}

func (c *MockCoverageCheck) Verify(ctx context.Context, zipCode string, carrierID string) (entity.CoverageCheckResponse, error) {
	args := c.Called(ctx, zipCode, carrierID)
	return args.Get(0).(entity.CoverageCheckResponse), errOrNil(args.Get(1))
}

func errOrNil(o interface{}) error {
	if o == nil {
		return nil
	}
	return o.(error)
}
