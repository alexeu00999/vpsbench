# План: Интерактивный TUI (Interactive Selector)

Этот план описывает реализацию интерактивного меню выбора компонентов бенчмарка с использованием BubbleTea.

- **Ветка:** `feature/interactive-tui`
- **Дата:** 2026-03-10
- **Статус:** Выполнено

## Settings
- **Testing:** Yes (Unit-тесты логики выбора в `internal/ui/selector_test.go`)
- **Logging:** Verbose (DEBUG логи событий TUI)
- **Docs:** Yes (Обновить `docs/cli.md` и `docs/getting-started.md`)

## Roadmap Linkage
- **Milestone:** "Интерактивный TUI"
- **Rationale:** Обеспечение удобного пользовательского интерфейса для выбора тестов перед их запуском.

## Tasks

### Фаза 1: Разработка BubbleTea модели
- [x] **Task 1: Create TUI Model and Views**
  - Создать `internal/ui/selector.go`.
  - Определить структуру `Model` для BubbleTea (список доступных бенчмарков, состояние выбора, таймер).
  - Реализовать метод `Init`, `Update` (обработка клавиш Space, Enter, Up/Down) и `View`.
  - Добавить таймер на 3 секунды, который автоматически завершает выбор с текущими настройками.
  - **Logging:** DEBUG логирование инициализации модели и выбора пользователя.

### Фаза 2: Интеграция в Main
- [x] **Task 2: Integrate UI into main.go**
  - Обновить `runBenchmark` в `cmd/vpsbench/main.go`.
  - Если не заданы флаги фильтрации и не установлен `--auto`, вызывать `ui.SelectComponents`.
  - Передавать список доступных бенчмарков в UI и получать список выбранных.
  - **Logging:** INFO лог о запуске интерактивного режима.

### Фаза 3: Поддержка выбора дисков (Advanced)
- [x] **Task 3: Implement Disk Selection**
  - Доработать модель, чтобы она позволяла выбирать конкретные диски из обнаруженных.
  - Обновить `internal/ui/selector.go` для отображения списка дисков под пунктом "DISK".
  - **Logging:** DEBUG логирование выбранных дисков.

### Фаза 4: Тестирование и документация
- [x] **Task 4: Unit Tests for UI Logic**
  - Создать `internal/ui/selector_test.go`.
  - Протестировать логику переключения чекбоксов (Update) без рендеринга.
- [x] **Task 5: Documentation Update**
  - Описать интерактивный режим в `docs/cli.md`.
  - Обновить `docs/getting-started.md` примером запуска.
  - Актуализировать `AGENTS.md`.

## Commit Plan
1. `feat(ui): implement basic BubbleTea selector with timer` (Task 1)
2. `feat(ui): integrate selector into main application flow` (Task 2)
3. `feat(ui): add per-disk selection support` (Task 3)
4. `test(ui): add unit tests for selector logic` (Task 4)
5. `docs(ui): update cli and getting started guides` (Task 5)
