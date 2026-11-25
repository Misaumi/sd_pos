package domain

import (
	"fmt"
	"time"
)

// Score - Value Object, балл уже проверен на диапазон
type Score struct {
	value float64
}

func NewScore(value, max float64) (Score, error) {
	if value < 0 || value > max {
		return Score{}, fmt.Errorf("%w: %.2f не в [0; %.2f]", ErrScoreOutOfRange, value, max)
	}
	return Score{value: value}, nil
}

func (s Score) Value() float64 { return s.value }

// Evaluation - сущность, оценка студента по критерию
type Evaluation struct {
	studentID   string
	criterionID string
	score       Score
	comment     string
	date        time.Time
}

func (e *Evaluation) StudentID() string   { return e.studentID }
func (e *Evaluation) CriterionID() string { return e.criterionID }
func (e *Evaluation) Score() Score        { return e.score }
func (e *Evaluation) Comment() string     { return e.comment }
func (e *Evaluation) Date() time.Time     { return e.date }

// Rescore исправляет балл, заново проверяя его по критерию из рубрикатора
func (e *Evaluation) Rescore(r *Rubric, value float64) error {
	c, ok := r.Find(e.criterionID)
	if !ok {
		return ErrCriterionNotInRubric
	}
	s, err := NewScore(value, c.MaxScore())
	if err != nil {
		return err
	}
	e.score = s
	return nil
}

// EvaluationFactory - паттерн Фабрика, единственное место создания Evaluation
type EvaluationFactory struct{}

func (EvaluationFactory) Create(st Student, r *Rubric, criterionID string, value float64, comment string, at time.Time) (*Evaluation, error) {
	c, ok := r.Find(criterionID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrCriterionNotInRubric, criterionID)
	}
	s, err := NewScore(value, c.MaxScore())
	if err != nil {
		return nil, err
	}
	return &Evaluation{
		studentID:   st.ID(),
		criterionID: criterionID,
		score:       s,
		comment:     comment,
		date:        at,
	}, nil
}
