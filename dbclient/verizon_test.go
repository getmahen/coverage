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

func TestVerizonVerifyCoverage(t *testing.T) {
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
				"zipcode":     "94105",
				"carriertype": "verizon",
				"vzelte":      "100",
				"vze_lte_ind": "Y",
				"state":       "CA",
			},
		},
		{
			desc:                 "happy path with Zip code that has No coverage when vzelte attribute less than threshold",
			zipCode:              "94105",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":     "94105",
				"carriertype": "sprint",
				"vzelte":      "50",
				"vze_lte_ind": "Y",
				"state":       "CA",
			},
		},
		{
			desc:                 "happy path with Zip code that has No coverage when vze_lte_ind attribute set to N",
			zipCode:              "94105",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":     "94105",
				"carriertype": "verizon",
				"vzelte":      "100",
				"vze_lte_ind": "N",
				"state":       "CA",
			},
		},
		{
			desc:                 "no coverage for a zip code with missing information in database",
			zipCode:              "94107",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":     "94105",
				"carriertype": "verizon",
				"vzelte":      "",
				"vze_lte_ind": "",
				"state":       "",
			},
		},
		{
			desc:                 "no coverage for a zip code with invalid vzelte attribute in database",
			zipCode:              "94107",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":     "94107",
				"carriertype": "verizon",
				"vzelte":      "abc",
				"vze_lte_ind": "Y",
				"state":       "CA",
			},
		},
		{
			desc:                 "no coverage for a zip code that does not exist in the database",
			zipCode:              "11111",
			expectZipCodeCovered: false,
			causeDynamoDbError:   false,
			dynamodbReturnPayload: map[string]string{
				"zipcode":     "",
				"carriertype": "",
				"vzelte":      "",
				"vze_lte_ind": "",
				"state":       "",
			},
		},
		{
			desc:               "Sad path with Zip code that has coverage but results in dynamodb error",
			zipCode:            "94105",
			causeDynamoDbError: true,
			dynamodbReturnPayload: map[string]string{
				"zipcode":     "94105",
				"carriertype": "verizon",
				"vzelte":      "100",
				"vze_lte_ind": "Y",
				"state":       "CA",
			},
		},
	}

	for _, tC := range testCases {
		tableName := aws.String("fakeCoverage")

		fakeDb := &fakeVerizonDynamoDB{}
		if tC.causeDynamoDbError {
			fakeDb = &fakeVerizonDynamoDB{
				t:               t,
				tableName:       tableName,
				payloadToReturn: tC.dynamodbReturnPayload,
				err:             errors.New("fake DB error"),
			}
		} else {
			fakeDb = &fakeVerizonDynamoDB{
				t:               t,
				tableName:       tableName,
				payloadToReturn: tC.dynamodbReturnPayload,
				err:             nil,
			}
		}

		sprintdbClient := NewVerizonClient(tableName, fakeDb)
		result, err := sprintdbClient.VerifyCoverage(context.Background(), tC.zipCode)

		if tC.causeDynamoDbError {
			assert.NotNil(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tC.expectZipCodeCovered, result)
		}

		assert.Equal(t, tC.zipCode, fakeDb.Keys["zipcode"])
		assert.Equal(t, "verizon", fakeDb.Keys["carriertype"])
	}
}

type fakeVerizonDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	tableName       *string
	Keys            map[string]string
	payloadToReturn map[string]string // Store fake return values
	err             error
	t               *testing.T
}

func (fd *fakeVerizonDynamoDB) GetItemWithContext(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	assert.Equal(fd.t, *fd.tableName, *input.TableName, "incorrect table name")

	expectedAttributes := map[string]*string{
		"#0": aws.String("zipcode"),
		"#1": aws.String("carriertype"),
		"#2": aws.String("vzelte"),
		"#3": aws.String("vze_lte_ind"),
		"#4": aws.String("state"),
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
			"zipcode":     &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["zipcode"])},
			"carriertype": &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["carriertype"])},
			"vzelte":      &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["vzelte"])},
			"vze_lte_ind": &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["vze_lte_ind"])},
			"state":       &dynamodb.AttributeValue{S: aws.String(fd.payloadToReturn["state"])},
		}}, nil
	}
}
