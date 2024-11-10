package services

import (
	"errors"
	"strconv"
	"x-menDetector/internal/models"
)

var validChars = map[rune]bool{
	'A': true,
	'T': true,
	'C': true,
	'G': true,
}

type MutanServicesImp struct {
	mutantRepository models.MutanRepository
}

func NewMutantServices(mutantRepository models.MutanRepository) models.MutanServices {
	return MutanServicesImp{mutantRepository: mutantRepository}
}

func validateDNA(dna []string) error {
	for _, row := range dna {
		for _, char := range row {
			if !validChars[char] {
				return errors.New("secuencia de ADN invalida: solo se permiten A, T, C, G")
			}
		}
	}
	return nil
}

func (ms MutanServicesImp) IsMutant(dna []string) (bool, error) {
	if err := validateDNA(dna); err != nil {
		return false, err
	}

	linesFound := 0
	dnaSize := len(dna)

	for rowIndex := 0; rowIndex < dnaSize; rowIndex++ {
		for colIndex := 0; colIndex < dnaSize; colIndex++ {
			if validateHorizontal(dna, rowIndex, colIndex, dnaSize) ||
				validateVertical(dna, rowIndex, colIndex, dnaSize) ||
				validateDiagonal(dna, rowIndex, colIndex, dnaSize) {
				linesFound++
				if linesFound >= 2 {
					err := ms.mutantRepository.IncrementCounter("mutant_count")
					if err != nil {
						return false, err
					}
					return true, nil
				}
			}
		}
	}
	err := ms.mutantRepository.IncrementCounter("human_count")
	if err != nil {
		return false, err
	}
	return false, nil
}

func validateHorizontal(dna []string, rowIndex, colIndex, dnaSize int) bool {
	if colIndex+3 < dnaSize {
		base := dna[rowIndex][colIndex]
		count := 1
		for k := colIndex + 1; k < dnaSize && dna[rowIndex][k] == base; k++ {
			count++
			if count >= 4 {
				return true
			}
		}
	}
	return false
}

func validateVertical(dna []string, rowIndex, colIndex, dnaSize int) bool {
	if rowIndex+3 < dnaSize {
		base := dna[rowIndex][colIndex]
		count := 1
		for k := rowIndex + 1; k < dnaSize && dna[k][colIndex] == base; k++ {
			count++
			if count >= 4 {
				return true
			}
		}
	}
	return false
}

func validateDiagonal(dna []string, rowIndex, colIndex, dnaSize int) bool {
	base := dna[rowIndex][colIndex]

	if rowIndex+3 < dnaSize && colIndex+3 < dnaSize {
		count := 1
		for k := 1; k < 4; k++ {
			if dna[rowIndex+k][colIndex+k] == base {
				count++
				if count == 4 {
					return true
				}
			} else {
				break
			}
		}
	}

	if rowIndex+3 < dnaSize && colIndex-3 >= 0 {
		count := 1
		for k := 1; k < 4; k++ {
			if dna[rowIndex+k][colIndex-k] == base {
				count++
				if count == 4 {
					return true
				}
			} else {
				break
			}
		}
	}

	return false
}

func (ms MutanServicesImp) GetMutantStats() (models.MutantStats, error) {
	var stats models.MutantStats

	mutantStr, err := ms.mutantRepository.GetCounter("mutant_count")
	if err != nil {
		return stats, err
	}
	humanStr, err := ms.mutantRepository.GetCounter("human_count")
	if err != nil {
		return stats, err
	}

	stats.MutantCount, _ = strconv.Atoi(mutantStr)
	stats.HumanCount, _ = strconv.Atoi(humanStr)

	if stats.HumanCount == 0 {
		stats.Ratio = float64(stats.MutantCount)
	} else {
		stats.Ratio = float64(stats.MutantCount) / float64(stats.HumanCount)
	}

	return stats, nil
}
