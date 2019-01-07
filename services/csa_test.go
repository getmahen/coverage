package services

import (
	"context"
	"errors"
	"os"
	"testing"

	"bitbucket.org/credomobile/coverage/dbclient"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCsa(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()
	csaService := NewCsa(&logger)

	assert.IsType(t, csa{}, csaService)
	assert.IsType(t, dbclient.SprintDbClient{}, csaService.dbClient)
}

func TestGetCsaHappyPathWithCsaFound(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()

	mockSprintCsaDbClient := mockSprintCsaDbClient{}
	mockSprintCsaDbClient.On("GetCsa", mock.Anything, mock.Anything).Return("fakeCsa", nil)

	csaService := csa{
		logger:   &logger,
		dbClient: mockSprintCsaDbClient,
	}

	response, err := csaService.GetCsa(context.Background(), "94105")

	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.CsaFound)
	assert.Equal(t, "fakeCsa", response.Csa)
	mockSprintCsaDbClient.AssertExpectations(t)
}

func TestGetCsaHappyPathWithCsaNotFound(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()

	mockSprintCsaDbClient := mockSprintCsaDbClient{}
	mockSprintCsaDbClient.On("GetCsa", mock.Anything, mock.Anything).Return("", nil)

	csaService := csa{
		logger:   &logger,
		dbClient: mockSprintCsaDbClient,
	}

	response, err := csaService.GetCsa(context.Background(), "94105")

	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, false, response.CsaFound)
	assert.Equal(t, "", response.Csa)
	mockSprintCsaDbClient.AssertExpectations(t)
}

func TestGetCsaWithDbClientError(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().Logger()

	mockSprintCsaDbClient := mockSprintCsaDbClient{}
	mockSprintCsaDbClient.On("GetCsa", mock.Anything, mock.Anything).Return("", errors.New("Fake Db error"))

	csaService := csa{
		logger:   &logger,
		dbClient: mockSprintCsaDbClient,
	}

	_, err := csaService.GetCsa(context.Background(), "94105")

	assert.Error(t, err)
	mockSprintCsaDbClient.AssertExpectations(t)
}

type mockSprintCsaDbClient struct {
	mock.Mock
}

func (m mockSprintCsaDbClient) GetCsa(ctx context.Context, zipCode string) (string, error) {
	args := m.Called(ctx, zipCode)
	return args.Get(0).(string), errOrNil(args.Get(1))
}
