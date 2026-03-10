package cpu

import (
	"context"
	"crypto/aes"
	"crypto/sha256"
	"log/slog"
	"math"
	"sort"
	"sync/atomic"
	"time"
)

// workloadResult содержит результат выполнения одного workload'а.
type workloadResult struct {
	Name string
	Ops  int64
}

// runMixedWorkload запускает все workload'ы последовательно в цикле на одной горутине
// до истечения duration. Возвращает суммарное количество операций.
// warmup — если true, результаты не учитываются (прогрев кэшей).
func runMixedWorkload(ctx context.Context, duration time.Duration, opsCounter *atomic.Int64) {
	slog.Debug("[cpu] mixed workload started", "duration", duration)
	deadline := time.Now().Add(duration)

	var localOps int64

	for time.Now().Before(deadline) {
		if ctx.Err() != nil {
			slog.Debug("[cpu] mixed workload cancelled by context")
			break
		}

		// Целочисленный workload
		localOps += workloadInteger()

		// Плавающая точка
		localOps += workloadFloat()

		// Крипто
		localOps += workloadCrypto()
	}

	opsCounter.Add(localOps)
	slog.Debug("[cpu] mixed workload finished", "local_ops", localOps)
}

// workloadInteger выполняет пакет целочисленных операций:
// сортировка массива + битовые операции.
// Возвращает количество выполненных операций.
func workloadInteger() int64 {
	const size = 256

	// Генерируем массив детерминированно (без rand для чистоты CPU теста)
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = (i*7 + 13) ^ (i << 3)
	}

	// Сортировка — классическая целочисленная нагрузка
	sort.Ints(data)

	// Битовые операции поверх отсортированного массива
	var acc int
	for i := 0; i < size; i++ {
		acc ^= data[i] << (uint(i) % 16)
		acc = (acc << 1) | (acc >> 31) // rotate left
	}

	// Предотвращаем оптимизацию компилятором
	sinkInt.Store(int64(acc))

	// size сортировки + size битовых операций
	return int64(size * 2)
}

// workloadFloat выполняет пакет операций с плавающей точкой:
// умножение матриц 4x4 + тригонометрические вычисления.
// Возвращает количество выполненных операций.
func workloadFloat() int64 {
	var ops int64

	// Умножение матриц 4x4 — типичная FP-нагрузка
	var a, b, c [4][4]float64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			a[i][j] = float64(i+1) * 1.1
			b[i][j] = float64(j+1) * 2.2
		}
	}

	// 10 итераций умножения
	for iter := 0; iter < 10; iter++ {
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				var sum float64
				for k := 0; k < 4; k++ {
					sum += a[i][k] * b[k][j]
				}
				c[i][j] = sum
			}
		}
		// Используем результат для следующей итерации
		a = c
		ops += 64 // 4*4*4 умножений и сложений
	}

	// Тригонометрические вычисления
	var trigResult float64
	for i := 0; i < 100; i++ {
		angle := float64(i) * 0.01 * math.Pi
		trigResult += math.Sin(angle) * math.Cos(angle) + math.Tan(angle*0.5)
	}
	ops += 300 // 100 итераций * 3 операции (sin, cos, tan)

	// Предотвращаем оптимизацию
	sinkFloat.Store(math.Float64bits(trigResult + c[0][0]))

	return ops
}

// workloadCrypto выполняет пакет криптографических операций:
// AES шифрование + SHA-256 хэширование.
// Возвращает количество выполненных операций.
func workloadCrypto() int64 {
	var ops int64

	// AES-128 шифрование блоков
	key := []byte("0123456789abcdef") // 16 байт = AES-128
	block, err := aes.NewCipher(key)
	if err != nil {
		slog.Error("[cpu] failed to create AES cipher", "error", err)
		return 0
	}

	plaintext := make([]byte, 16) // один AES блок
	ciphertext := make([]byte, 16)
	for i := range plaintext {
		plaintext[i] = byte(i * 7)
	}

	// 50 шифрований
	for i := 0; i < 50; i++ {
		block.Encrypt(ciphertext, plaintext)
		// Используем результат для следующей итерации
		copy(plaintext, ciphertext)
	}
	ops += 50

	// SHA-256 хэширование
	data := make([]byte, 256)
	for i := range data {
		data[i] = byte(i)
	}

	hash := data
	for i := 0; i < 50; i++ {
		h := sha256.Sum256(hash)
		hash = h[:]
	}
	ops += 50

	// Предотвращаем оптимизацию
	sinkInt.Store(int64(ciphertext[0]) + int64(hash[0]))

	return ops
}

// Sink-переменные предотвращают оптимизацию компилятором (dead code elimination).
var (
	sinkInt   atomic.Int64
	sinkFloat atomic.Uint64
)

// storeSinkFloat сохраняет float64 в atomic переменную через битовое преобразование.
func init() {
	// Инициализация sink — ничего не делает, но гарантирует что переменные используются.
	sinkInt.Store(0)
	sinkFloat.Store(0)
}
