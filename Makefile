NO_COLOR=\033[0m
OK_COLOR=\033[0;32m
TS_COLOR=\033[1;30m
BUILD_DATE_VER=""
TIMESTAMP := $(shell date -u +%Y%m%d%H%M%S)
PORT="8002"

BASE_ENV_VALS := ENVIRONMENT="LOCAL_DEV" \
		VAULT_TOKEN="cbc762a7-1eba-0aac-49f8-0deff1bcdfed" \
		VAULT_URL="https://localhost:8243" \
		CONSUL_URL="http://localhost:8500" \
		DATADOG_APM_URL="localhost:8126" \
		DATADOG_STATS_URL="localhost:8125" \
		LISTEN_IP=127.0.0.1 \
		NOMAD_PORT_http=3999 \
		LOG_LEVEL="debug" \
		DYNAMODB_ARN="arn:aws:dynamodb:us-east-2:674346455231:table/coverage" \
		AWS_REGION="us-east-2"

.PHONY: test
test:
	@echo "$(OK_COLOR)==> Testing$(NO_COLOR)"
	go test -cover -race ./...

.PHONY: clean
clean:
	@echo "$(TS_COLOR)$(shell date "+%Y/%m/%d %H:%M:%S")$(NO_COLOR)$(OK_COLOR) ==> Cleaning$(NO_COLOR)"
	@go clean
	@rm -f coverage coverage.zip

.PHONY: build
build: clean
	@echo "$(TS_COLOR)$(shell date "+%Y/%m/%d %H:%M:%S")$(NO_COLOR)$(OK_COLOR)==> Building$(NO_COLOR)"
	GOOS=linux go build --ldflags "-X bitbucket.org/credomobile/coverage/handler.version=`git rev-parse HEAD`" -o coverage
	zip coverage.zip coverage vault-cas.crt

.PHONY: upload
upload:
	@echo "$(TS_COLOR)$(shell date "+%Y/%m/%d %H:%M:%S")$(NO_COLOR)$(OK_COLOR)==> Deploying Zip to s3$(NO_COLOR)"
	aws s3 cp coverage.zip s3://credo-dev-lambdas/coverage.zip --metadata GitHash=`git rev-parse HEAD`

.PHONY: run
run:
	@echo "$(TS_COLOR)$(shell date "+%Y/%m/%d %H:%M:%S")$(NO_COLOR)$(OK_COLOR)==> Running Lambda locally on PORT:$(PORT) $(NO_COLOR)"
	$(BASE_ENV_VALS) _LAMBDA_SERVER_PORT=$(PORT) go run main.go