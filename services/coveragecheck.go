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

//ORIGINAL CODE *****
// type coverageCheck struct {
// 	tableName  *string
// 	logger     *zerolog.Logger
// 	connection dynamodbiface.DynamoDBAPI
// }

// func NewCoverageCheck(dynamodbARN string, logger *zerolog.Logger) CoverageCheck {
// 	awsSession, err := session.NewSession()
// 	if err != nil {
// 		logger.Fatal().Err(err).Msg("unable to create connection to dynamodb")
// 	}
// 	dynamo := dynamodb.New(awsSession)
// 	//xray.AWS(dynamo.Client)

// 	return coverageCheck{
// 		tableName:  aws.String(strings.Split(dynamodbARN, "/")[1]),
// 		connection: dynamodbiface.DynamoDBAPI(dynamo),
// 		logger:     logger,
// 	}
// }

// func (c coverageCheck) Verify(ctx context.Context, zipCode string, carrierID string) (entity.CoverageCheckResponse, error) {
// 	c.logger.Info().Msgf("Verifying coverage for zipcode: %s and carrierID: %s", zipCode, carrierID)

// 	return entity.CoverageCheckResponse{IsCovered: true}, nil

// }
//ORIGINAL CODE *****

//NEW CODE WITH FACTORY PATTERN
type coverageCheck struct {
	logger          *zerolog.Logger
	dbclientFactory dbclient.ClientFactory
}

func NewCoverageCheck(dbclientFactory dbclient.ClientFactory, logger *zerolog.Logger) CoverageCheck {
	//xray.AWS(dynamo.Client)
	return coverageCheck{
		logger:          logger,
		dbclientFactory: dbclientFactory,
	}
}

func (c coverageCheck) Verify(ctx context.Context, zipCode string, carrierID string) (entity.CoverageCheckResponse, error) {
	c.logger.Info().Msgf("Verifying coverage for zipcode: %s and carrierID: %s", zipCode, carrierID)

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
