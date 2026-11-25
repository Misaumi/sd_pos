package main

import (
	"fmt"
	"time"

	"pos/include"
)

type TotalRow struct {
	StudentID string
	Name      string
	Total     float64
}

// TotalsExporter - порт выгрузки итогов (CSV, JSON и т.д.)
type TotalsExporter interface {
	Export(path string, rows []TotalRow) error
}

// Facade - паттерн Фасад, единая точка входа сценариев для консоли
type Facade struct {
	students    domain.StudentRepository
	rubrics     domain.RubricRepository
	evaluations domain.EvaluationRepository
	factory     domain.EvaluationFactory
	calculator  *domain.ScoreCalculator
	exporter    TotalsExporter
}

func NewFacade(s domain.StudentRepository, r domain.RubricRepository, e domain.EvaluationRepository,
	calc *domain.ScoreCalculator, exp TotalsExporter) *Facade {
	return &Facade{students: s, rubrics: r, evaluations: e, calculator: calc, exporter: exp}
}

func (f *Facade) AddStudent(name string) (domain.Student, error) {
	id := fmt.Sprintf("s%d", len(f.students.List())+1)
	st, err := domain.NewStudent(id, name)
	if err != nil {
		return domain.Student{}, err
	}
	return st, f.students.Save(st)
}

func (f *Facade) ListStudents() []domain.Student {
	return f.students.List()
}

func (f *Facade) AddCriterion(name string, maxScore, weight float64) (domain.Criterion, error) {
	r := f.rubrics.Get()
	id := fmt.Sprintf("c%d", len(r.Criteria())+1)
	c, err := domain.NewCriterion(id, name, maxScore, weight)
	if err != nil {
		return domain.Criterion{}, err
	}
	r.Add(c)
	return c, nil
}

func (f *Facade) ListCriteria() []domain.Criterion {
	return f.rubrics.Get().Criteria()
}

func (f *Facade) SetGrade(studentID, criterionID string, value float64, comment string) error {
	st, err := f.students.Get(studentID)
	if err != nil {
		return err
	}
	if _, err := f.evaluations.Find(studentID, criterionID); err == nil {
		return domain.ErrDuplicateEvaluation
	}
	e, err := f.factory.Create(st, f.rubrics.Get(), criterionID, value, comment, time.Now())
	if err != nil {
		return err
	}
	return f.evaluations.Save(e)
}

func (f *Facade) ChangeGrade(studentID, criterionID string, value float64) error {
	e, err := f.evaluations.Find(studentID, criterionID)
	if err != nil {
		return err
	}
	return e.Rescore(f.rubrics.Get(), value)
}

func (f *Facade) CancelGrade(studentID, criterionID string) error {
	return f.evaluations.Delete(studentID, criterionID)
}

func (f *Facade) ListGrades(studentID string) ([]*domain.Evaluation, error) {
	if _, err := f.students.Get(studentID); err != nil {
		return nil, err
	}
	return f.evaluations.ListByStudent(studentID), nil
}

func (f *Facade) Total(studentID string) (float64, error) {
	if _, err := f.students.Get(studentID); err != nil {
		return 0, err
	}
	return f.calculator.Total(f.rubrics.Get(), f.evaluations.ListByStudent(studentID)), nil
}

// SetStrategy переключает правило подсчёта итога
func (f *Facade) SetStrategy(name string) error {
	switch name {
	case "all":
		f.calculator.SetStrategy(domain.RoundedStrategy{Inner: domain.WeightedAll{}, Digits: 2})
	case "rated":
		f.calculator.SetStrategy(domain.RoundedStrategy{Inner: domain.WeightedRatedOnly{}, Digits: 2})
	default:
		return fmt.Errorf("неизвестная стратегия %q (all или rated)", name)
	}
	return nil
}

func (f *Facade) ExportTotals(path string) error {
	var rows []TotalRow
	for _, st := range f.students.List() {
		t, _ := f.Total(st.ID())
		rows = append(rows, TotalRow{StudentID: st.ID(), Name: st.Name(), Total: t})
	}
	return f.exporter.Export(path, rows)
}
