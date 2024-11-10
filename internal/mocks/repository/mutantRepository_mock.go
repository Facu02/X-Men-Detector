package mock

import (
	"github.com/stretchr/testify/mock"
)

type MockMutanRepository struct {
	mock.Mock
}

func (m *MockMutanRepository) IncrementCounter(counterName string) error {
	args := m.Called(counterName)
	return args.Error(0)
}

func (m *MockMutanRepository) GetCounter(counterName string) (string, error) {
	args := m.Called(counterName)
	return args.String(0), args.Error(1)
}
