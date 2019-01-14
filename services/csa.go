package services

import (
	"context"
	"errors"
	"strings"

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
	dbClient dbclient.SprintCsaDbClient
}

//NewCsa constructs and gives back csa service
func NewCsa(dynamodbARN string, logger *zerolog.Logger) (csa, error) {
	//xray.AWS(dynamo.Client)

	//awsSession, err := session.NewSession()
	config := &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String("http://localhost:8000"),
	}
	awsSession, err := session.NewSession(config)

	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create connection to dynamodb")
		return csa{}, err
	}
	dynamo := dynamodb.New(awsSession)

	if len(strings.Split(dynamodbARN, "/")) < 2 {
		return csa{}, errors.New("Invalid dynamodbARN")
	}

	return csa{
		dbClient: dbclient.NewSprintClient(aws.String(strings.Split(dynamodbARN, "/")[1]), dynamodbiface.DynamoDBAPI(dynamo)),
	}, nil
}

func (c csa) GetCsa(ctx context.Context, zipCode string) (entity.CsaResponse, error) {
	zerolog.Ctx(ctx).Info().Msgf("Getting Csa for zipcode: %s", zipCode)

	csa, err := c.dbClient.GetCsa(ctx, zipCode)
	if err != nil {
		return entity.CsaResponse{}, err
	}

	return entity.CsaResponse{CsaFound: len(csa) > 0, Csa: csa}, nil
}
