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

type SprintCsaDbClient interface {
	GetCsa(ctx context.Context, zipCode string) (string, error)
}

type sprintDbClient struct {
	tableName  string
	logger     *zerolog.Logger
	connection dynamodbiface.DynamoDBAPI
}

type sprintCoverageData struct {
	CsaLeaf     string `json:"CSA_Leaf"`
	CurPctCov   string `json:"Cur_Pct_Cov"`
	LTE4GPctCov string `json:"LTE_4G_PctCov"`
}

func NewSprintClient(tableName string, logger *zerolog.Logger, connection dynamodbiface.DynamoDBAPI) sprintDbClient {
	return sprintDbClient{tableName: tableName, logger: logger, connection: connection}
}

func (s sprintDbClient) VerifyCoverage(ctx context.Context, zipCode string, carrierID string) (bool, error) {
	s.logger.Info().Msgf("*** IN SPRINT DB CLIENT VerifyCoverage() ***")

	data, err := s.getData(ctx, zipCode)
	if err != nil {
		return false, nil
	}
	fmt.Println("CsaLeaf: ", data.CsaLeaf)
	fmt.Println("CurPctCov: ", data.CurPctCov)
	fmt.Println("LTE4GPctCov: ", data.LTE4GPctCov)

	covered := s.isZipCovered(zipCode, data)
	return covered, nil
}

func (s sprintDbClient) GetCsa(ctx context.Context, zipCode string) (string, error) {
	s.logger.Info().Msgf("*** IN SPRINT DB CLIENT GetCsa() for zipcode %s***", zipCode)

	data, err := s.getData(ctx, zipCode)
	if err != nil {
		return "", nil
	}
	return data.CsaLeaf, nil
}

func (s sprintDbClient) getData(ctx context.Context, zipCode string) (sprintCoverageData, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ZIP": {
				S: aws.String(zipCode),
			},
		},
	}

	result, err := s.connection.GetItemWithContext(ctx, input)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to query dynamodb")
		return sprintCoverageData{}, err
	}

	item := item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to UnmarshalMap data from dynamodb")
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.Zip == "" {
		s.logger.Debug().Msgf("Could not find coverage data for zipcode: %s", zipCode)
		return sprintCoverageData{}, nil
	}

	var data sprintCoverageData
	json.Unmarshal([]byte(item.JsonData), &data)
	return data, nil
}
func (s sprintDbClient) isZipCovered(zipCode string, data sprintCoverageData) bool {
	if len(data.CurPctCov) == 0 || len(data.LTE4GPctCov) == 0 {
		s.logger.Debug().Msgf("zipcode: %s not covered as either Cur_Pct_Cov, LTE_4G_PctCov fields are empty", zipCode)
		return false
	}

	curPctCov, err := strconv.ParseFloat(data.CurPctCov, 64)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Illegal value in Cur_Pct_Cov")
		return false
	}

	lTE4GPctCov, err := strconv.ParseFloat(data.LTE4GPctCov, 64)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Illegal value in LTE_4G_PctCov")
		return false
	}

	if lTE4GPctCov > 50 && curPctCov > 50 {
		return true
	}
	return false
}
