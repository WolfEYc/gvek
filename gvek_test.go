package gvek_test

import (
	"os"
	"testing"

	"github.com/WolfEYc/gvek"
)

func TestMain(m *testing.M) {
	gvek.Init()
	code := m.Run()
	os.Exit(code)
}

func TestAddf32(t *testing.T) {
	a_slice := []float32{2, 3, 8, 9, 8, 7, 3, 5, 12}
	b_slice := []float32{3, 8, 7, 2, 4, 5, 9, 6, -7}
	c_actual := make([]float32, len(a_slice))
	gvek.Add_f32(c_actual, a_slice, b_slice)

	c_expect := []float32{5, 11, 15, 11, 12, 12, 12, 11, 5}
	if len(c_expect) != len(c_actual) {
		t.Errorf("len(c_expect)=%d len(c_actual)=%d", len(c_expect), len(c_actual))
	}
	for i, expected := range c_expect {
		actual := c_actual[i]
		diff := actual - expected
		if diff < 0 {
			diff = -diff
		}
		if diff > 0.01 {
			t.Errorf("i=%d, expected=%.2f, actual=%.2f", i, expected, actual)
			t.Logf("actual=%v", c_actual)
		}
	}
}
