package models

type MutantStats struct {
	MutantCount int     `json:"mutant_count"`
	HumanCount  int     `json:"human_count"`
	Ratio       float64 `json:"ratio"`
}
