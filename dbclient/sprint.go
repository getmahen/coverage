package dbclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rs/zerolog"
)

type SprintCsaDbClient interface {
	GetCsa(ctx context.Context, zipCode string) (string, error)
}

type sprintDbClient struct {
	logger     *zerolog.Logger
	tableName  *string
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

//NewSprintClient construts and returns Sprint's db client
func NewSprintClient(tableName *string, connection dynamodbiface.DynamoDBAPI) sprintDbClient {
	return sprintDbClient{tableName: tableName, connection: connection}
}

func (s sprintDbClient) VerifyCoverage(ctx context.Context, zipCode string) (bool, error) {
	zerolog.Ctx(ctx).Info().Msgf("*** IN SPRINT DB CLIENT VerifyCoverage() ***")

	data, err := s.getData(ctx, zipCode)
	if err != nil {
		return false, err
	}
	fmt.Println("CsaLeaf: ", data.CsaLeaf)
	fmt.Println("CurPctCov: ", data.CurPctCov)
	fmt.Println("LTE4GPctCov: ", data.Lte4GPctCov)

	covered := s.isZipCovered(ctx, zipCode, data)
	return covered, nil
}

func (s sprintDbClient) GetCsa(ctx context.Context, zipCode string) (string, error) {
	zerolog.Ctx(ctx).Info().Msgf("*** IN SPRINT DB CLIENT GetCsa() for zipcode %s***", zipCode)

	data, err := s.getData(ctx, zipCode)
	if err != nil {
		return "", err
	}
	return data.CsaLeaf, nil
}

func (s sprintDbClient) getData(ctx context.Context, zipCode string) (sprintCoverageData, error) {
	// Get cur_pct_cov and lte_4g_pctcov attributes
	proj := expression.NamesList(expression.Name("zipcode"), expression.Name("carriertype"), expression.Name("csa_leaf"), expression.Name("cur_pct_cov"), expression.Name("lte_4g_pctcov"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to build projection expression to query dynamodb table for Sprint coverage")
		return sprintCoverageData{}, err
	}

	input := &dynamodb.GetItemInput{
		TableName: s.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"zipcode": {
				S: aws.String(zipCode),
			},
			"carriertype": {
				S: aws.String("sprint"),
			},
		},
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}

	result, err := s.connection.GetItemWithContext(ctx, input)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to query coverage dynamodb table")
		return sprintCoverageData{}, err
	}

	item := sprintCoverageData{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to UnmarshalMap Sprint coverage data from dynamodb")
		return sprintCoverageData{}, err
	}

	if item.ZipCode == "" {
		zerolog.Ctx(ctx).Debug().Msgf("Could not find coverage data for zipcode: %s", zipCode)
		return sprintCoverageData{}, nil
	}

	// var data sprintCoverageData
	// json.Unmarshal([]byte(item.JsonData), &data)
	return item, nil
}
func (s sprintDbClient) isZipCovered(ctx context.Context, zipCode string, data sprintCoverageData) bool {
	if len(data.CurPctCov) == 0 || len(data.Lte4GPctCov) == 0 {
		zerolog.Ctx(ctx).Debug().Msgf("zipcode: %s not covered as either Cur_Pct_Cov, LTE_4G_PctCov fields are empty", zipCode)
		return false
	}

	curPctCov, err := strconv.ParseFloat(data.CurPctCov, 64)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Illegal value in Cur_Pct_Cov")
		return false
	}

	lTE4GPctCov, err := strconv.ParseFloat(data.Lte4GPctCov, 64)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Illegal value in LTE_4G_PctCov")
		return false
	}

	if lTE4GPctCov > 50 && curPctCov > 50 {
		return true
	}
	return false
}
