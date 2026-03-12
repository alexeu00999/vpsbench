# План: Сетевой бенчмарк (Network Benchmark)

Этот план описывает реализацию модуля сетевого бенчмаркинга, включая измерение пинга, скорости загрузки и отдачи до серверов в Европе (EU), США (US) и Азии (Asia).

- **Ветка:** `feature/network-benchmark`
- **Дата:** 2026-03-10
- **Статус:** Выполнено

## Settings
- **Testing:** Yes (Unit tests in `internal/network/network_test.go`)
- **Logging:** Verbose (DEBUG logs for each measurement step)
- **Docs:** Yes (Update `docs/cli.md` and `docs/architecture.md`)

## Roadmap Linkage
- **Milestone:** "Сетевой бенчмарк"
- **Rationale:** Реализация основной функциональности тестирования сети для достижения паритета с заявленными возможностями VPSBENCH.

## Tasks

### Фаза 1: Определение серверов и структур данных
- [x] **Task 1: Define Network Measurement Structures**
  - Добавить структуры `Server` и `RegionResult` в `internal/network/network.go`.
  - Определить список тестовых серверов для EU, US, ASIA (например, Cloudflare, AWS, Hetzner).
  - **Logging:** Спискок серверов при инициализации.

### Фаза 2: Реализация измерений (Download/Upload)
- [x] **Task 2: Implement Download Measurement**
  - Реализовать функцию `measureDownload(ctx, url)`. Использовать `http.Client` для загрузки фиксированного объема данных (например, 10MB) и расчета скорости (Mbps).
  - **Logging:** DEBUG логи начала загрузки, времени завершения и рассчитанной скорости.
- [x] **Task 3: Implement Upload Measurement**
  - Реализовать функцию `measureUpload(ctx, url)`. Использовать `http.MethodPost` для отправки данных и расчета скорости (Mbps).
  - **Logging:** DEBUG логи начала отправки и финальной скорости.

### Фаза 3: Интеграция в модуль NetworkBench
- [x] **Task 4: Update NetworkBench.Run**
  - Обновить метод `Run` для последовательного (или параллельного) выполнения тестов для всех регионов.
  - Собрать результаты в `bench.ModuleResult`.
  - Реализовать автодетект публичного IP (если не сделано в sysinfo).
  - **Logging:** INFO логи старта/завершения модуля, DEBUG логи по каждому региону.

### Фаза 4: Тестирование и документация
- [x] **Task 5: Write Unit Tests**
  - Создать `internal/network/network_test.go`. Протестировать расчеты скорости и обработку ошибок (timeout, unreachable).
  - **Logging:** Использовать мок-серверы для тестов.
- [x] **Task 6: Documentation Checkpoint**
  - Обновить `docs/cli.md` (удалить "В разработке" для сети).
  - Обновить `docs/architecture.md` (удалить "В разработке").
  - Убедиться, что `AGENTS.md` актуален.

## Commit Plan
1. `feat(network): define test servers and structures` (Task 1)
2. `feat(network): implement download and upload measurements` (Tasks 2, 3)
3. `feat(network): integrate measurements into NetworkBench.Run` (Task 4)
4. `test(network): add unit tests and update documentation` (Tasks 5, 6)
