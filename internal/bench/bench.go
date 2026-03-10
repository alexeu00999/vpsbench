package bench

import "context"

// Result — результат одного измерения внутри модуля бенчмарка.
type Result struct {
	Name    string  `json:"name"`    // Название метрики: "Seq. Read", "Single-core"
	Value   float64 `json:"value"`   // Измеренное значение: 850, 1240
	Unit    string  `json:"unit"`    // Единица измерения: "MB/s", "ops/s", "IOPS", "ms"
	Percent int     `json:"percent"` // 0-100, рассчитывается модулем rating
}

// ModuleResult — агрегированный результат одного модуля бенчмарка.
type ModuleResult struct {
	Module  string   `json:"module"`  // Название модуля: "CPU", "RAM", "DISK", "NETWORK"
	Info    string   `json:"info"`    // Описание: "Intel Xeon E5-2690 v3 (4 Cores)"
	Results []Result `json:"results"` // Список измерений
	Err     error    `json:"-"`       // Игнорируем Err в JSON (через Results.Err можно обработать)
}

// Benchmark — интерфейс, который реализуют все модули бенчмарка.
type Benchmark interface {
	// Name возвращает имя модуля (используется для фильтрации и вывода).
	Name() string

	// Run запускает бенчмарк и возвращает результаты.
	// Контекст позволяет отменить выполнение.
	Run(ctx context.Context) ModuleResult
}

// Report — полная структура отчёта бенчмарка для JSON-вывода.
type Report struct {
	Timestamp     string         `json:"timestamp"`
	SystemInfo    SystemInfo     `json:"system_info"`
	ModuleResults []ModuleResult `json:"module_results"`
	OverallRating int            `json:"overall_rating"`
}

// SystemInfo — информация о системе для отчёта.
type SystemInfo struct {
	OSVersion string `json:"os_version"`
	Kernel    string `json:"kernel"`
	Arch      string `json:"arch"`
	CPUModel  string `json:"cpu_model"`
	CPUCores  int    `json:"cpu_cores"`
	RAM       string `json:"ram"`
	Disks     string `json:"disks"`
	Location  string `json:"location"`
	PublicIP  string `json:"public_ip"`
}
