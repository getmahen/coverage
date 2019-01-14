package dbclient

import (
	"context"
	"errors"
	"strings"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type CoverageCheckClient interface {
	VerifyCoverage(ctx context.Context, zipCode string) (bool, error)
}

type ClientFactory interface {
	GetDbClient(t entity.CarrierType) (CoverageCheckClient, error)
}

type clientFactoryImpl struct {
	tableName  *string
	connection dynamodbiface.DynamoDBAPI
}

// NewDbClientFactory constructs and gives back a db client factory that can be used to retrieve carrier specfic db client.
func NewDbClientFactory(dynamodbARN string, logger *zerolog.Logger) (ClientFactory, error) {
	//awsSession, err := session.NewSession()
	config := &aws.Config{
		Region:   aws.String("us-east-2"),
		Endpoint: aws.String("http://localhost:8000"),
	}
	awsSession, err := session.NewSession(config)

	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create connection to dynamodb")
		return nil, err
	}
	dynamo := dynamodb.New(awsSession)

	if len(strings.Split(dynamodbARN, "/")) < 2 {
		return nil, errors.New("Invalid dynamodbARN")
	}

	return clientFactoryImpl{
		tableName:  aws.String(strings.Split(dynamodbARN, "/")[1]),
		connection: dynamodbiface.DynamoDBAPI(dynamo),
	}, nil
}

func (c clientFactoryImpl) GetDbClient(t entity.CarrierType) (CoverageCheckClient, error) {
	switch t {
	case entity.Sprint:
		return NewSprintClient(c.tableName, c.connection), nil
	case entity.Verizon:
		return NewVerizonClient(c.tableName, c.connection), nil
	default:
		//if type is invalid, return an error
		return nil, errors.New("Invalid Carrier Type")
	}
}
