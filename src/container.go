package main

import "pos/include"

// buildFacade - composition root, здесь (и только здесь) собираются зависимости
func buildFacade() *Facade {
	students := NewMemStudentRepo()
	rubrics := NewMemRubricRepo()
	evals := NewMemEvaluationRepo()
	exporter := CSVExporter{}

	strategy := domain.RoundedStrategy{Inner: domain.WeightedAll{}, Digits: 2}
	calc := domain.NewScoreCalculator(strategy)

	return NewFacade(students, rubrics, evals, calc, exporter)
}
