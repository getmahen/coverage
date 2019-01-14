package main

import (
	"log"

	"bitbucket.org/credomobile/coverage/dbclient"
	"bitbucket.org/credomobile/coverage/handlers"
	"bitbucket.org/credomobile/coverage/services"
	"bitbucket.org/credomobile/coverage/validators"
	"bitbucket.org/credomobile/frink"
	"bitbucket.org/credomobile/frink/flambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Config struct {
	frink.BaseConfig
	DynamoDBArn string `env:"DYNAMODB_ARN"`
}

var initialized = false
var frinkLambda *flambda.FrinkAdapter

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		log.Println("Lambda COLD START")

		config := &Config{}
		opts, err := frink.NewDefaultLambdaOptions()
		if err != nil {
			log.Fatal("unable to configure options: ", err)
		}
		app, err := frink.New("coverage", opts, config)
		if err != nil {
			app.Logger.Fatal().Err(err).Msg("unable to configure application")
		}

		dbclientFactory, err := dbclient.NewDbClientFactory(config.DynamoDBArn, app.Logger)
		if err != nil {
			app.Logger.Fatal().Err(err).Msg("unable to configure Db Client")
		}

		coverageCheckService := services.NewCoverageCheck(dbclientFactory)

		csaService, err := services.NewCsa(config.DynamoDBArn, app.Logger)
		if err != nil {
			app.Logger.Fatal().Err(err).Msg("unable to configure Csa service")
		}

		coverageCheckValidator := validators.NewCoverageCheckValidator()
		csaValidator := validators.NewCsaValidator()

		app.Router.Get("/v1/coveragecheck", handlers.CheckCoverage(coverageCheckValidator, coverageCheckService))
		app.Router.Get("/v1/csa", handlers.GetCsa(csaValidator, csaService))

		frinkLambda = flambda.New(app)
		initialized = true
	}
	return frinkLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
