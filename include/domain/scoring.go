package domain

import "math"

// ScoringStrategy - паттерн Стратегия, правило расчёта итога (в процентах)
type ScoringStrategy interface {
	Calculate(r *Rubric, evals []*Evaluation) float64
}

// WeightedAll - неоценённые критерии считаются как 0
type WeightedAll struct{}

func (WeightedAll) Calculate(r *Rubric, evals []*Evaluation) float64 {
	total := r.TotalWeight()
	if total == 0 {
		return 0
	}
	return weightedSum(r, evals) / total * 100
}

// WeightedRatedOnly - считаются только критерии, по которым есть оценка
type WeightedRatedOnly struct{}

func (WeightedRatedOnly) Calculate(r *Rubric, evals []*Evaluation) float64 {
	ratedWeight := 0.0
	for _, e := range evals {
		if c, ok := r.Find(e.criterionID); ok {
			ratedWeight += c.Weight()
		}
	}
	if ratedWeight == 0 {
		return 0
	}
	return weightedSum(r, evals) / ratedWeight * 100
}

func weightedSum(r *Rubric, evals []*Evaluation) float64 {
	sum := 0.0
	for _, e := range evals {
		c, ok := r.Find(e.criterionID)
		if !ok {
			continue
		}
		sum += e.score.Value() / c.MaxScore() * c.Weight()
	}
	return sum
}

// RoundedStrategy - паттерн Декоратор, округляет результат любой стратегии
type RoundedStrategy struct {
	Inner  ScoringStrategy
	Digits int
}

func (d RoundedStrategy) Calculate(r *Rubric, evals []*Evaluation) float64 {
	k := math.Pow(10, float64(d.Digits))
	return math.Round(d.Inner.Calculate(r, evals)*k) / k
}

// ScoreCalculator - Domain Service, итог не принадлежит одной сущности
type ScoreCalculator struct {
	strategy ScoringStrategy
}

func NewScoreCalculator(s ScoringStrategy) *ScoreCalculator {
	return &ScoreCalculator{strategy: s}
}

func (c *ScoreCalculator) SetStrategy(s ScoringStrategy) { c.strategy = s }

func (c *ScoreCalculator) Total(r *Rubric, evals []*Evaluation) float64 {
	return c.strategy.Calculate(r, evals)
}
