# Implementation Plan: Скелет приложения

Branch: feature/app-skeleton
Created: 2026-03-10

## Settings
- Testing: yes
- Logging: verbose
- Docs: no

## Roadmap Linkage
Milestone: "Скелет приложения"
Rationale: Базовая структура Go-проекта, интерфейсы и оркестратор — фундамент для всех модулей

## Commit Plan
- **Commit 1** (после задач 1-2): "feat: init Go module with bench interface and types"
- **Commit 2** (после задач 3-7): "feat: add runner, module stubs, sysinfo, rating, output"
- **Commit 3** (после задачи 8): "feat: add Cobra CLI entry point with flag parsing"
- **Commit 4** (после задачи 9): "test: add unit tests for bench, runner, rating, output"

## Tasks

### Phase 1: Фундамент
- [x] Task 1: Инициализация Go-модуля и зависимостей (go mod init, cobra, bubbletea, lipgloss, структура директорий)
- [x] Task 2: Интерфейс Benchmark и типы данных (bench.go — Result, ModuleResult, Benchmark interface)

### Phase 2: Модули
- [x] Task 3: Оркестратор Runner (runner.go — RunAll параллельно, RunSelected) (depends on 2)
- [x] Task 4: Заглушки модулей бенчмарков — cpu, ram, disk, network с фейковыми данными (depends on 2)
- [x] Task 5: Заглушка sysinfo — определение OS, Arch, NumCPU через runtime
- [x] Task 6: Заглушка rating — расчёт процентов vs базовая конфигурация (depends on 2)
- [x] Task 7: Заглушка output — Lipgloss прогресс-бары, цветовая шкала, рендер результатов (depends on 2, 5)
<!-- Commit checkpoint: tasks 1-7 -->

### Phase 3: Интеграция
- [x] Task 8: Cobra CLI и точка входа main.go — флаги, инициализация, запуск пайплайна (depends on 3, 4, 5, 6, 7)
<!-- Commit checkpoint: task 8 -->

### Phase 4: Тесты
- [x] Task 9: Unit-тесты для bench, runner, rating, output (depends on 2, 3, 6, 7)
<!-- Commit checkpoint: task 9 -->
