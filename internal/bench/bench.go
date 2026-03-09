package bench

import "context"

// Result — результат одного измерения внутри модуля бенчмарка.
type Result struct {
	Name    string  // Название метрики: "Seq. Read", "Single-core"
	Value   float64 // Измеренное значение: 850, 1240
	Unit    string  // Единица измерения: "MB/s", "ops/s", "IOPS", "ms"
	Percent int     // 0-100, рассчитывается модулем rating
}

// ModuleResult — агрегированный результат одного модуля бенчмарка.
type ModuleResult struct {
	Module  string   // Название модуля: "CPU", "RAM", "DISK", "NETWORK"
	Info    string   // Описание: "Intel Xeon E5-2690 v3 (4 Cores)"
	Results []Result // Список измерений
	Err     error    // nil если модуль отработал без ошибок
}

// Benchmark — интерфейс, который реализуют все модули бенчмарка.
type Benchmark interface {
	// Name возвращает имя модуля (используется для фильтрации и вывода).
	Name() string

	// Run запускает бенчмарк и возвращает результаты.
	// Контекст позволяет отменить выполнение.
	Run(ctx context.Context) ModuleResult
}
