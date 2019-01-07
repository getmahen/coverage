package validators

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/stretchr/testify/assert"
)

func TestCsaValidator(t *testing.T) {
	testCases := []struct {
		desc             string
		zipCode          string
		expectError      bool
		expectedResponse []entity.Error
	}{
		{
			desc:             "Validates a valid zipcode",
			zipCode:          "94105",
			expectError:      false,
			expectedResponse: nil,
		},
		{
			desc:             "Validates a missing zipcode",
			zipCode:          "",
			expectError:      true,
			expectedResponse: []entity.Error{entity.Error{Message: "Missing required property", Path: "zipcode"}},
		},
		{
			desc:             "Validates an Invalid zipcode",
			zipCode:          "abc",
			expectError:      true,
			expectedResponse: []entity.Error{entity.Error{Message: "Illegal value for property", Path: "zipcode"}},
		},
		{
			desc:             "Validates an Invalid zipcode with exceeded length",
			zipCode:          "941050",
			expectError:      true,
			expectedResponse: []entity.Error{entity.Error{Message: "Illegal value for property", Path: "zipcode"}},
		},
		{
			desc:             "Validates an Invalid zipcode with short length",
			zipCode:          "9410",
			expectError:      true,
			expectedResponse: []entity.Error{entity.Error{Message: "Illegal value for property", Path: "zipcode"}},
		},
	}

	for _, tC := range testCases {

		t.Run(tC.desc, func(t *testing.T) {

			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/csa?zipcode=%s", "fakeUrlBasePath", tC.zipCode), nil)

			validator := NewCsaValidator()
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
