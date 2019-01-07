package handlers

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/credomobile/coverage/entity"
	"bitbucket.org/credomobile/coverage/services"
	"bitbucket.org/credomobile/coverage/validators"
	"github.com/rs/zerolog/log"
)

func CheckCoverage(validator validators.CoverageCheckValidator, coverageCheckService services.CoverageCheck) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var validationErrors []entity.Error

		validationErrors = validator.Validate(r.Context(), r)
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.Response{Errors: validationErrors})
			return
		}

		ctx := r.Context()
		zipCode := r.URL.Query().Get("zipcode")
		carrierID := r.URL.Query().Get("carrierid")
		response, err := coverageCheckService.Verify(ctx, zipCode, carrierID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("Error occurred checking coverage for zipcode: %s and carrierID: %s", zipCode, carrierID)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.Error{Message: "There is a problem on the server. Please try again later"})
			return
		}

		result, _ := json.Marshal(entity.Response{Result: response})
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
