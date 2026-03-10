package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/vpsbench/internal/bench"
	"github.com/user/vpsbench/internal/cpu"
	"github.com/user/vpsbench/internal/disk"
	"github.com/user/vpsbench/internal/network"
	"github.com/user/vpsbench/internal/output"
	"github.com/user/vpsbench/internal/ram"
	"github.com/user/vpsbench/internal/rating"
	"github.com/user/vpsbench/internal/sysinfo"
)

var (
	flagAuto    bool
	flagJSON    bool
	flagNoColor bool
	flagCPU     bool
	flagRAM     bool
	flagDisk    bool
	flagNetwork bool
)

func main() {
	setupLogger()

	rootCmd := &cobra.Command{
		Use:   "vpsbench",
		Short: "SUPER-BENCH — комплексный бенчмарк VPS и серверов",
		Long:  "Тестирование производительности CPU, RAM, Disk I/O и сети с красивым цветным выводом.",
		RunE:  runBenchmark,
	}

	rootCmd.Flags().BoolVar(&flagAuto, "auto", false, "Пропустить интерактивное меню, запустить всё по умолчанию")
	rootCmd.Flags().BoolVar(&flagJSON, "json", false, "Вывод в формате JSON")
	rootCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "Без цветов (для перенаправления в файл)")
	rootCmd.Flags().BoolVar(&flagCPU, "cpu", false, "Только CPU тест")
	rootCmd.Flags().BoolVar(&flagRAM, "ram", false, "Только RAM тест")
	rootCmd.Flags().BoolVar(&flagDisk, "disk", false, "Только Disk I/O тест")
	rootCmd.Flags().BoolVar(&flagNetwork, "network", false, "Только сетевой тест")

	if err := rootCmd.Execute(); err != nil {
		slog.Error("[main] command failed", "error", err)
		os.Exit(1)
	}
}

func runBenchmark(cmd *cobra.Command, args []string) error {
	slog.Info("[main] starting SUPER-BENCH")

	// Определяем систему
	ctx := context.Background()
	info, err := sysinfo.Detect(ctx)
	if err != nil {
		slog.Error("[main] system detection failed", "error", err)
		return fmt.Errorf("system detection: %w", err)
	}

	// Конвертируем sysinfo.DiskInfo в disk.DiskTarget для бенчмарка
	var diskTargets []disk.DiskTarget
	for _, d := range info.Disks {
		diskTargets = append(diskTargets, disk.DiskTarget{
			Device: d.Device,
			Type:   d.Type,
			Path:   "", // empty = use default temp dir (mount point resolution is TODO)
		})
	}
	if len(diskTargets) > 0 {
		slog.Debug("[main] disk targets from sysinfo", "count", len(diskTargets))
	} else {
		slog.Debug("[main] no disks detected, disk benchmark will use default temp dir")
	}

	// Собираем бенчмарки
	allBenchmarks := []bench.Benchmark{
		cpu.New(),
		ram.New(),
		disk.New(disk.WithTargets(diskTargets)),
		network.New(),
	}

	// Фильтрация по флагам
	var selected []bench.Benchmark
	hasFilter := flagCPU || flagRAM || flagDisk || flagNetwork

	if hasFilter {
		slog.Debug("[main] filtering benchmarks by flags", "cpu", flagCPU, "ram", flagRAM, "disk", flagDisk, "network", flagNetwork)
		for _, b := range allBenchmarks {
			switch b.Name() {
			case "CPU":
				if flagCPU {
					selected = append(selected, b)
				}
			case "RAM":
				if flagRAM {
					selected = append(selected, b)
				}
			case "DISK":
				if flagDisk {
					selected = append(selected, b)
				}
			case "NETWORK":
				if flagNetwork {
					selected = append(selected, b)
				}
			}
		}
	} else {
		selected = allBenchmarks
	}

	slog.Info("[main] running benchmarks", "count", len(selected))

	// Запускаем
	runner := bench.NewRunner(selected)
	results := runner.RunAll(ctx)

	// Обогащаем Info из sysinfo
	for i := range results {
		switch results[i].Module {
		case "CPU":
			results[i].Info = fmt.Sprintf("%s (%d Cores)", info.CPUModel, info.CPUCores)
		case "RAM":
			results[i].Info = sysinfo.FormatRAM(info.RAMTotal)
		case "DISK":
			// Info is set by DiskBench itself based on targets
			if results[i].Info == "" || results[i].Info == "Default" {
				results[i].Info = sysinfo.FormatDisks(info.Disks)
			}
		case "NETWORK":
			if info.PublicIP != "" {
				results[i].Info = fmt.Sprintf("IP: %s | %s", info.PublicIP, info.Location)
			} else {
				results[i].Info = info.Location
			}
		}
	}

	// Рассчитываем рейтинг
	baseline := rating.DefaultBaseline()
	results = rating.Calculate(results, baseline)
	overall := rating.OverallRating(results)

	// Выводим результаты
	if flagJSON {
		// TODO: JSON вывод
		slog.Info("[main] JSON output not yet implemented")
		fmt.Println("{\"error\": \"JSON output not yet implemented\"}")
		return nil
	}

	headerData := output.HeaderData{
		OSVersion: info.OSVersion,
		Kernel:    info.Kernel,
		Arch:      info.Arch,
		CPUModel:  info.CPUModel,
		CPUCores:  info.CPUCores,
		RAM:       sysinfo.FormatRAM(info.RAMTotal),
		Disks:     sysinfo.FormatDisks(info.Disks),
		Location:  info.Location,
		PublicIP:  info.PublicIP,
	}

	fmt.Println()
	fmt.Println(output.RenderHeader(headerData, overall))
	for _, r := range results {
		fmt.Print(output.RenderModuleResult(r))
	}
	fmt.Println(strings.Repeat("-", 70))
	fmt.Println(output.RenderFooter())
	fmt.Println()

	slog.Info("[main] benchmark complete", "overall_rating", overall)
	return nil
}

func setupLogger() {
	level := slog.LevelInfo

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
	slog.Debug("[main] logger initialized", "level", level)
}
