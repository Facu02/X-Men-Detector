package utils

import "fmt"

func ValidateDNA(dna []string) error {
	numRows := len(dna)
	if numRows < 4 || numRows > 50 {
		return fmt.Errorf("el tamaño de la matriz debe estar entre 4x4 y 50x50")
	}

	for _, row := range dna {
		if len(row) != numRows {
			return fmt.Errorf("cada fila debe tener la misma longitud que el numero total de filas")
		}
	}

	return nil
}
