package domain

import "errors"

var (
	ErrEmptyName            = errors.New("название не может быть пустым")
	ErrBadMaxScore          = errors.New("max_score должен быть больше 0")
	ErrBadWeight            = errors.New("вес должен быть больше 0")
	ErrScoreOutOfRange      = errors.New("балл вне диапазона критерия")
	ErrCriterionNotInRubric = errors.New("критерия нет в рубрикаторе")
	ErrStudentNotFound      = errors.New("студент не найден")
	ErrEvaluationNotFound   = errors.New("оценка не найдена")
	ErrDuplicateEvaluation  = errors.New("оценка уже выставлена (используйте исправление)")
)
