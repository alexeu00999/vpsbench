package cpu

import (
	"testing"
)

func TestWorkloadInteger(t *testing.T) {
	ops := workloadInteger()
	if ops <= 0 {
		t.Errorf("workloadInteger() returned %d, want > 0", ops)
	}
	t.Logf("workloadInteger: %d ops", ops)
}

func TestWorkloadFloat(t *testing.T) {
	ops := workloadFloat()
	if ops <= 0 {
		t.Errorf("workloadFloat() returned %d, want > 0", ops)
	}
	t.Logf("workloadFloat: %d ops", ops)
}

func TestWorkloadCrypto(t *testing.T) {
	ops := workloadCrypto()
	if ops <= 0 {
		t.Errorf("workloadCrypto() returned %d, want > 0", ops)
	}
	t.Logf("workloadCrypto: %d ops", ops)
}

func TestWorkloadsDeterministic(t *testing.T) {
	// Запускаем каждый workload дважды — результат должен быть одинаковым
	// (детерминированные данные на входе).
	int1 := workloadInteger()
	int2 := workloadInteger()
	if int1 != int2 {
		t.Errorf("workloadInteger not deterministic: %d vs %d", int1, int2)
	}

	float1 := workloadFloat()
	float2 := workloadFloat()
	if float1 != float2 {
		t.Errorf("workloadFloat not deterministic: %d vs %d", float1, float2)
	}

	crypto1 := workloadCrypto()
	crypto2 := workloadCrypto()
	if crypto1 != crypto2 {
		t.Errorf("workloadCrypto not deterministic: %d vs %d", crypto1, crypto2)
	}
}
