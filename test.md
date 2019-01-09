/// TO GET JWT
getpostmantoken --expire=1d --service=coverage | pbcopy

# RUN LOCAL DYNAMODB DOCKER 
docker run -d -p 8000:8000 amazon/dynamodb-local


curl -X GET \
  'https://qa-api.credomobile.com/coverage/v1/coveragecheck?zipcode=94538&carrierid=2' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwb3N0bWFuIiwiY2FsbGVlQWdlbnRBdWQiOiJjb3ZlcmFnZSIsImlhdCI6MTUzOTExNDAxNiwiZXhwIjoxNTcwNjcxNjE2fQ.G_Y-Mae6GNvlf9Jly6SLjlE3BTVQxqllT4WLK4TccPNPm2B4AYYQ0g_Ep3JIWym4YDf8rgcS6hqGrWCTglfeUeynNzG60TQB73h6NZBQ5wVngLM6yTdxUMZIBkmNwFIO0-pBnWY7fMBB4IrxZWyCNKq88kRk7XMDAJpl0LtHQ8Dkw9S7AMnsjH1yPQU9auXKtirKpAZiEHVgCIEAZ81ZRNOYvidHRrdejGA6Ar2ry3RVebAKO_Kpw9PSyhkiEA2IHuILkjrexj9Hbf0x_rBHBOXvEtGMJ5qhl0X5kACk5IeqJ31nbatexxv7SexNGrlki0_lU5ZBYJHXrF8OZT1LGA' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: f767d9ab-b650-4d57-9fd7-0e13427a1a77' \
  -H 'cache-control: no-cache' \
  -H 'x-api-key: JuT33brJpR4shKXBuVIIf4xeGnvTTvx6bZGLkvP6' \
  -H 'x-trackingid: 565242ac-a5d7-4377-a88f-adc0d8cdd3ea'


  curl -X GET \
  'https://credoqa.dev/coverage/v1/coveragecheck?zipcode=94538&carrierid=2' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJwb3N0bWFuIiwiY2FsbGVlQWdlbnRBdWQiOiJjb3ZlcmFnZSIsImlhdCI6MTUzOTExNDAxNiwiZXhwIjoxNTcwNjcxNjE2fQ.G_Y-Mae6GNvlf9Jly6SLjlE3BTVQxqllT4WLK4TccPNPm2B4AYYQ0g_Ep3JIWym4YDf8rgcS6hqGrWCTglfeUeynNzG60TQB73h6NZBQ5wVngLM6yTdxUMZIBkmNwFIO0-pBnWY7fMBB4IrxZWyCNKq88kRk7XMDAJpl0LtHQ8Dkw9S7AMnsjH1yPQU9auXKtirKpAZiEHVgCIEAZ81ZRNOYvidHRrdejGA6Ar2ry3RVebAKO_Kpw9PSyhkiEA2IHuILkjrexj9Hbf0x_rBHBOXvEtGMJ5qhl0X5kACk5IeqJ31nbatexxv7SexNGrlki0_lU5ZBYJHXrF8OZT1LGA' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: f767d9ab-b650-4d57-9fd7-0e13427a1a77' \
  -H 'cache-control: no-cache' \
  -H 'x-api-key: JuT33brJpR4shKXBuVIIf4xeGnvTTvx6bZGLkvP6' \
  -H 'x-trackingid: 565242ac-a5d7-4377-a88f-adc0d8cdd3ea'


//TODO
- Move the Valdation logic to the handler level from Service layer - (DONE)
- Fix the Validator to accept (w http.ResponseWriter, r *http.Request) as the func signature - (DONE)
- Add more table driven tests for Validator - (DONE)
- Check if the Frink's Validator can be used instead
- JWT Auth
- Implement Datadog/X-ray


//DYNAMODB ARN
arn::aws::dynamodb/coverage


///DYNAMODB LOCAL CLI COMMANDS
aws dynamodb create-table \
    --table-name Music \
    --attribute-definitions \
        AttributeName=Artist,AttributeType=S \
        AttributeName=SongTitle,AttributeType=S \
    --key-schema AttributeName=Artist,KeyType=HASH AttributeName=SongTitle,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url http://localhost:8000


aws dynamodb create-table \
    --cli-input-json file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/verizon_table_coverage.json \
    --endpoint-url http://localhost:8000

aws dynamodb create-table \
--cli-input-json file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/sprint_table_coverage.json \
--endpoint-url http://localhost:8000


// CREATE ZIP CODE DATA RECORD FOR VERIZON

    aws dynamodb put-item \
    --table-name verizon_coverage \
    --item '{"ZIP": {"S": "94105"}, "env": {"S": "qa/prod"}}' \
    --endpoint-url http://localhost:8000

    aws dynamodb put-item \
    --table-name verizon_coverage \
    --item file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/verizon_coverage_data.json \
    --endpoint-url http://localhost:8000

// CREATE ZIP CODE DATA RECORD FOR SPRINT
aws dynamodb put-item \
    --table-name sprint_coverage \
    --item file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/sprint_coverage_data.json \
    --endpoint-url http://localhost:8000


aws dynamodb put-item \
--table-name Music  \
--item \
    '{"Artist": {"S": "No One You Know"}, "SongTitle": {"S": "Call Me Today"}, "AlbumTitle": {"S": "Somewhat Famous"}}' \
--return-consumed-capacity TOTAL \
--endpoint-url http://localhost:8000

aws dynamodb put-item \
    --table-name Music \
    --item '{ \
        "Artist": {"S": "Acme Band"}, \
        "SongTitle": {"S": "Happy Day"}, \
        "AlbumTitle": {"S": "Songs About Life"} }' \
    --return-consumed-capacity TOTAL \
    --endpoint-url http://localhost:8000

aws dynamodb query --table-name Music --key-conditions file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/dynamoquery.json --endpoint-url http://localhost:8000


aws dynamodb describe-table --table-name Music --endpoint-url http://localhost:8000


aws dynamodb list-tables --endpoint-url http://localhost:8000


aws dynamodb delete-table --table-name sprint_coverage --endpoint-url http://localhost:8000

aws dynamodb get-item --table-name verizon_coverage --key file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/coveragequery.json --endpoint-url http://localhost:8000

aws dynamodb get-item --table-name sprint_coverage --key file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/coveragequery.json --endpoint-url http://localhost:8000

aws dynamodb batch-write-item --request-items file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/verizon_coverage_batch_data.json --endpoint-url http://localhost:8000

aws dynamodb batch-write-item --request-items file:///Users/mrekapally/go/src/bitbucket.org/credomobile/coverage/sprint_coverage_batch_data.json --endpoint-url http://localhost:8000