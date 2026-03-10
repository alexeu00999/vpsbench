# План: CPU бенчмарк

**Ветка:** `feature/cpu-benchmark`
**Дата:** 2026-03-10

## Settings

- **Testing:** yes
- **Logging:** verbose (DEBUG уровень для всех операций)
- **Docs:** no

## Roadmap Linkage

- **Milestone:** "CPU бенчмарк"
- **Rationale:** Реализация полноценного CPU бенчмарка с single-core и multi-core тестами — четвёртая веха проекта

## Контекст

Текущая реализация в `internal/cpu/cpu.go` — заглушка с примитивным `math.Sqrt` в цикле.
Необходимо заменить на реалистичные вычислительные нагрузки: целочисленная арифметика,
операции с плавающей точкой, криптографические вычисления. Добавить warmup-фазу,
улучшить точность замеров, калибровать baseline.

## Задачи

### Фаза 1: Вычислительные нагрузки

- [x] **Task #18:** Реализовать вычислительные workload'ы для CPU бенчмарка
  - `internal/cpu/workloads.go` — три workload'а: integer (сортировка, битовые операции), float (матрицы, тригонометрия), crypto (AES, SHA-256)
  - Каждый workload возвращает количество операций
  - Verbose slog логирование

### Фаза 2: Интеграция

- [x] **Task #19:** Переписать основной CPU бенчмарк с новыми workload'ами _(blocked by #18)_
  - `internal/cpu/cpu.go` — warmup-фаза, atomic счётчики, 3s длительность
  - Single-core и Multi-core тесты со смешанной нагрузкой

### Фаза 3: Калибровка и тесты (параллельно)

- [x] **Task #20:** Калибровать baseline значения для CPU в rating _(blocked by #19)_
  - `internal/rating/rating.go` — реалистичные baseline для Ryzen 9 5900X

- [x] **Task #21:** Написать тесты для CPU бенчмарка _(blocked by #19)_
  - `internal/cpu/cpu_test.go` — тесты бенчмарка (Name, Run, Cancel, scaling)
  - `internal/cpu/workloads_test.go` — тесты workload'ов

## Commit Plan

- **После Task #19:** `feat(cpu): implement realistic CPU benchmark workloads`
- **После Task #20 + #21:** `test(cpu): add tests and calibrate baseline values`
