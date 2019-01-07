package dbclient

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type verizonDbClient struct {
	tableName  string
	logger     *zerolog.Logger
	connection dynamodbiface.DynamoDBAPI
}

type item struct {
	Zip      string `json:"ZIP"`
	JsonData string `json:"JSON_DATA"`
}

type verizonCoverageData struct {
	VzwLte    string `json:"VZW_LTE"`
	VzwLteInd string `json:"VZW_LTE_IND"`
	State     string `json:"STATE"`
	PoName    string `json:"PO_NAME"`
}

func NewVerizonClient(tableName string, logger *zerolog.Logger, connection dynamodbiface.DynamoDBAPI) verizonDbClient {
	return verizonDbClient{tableName: tableName, logger: logger, connection: connection}
}

func (v verizonDbClient) VerifyCoverage(ctx context.Context, zipCode string, carrierID string) (bool, error) {
	v.logger.Info().Msgf("*** IN VERIZON DB CLIENT ***")

	input := &dynamodb.GetItemInput{
		TableName: aws.String(v.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ZIP": {
				S: aws.String(zipCode),
			},
		},
	}

	result, err := v.connection.GetItemWithContext(ctx, input)
	if err != nil {
		v.logger.Fatal().Err(err).Msg("failed to query dynamodb")
		return false, err
	}

	item := item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		v.logger.Fatal().Err(err).Msg("failed to UnmarshalMap data from dynamodb")
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.Zip == "" {
		v.logger.Debug().Msgf("Could not find coverage for zipcode: %s", zipCode)
		return false, nil
	}

	var data verizonCoverageData
	json.Unmarshal([]byte(item.JsonData), &data)

	covered := v.isZipCovered(zipCode, data)
	return covered, nil
}

func (v verizonDbClient) isZipCovered(zipCode string, data verizonCoverageData) bool {
	if len(data.VzwLte) == 0 || len(data.VzwLteInd) == 0 || len(data.State) == 0 {
		v.logger.Debug().Msgf("zipcode: %s not covered as either Vzw_Lte, Vzw_Lte_Ind or State fields are empty", zipCode)
		return false
	}

	vzwlte, err := strconv.ParseFloat(data.VzwLte, 64)
	if err != nil {
		v.logger.Fatal().Err(err).Msg("Illegal value in VZW_LTE")
		return false
	}
	if vzwlte > 50 && data.VzwLteInd == "Y" {
		return true
	}
	return false
}
