package dbclient

import (
	"context"
	"strconv"

	"bitbucket.org/credomobile/coverage/entity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rs/zerolog"
)

type verizonDbClient struct {
	logger     *zerolog.Logger
	connection dynamodbiface.DynamoDBAPI
}

type verizonCoverageData struct {
	ZipCode         string `json:"zipcode"`
	CarrierType     string `json:"carriertype"`
	EncZip          string `json:"enc_zip"`
	State           string `json:"state"`
	PoName          string `json:"po_name"`
	VzwVoiceOr1x    string `json:"vzwvoiceor1x"`
	VzwVoiceOr1xInd string `json:"vzw_voice_or_1x_ind"`
	VzwEvdo         string `json:"vzwevdo"`
	VzwEvdoInd      string `json:"vzw_evdo_ind"`
	VzwLte          string `json:"vzelte"`
	VzwLteInd       string `json:"vze_lte_ind"`
	AllLte          string `json:"alltle"`
	AllLteInd       string `json:"all_tle_ind"`
	LoadDate        string `json:"load_date"`
	County          string `json:"county"`
	MtaCode         string `json:"mtacode"`
	MtaName         string `json:"mtaname"`
	BtaCode         string `json:"btacode"`
	BtaName         string `json:"btaname"`
	MsaRsaCode      string `json:"msarsacode"`
	MsaRsaName      string `json:"msarsaname"`
}

func NewVerizonClient(logger *zerolog.Logger, connection dynamodbiface.DynamoDBAPI) verizonDbClient {
	return verizonDbClient{logger: logger, connection: connection}
}

func (v verizonDbClient) VerifyCoverage(ctx context.Context, zipCode string) (bool, error) {
	v.logger.Info().Msgf("*** IN VERIZON DB CLIENT ***")

	proj := expression.NamesList(expression.Name("zipcode"), expression.Name("carriertype"), expression.Name("vzelte"), expression.Name("vze_lte_ind"), expression.Name("state"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		v.logger.Error().Err(err).Msg("failed to build projection expression to query dynamodb table for Verizon coverage")
		return false, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(entity.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"zipcode": {
				S: aws.String(zipCode),
			},
			"carriertype": {
				S: aws.String("verizon"),
			},
		},
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	}

	result, err := v.connection.GetItemWithContext(ctx, input)
	if err != nil {
		v.logger.Error().Err(err).Msg("failed to query dynamodb")
		return false, err
	}

	item := verizonCoverageData{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		v.logger.Error().Err(err).Msg("failed to UnmarshalMap Verizon coverage data from dynamodb")
		return false, err
	}

	if item.ZipCode == "" {
		v.logger.Debug().Msgf("Could not find coverage for zipcode: %s", zipCode)
		return false, nil
	}

	covered := v.isZipCovered(zipCode, item)
	return covered, nil
}

func (v verizonDbClient) isZipCovered(zipCode string, data verizonCoverageData) bool {
	if len(data.VzwLte) == 0 || len(data.VzwLteInd) == 0 || len(data.State) == 0 {
		v.logger.Debug().Msgf("zipcode: %s not covered as either vzwlte, vzw_lte_ind or state fields are empty", zipCode)
		return false
	}

	vzwlte, err := strconv.ParseFloat(data.VzwLte, 64)
	if err != nil {
		v.logger.Error().Err(err).Msg("Illegal value in vzwlte")
		return false
	}
	if vzwlte > 50 && data.VzwLteInd == "Y" {
		return true
	}
	return false
}
