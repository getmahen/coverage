package dbclient

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

func TestSprintVerifyCoverage(t *testing.T) {
	testCases := []struct {
		desc                  string
		zipCode               string
		expectZipCodeCovered  bool
		causeDynamoDbError    bool
		dynamodbReturnPayload map[string]string
	}{
		{
			desc:                 "happy path with Zip code that has coverage",
			zipCode:              "94105",
			expectZipCodeCovered: true,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94105",
				"carriertype":   "sprint",
				"csa_leaf":      "fakeCsa",
				"cur_pct_cov":   "100",
				"lte_4g_pctcov": "100",
			},
		},
		{
			desc:                 "happy path with Zip code that has No coverage",
			zipCode:              "94105",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94105",
				"carriertype":   "sprint",
				"csa_leaf":      "fakeCsa",
				"cur_pct_cov":   "50",
				"lte_4g_pctcov": "50",
			},
		},
		{
			desc:                 "no coverage for a zip code with missing information in database",
			zipCode:              "94107",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94107",
				"carriertype":   "sprint",
				"csa_leaf":      "",
				"cur_pct_cov":   "",
				"lte_4g_pctcov": "",
			},
		},
		{
			desc:                 "no coverage for a zip code with invalid cur_pct_cov attribute in database",
			zipCode:              "94107",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94107",
				"carriertype":   "sprint",
				"csa_leaf":      "",
				"cur_pct_cov":   "abc",
				"lte_4g_pctcov": "100",
			},
		},
		{
			desc:                 "results in no coverage for a zip code with invalid lte_4g_pctcov attribute in database",
			zipCode:              "94107",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94107",
				"carriertype":   "sprint",
				"csa_leaf":      "",
				"cur_pct_cov":   "100",
				"lte_4g_pctcov": "abc",
			},
		},
		{
			desc:                 "results in no coverage for a zip code that does not exist in the database",
			zipCode:              "11111",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "",
				"carriertype":   "",
				"csa_leaf":      "",
				"cur_pct_cov":   "",
				"lte_4g_pctcov": "",
			},
		},
		{
			desc:               "Sad path with Zip code that has coverage but results in dynamodb error",
			zipCode:            "94105",
			causeDynamoDbError: true,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94105",
				"carriertype":   "sprint",
				"csa_leaf":      "fakeCsa",
				"cur_pct_cov":   "100",
				"lte_4g_pctcov": "100",
			},
		},
	}

	for _, tC := range testCases {
		tableName := aws.String("fakeCoverage")

		fakeDb := &fakeSprintDynamoDB{}
		if tC.causeDynamoDbError {
			fakeDb = &fakeSprintDynamoDB{
				t:               t,
				tableName:       tableName,
				payloadToReturn: tC.dynamodbReturnPayload,
				err:             errors.New("fake DB error"),
			}
		} else {
			fakeDb = &fakeSprintDynamoDB{
				t:               t,
				tableName:       tableName,
				payloadToReturn: tC.dynamodbReturnPayload,
				err:             nil,
			}
		}

		sprintdbClient := NewSprintClient(tableName, fakeDb)
		result, err := sprintdbClient.VerifyCoverage(context.Background(), tC.zipCode)

		if tC.causeDynamoDbError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tC.expectZipCodeCovered, result)
		}

		assert.Equal(t, tC.zipCode, fakeDb.Keys["zipcode"])
		assert.Equal(t, "sprint", fakeDb.Keys["carriertype"])
	}
}

func TestGetCsa(t *testing.T) {
	testCases := []struct {
		desc                  string
		zipCode               string
		expectCsa             bool
		causeDynamoDbError    bool
		zipCodeExistInDb      bool
		dynamodbReturnPayload map[string]string
	}{
		{
			desc:               "happy path with Zip code that has Csa",
			zipCode:            "94105",
			expectCsa:          true,
			causeDynamoDbError: false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94105",
				"carriertype":   "sprint",
				"csa_leaf":      "fakeCsa",
				"cur_pct_cov":   "100",
				"lte_4g_pctcov": "100",
			},
		},
		{
			desc:               "happy path with Zip code that has no Csa",
			zipCode:            "94105",
			expectCsa:          true,
			causeDynamoDbError: false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94105",
				"carriertype":   "sprint",
				"csa_leaf":      "",
				"cur_pct_cov":   "100",
				"lte_4g_pctcov": "100",
			},
		},
		{
			desc:               "Sad path when there is a failure to retrieve data from database",
			zipCode:            "94105",
			expectCsa:          true,
			causeDynamoDbError: true,
			dynamodbReturnPayload: map[string]string{
				"zipcode":       "94105",
				"carriertype":   "sprint",
				"csa_leaf":      "",
				"cur_pct_cov":   "100",
				"lte_4g_pctcov": "100",
			},
		},
	}
	for _, tC := range testCases {
		tableName := aws.String("fakeCoverage")

		fakeDb := &fakeSprintDynamoDB{}
		if tC.causeDynamoDbError {
			fakeDb = &fakeSprintDynamoDB{
				t:               t,
				tableName:       tableName,
				payloadToReturn: tC.dynamodbReturnPayload,
				err:             errors.New("fake DB error"),
			}
		} else {
			fakeDb = &fakeSprintDynamoDB{
				t:               t,
				tableName:       tableName,
				payloadToReturn: tC.dynamodbReturnPayload,
				err:             nil,
			}
		}

		sprintdbClient := NewSprintClient(aws.String("fakeCoverage"), fakeDb)
		result, err := sprintdbClient.GetCsa(context.Background(), tC.zipCode)

		if tC.causeDynamoDbError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tC.dynamodbReturnPayload["csa_leaf"], result)
		}

		assert.Equal(t, tC.zipCode, fakeDb.Keys["zipcode"])
		assert.Equal(t, "sprint", fakeDb.Keys["carriertype"])
	}
}

type fakeSprintDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	tableName       *string
	Keys            map[string]string
	payloadToReturn map[string]string // Store fake return values
	err             error
	t               *testing.T
}

func (fd *fakeSprintDynamoDB) GetItemWithContext(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	assert.Equal(fd.t, *fd.tableName, *input.TableName, "incorrect table name")

	expectedAttributes := map[string]*string{
		"#0": aws.String("zipcode"),
		"#1": aws.String("carriertype"),
		"#2": aws.String("csa_leaf"),
		"#3": aws.String("cur_pct_cov"),
		"#4": aws.String("lte_4g_pctcov"),
	}

	actual := input.ExpressionAttributeNames
	if e, a := expectedAttributes, actual; !reflect.DeepEqual(a, e) {
		fd.t.Errorf("expect %v, got %v", e, a)
	}

	fd.Keys = make(map[string]string)
	for k, v := range input.Key {
		if v == nil {
			continue
		}

		if v.S != nil {
			fd.Keys[k] = *v.S
		} else if v.N != nil {
			fd.Keys[k] = *v.N
		}
	}

	if fd.err != nil {
		return &dynamodb.GetItemOutput{}, errors.New("Fake dynamodb error")
	} else {
		return &dynamodb.GetItemOutput{Item: map[string]*dynamodb.AttributeValue{
			"zipcode":       &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["zipcode"])},
			"carriertype":   &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["carriertype"])},
			"csa_leaf":      &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["csa_leaf"])},
			"cur_pct_cov":   &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["cur_pct_cov"])},
			"lte_4g_pctcov": &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["lte_4g_pctcov"])},
		}}, nil
	}
}
