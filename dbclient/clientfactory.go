package dbclient

import (
	"context"
	"errors"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type Client interface {
	VerifyCoverage(ctx context.Context, zipCode string, carrierID string) (bool, error)
}

type ClientFactory interface {
	GetDbClient(t entity.CarrierType) (Client, error)
}

type clientFactoryImpl struct {
	logger     *zerolog.Logger
	connection dynamodbiface.DynamoDBAPI
}

func NewDbClientFactory(dynamodbARN string, logger *zerolog.Logger) ClientFactory {
	//awsSession, err := session.NewSession()
	config := &aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	}
	awsSession, err := session.NewSession(config)

	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create connection to dynamodb")
	}
	dynamo := dynamodb.New(awsSession)
	return clientFactoryImpl{
		connection: dynamodbiface.DynamoDBAPI(dynamo),
		logger:     logger,
	}
}

func (c clientFactoryImpl) GetDbClient(t entity.CarrierType) (Client, error) {
	switch t {
	case entity.Sprint:
		//return SprintDbClient{TableName: "sprint_coverage", Logger: c.logger, Connection: c.connection}, nil
		return NewSprintClient("Music", c.logger, c.connection), nil
	case entity.Verizon:
		return VerizonDbClient{TableName: "verizon_coverage", Logger: c.logger, Connection: c.connection}, nil
	default:
		//if type is invalid, return an error
		return nil, errors.New("Invalid Carrier Type")
	}
}
