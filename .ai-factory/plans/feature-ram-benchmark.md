# План: RAM бенчмарк

**Ветка:** `feature/ram-benchmark`
**Дата:** 2026-03-10

## Settings

- **Testing:** yes
- **Logging:** verbose (DEBUG уровень для всех операций)
- **Docs:** no

## Roadmap Linkage

- **Milestone:** "RAM бенчмарк"
- **Rationale:** Реализация полноценного RAM бенчмарка с измерением скорости чтения/записи — пятая веха проекта

## Контекст

Текущая реализация в `internal/ram/ram.go` — заглушка с багом: пишет только 12.5% буфера
из-за stride `unsafe.Sizeof(uint64(0))`. Необходимо переписать с нуля: отдельные workload'ы
для sequential write и read, правильная работа через []uint64, калибровка baseline для DDR5.

## Задачи

### Фаза 1: Workload'ы

- [x] **Task #22:** Реализовать RAM workload'ы для чтения и записи
  - `internal/ram/workloads.go` — workloadWrite() и workloadRead() через []uint64
  - Буфер 128 MB, verbose slog логирование

### Фаза 2: Интеграция

- [x] **Task #23:** Переписать основной RAM бенчмарк _(blocked by #22)_
  - `internal/ram/ram.go` — два Result'а: Write и Read в MB/s
  - 3s на каждый workload, убрать старый measureSpeed()

### Фаза 3: Калибровка и тесты (параллельно)

- [x] **Task #24:** Калибровать baseline и обновить rating для RAM _(blocked by #23)_
  - `internal/rating/rating.go` — заменить "RAM:Speed (R/W)" на "RAM:Write" + "RAM:Read"
  - Baseline для DDR5-5600: ~40-50 GB/s

- [x] **Task #25:** Написать тесты для RAM бенчмарка _(blocked by #23)_
  - `internal/ram/ram_test.go` + `internal/ram/workloads_test.go`

## Commit Plan

- **После Task #23:** `feat(ram): implement RAM benchmark with sequential read/write`
- **После Task #24 + #25:** `test(ram): add tests and calibrate baseline values`
