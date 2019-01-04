package handlers

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/credomobile/coverage/entity"
	"bitbucket.org/credomobile/coverage/services"
	"bitbucket.org/credomobile/coverage/validators"
	"github.com/rs/zerolog/log"
)

func GetCsa(validator validators.CsaValidator, csaService services.Csa) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var validationErrors []entity.Error

		// zipCode := r.URL.Query().Get("zipcode")
		// if zipCode == "" {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(entity.Error{Message: "Missing required property", Path: "zipcode"})
		// }

		// validationErrors = validator.Validate(r.Context(), zipCode)
		// if len(validationErrors) > 0 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(entity.Response{Errors: validationErrors})
		// 	return
		// }

		validationErrors = validator.Validate(r.Context(), r)
		if len(validationErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.Response{Errors: validationErrors})
			return
		}

		ctx := r.Context()
		zipCode := r.URL.Query().Get("zipcode")
		response, err := csaService.GetCsa(ctx, zipCode)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("Error occurred getting csa for zipcode: %s", zipCode)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.Error{Message: "There is a problem on the server. Please try again later"})
			return
		}

		result, _ := json.Marshal(entity.Response{Result: response})
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}
