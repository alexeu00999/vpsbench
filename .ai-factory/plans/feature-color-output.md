# План: Цветной вывод результатов (Color Output)

Этот план описывает улучшение визуального вывода бенчмарка с использованием Lipgloss, настройку цветовой шкалы и поддержку режима без цветов (--no-color).

- **Ветка:** `feature/color-output`
- **Дата:** 2026-03-10
- **Статус:** Выполнено

## Settings
- **Testing:** Yes (Unit-тесты логики форматирования в `internal/output/output_test.go`)
- **Logging:** Verbose (DEBUG логи процесса рендеринга)
- **Docs:** Yes (Обновить `docs/cli.md`)

## Roadmap Linkage
- **Milestone:** "Цветной вывод результатов"
- **Rationale:** Обеспечение наглядного и профессионального представления результатов тестирования.

## Tasks

### Фаза 1: Поддержка --no-color
- [x] **Task 1: Implement No-Color Support**
  - Добавить глобальную (в рамках пакета) или передаваемую настройку `SetNoColor(bool)` в `internal/output/output.go`.
  - Обновить все функции рендеринга (`ColorForPercent`, `RenderProgressBar`, `RenderHeader`, `RenderModuleResult`), чтобы они игнорировали цвета, если установлен флаг.
  - **Logging:** DEBUG логирование статуса раскраски при инициализации.

### Фаза 2: Улучшение визуального стиля
- [x] **Task 2: Refine Layout and Borders**
  - Использовать `lipgloss.Border` для создания более красивых рамок в шапке и подвале отчета.
  - Настроить выравнивание колонок для более четкого отображения метрик.
  - Добавить жирный шрифт для заголовков модулей.
  - **Logging:** DEBUG логирование параметров макета (ширина, отступы).

### Фаза 3: Интеграция флага из Main
- [x] **Task 3: Connect Flag to Output**
  - Обновить `cmd/vpsbench/main.go`, чтобы он вызывал `output.SetNoColor(flagNoColor)` перед рендерингом результатов.
  - **Logging:** INFO лог о применении цветовой схемы.

### Фаза 4: Тестирование и документация
- [x] **Task 4: Update Unit Tests**
  - Обновить `internal/output/output_test.go`.
  - Добавить тесты на рендеринг с цветами и без них (проверка отсутствия ANSI-последовательностей).
- [x] **Task 5: Documentation Update**
  - Актуализировать `docs/cli.md` описанием флага `--no-color`.
  - Убедиться, что `AGENTS.md` актуален.

## Commit Plan
1. `feat(output): add support for no-color mode` (Task 1, 3)
2. `feat(output): improve layout and borders using lipgloss` (Task 2)
3. `test(output): add tests for no-color rendering` (Task 4)
4. `docs(output): update cli documentation` (Task 5)
