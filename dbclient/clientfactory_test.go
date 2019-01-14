package dbclient

import (
	"os"
	"testing"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewDbClientFactory(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()
	dbClientFactory, err := NewDbClientFactory("abc/fakeDyanamoDbArn", &logger)

	assert.NoError(t, err)
	assert.Implements(t, (*ClientFactory)(nil), dbClientFactory)
}

func TestNewDbClientFactoryForInvalidDynamodbArn(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()
	_, err := NewDbClientFactory("fakeDyanamoDbArn", &logger)

	assert.Error(t, err)
}

func TestGetDbClientForSprint(t *testing.T) {
	dbClientFactory := clientFactoryImpl{tableName: aws.String("fakeCoverage"), connection: &fakeDynamoDB{}}

	//Sprint Db Client
	dbClient, err := dbClientFactory.GetDbClient(entity.CarrierType("1"))
	assert.NoError(t, err)
	assert.IsType(t, sprintDbClient{}, dbClient)
}

func TestGetDbClientForVerizon(t *testing.T) {
	dbClientFactory := clientFactoryImpl{tableName: aws.String("fakeCoverage"), connection: &fakeDynamoDB{}}

	//Verizon Db Client
	dbClient, err := dbClientFactory.GetDbClient(entity.CarrierType("2"))
	assert.NoError(t, err)
	assert.IsType(t, verizonDbClient{}, dbClient)
}

func TestGetDbClientForInvalidCarrierID(t *testing.T) {
	dbClientFactory := clientFactoryImpl{tableName: aws.String("fakeCoverage"), connection: &fakeDynamoDB{}}

	//Verizon Db Client
	dbClient, err := dbClientFactory.GetDbClient(entity.CarrierType("5"))
	assert.Error(t, err)
	assert.Nil(t, dbClient)
}

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	Keys            map[string]string
	payloadToReturn map[string]string // Store fake return values
	err             error
	t               *testing.T
}
