package domain

import "strings"

// Criterion - критерий оценивания: название, максимум и вес (вклад в итог)
type Criterion struct {
	id       string
	name     string
	maxScore float64
	weight   float64
}

func NewCriterion(id, name string, maxScore, weight float64) (Criterion, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Criterion{}, ErrEmptyName
	}
	if maxScore <= 0 {
		return Criterion{}, ErrBadMaxScore
	}
	if weight <= 0 {
		return Criterion{}, ErrBadWeight
	}
	return Criterion{id: id, name: name, maxScore: maxScore, weight: weight}, nil
}

func (c Criterion) ID() string        { return c.id }
func (c Criterion) Name() string      { return c.name }
func (c Criterion) MaxScore() float64 { return c.maxScore }
func (c Criterion) Weight() float64   { return c.weight }

// Rubric - рубрикатор (корень агрегата), хранит критерии курса
type Rubric struct {
	criteria []Criterion
}

func NewRubric() *Rubric {
	return &Rubric{}
}

func (r *Rubric) Add(c Criterion) {
	r.criteria = append(r.criteria, c)
}

func (r *Rubric) Find(id string) (Criterion, bool) {
	for _, c := range r.criteria {
		if c.id == id {
			return c, true
		}
	}
	return Criterion{}, false
}

// Criteria возвращает копию, чтобы снаружи нельзя было испортить список
func (r *Rubric) Criteria() []Criterion {
	res := make([]Criterion, len(r.criteria))
	copy(res, r.criteria)
	return res
}

func (r *Rubric) TotalWeight() float64 {
	sum := 0.0
	for _, c := range r.criteria {
		sum += c.weight
	}
	return sum
}
