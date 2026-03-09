# Архитектура: Modular Monolith

## Обзор
Проект использует модульный монолит — единый бинарник Go с чёткими границами между модулями. Каждый бенчмарк (CPU, RAM, Disk, Network) — это независимый модуль, реализующий общий интерфейс `Benchmark`. Модули не знают друг о друге, общаются только через общие типы данных.

Этот паттерн выбран потому что:
- CLI-инструмент собирается в один бинарник — никакой микросервисной оркестрации
- Каждый бенчмарк логически независим — идеально ложится на модули
- Простота добавления новых тестов — реализуешь интерфейс, регистрируешь модуль
- В V2 модуль отправки результатов на сервер добавляется как ещё один модуль

## Решение
- **Тип проекта:** CLI-инструмент для бенчмаркинга
- **Стек:** Go + Cobra + BubbleTea/Lipgloss
- **Ключевой фактор:** Независимые модули бенчмарков + единый бинарник

## Структура папок
```
vpsbench/
├── cmd/
│   └── vpsbench/
│       └── main.go              # Точка входа, инициализация Cobra
├── internal/
│   ├── bench/                   # Общий интерфейс и типы
│   │   ├── bench.go             # Интерфейс Benchmark, типы Result
│   │   └── runner.go            # Оркестратор запуска бенчмарков
│   ├── cpu/                     # Модуль CPU бенчмарка
│   │   └── cpu.go
│   ├── ram/                     # Модуль RAM бенчмарка
│   │   └── ram.go
│   ├── disk/                    # Модуль Disk I/O бенчмарка
│   │   └── disk.go
│   ├── network/                 # Модуль сетевого бенчмарка
│   │   └── network.go
│   ├── sysinfo/                 # Детект системы (CPU модель, ОС, локация)
│   │   └── sysinfo.go
│   ├── rating/                  # Расчёт процентов относительно базовой конфигурации
│   │   └── rating.go
│   ├── ui/                      # Интерактивный TUI (BubbleTea)
│   │   ├── selector.go          # Выбор компонентов чекбоксами + таймер
│   │   └── results.go           # Отрисовка результатов с прогресс-барами
│   └── output/                  # Форматирование и цвета (Lipgloss)
│       └── output.go
├── Dockerfile
├── go.mod
├── go.sum
└── Makefile
```

## Правила зависимостей

Все модули бенчмарков зависят только от `bench/` (общий интерфейс).

```
cmd/vpsbench/main.go
    └── internal/bench/runner.go     (оркестратор)
            ├── internal/cpu/        (реализует bench.Benchmark)
            ├── internal/ram/        (реализует bench.Benchmark)
            ├── internal/disk/       (реализует bench.Benchmark)
            ├── internal/network/    (реализует bench.Benchmark)
            └── internal/sysinfo/    (системная информация)
    └── internal/ui/                 (TUI: выбор → запуск → результаты)
            ├── internal/bench/      (типы и интерфейсы)
            ├── internal/rating/     (расчёт процентов)
            └── internal/output/     (форматирование)
```

- ✅ Модули бенчмарков → `bench/` (общие типы и интерфейс)
- ✅ `ui/` → `bench/`, `rating/`, `output/`
- ✅ `rating/` → `bench/` (типы результатов)
- ✅ `output/` → `bench/` (типы результатов)
- ❌ Модули бенчмарков НЕ зависят друг от друга (cpu не знает о ram)
- ❌ `bench/` НЕ зависит ни от чего внутреннего
- ❌ `rating/` и `output/` НЕ зависят от конкретных модулей бенчмарков

## Коммуникация между модулями

Модули не общаются напрямую. Оркестратор (`bench/runner.go`) управляет потоком:

1. `sysinfo` собирает информацию о системе
2. `ui/selector` показывает выбор компонентов
3. `runner` запускает выбранные бенчмарки параллельно
4. Каждый бенчмарк возвращает `bench.Result`
5. `rating` рассчитывает проценты
6. `ui/results` + `output` рендерят результаты

## Ключевые принципы

1. **Один интерфейс — много реализаций.** Все бенчмарки реализуют `bench.Benchmark`. Добавление нового теста = новый пакет + регистрация.
2. **Нет глобального состояния.** Конфигурация передаётся явно через структуры. Никаких `init()` с побочными эффектами.
3. **Graceful degradation.** Если модуль не может запуститься (нет сети, нет диска) — он возвращает ошибку, остальные работают.
4. **internal/ — всё приватно.** Только `cmd/` является публичной точкой входа. Никаких `pkg/` — инструмент не библиотека.

## Примеры кода

### Интерфейс бенчмарка
```go
// internal/bench/bench.go
package bench

import "context"

// Result — результат одного измерения
type Result struct {
    Name     string  // "Seq. Read", "Single-core"
    Value    float64 // 850, 1240
    Unit     string  // "MB/s", "ops/s", "IOPS", "ms"
    Percent  int     // 0-100, рассчитывается rating
}

// ModuleResult — результаты одного модуля
type ModuleResult struct {
    Module  string   // "CPU", "RAM", "DISK", "NETWORK"
    Info    string   // "Intel Xeon E5-2690 v3 (4 Cores)"
    Results []Result
    Err     error    // nil если всё ок
}

// Benchmark — интерфейс для всех модулей бенчмарка
type Benchmark interface {
    Name() string
    Run(ctx context.Context) ModuleResult
}
```

### Реализация модуля
```go
// internal/cpu/cpu.go
package cpu

import (
    "context"
    "vpsbench/internal/bench"
)

type CPUBench struct{}

func New() *CPUBench {
    return &CPUBench{}
}

func (c *CPUBench) Name() string {
    return "CPU"
}

func (c *CPUBench) Run(ctx context.Context) bench.ModuleResult {
    result := bench.ModuleResult{Module: "CPU"}

    // Детект CPU
    result.Info = detectCPUModel()

    // Single-core тест
    single := runSingleCore(ctx)
    result.Results = append(result.Results, bench.Result{
        Name:  "Single-core",
        Value: single,
        Unit:  "ops/s",
    })

    // Multi-core тест
    multi := runMultiCore(ctx)
    result.Results = append(result.Results, bench.Result{
        Name:  "Multi-core",
        Value: multi,
        Unit:  "ops/s",
    })

    return result
}
```

### Оркестратор
```go
// internal/bench/runner.go
package bench

import (
    "context"
    "sync"
)

type Runner struct {
    benchmarks []Benchmark
}

func NewRunner(benchmarks []Benchmark) *Runner {
    return &Runner{benchmarks: benchmarks}
}

func (r *Runner) RunAll(ctx context.Context) []ModuleResult {
    results := make([]ModuleResult, len(r.benchmarks))
    var wg sync.WaitGroup

    for i, b := range r.benchmarks {
        wg.Add(1)
        go func(idx int, bm Benchmark) {
            defer wg.Done()
            results[idx] = bm.Run(ctx)
        }(i, b)
    }

    wg.Wait()
    return results
}
```

## Анти-паттерны

- ❌ **Не импортировать модули бенчмарков друг в друга.** `cpu` не должен знать о `disk`. Если нужна общая логика — выносить в `bench/`.
- ❌ **Не использовать глобальные переменные.** Никаких `var globalConfig Config`. Передавать через параметры.
- ❌ **Не класть бизнес-логику в `cmd/`.** `main.go` только инициализирует и связывает модули.
- ❌ **Не делать `pkg/`.** Это CLI-инструмент, а не библиотека. Всё в `internal/`.
