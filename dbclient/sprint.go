package dbclient

import (
	"context"
	"fmt"
	"strconv"

	"bitbucket.org/credomobile/coverage/entity"
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
	logger     *zerolog.Logger
	connection dynamodbiface.DynamoDBAPI
}

type sprintCoverageData struct {
	ZipCode        string `json:"zipcode"`
	CarrierType    string `json:"carriertype"`
	ZipPostalCity  string `json:"zip_postal_city"`
	State          string `json:"state"`
	ZipCodeArea    string `json:"zipcode_area"`
	MktName        string `json:"mkt_name"`
	CsaLeaf        string `json:"csa_leaf"`
	CurPctCov      string `json:"cur_pct_cov"`
	CurEvdoPctCov  string `json:"cur_evdo_pct_cov"`
	Roam1XPctCov   string `json:"roam1x_pct_cov"`
	EVDORoamPctCov string `json:"evdoroam_pct_cov"`
	CDMARoamPctCov string `json:"cdmaroam_pct_cov"`
	Lte4GPctCov    string `json:"lte_4g_pctcov"`
	Lte2500PctCov  string `json:"lte_2500_PctCov"`
	ZipCenterLon   string `json:"zip_center_lon"`
	ZipCenterLat   string `json:"zip_center_lat"`
}

func NewSprintClient(logger *zerolog.Logger, connection dynamodbiface.DynamoDBAPI) sprintDbClient {
	return sprintDbClient{logger: logger, connection: connection}
}

func (s sprintDbClient) VerifyCoverage(ctx context.Context, zipCode string) (bool, error) {
	s.logger.Info().Msgf("*** IN SPRINT DB CLIENT VerifyCoverage() ***")

	data, err := s.getData(ctx, zipCode)
	if err != nil {
		return false, err
	}
	fmt.Println("CsaLeaf: ", data.CsaLeaf)
	fmt.Println("CurPctCov: ", data.CurPctCov)
	fmt.Println("LTE4GPctCov: ", data.Lte4GPctCov)

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
		TableName: aws.String(entity.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"zipcode": {
				S: aws.String(zipCode),
			},
			"carriertype": {
				S: aws.String("sprint"),
			},
		},
	}

	result, err := s.connection.GetItemWithContext(ctx, input)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to query coverage dynamodb table")
		return sprintCoverageData{}, err
	}

	item := sprintCoverageData{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to UnmarshalMap data from dynamodb")
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.ZipCode == "" {
		s.logger.Debug().Msgf("Could not find coverage data for zipcode: %s", zipCode)
		return sprintCoverageData{}, nil
	}

	// var data sprintCoverageData
	// json.Unmarshal([]byte(item.JsonData), &data)
	return item, nil
}
func (s sprintDbClient) isZipCovered(zipCode string, data sprintCoverageData) bool {
	if len(data.CurPctCov) == 0 || len(data.Lte4GPctCov) == 0 {
		s.logger.Debug().Msgf("zipcode: %s not covered as either Cur_Pct_Cov, LTE_4G_PctCov fields are empty", zipCode)
		return false
	}

	curPctCov, err := strconv.ParseFloat(data.CurPctCov, 64)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Illegal value in Cur_Pct_Cov")
		return false
	}

	lTE4GPctCov, err := strconv.ParseFloat(data.Lte4GPctCov, 64)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Illegal value in LTE_4G_PctCov")
		return false
	}

	if lTE4GPctCov > 50 && curPctCov > 50 {
		return true
	}
	return false
}
