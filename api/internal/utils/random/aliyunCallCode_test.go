package random

import (
	"testing"

	"github.com/alibabacloud-go/dyvmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDyvmsapiClient is a mock of the DyvmsapiClientInterface
type MockDyvmsapiClient struct {
	mock.Mock
}

func (m *MockDyvmsapiClient) SendVerification(request *client.SendVerificationRequest) (*client.SendVerificationResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*client.SendVerificationResponse), args.Error(1)
}

func TestCreateDyvmsapiClient(t *testing.T) {
	client, err := CreateDyvmsapiClient()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestSendVerificationCode(t *testing.T) {
	// Create a mock client
	mockClient := new(MockDyvmsapiClient)

	// Set up the expected behavior
	mockClient.On("SendVerification", mock.Anything).Return(&client.SendVerificationResponse{
		Body: &client.SendVerificationResponseBody{
			Code:    tea.String("OK"),
			Message: tea.String("Success"),
		},
	}, nil)

	// Call the function we want to test
	err := SendVerificationCode("13682233421", mockClient)

	// Assert that the error is nil
	assert.NoError(t, err)

	// Assert that our expectations were met
	mockClient.AssertExpectations(t)
}

func TestSendVerificationCodeWithNilClient(t *testing.T) {
	// This test will use the actual CreateDyvmsapiClient function
	// Note: This test may fail if it cannot connect to the actual service
	assert.NotPanics(t, func() {
		_ = SendVerificationCode("13682233421", nil)
	})
}
