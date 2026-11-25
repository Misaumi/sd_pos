package domain

import "strings"

// Student - сущность (Entity), идентифицируется по id, а не по имени
type Student struct {
	id   string
	name string
}

func NewStudent(id, name string) (Student, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Student{}, ErrEmptyName
	}
	return Student{id: id, name: name}, nil
}

func (s Student) ID() string   { return s.id }
func (s Student) Name() string { return s.name }
