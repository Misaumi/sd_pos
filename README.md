# ПОС: Преподавательская Оценка Студентов (Go)

## Запуск
```
go run ./src          # интерактивное меню
go run ./src demo     # сквозной сценарий
go test ./...         # тесты
```

## Архитектура: Ports and Adapters
```
include/domain       ядро: сущности, VO, инварианты, Domain Service, порты репозиториев
src/facade.go        сценарии (Facade) и порт TotalsExporter
src/memory.go        адаптер: хранилище в памяти
src/csv_exporter.go  адаптер: экспорт в CSV
src/console.go       адаптер: консольный интерфейс
src/container.go     composition root (DI)
src/main.go          точка входа
```
Зависимости направлены внутрь: файлы в `src` зависят от `include/domain`, домен не зависит ни от чего.

## Карта DDD
- Ubiquitous Language: Student, Criterion, Rubric, Evaluation, Score, итог
- Entity: Student, Criterion, Evaluation (у оценки идентичность — пара студент + критерий)
- Value Object: Score, неизменяем и создаётся только внутри [0; max]
- Агрегат: Rubric (корень) владеет критериями, Evaluation ссылается на студента и критерий по id
- Инварианты: балл в диапазоне, критерий из рубрикатора, max_score > 0, вес > 0, имя не пустое, одна оценка на пару студент-критерий
- Repository (порты): StudentRepository, RubricRepository, EvaluationRepository
- Domain Service: ScoreCalculator, итог = сумма(score/max * вес) / сумма весов * 100

## Паттерны GoF
- Фабрика (порождающий): `EvaluationFactory.Create`, единственное место создания Evaluation с проверкой критерия и диапазона
- Декоратор (структурный): `RoundedStrategy` округляет результат любой стратегии
- Фасад (структурный): `Facade`, консоль не знает про репозитории и фабрику
- Стратегия (поведенческий): `ScoringStrategy` (`WeightedAll`, `WeightedRatedOnly`), правило меняется в меню (пункт 11)

## Роль ИИ
См. AI_USAGE.md
