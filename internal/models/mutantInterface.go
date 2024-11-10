package models

type MutanRepository interface {
	IncrementCounter(counterName string) error
	GetCounter(counterName string) (string, error)
}

type MutanServices interface {
	IsMutant([]string) (bool, error)
	GetMutantStats() (MutantStats, error)
}
