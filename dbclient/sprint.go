package dbclient

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog"
)

type SprintDbClient struct {
	TableName  string
	Logger     *zerolog.Logger
	Connection dynamodbiface.DynamoDBAPI
}

type Item struct {
	Artist    string `json:"artist"`
	SongTitle string `json:"songtitle"`
}

func NewSprintClient(tableName string, logger *zerolog.Logger, connection dynamodbiface.DynamoDBAPI) SprintDbClient {
	return SprintDbClient{TableName: tableName, Logger: logger, Connection: connection}
}

func (s SprintDbClient) VerifyCoverage(ctx context.Context, zipCode string, carrierID string) (bool, error) {
	s.Logger.Info().Msgf("*** IN SPRINT DB CLIENT VerifyCoverage() ***")

	result, err := s.Connection.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		fmt.Println("Number of tables:", len(result.TableNames))
		for _, table := range result.TableNames {
			fmt.Println("Name:", table)
		}
	}

	// req := &dynamodb.DescribeTableInput{
	// 	TableName: aws.String("music"),
	// }
	// result, err := s.Connection.DescribeTable(req)
	// if err != nil {
	// 	fmt.Printf("%s", err)
	// }
	// table := result.Table
	// fmt.Printf("done", table)

	// input := &dynamodb.GetItemInput{
	// 	TableName: aws.String("Music"),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"artist": {S: aws.String("No One You Know")},
	// 	},
	// }

	// _, err := s.Connection.GetItemWithContext(ctx, input)
	// if err != nil {
	// 	s.Logger.Fatal().Err(err).Msg("failed to query dynamodb")

	// 	return false, err
	// }

	//item := Item{}
	// err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	// }

	// if item.Artist == "" {
	// 	fmt.Println("Could not find 'The Big New Movie' (2015)")
	// 	return false, nil
	// }

	// fmt.Println("Found item:")
	// fmt.Println("Artist:  ", item.Artist)
	// fmt.Println("SongTitle: ", item.SongTitle)

	return true, nil
}

func (s SprintDbClient) GetCsa(ctx context.Context, zipCode string) (string, error) {
	s.Logger.Info().Msgf("*** IN SPRINT DB CLIENT GetCsa() ***")
	return "fakeCsa", nil
}
