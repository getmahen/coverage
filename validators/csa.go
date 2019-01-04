package validators

import (
	"context"
	"net/http"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/rs/zerolog/log"
)

type CsaValidator interface {
	//Validate(ctx context.Context, zipCode string) []entity.Error
	Validate(ctx context.Context, r *http.Request) []entity.Error
}
type csaValidator struct {
}

func NewCsaValidator() CsaValidator {
	return csaValidator{}
}

//** OLD CODE
// func (v csaValidator) Validate(ctx context.Context, zipCode string) []entity.Error {
// 	var validationErrors []entity.Error

// 	isZipCodeValid := zipCodeRegex.MatchString(zipCode)
// 	log.Ctx(ctx).Debug().Bool("regexCheck", isZipCodeValid).Str("zipCode", zipCode)

// 	if !isZipCodeValid {
// 		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "zipcode"})
// 	}

// 	return validationErrors
// }

func (v csaValidator) Validate(ctx context.Context, r *http.Request) []entity.Error {
	var validationErrors []entity.Error

	zipCode := r.URL.Query().Get("zipcode")
	if zipCode == "" {
		validationErrors = append(validationErrors, entity.Error{Message: "Missing required property", Path: "zipcode"})
		return validationErrors
	}

	isZipCodeValid := zipCodeRegex.MatchString(zipCode)
	if !isZipCodeValid {
		log.Ctx(ctx).Debug().Bool("regexCheck", isZipCodeValid).Str("zipCode", zipCode)
		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "zipcode"})
		return validationErrors
	}
	return validationErrors
}
