# План: CLI флаги и режимы

Ветка: `feature/cli-flags-modes`
Дата: 2026-03-10

## Обзор
Реализация полноценной поддержки CLI флагов для автоматизации и интеграции. Добавление вывода в формате JSON, корректное управление цветами и возможность точечного запуска тестов.

## Настройки (Settings)
- Тесты: Да (автоматические тесты для JSON и фильтрации)
- Логи: Verbose (через `slog.Debug`)
- Документация: Обновление README и docs/cli.md
- Роадмап: Связано с вехой "CLI флаги и режимы"

## Roadmap Linkage
- Веха: "CLI флаги и режимы"
- Обоснование: Переход от чисто интерактивного инструмента к полноценной CLI-утилите для автоматизации.

## Задачи (Tasks)

### Фаза 1: Инфраструктура вывода
- [x] **Задача 1.1: Реализация JSON-структур в `internal/bench`**
    - Создать структуру `Report` (System info + ModuleResults + OverallRating).
    - Файлы: `internal/bench/bench.go`
    - Логи: `slog.Debug("Creating report structure", "module", "bench")`
- [x] **Задача 1.2: Поддержка вывода JSON в `main.go`**
    - Реализовать сериализацию `Report` при установленном флаге `--json`.
    - Убедиться, что при `--json` не выводится лишний текст в stdout.
    - Файлы: `cmd/vpsbench/main.go`
    - Логи: `slog.Info("JSON output enabled", "flag", "--json")`
- [x] **Задача 1.3: Улучшение `internal/output` для `--no-color`**
    - Убедиться, что все Lipgloss-стили корректно обрабатывают `SetNoColor`.
    - Файлы: `internal/output/output.go`
    - Логи: `slog.Debug("Setting color mode", "no_color", noColor)`

### Фаза 2: Логика фильтрации и режимы
- [x] **Задача 2.1: Рефакторинг логики фильтрации бенчмарков**
    - Выделить логику выбора тестов на основе флагов и интерактивного ввода в отдельную функцию или компонент.
    - Файлы: `cmd/vpsbench/main.go`, `internal/bench/runner.go`
    - Логи: `slog.Debug("Filtering benchmarks", "count", len(selected))`
- [x] **Задача 2.2: Реализация флага `--auto`**
    - Полный пропуск интерактивного меню и запуск всех тестов с дефолтными настройками дисков.
    - Файлы: `cmd/vpsbench/main.go`
    - Логи: `slog.Info("Auto mode enabled, skipping interactive TUI")`

### Фаза 3: Тестирование и Валидация
- [x] **Задача 3.1: Тесты для JSON-маршалинга**
    - Написать тесты, проверяющие корректность структуры JSON.
    - Файлы: `internal/bench/bench_test.go`
- [x] **Задача 3.2: Тесты для логики фильтрации**
    - Проверить, что флаги `--cpu`, `--ram` и т.д. правильно выбирают нужные модули.
    - Файлы: `internal/bench/runner_test.go`

### Фаза 4: Документация и Финализация
- [x] **Задача 4.1: Обновление `README.md` и `docs/cli.md`**
    - Описать все доступные флаги и примеры их использования.
    - Файлы: `README.md`, `docs/cli.md`
- [x] **Задача 4.2: Финальная проверка и обновление `ROADMAP.md`**
    - Отметить веху "CLI флаги и режимы" as завершенную.
    - Файлы: `ROADMAP.md`

## План коммитов (Commit Plan)
1. `feat(bench): add Report structure for JSON output` (Задачи 1.1, 1.2)
2. `feat(output): support no-color mode in Lipgloss` (Задача 1.3)
3. `refactor(cli): extract benchmark filtering logic` (Задачи 2.1, 2.2)
4. `test(cli): add tests for JSON output and benchmark filtering` (Задачи 3.1, 3.2)
5. `docs(cli): document flags and auto mode` (Задачи 4.1, 4.2)
