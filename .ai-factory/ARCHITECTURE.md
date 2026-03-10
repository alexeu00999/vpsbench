# Архитектура: Modular Monolith

## Обзор
Проект использует паттерн **Modular Monolith** (Модульный Монолит). Весь функционал собран в единый бинарник Go, но внутри разделен на четко очерченные модули (CPU, RAM, Disk, Network). Это обеспечивает простоту деплоя и при этом позволяет поддерживать высокую степень независимости компонентов.

## Решение
- **Тип проекта:** CLI-инструмент для бенчмаркинга
- **Стек:** Go + Cobra + BubbleTea/Lipgloss
- **Ключевой фактор:** Независимость бенчмарков при едином рантайме.

## Folder Structure
```
vpsbench/
├── cmd/
│   └── vpsbench/
│       └── main.go              # Точка входа, инициализация Cobra
├── internal/
│   ├── bench/                   # Общий интерфейс и типы (Ядро системы)
│   │   ├── bench.go             # Интерфейс Benchmark, типы Result
│   │   └── runner.go            # Оркестратор запуска бенчмарков
│   ├── cpu/                     # Модуль CPU бенчмарка
│   ├── ram/                     # Модуль RAM бенчмарка
│   ├── disk/                    # Модуль Disk I/O бенчмарка
│   ├── network/                 # Модуль сетевого бенчмарка
│   ├── sysinfo/                 # Детект системы (информация об окружении)
│   ├── rating/                  # Расчёт баллов (логика оценки)
│   ├── ui/                      # Интерактивный TUI (BubbleTea)
│   └── output/                  # Форматирование и цвета (Lipgloss)
├── .ai-factory/                 # AI-контекст и планы
├── Dockerfile                   # Контейнеризация
└── Makefile                     # Автоматизация сборки
```

## Dependency Rules
Все модули бенчмарков зависят только от `bench/` (общие интерфейсы).

- ✅ `internal/cpu` → `internal/bench`
- ✅ `internal/ui` → `internal/bench`, `internal/rating`, `internal/output`
- ✅ `internal/rating` → `internal/bench` (типы результатов)
- ❌ Модули бенчмарков (cpu, ram, disk) НЕ зависят друг от друга
- ❌ `internal/bench` НЕ зависит от конкретных реализаций бенчмарков

## Layer/Module Communication
- **Интерфейс Benchmark:** Все тесты реализуют `bench.Benchmark`.
- **Оркестратор:** `bench.Runner` запускает выбранные тесты параллельно через Goroutines.
- **Результаты:** Каждый модуль возвращает `bench.ModuleResult`, который затем обрабатывается `rating` и `ui`.

## Key Principles
1. **Low Coupling:** Модули бенчмарков не знают о существовании друг друга.
2. **High Cohesion:** Каждый модуль отвечает только за один тип тестов.
3. **Single Responsibility:** `runner` только запускает, `rating` только считает, `ui` только отображает.

## Code Examples

### Интерфейс бенчмарка
```go
// internal/bench/bench.go
package bench

import "context"

type Benchmark interface {
    Name() string
    Run(ctx context.Context) ModuleResult
}
```

### Пример реализации (CPU)
```go
// internal/cpu/cpu.go
package cpu

import (
    "context"
    "vpsbench/internal/bench"
)

type CPUBench struct{}

func (c *CPUBench) Name() string { return "CPU" }
func (c *CPUBench) Run(ctx context.Context) bench.ModuleResult {
    // Логика теста...
    return bench.ModuleResult{Module: "CPU", Results: results}
}
```

## Anti-Patterns
- ❌ **Циклические зависимости:** Импорт `cpu` в `ram` или наоборот строго запрещен.
- ❌ **Глобальное состояние:** Конфигурации должны передаваться через конструкторы (`New()`).
- ❌ **Бизнес-логика в main:** Вся логика должна находиться в `internal/`.
