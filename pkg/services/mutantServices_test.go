package services

import (
	"errors"
	"testing"
	mock "x-menDetector/internal/mocks/repository"

	"github.com/stretchr/testify/assert"
)

func TestValidateDNA(t *testing.T) {
	tests := []struct {
		name    string
		dna     []string
		wantErr bool
	}{
		{
			name:    "secuencia valida",
			dna:     []string{"ATCG", "ATCG", "ATCG", "ATCG"},
			wantErr: false,
		},
		{
			name:    "secuencia invalida con carácter no permitido",
			dna:     []string{"ATCX", "ATCG", "ATCG", "ATCG"},
			wantErr: true,
		},
		{
			name:    "secuencia con solo caracteres validos",
			dna:     []string{"AATT", "CCGG", "TTAA", "GGCC"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDNA(tt.dna)
			if tt.wantErr {
				assert.Error(t, err, "se esperaba un error por secuencia invalida")
			} else {
				assert.NoError(t, err, "no se esperaba error por secuencia valida")
			}
		})
	}
}

func TestValidateDiagonal(t *testing.T) {
	tests := []struct {
		name     string
		dna      []string
		rowIndex int
		colIndex int
		expected bool
	}{
		{
			name: "diagonal hacia abajo a la derecha con 4 caracteres iguales",
			dna: []string{
				"ATCG",
				"CAGA",
				"TGAC",
				"GACA",
			},
			rowIndex: 0,
			colIndex: 0,
			expected: true,
		},
		{
			name: "diagonal hacia abajo a la izquierda con 4 caracteres iguales",
			dna: []string{
				"ATCT",
				"CGTA",
				"GTAC",
				"TACG",
			},
			rowIndex: 0,
			colIndex: 3,
			expected: true,
		},
		{
			name: "secuencia muy pequeña para diagonal de 4 caracteres",
			dna: []string{
				"AT",
				"CG",
			},
			rowIndex: 0,
			colIndex: 0,
			expected: false,
		},
		{
			name: "no hay secuencias diagonales de 4 caracteres",
			dna: []string{
				"ATCG",
				"GCTA",
				"TCAG",
				"ATCG",
			},
			rowIndex: 0,
			colIndex: 0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateDiagonal(tt.dna, tt.rowIndex, tt.colIndex, len(tt.dna))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateVertical(t *testing.T) {
	tests := []struct {
		name     string
		dna      []string
		rowIndex int
		colIndex int
		expected bool
	}{
		{
			name: "secuencia vertical con 4 caracteres iguales",
			dna: []string{
				"ATCG",
				"ATCG",
				"ATCG",
				"ATCG",
			},
			rowIndex: 0,
			colIndex: 0,
			expected: true,
		},
		{
			name: "sin secuencia vertical de 4 caracteres iguales",
			dna: []string{
				"ATCG",
				"CTGA",
				"TGAC",
				"GACT",
			},
			rowIndex: 0,
			colIndex: 0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateVertical(tt.dna, tt.rowIndex, tt.colIndex, len(tt.dna))
			if result != tt.expected {
				t.Errorf("validateVertical() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateHorizontal(t *testing.T) {
	tests := []struct {
		name     string
		dna      []string
		rowIndex int
		colIndex int
		expected bool
	}{
		{
			name: "secuencia horizontal con 4 caracteres iguales",
			dna: []string{
				"AAAA",
				"ATCG",
				"ATCG",
				"ATCG",
			},
			rowIndex: 0,
			colIndex: 0,
			expected: true,
		},
		{
			name: "sin secuencia horizontal de 4 caracteres iguales",
			dna: []string{
				"ATCG",
				"CTGA",
				"TGAC",
				"GACT",
			},
			rowIndex: 0,
			colIndex: 0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateHorizontal(tt.dna, tt.rowIndex, tt.colIndex, len(tt.dna))
			if result != tt.expected {
				t.Errorf("validateHorizontal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetMutantStats(t *testing.T) {
	mockRepo := new(mock.MockMutanRepository)
	service := MutanServicesImp{mutantRepository: mockRepo}

	t.Run("obtiene GetMutantStats exitosamente con valores validos", func(t *testing.T) {
		mockRepo.On("GetCounter", "mutant_count").Return("40", nil).Once()
		mockRepo.On("GetCounter", "human_count").Return("100", nil).Once()

		stats, err := service.GetMutantStats()

		assert.NoError(t, err, "error inesperado al obtener estadísticas de mutantes")
		assert.Equal(t, 40, stats.MutantCount)
		assert.Equal(t, 100, stats.HumanCount)
		assert.Equal(t, 0.4, stats.Ratio)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetMutantStatsErrorMutant(t *testing.T) {
	mockRepo := new(mock.MockMutanRepository)
	service := MutanServicesImp{mutantRepository: mockRepo}

	t.Run("error al obtener mutant_count", func(t *testing.T) {
		mockRepo.On("GetCounter", "mutant_count").Return("", errors.New("error al obtener mutant_count")).Once()

		stats, err := service.GetMutantStats()

		assert.Error(t, err, "se esperaba un error al obtener mutant_count")
		assert.EqualError(t, err, "error al obtener mutant_count")
		assert.Equal(t, 0, stats.MutantCount)
		assert.Equal(t, 0, stats.HumanCount)
		assert.Equal(t, 0.0, stats.Ratio)

		mockRepo.AssertExpectations(t)
	})

}

func TestGetMutantStatsHuman(t *testing.T) {
	mockRepo := new(mock.MockMutanRepository)
	service := MutanServicesImp{mutantRepository: mockRepo}

	t.Run("error al obtener human_count", func(t *testing.T) {
		mockRepo.On("GetCounter", "mutant_count").Return("40", nil).Once()
		mockRepo.On("GetCounter", "human_count").Return("", errors.New("error al obtener human_count")).Once()

		stats, err := service.GetMutantStats()

		assert.Error(t, err, "se esperaba un error al obtener human_count")
		assert.EqualError(t, err, "error al obtener human_count")
		assert.Equal(t, 0, stats.HumanCount)
		assert.Equal(t, 0, stats.MutantCount)
		assert.Equal(t, 0.0, stats.Ratio)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetMutantStatsMutantCount(t *testing.T) {
	mockRepo := new(mock.MockMutanRepository)
	service := MutanServicesImp{mutantRepository: mockRepo}

	t.Run("HumanCount es 0, retorna MutantCount como Ratio", func(t *testing.T) {
		mockRepo.On("GetCounter", "mutant_count").Return("40", nil).Once()
		mockRepo.On("GetCounter", "human_count").Return("0", nil).Once()

		stats, err := service.GetMutantStats()

		assert.NoError(t, err, "error inesperado al obtener estadísticas cuando human_count es 0")
		assert.Equal(t, 40, stats.MutantCount)
		assert.Equal(t, 0, stats.HumanCount)
		assert.Equal(t, 40.0, stats.Ratio)

		mockRepo.AssertExpectations(t)
	})
}

func TestIsMutant(t *testing.T) {
	t.Run("DNA es mutante, incrementa mutant_count y retorna true", func(t *testing.T) {
		mockRepo := new(mock.MockMutanRepository)
		mockRepo.On("IncrementCounter", "mutant_count").Return(nil).Once()
		service := MutanServicesImp{mutantRepository: mockRepo}

		dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}
		isMutant, err := service.IsMutant(dna)

		assert.NoError(t, err)
		assert.True(t, isMutant)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DNA no es mutante, incrementa human_count y retorna false", func(t *testing.T) {
		mockRepo := new(mock.MockMutanRepository)
		mockRepo.On("IncrementCounter", "human_count").Return(nil).Once()
		service := MutanServicesImp{mutantRepository: mockRepo}

		dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGACGG", "GCGTCA", "TCACTG"}
		isMutant, err := service.IsMutant(dna)

		assert.NoError(t, err)
		assert.False(t, isMutant)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error en validateDNA, retorna false con error", func(t *testing.T) {
		mockRepo := new(mock.MockMutanRepository)
		service := MutanServicesImp{mutantRepository: mockRepo}

		dna := []string{"ATE", "CAG", "TTA"}
		isMutant, err := service.IsMutant(dna)

		assert.Error(t, err)
		assert.False(t, isMutant)
	})

	t.Run("error al incrementar mutant_count, retorna false con error", func(t *testing.T) {
		mockRepo := new(mock.MockMutanRepository)
		mockRepo.On("IncrementCounter", "mutant_count").Return(errors.New("error al incrementar mutant_count")).Once()
		service := MutanServicesImp{mutantRepository: mockRepo}

		dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}
		isMutant, err := service.IsMutant(dna)

		assert.Error(t, err)
		assert.False(t, isMutant)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error al incrementar human_count, retorna false con error", func(t *testing.T) {
		mockRepo := new(mock.MockMutanRepository)
		mockRepo.On("IncrementCounter", "human_count").Return(errors.New("error al incrementar human_count")).Once()
		service := MutanServicesImp{mutantRepository: mockRepo}

		dna := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGACGG", "GCGTCA", "TCACTG"}
		isMutant, err := service.IsMutant(dna)

		assert.Error(t, err)
		assert.False(t, isMutant)
		mockRepo.AssertExpectations(t)
	})
}
