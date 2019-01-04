package dbclient

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type VerizonDbClient struct {
	TableName  string
	Logger     *zerolog.Logger
	Connection dynamodbiface.DynamoDBAPI
}

func (v VerizonDbClient) VerifyCoverage(ctx context.Context, zipCode string, carrierID string) (bool, error) {
	v.Logger.Info().Msgf("*** IN VERIZON DB CLIENT ***")
	return true, nil
}
