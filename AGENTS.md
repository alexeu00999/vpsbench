# AGENTS.md

> Project map for AI agents. Keep this file up-to-date as the project evolves.

## Project Overview
Комплексный инструмент бенчмаркинга VPS и серверов. Запускает тесты производительности и выводит результаты в виде красивой цветной ASCII-графики в терминале.

## Tech Stack
- **Language:** Go
- **Framework:** Cobra (CLI) + BubbleTea/Lipgloss (TUI)
- **Architecture:** Modular Monolith

## Project Structure
```
vpsbench/
├── cmd/
│   └── vpsbench/
│       └── main.go              # Точка входа, инициализация Cobra
├── internal/
│   ├── bench/                   # Общий интерфейс и типы (Benchmark, Result)
│   │   ├── bench.go             # Интерфейс и основные структуры
│   │   └── runner.go            # Оркестратор запуска бенчмарков
│   ├── cpu/                     # Модуль CPU бенчмарка (single/multi core)
│   ├── ram/                     # Модуль RAM бенчмарка (read/write speed)
│   ├── disk/                    # Модуль Disk I/O бенчмарка (seq/random)
│   ├── network/                 # Модуль сетевого бенчмарка
│   ├── sysinfo/                 # Детект системы (CPU, RAM, OS, Location)
│   ├── rating/                  # Расчёт рейтинга (в разработке)
│   ├── ui/                      # TUI компоненты (BubbleTea)
│   └── output/                  # Форматирование вывода и цвета (Lipgloss)
├── docs/                        # Подробная документация (architecture, cli, etc.)
├── .ai-factory/                 # AI-контекст (DESCRIPTION, ARCHITECTURE, ROADMAP)
├── Dockerfile                   # Контейнеризация для универсального запуска
├── Makefile                     # Сборка, тесты, линтинг
├── go.mod                       # Зависимости Go
└── README.md                    # Описание проекта и инструкции
```

## Key Entry Points
| File | Purpose |
|------|---------|
| `cmd/vpsbench/main.go` | Main application entry point |
| `internal/bench/bench.go` | Central benchmark interface definition |
| `internal/bench/runner.go` | Orchestrator for concurrent benchmark execution |
| `.ai-factory/ARCHITECTURE.md` | Detailed architectural guidelines |

## Documentation
| Document | Path | Description |
|----------|------|-------------|
| README | README.md | Project landing page |
| Architecture | docs/architecture.md | High-level system design |
| CLI Interface | docs/cli.md | Commands, flags, and interaction |
| Configuration | docs/configuration.md | Environment variables and config |
| Getting Started | docs/getting-started.md | Installation and usage |

## AI Context Files
| File | Purpose |
|------|---------|
| AGENTS.md | This file — project structure map |
| .ai-factory/DESCRIPTION.md | Project specification and tech stack |
| .ai-factory/ARCHITECTURE.md | Architecture decisions and guidelines |
| .ai-factory/ROADMAP.md | Implementation progress tracking |
| .ai-factory/plans/ | Specific feature implementation plans |
