package services

import (
	"context"

	"bitbucket.org/credomobile/coverage/dbclient"
	"bitbucket.org/credomobile/coverage/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type Csa interface {
	GetCsa(ctx context.Context, zipCode string) (entity.CsaResponse, error)
}

type csa struct {
	logger   *zerolog.Logger
	dbClient dbclient.SprintCsaDbClient
}

func NewCsa(logger *zerolog.Logger) csa {
	//xray.AWS(dynamo.Client)

	//awsSession, err := session.NewSession()
	config := &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String("http://localhost:8000"),
	}
	awsSession, err := session.NewSession(config)

	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create connection to dynamodb")
	}
	dynamo := dynamodb.New(awsSession)

	return csa{
		logger:   logger,
		dbClient: dbclient.NewSprintClient(logger, dynamodbiface.DynamoDBAPI(dynamo)),
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
