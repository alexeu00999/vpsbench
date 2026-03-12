# План: Система рейтинга (Rating System)

Этот план описывает доработку системы рейтинга для VPSBENCH, включая актуализацию эталонных значений (baseline), поддержку новых метрик и расчет цветовой шкалы для визуального вывода.

- **Ветка:** `feature/rating-system`
- **Дата:** 2026-03-10
- **Статус:** Выполнено

## Settings
- **Testing:** Yes (Unit tests in `internal/rating/rating_test.go`)
- **Logging:** Verbose (DEBUG логи процесса расчета)
- **Docs:** Yes (Обновить `docs/architecture.md`)

## Roadmap Linkage
- **Milestone:** "Система рейтинга"
- **Rationale:** Реализация логики оценки производительности относительно базовой конфигурации для последующего визуального вывода.

## Tasks

### Фаза 1: Актуализация Baseline
- [x] **Task 1: Update Baseline Values**
  - Обновить `DefaultBaseline` в `internal/rating/rating.go`.
  - Добавить точные значения для Disk (Seq. Read/Write, Rand 4K IOPS), Network (Ping, Download, Upload для регионов EU/US/ASIA) и RAM.
  - Использовать значения из `DESCRIPTION.md` как ориентир.
  - **Logging:** Вывод загруженного baseline в DEBUG.

### Фаза 2: Доработка алгоритма расчета
- [x] **Task 2: Refine Calculation Logic**
  - Разрешить значения > 100% (убрать жесткий Clamp до 100 в `Calculate`).
  - Добавить обработку новых сетевых метрик по регионам (Ping/Download/Upload).
  - Убедиться, что латентность (ms) рассчитывается корректно (инвертированная шкала).
  - **Logging:** Логирование рассчитанного процента для каждого `Result`.

### Фаза 3: Цветовая шкала и визуализация
- [x] **Task 3: Implement Color Rating**
  - Добавить функцию `GetColor(percent int) string` (или аналогичную) в `internal/rating/rating.go` или `internal/output/output.go`.
  - Реализовать шкалу из `DESCRIPTION.md`:
    - 0-20%: Красный
    - 21-40%: Оранжевый
    - 41-60%: Белый
    - 61-80%: Жёлтый
    - 81-100%+: Зелёный
  - **Logging:** DEBUG логирование выбранного цвета для рейтинга.

### Фаза 4: Тестирование и документация
- [x] **Task 4: Comprehensive Unit Tests**
  - Обновить `internal/rating/rating_test.go`.
  - Протестировать расчеты для всех типов метрик (скорость, латентность).
  - Протестировать расчет `OverallRating`.
  - Протестировать логику цветовой шкалы.
- [x] **Task 5: Documentation Update**
  - Отразить изменения в `docs/architecture.md`.
  - Убедиться, что `AGENTS.md` актуален.

## Commit Plan
1. `feat(rating): update baseline values and remove 100% clamp` (Task 1, 2)
2. `feat(rating): implement color scale logic` (Task 3)
3. `test(rating): add comprehensive unit tests` (Task 4)
4. `docs(rating): update architecture documentation` (Task 5)
