package domain

// Порты хранилищ, реализации лежат в src

type StudentRepository interface {
	Save(s Student) error
	Get(id string) (Student, error)
	List() []Student
}

type RubricRepository interface {
	Get() *Rubric
}

type EvaluationRepository interface {
	Save(e *Evaluation) error
	Find(studentID, criterionID string) (*Evaluation, error)
	ListByStudent(studentID string) []*Evaluation
	Delete(studentID, criterionID string) error
}
