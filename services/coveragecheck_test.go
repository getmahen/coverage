package services

import (
	"context"
	"errors"
	"testing"

	"bitbucket.org/credomobile/coverage/dbclient"
	"bitbucket.org/credomobile/coverage/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCoverageCheckHappyPathForSprint(t *testing.T) {
	dbClientFactory := mockClientFactory{}

	mockSprintClient := mockSprintClient{}
	mockSprintClient.On("VerifyCoverage", mock.Anything, mock.Anything).Return(true, nil)
	dbClientFactory.On("GetDbClient", mock.Anything).Return(mockSprintClient, nil)

	service := NewCoverageCheck(dbClientFactory)
	response, err := service.Verify(context.Background(), "94105", "1")

	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.IsCovered)
	mockSprintClient.AssertExpectations(t)
	dbClientFactory.AssertExpectations(t)
}

func TestCoverageCheckHappyPathForVerizon(t *testing.T) {
	dbClientFactory := mockClientFactory{}
	mockVerizonClient := mockVerizonClient{}

	mockVerizonClient.On("VerifyCoverage", mock.Anything, mock.Anything).Return(true, nil)
	dbClientFactory.On("GetDbClient", mock.Anything).Return(mockVerizonClient, nil)

	service := NewCoverageCheck(dbClientFactory)
	response, err := service.Verify(context.Background(), "94105", "2")

	assert.NotNil(t, response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.IsCovered)
	mockVerizonClient.AssertExpectations(t)
	dbClientFactory.AssertExpectations(t)
}

func TestCoverageCheckWithDbClientFactoryError(t *testing.T) {
	dbClientFactory := mockClientFactory{}
	mockVerizonClient := mockVerizonClient{}
	mockVerizonClient.On("VerifyCoverage", mock.Anything, mock.Anything).Return(true, nil)

	dbClientFactory.On("GetDbClient", mock.Anything).Return(mockVerizonClient, errors.New("Fake db Client Factory error"))

	service := NewCoverageCheck(dbClientFactory)
	_, err := service.Verify(context.Background(), "94105", "2")

	assert.Error(t, err)
	dbClientFactory.AssertExpectations(t)
}

func TestCoverageCheckWithDbClientError(t *testing.T) {
	dbClientFactory := mockClientFactory{}

	mockVerizonClient := mockVerizonClient{}
	mockVerizonClient.On("VerifyCoverage", mock.Anything, mock.Anything).Return(false, errors.New("Fake db Client error"))
	dbClientFactory.On("GetDbClient", mock.Anything).Return(mockVerizonClient, nil)

	service := NewCoverageCheck(dbClientFactory)
	_, err := service.Verify(context.Background(), "94105", "2")

	assert.Error(t, err)
	mockVerizonClient.AssertExpectations(t)
	dbClientFactory.AssertExpectations(t)
}

type mockClientFactory struct {
	mock.Mock
}

func (m mockClientFactory) GetDbClient(t entity.CarrierType) (dbclient.CoverageCheckClient, error) {
	args := m.Called(t)
	return args.Get(0).(dbclient.CoverageCheckClient), errOrNil(args.Get(1))
}

type mockSprintClient struct {
	mock.Mock
}

type mockVerizonClient struct {
	mock.Mock
}

func (m mockSprintClient) VerifyCoverage(ctx context.Context, zipCode string) (bool, error) {
	args := m.Called(ctx, zipCode)
	return args.Get(0).(bool), errOrNil(args.Get(1))
}

func (m mockVerizonClient) VerifyCoverage(ctx context.Context, zipCode string) (bool, error) {
	args := m.Called(ctx, zipCode)
	return args.Get(0).(bool), errOrNil(args.Get(1))
}

func errOrNil(o interface{}) error {
	if o == nil {
		return nil
	}
	return o.(error)
}
