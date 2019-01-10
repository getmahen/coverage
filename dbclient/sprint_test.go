package dbclient

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// var sprintCoverageRequiredFields = []string{
// 	"zipcode",
// 	"carriertype",
// 	"csa_leaf",
// 	"cur_pct_cov",
// 	"lte_4g_pctcov",
// }

func TestVerifyCoverageHappyPath(t *testing.T) {
	zipCode := "94105"
	logger := zerolog.New(os.Stdout).With().Logger()

	fakeDb := &fakeDynamoDB{t: t}
	sprintdbClient := NewSprintClient(&logger, fakeDb)
	_, err := sprintdbClient.VerifyCoverage(context.Background(), zipCode)
	assert.NoError(t, err)
	assert.Equal(t, zipCode, fakeDb.payload["zipcode"])
	assert.Equal(t, "sprint", fakeDb.payload["carriertype"])
}

func (fd *fakeDynamoDB) GetItemWithContext(ctx aws.Context, input *dynamodb.GetItemInput, opts ...request.Option) (*dynamodb.GetItemOutput, error) {
	assert.Equal(fd.t, "coverage", *input.TableName, "incorrect table name")

	//for _, field := range sprintCoverageRequiredFields {
	//expression, found := input.ExpressionAttributeNames[":name"]
	//expression := input.ProjectionExpression
	//fmt.Printf("expression: %v", *expression)
	//assert.Contains(fd.t, expression.String(), `["testprocessor"]`)
	//assert.True(fd.t, found, fmt.Sprintf("%s is a required field and is missing", field))
	//}

	expected := map[string]*string{
		"#0": aws.String("zipcode"),
		"#1": aws.String("carriertype"),
		"#2": aws.String("csa_leaf"),
		"#3": aws.String("cur_pct_cov"),
		"#4": aws.String("lte_4g_pctcov"),
	}

	actual := input.ExpressionAttributeNames
	if e, a := expected, actual; !reflect.DeepEqual(a, e) {
		fd.t.Errorf("expect %v, got %v", e, a)
	}

	fd.payload = make(map[string]string)
	for k, v := range input.Key {
		if v == nil {
			continue
		}

		if v.S != nil {
			fd.payload[k] = *v.S
		} else if v.N != nil {
			fd.payload[k] = *v.N
		}
	}

	return &dynamodb.GetItemOutput{}, nil
}
