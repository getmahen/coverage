package validators

import (
	"context"
	"net/http"
	"regexp"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/rs/zerolog/log"
)

var zipCodeRegex = regexp.MustCompile("^(\\d{5})?$")

type CoverageCheckValidator interface {
	Validate(ctx context.Context, r *http.Request) []entity.Error
}
type coverageCheckValidator struct {
}

func NewCoverageCheckValidator() CoverageCheckValidator {
	return coverageCheckValidator{}
}

//*** OLD CODE
// func (v coverageCheckValidator) Validate(ctx context.Context, zipCode string, carrierID CarrierIDType) []entity.Error {
// 	var validationErrors []entity.Error

// 	isZipCodeValid := zipCodeRegex.MatchString(zipCode)
// 	log.Ctx(ctx).Debug().Bool("regexCheck", isZipCodeValid).Str("zipCode", zipCode)

// 	if !isZipCodeValid {
// 		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "zipcode"})
// 	}

// 	isValidCarrierID := false
// 	switch carrierID {
// 	case SPRINT:
// 		isValidCarrierID = true
// 		log.Ctx(ctx).Debug().Bool("carrierIDCheck", isValidCarrierID).Interface("carrierID", carrierID).Interface("SPRINT", carrierID)
// 	case VERIZON:
// 		isValidCarrierID = true
// 		log.Ctx(ctx).Debug().Bool("carrierIDCheck", isValidCarrierID).Interface("carrierID", carrierID).Interface("VERIZON", carrierID)
// 	default:
// 		log.Ctx(ctx).Debug().Interface("Invalid Carrier ID", carrierID)
// 		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "carrierid"})
// 	}
// 	return validationErrors
// }
//*** OLD CODE

// func (v coverageCheckValidator) Validate(ctx context.Context, zipCode string, carrierID string) []entity.Error {
// 	var validationErrors []entity.Error

// 	isZipCodeValid := zipCodeRegex.MatchString(zipCode)
// 	log.Ctx(ctx).Debug().Bool("regexCheck", isZipCodeValid).Str("zipCode", zipCode)

// 	if !isZipCodeValid {
// 		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "zipcode"})
// 	}

// 	isValidCarrierID := false
// 	switch carrierID {
// 	case "1":
// 		isValidCarrierID = true
// 		log.Ctx(ctx).Debug().Bool("carrierIDCheck", isValidCarrierID).Interface("carrierID", carrierID).Interface("SPRINT", carrierID)
// 	case "2":
// 		isValidCarrierID = true
// 		log.Ctx(ctx).Debug().Bool("carrierIDCheck", isValidCarrierID).Interface("carrierID", carrierID).Interface("VERIZON", carrierID)
// 	default:
// 		log.Ctx(ctx).Debug().Interface("Invalid Carrier ID", carrierID)
// 		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "carrierid"})
// 	}
// 	return validationErrors
// }

func (v coverageCheckValidator) Validate(ctx context.Context, r *http.Request) []entity.Error {
	var validationErrors []entity.Error

	zipCode := r.URL.Query().Get("zipcode")
	if zipCode == "" {
		validationErrors = append(validationErrors, entity.Error{Message: "Missing required property", Path: "zipcode"})
	}

	carrierID := r.URL.Query().Get("carrierid")
	if carrierID == "" {
		validationErrors = append(validationErrors, entity.Error{Message: "Missing required property", Path: "carrierid"})
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	isZipCodeValid := zipCodeRegex.MatchString(zipCode)
	log.Ctx(ctx).Debug().Bool("regexCheck", isZipCodeValid).Str("zipCode", zipCode)

	if !isZipCodeValid {
		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "zipcode"})
	}

	isValidCarrierID := false
	switch entity.CarrierType(carrierID) {
	case entity.Sprint:
		isValidCarrierID = true
		log.Ctx(ctx).Debug().Bool("carrierIDCheck", isValidCarrierID).Interface("carrierID", carrierID).Interface("SPRINT", carrierID)
	case entity.Verizon:
		isValidCarrierID = true
		log.Ctx(ctx).Debug().Bool("carrierIDCheck", isValidCarrierID).Interface("carrierID", carrierID).Interface("VERIZON", carrierID)
	default:
		log.Ctx(ctx).Debug().Interface("Invalid Carrier ID", carrierID)
		validationErrors = append(validationErrors, entity.Error{Message: "Illegal value for property", Path: "carrierid"})
	}
	return validationErrors
}
