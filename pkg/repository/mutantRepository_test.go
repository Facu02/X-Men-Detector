package repository_test

import (
	"errors"
	"testing"
	"x-menDetector/pkg/repository"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name            string
	expectedCounter string
	counterName     string
	incrementCount  int
	setupMocks      func(mock redismock.ClientMock)
	expectedError   error
}

func TestMutantRepository(t *testing.T) {
	testCases := []testCase{
		{
			name:            "incrementar contador",
			counterName:     "mutant_count",
			incrementCount:  1,
			expectedCounter: "1",
			setupMocks: func(mock redismock.ClientMock) {
				mock.ExpectIncr("mutant_count").SetVal(1)
				mock.ExpectGet("mutant_count").SetVal("1")
			},
			expectedError: nil,
		},
		{
			name:            "obtener contador sin valor previo",
			counterName:     "human_count",
			incrementCount:  0,
			expectedCounter: "0",
			setupMocks: func(mock redismock.ClientMock) {
				mock.ExpectGet("human_count").SetVal("0")
			},
			expectedError: nil,
		},
		{
			name:            "obtener contador después de múltiples incrementos",
			counterName:     "mutant_count",
			incrementCount:  3,
			expectedCounter: "3",
			setupMocks: func(mock redismock.ClientMock) {
				mock.ExpectIncr("mutant_count").SetVal(1)
				mock.ExpectIncr("mutant_count").SetVal(2)
				mock.ExpectIncr("mutant_count").SetVal(3)
				mock.ExpectGet("mutant_count").SetVal("3")
			},
			expectedError: nil,
		},
		{
			name:            "obtener contador con valor no encontrado (redis.Nil)",
			counterName:     "unknown_count",
			incrementCount:  0,
			expectedCounter: "0",
			setupMocks: func(mock redismock.ClientMock) {
				mock.ExpectGet("unknown_count").SetErr(redis.Nil)
			},
			expectedError: nil,
		},
		{
			name:            "obtener contador con error en redis",
			counterName:     "mutant_count",
			incrementCount:  0,
			expectedCounter: "",
			setupMocks: func(mock redismock.ClientMock) {
				mock.ExpectGet("mutant_count").SetErr(errors.New("redis error"))
			},
			expectedError: errors.New("redis error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			repo := repository.NewMutanRepository(client)

			tc.setupMocks(mock)

			for i := 0; i < tc.incrementCount; i++ {
				err := repo.IncrementCounter(tc.counterName)
				assert.NoError(t, err, "no deberia haber error al incrementar el contador")
			}

			value, err := repo.GetCounter(tc.counterName)
			assert.Equal(t, tc.expectedCounter, value, "el valor del contador no coincide")

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error(), "el error no es el esperado")
			} else {
				assert.NoError(t, err, "no deberia haber error al obtener el contador")
			}

			assert.NoError(t, mock.ExpectationsWereMet(), "hay expectativas no cumplidas")
		})
	}
}
