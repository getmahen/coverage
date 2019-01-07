package dbclient

import (
	"os"
	"testing"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetDbClientForSprint(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()
	dbClientFactory := clientFactoryImpl{connection: &fakeDynamoDB{}, logger: &logger}

	//Sprint Db Client
	dbClient, err := dbClientFactory.GetDbClient(entity.CarrierType("1"))
	assert.NoError(t, err)
	assert.IsType(t, sprintDbClient{}, dbClient)
}

func TestGetDbClientForVerizon(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()
	dbClientFactory := clientFactoryImpl{connection: &fakeDynamoDB{}, logger: &logger}

	//Verizon Db Client
	dbClient, err := dbClientFactory.GetDbClient(entity.CarrierType("2"))
	assert.NoError(t, err)
	assert.IsType(t, verizonDbClient{}, dbClient)
}

func TestGetDbClientForInvalidCarrierID(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()
	dbClientFactory := clientFactoryImpl{connection: &fakeDynamoDB{}, logger: &logger}

	//Verizon Db Client
	dbClient, err := dbClientFactory.GetDbClient(entity.CarrierType("5"))
	assert.Error(t, err)
	assert.Nil(t, dbClient)
}

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}
