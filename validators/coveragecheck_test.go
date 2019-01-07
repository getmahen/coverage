package validators

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/stretchr/testify/assert"
)

func TestCoverageCheckValidator(t *testing.T) {
	testCases := []struct {
		desc             string
		zipCode          string
		carrierID        string
		expectError      bool
		expectedResponse []entity.Error
	}{
		{
			desc:             "Validates a valid zipcode and carrierid for Sprint",
			zipCode:          "94105",
			carrierID:        "1",
			expectError:      false,
			expectedResponse: nil,
		},
		{
			desc:             "Validates a valid zipcode and carrierid for Verizon",
			zipCode:          "94105",
			carrierID:        "2",
			expectError:      false,
			expectedResponse: nil,
		},
		{
			desc:        "Validates a missing zipcode and carrierid",
			zipCode:     "",
			carrierID:   "",
			expectError: true,
			expectedResponse: []entity.Error{
				entity.Error{Message: "Missing required property", Path: "zipcode"},
				entity.Error{Message: "Missing required property", Path: "carrierid"},
			},
		},
		{
			desc:        "Validates a valid zipcode and missing carrierid",
			zipCode:     "94105",
			carrierID:   "",
			expectError: true,
			expectedResponse: []entity.Error{
				entity.Error{Message: "Missing required property", Path: "carrierid"},
			},
		},
		{
			desc:        "Validates a missing zipcode and valid carrierid",
			zipCode:     "",
			carrierID:   "2",
			expectError: true,
			expectedResponse: []entity.Error{
				entity.Error{Message: "Missing required property", Path: "zipcode"},
			},
		},
		{
			desc:        "Validates a invalid zipcode and an invalid carrierid",
			zipCode:     "941ab",
			carrierID:   "3",
			expectError: true,
			expectedResponse: []entity.Error{
				entity.Error{Message: "Illegal value for property", Path: "zipcode"},
				entity.Error{Message: "Illegal value for property", Path: "carrierid"},
			},
		},
	}

	for _, tC := range testCases {

		t.Run(tC.desc, func(t *testing.T) {

			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/coveragecheck?zipcode=%s&carrierid=%s", "fakeUrlBasePath", tC.zipCode, tC.carrierID), nil)

			validator := NewCoverageCheckValidator()
			response := validator.Validate(context.Background(), req)

			if !tC.expectError {
				assert.Nil(t, response)
			} else {
				assert.NotNil(t, response)
				assert.Equal(t, tC.expectedResponse, response)
			}
		})
	}
}

func TestInvalidZipCodeRegex(t *testing.T) {
	testZipCodes := []string{
		"9410a",
		"abcde",
		"abc",
		"941",
		"941055",
	}
	for _, zipCode := range testZipCodes {
		t.Run(zipCode, func(t *testing.T) {
			result := zipCodeRegex.MatchString(zipCode)
			assert.Equal(t, false, result)
		})
	}
}

func TestValidZipCodeRegex(t *testing.T) {
	result := zipCodeRegex.MatchString("94105")
	assert.Equal(t, true, result)
}
