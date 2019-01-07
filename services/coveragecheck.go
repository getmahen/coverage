package services

import (
	"context"

	"bitbucket.org/credomobile/coverage/dbclient"
	"bitbucket.org/credomobile/coverage/entity"
	"github.com/rs/zerolog"
)

type CoverageCheck interface {
	Verify(ctx context.Context, zipCode string, carrierID string) (entity.CoverageCheckResponse, error)
}

type coverageCheck struct {
	dbclientFactory dbclient.ClientFactory
}

func NewCoverageCheck(dbclientFactory dbclient.ClientFactory) CoverageCheck {
	//xray.AWS(dynamo.Client)
	return coverageCheck{
		dbclientFactory: dbclientFactory,
	}
}

func (c coverageCheck) Verify(ctx context.Context, zipCode string, carrierID string) (entity.CoverageCheckResponse, error) {
	zerolog.Ctx(ctx).Info().Msgf("Verifying coverage for zipcode: %s and carrierID: %s", zipCode, carrierID)

	dbclient, err := c.dbclientFactory.GetDbClient(entity.CarrierType(carrierID))
	if err != nil {
		return entity.CoverageCheckResponse{}, err
	}

	isCovered, err := dbclient.VerifyCoverage(ctx, zipCode, carrierID)
	if err != nil {
		return entity.CoverageCheckResponse{}, err
	}

	return entity.CoverageCheckResponse{IsCovered: isCovered}, nil
}
