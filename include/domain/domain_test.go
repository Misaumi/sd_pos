package domain

import (
	"errors"
	"testing"
	"time"
)

func setup(t *testing.T) (Student, *Rubric) {
	st, _ := NewStudent("s1", "Иван")
	r := NewRubric()
	c1, _ := NewCriterion("c1", "Лекции", 10, 1)
	c2, _ := NewCriterion("c2", "Семинары", 10, 3)
	r.Add(c1)
	r.Add(c2)
	return st, r
}

func TestScoreOutOfRange(t *testing.T) {
	st, r := setup(t)
	_, err := EvaluationFactory{}.Create(st, r, "c1", 12, "", time.Now())
	if !errors.Is(err, ErrScoreOutOfRange) {
		t.Fatalf("ожидали ErrScoreOutOfRange, получили %v", err)
	}
	_, err = EvaluationFactory{}.Create(st, r, "c1", -1, "", time.Now())
	if !errors.Is(err, ErrScoreOutOfRange) {
		t.Fatalf("отрицательный балл должен быть ошибкой")
	}
}

func TestUnknownCriterion(t *testing.T) {
	st, r := setup(t)
	_, err := EvaluationFactory{}.Create(st, r, "zzz", 5, "", time.Now())
	if !errors.Is(err, ErrCriterionNotInRubric) {
		t.Fatalf("ожидали ErrCriterionNotInRubric, получили %v", err)
	}
}

func TestRescoreChecksRange(t *testing.T) {
	st, r := setup(t)
	e, _ := EvaluationFactory{}.Create(st, r, "c1", 5, "", time.Now())
	if err := e.Rescore(r, 11); err == nil {
		t.Fatal("Rescore должен отклонить 11 из 10")
	}
	if err := e.Rescore(r, 7); err != nil || e.Score().Value() != 7 {
		t.Fatal("Rescore должен принять 7")
	}
}

func TestTotalWeighted(t *testing.T) {
	st, r := setup(t)
	e1, _ := EvaluationFactory{}.Create(st, r, "c1", 10, "", time.Now())
	e2, _ := EvaluationFactory{}.Create(st, r, "c2", 5, "", time.Now())
	evals := []*Evaluation{e1, e2}

	if got := NewScoreCalculator(WeightedAll{}).Total(r, evals); got != 62.5 {
		t.Fatalf("ожидали 62.5, получили %v", got)
	}
	if got := NewScoreCalculator(WeightedAll{}).Total(r, evals[:1]); got != 25 {
		t.Fatalf("ожидали 25, получили %v", got)
	}
	if got := NewScoreCalculator(WeightedRatedOnly{}).Total(r, evals[:1]); got != 100 {
		t.Fatalf("ожидали 100, получили %v", got)
	}
}

func TestRoundingDecorator(t *testing.T) {
	st, r := setup(t)
	e1, _ := EvaluationFactory{}.Create(st, r, "c1", 10, "", time.Now())
	e2, _ := EvaluationFactory{}.Create(st, r, "c2", 3, "", time.Now())
	s := RoundedStrategy{Inner: WeightedAll{}, Digits: 1}
	if got := s.Calculate(r, []*Evaluation{e1, e2}); got != 47.5 {
		t.Fatalf("получили %v", got)
	}
}
