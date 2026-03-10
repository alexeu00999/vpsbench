# План: Детект системы

**Ветка:** `feature/system-detect`
**Дата:** 2026-03-10

## Настройки

- **Тесты:** Да
- **Логирование:** Verbose (DEBUG)
- **Документация:** Нет (warn-only)

## Привязка к роадмапу

- **Веха:** "Детект системы"
- **Обоснование:** Реализация полного определения характеристик системы — CPU модель, RAM, диски, ОС, GeoIP локация

## Контекст

Текущий `internal/sysinfo/sysinfo.go` содержит заглушки:
- CPUModel = "Detection pending"
- RAMTotal = 0
- Location = "Unknown"

Нужно реализовать реальный детект через системные файлы (/proc/cpuinfo, /proc/meminfo, /sys/block/, /etc/os-release) на Linux и sysctl/diskutil/sw_vers на macOS. Локация определяется через GeoIP API.

Также нужно исправить архитектурное нарушение: output.go импортирует sysinfo напрямую.

## Задачи

### Фаза 1: Платформенные детекторы (параллельные задачи)

- [x] **Задача 1:** Детект CPU — модель и характеристики
  - `internal/sysinfo/detect_linux.go`, `internal/sysinfo/detect_darwin.go`
  - Linux: /proc/cpuinfo, macOS: sysctl machdep.cpu.brand_string

- [x] **Задача 2:** Детект RAM — общий объём памяти
  - `internal/sysinfo/detect_linux.go`, `internal/sysinfo/detect_darwin.go`
  - Linux: /proc/meminfo, macOS: sysctl hw.memsize

- [x] **Задача 3:** Детект дисков — список и характеристики
  - `internal/sysinfo/sysinfo.go` (struct DiskInfo), `detect_linux.go`, `detect_darwin.go`
  - Linux: /sys/block/, macOS: diskutil

- [x] **Задача 4:** Детект ОС — дистрибутив и версия
  - `internal/sysinfo/sysinfo.go` (поля OSVersion, Kernel), `detect_linux.go`, `detect_darwin.go`
  - Linux: /etc/os-release + uname, macOS: sw_vers + uname

- [x] **Задача 5:** Детект локации — GeoIP по внешнему IP
  - `internal/sysinfo/geoip.go`
  - HTTP GET к ip-api.com, таймаут 5с, graceful fallback

### Фаза 2: Интеграция

- [x] **Задача 6:** Обновить Detect() и интеграция с main.go
  - Вызвать все детекторы, передать ctx, прокинуть sysinfo в бенчмарк-модули
  - Блокируется задачами 1-5

- [x] **Задача 7:** Обновить output для новых полей SystemInfo
  - Отобразить все новые поля в хедере, убрать прямой импорт sysinfo
  - Блокируется задачей 6

### Фаза 3: Тесты

- [x] **Задача 8:** Тесты для sysinfo детекторов
  - `internal/sysinfo/sysinfo_test.go`
  - Блокируется задачей 6

## План коммитов

1. **После задач 1-5:** `feat(sysinfo): add platform-specific system detection`
2. **После задач 6-7:** `feat(sysinfo): integrate detectors with main pipeline and output`
3. **После задачи 8:** `test(sysinfo): add tests for system detection`
