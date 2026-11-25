package main

import "pos/include"

type MemStudentRepo struct {
	items []domain.Student
}

func NewMemStudentRepo() *MemStudentRepo { return &MemStudentRepo{} }

func (r *MemStudentRepo) Save(s domain.Student) error {
	for i, x := range r.items {
		if x.ID() == s.ID() {
			r.items[i] = s
			return nil
		}
	}
	r.items = append(r.items, s)
	return nil
}

func (r *MemStudentRepo) Get(id string) (domain.Student, error) {
	for _, x := range r.items {
		if x.ID() == id {
			return x, nil
		}
	}
	return domain.Student{}, domain.ErrStudentNotFound
}

func (r *MemStudentRepo) List() []domain.Student {
	res := make([]domain.Student, len(r.items))
	copy(res, r.items)
	return res
}

type MemRubricRepo struct {
	rubric *domain.Rubric
}

func NewMemRubricRepo() *MemRubricRepo { return &MemRubricRepo{rubric: domain.NewRubric()} }

func (r *MemRubricRepo) Get() *domain.Rubric { return r.rubric }

type MemEvaluationRepo struct {
	items []*domain.Evaluation
}

func NewMemEvaluationRepo() *MemEvaluationRepo { return &MemEvaluationRepo{} }

func (r *MemEvaluationRepo) Save(e *domain.Evaluation) error {
	for i, x := range r.items {
		if x.StudentID() == e.StudentID() && x.CriterionID() == e.CriterionID() {
			r.items[i] = e
			return nil
		}
	}
	r.items = append(r.items, e)
	return nil
}

func (r *MemEvaluationRepo) Find(studentID, criterionID string) (*domain.Evaluation, error) {
	for _, x := range r.items {
		if x.StudentID() == studentID && x.CriterionID() == criterionID {
			return x, nil
		}
	}
	return nil, domain.ErrEvaluationNotFound
}

func (r *MemEvaluationRepo) ListByStudent(studentID string) []*domain.Evaluation {
	var res []*domain.Evaluation
	for _, x := range r.items {
		if x.StudentID() == studentID {
			res = append(res, x)
		}
	}
	return res
}

func (r *MemEvaluationRepo) Delete(studentID, criterionID string) error {
	for i, x := range r.items {
		if x.StudentID() == studentID && x.CriterionID() == criterionID {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return domain.ErrEvaluationNotFound
}
