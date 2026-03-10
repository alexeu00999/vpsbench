[← Быстрый старт](getting-started.md) · [Назад к README](../README.md) · [CLI-интерфейс →](cli.md)

# Архитектура

## Обзор

SUPER-BENCH построен как модульный монолит на Go. Каждый бенчмарк — независимый модуль, реализующий общий интерфейс `Benchmark`. Единый бинарник без зависимостей.

## Структура проекта

```
vpsbench/
├── cmd/
│   └── vpsbench/
│       └── main.go              # Точка входа, Cobra CLI
├── internal/
│   ├── bench/                   # Общий интерфейс и оркестратор
│   │   ├── bench.go             # Интерфейс Benchmark, типы Result
│   │   └── runner.go            # Параллельный запуск бенчмарков
│   ├── cpu/                     # CPU бенчмарк (single/multi-core)
│   ├── ram/                     # RAM бенчмарк (R/W speed)
│   ├── disk/                    # Disk I/O (seq/rand read/write)
│   ├── network/                 # Сеть (ping, download/upload)
│   ├── sysinfo/                 # Детект системы (CPU, OS, локация)
│   ├── rating/                  # Расчёт процентов vs базовая конфигурация
│   ├── ui/                      # TUI: выбор компонентов + вывод результатов
│   │   ├── selector.go          # Интерактивные чекбоксы + таймер
│   │   └── results.go           # Отрисовка прогресс-баров
│   └── output/                  # Форматирование и цветовая палитра
├── Dockerfile
├── go.mod
└── Makefile
```

## Ключевые модули

| Модуль | Пакет | Назначение |
|--------|-------|-----------|
| Интерфейс | `internal/bench` | Общий интерфейс `Benchmark` и типы результатов |
| CPU | `internal/cpu` | Замер ops/s single-core и multi-core |
| RAM | `internal/ram` | Замер скорости чтения/записи памяти |
| Disk | `internal/disk` | Последовательный и случайный I/O |
| Network | `internal/network` | Пинг, скорость до EU/US/Asia (В разработке) |
| Sysinfo | `internal/sysinfo` | Информация о системе |
| Rating | `internal/rating` | Процентная оценка (В разработке) |
| UI | `internal/ui` | BubbleTea TUI — выбор и отображение (В разработке) |
| Output | `internal/output` | Lipgloss — цвета и форматирование (В разработке) |

## Поток данных (В разработке)

```
main.go
  → sysinfo.Detect()          # Собрать информацию о системе
  → ui.SelectComponents()     # Пользователь выбирает тесты (3 сек таймер)
  → runner.RunAll(selected)   # Параллельный запуск выбранных бенчмарков
  → rating.Calculate(results) # Расчёт процентов
  → ui.ShowResults(results)   # Отрисовка цветных прогресс-барок
```

## Правила зависимостей

- ✅ Модули бенчмарков зависят только от `bench/` (общие типы)
- ✅ `ui/` зависит от `bench/`, `rating/`, `output/`
- ❌ Модули бенчмарков НЕ зависят друг от друга
- ❌ `bench/` НЕ зависит ни от чего внутреннего

## Базовая конфигурация (100%)

| Компонент | Эталон |
|-----------|--------|
| CPU | 8-ядерный Ryzen 9 |
| Диск | NVMe Gen4 |
| Сеть | 10 Gbps |
| RAM | 16 GB DDR5 |

## Смотри также

- [Быстрый старт](getting-started.md) — установка и первый запуск
- [CLI-интерфейс](cli.md) — команды и флаги
