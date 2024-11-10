package mocks

import (
	"x-menDetector/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockMutanServices struct {
	mock.Mock
}

func (m *MockMutanServices) IsMutant(dna []string) (bool, error) {
	args := m.Called(dna)
	return args.Bool(0), args.Error(1)
}

func (m *MockMutanServices) GetMutantStats() (models.MutantStats, error) {
	args := m.Called()
	return args.Get(0).(models.MutantStats), args.Error(1)
}
