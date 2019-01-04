package services

import (
	"context"

	"bitbucket.org/credomobile/coverage/dbclient"
	"bitbucket.org/credomobile/coverage/entity"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type Csa interface {
	GetCsa(ctx context.Context, zipCode string) (entity.CsaResponse, error)
}

//*** ORIGINAL CODE
// type csa struct {
// 	tableName  *string
// 	logger     *zerolog.Logger
// 	connection dynamodbiface.DynamoDBAPI
// }

// func NewCsa(dynamodbARN string, logger *zerolog.Logger) csa {
// 	awsSession, err := session.NewSession()
// 	if err != nil {
// 		logger.Fatal().Err(err).Msg("unable to create connection to dynamodb")
// 	}
// 	dynamo := dynamodb.New(awsSession)
// 	//xray.AWS(dynamo.Client)

// 	return csa{
// 		tableName:  aws.String(strings.Split(dynamodbARN, "/")[1]),
// 		connection: dynamodbiface.DynamoDBAPI(dynamo),
// 		logger:     logger,
// 	}
// }

// func (c csa) GetCsa(ctx context.Context, zipCode string) (entity.CsaResponse, error) {
// 	c.logger.Info().Msgf("Getting Csa for zipcode: %s", zipCode)

// 	return entity.CsaResponse{CsaFound: true, Csa: "fakeCsa"}, nil
// }
//*** ORIGINAL CODE

type csa struct {
	logger   *zerolog.Logger
	dbClient dbclient.SprintDbClient
}

func NewCsa(logger *zerolog.Logger) csa {
	//xray.AWS(dynamo.Client)

	awsSession, err := session.NewSession()
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create connection to dynamodb")
	}
	dynamo := dynamodb.New(awsSession)

	return csa{
		logger: logger,
		dbClient: dbclient.SprintDbClient{
			TableName:  "sprint_coverage",
			Logger:     logger,
			Connection: dynamodbiface.DynamoDBAPI(dynamo),
		},
	}
}

func (c csa) GetCsa(ctx context.Context, zipCode string) (entity.CsaResponse, error) {
	c.logger.Info().Msgf("Getting Csa for zipcode: %s", zipCode)

	csa, err := c.dbClient.GetCsa(ctx, zipCode)
	if err != nil {
		return entity.CsaResponse{}, err
	}

	return entity.CsaResponse{CsaFound: len(csa) > 0, Csa: csa}, nil
}
