package gvek_test

import (
	"os"
	"testing"

	"github.com/WolfEYc/gvek"
)

func TestMain(m *testing.M) {
	gvek.Init(true)
	code := m.Run()
	num_active := gvek.List_active_ctxs()
	if num_active != 0 {
		os.Exit(1)
	}
	os.Exit(code)
}

func TestAddf32(t *testing.T) {
	a_slice := []float32{2, 3, 8, 9, 8, 7, 3, 5, 12}
	b_slice := []float32{3, 8, 7, 2, 4, 5, 9, 6, -7}
	ctx := gvek.Make_Ctx_Smart(4, uint(len(a_slice)), 3)
	defer gvek.Free_Ctx(ctx)
	var args gvek.Apply_Args[float32]
	args.A = gvek.New_f32(ctx, uint(len(a_slice)))
	copy(args.A, a_slice)
	args.B = gvek.New_f32(ctx, uint(len(b_slice)))
	copy(args.B, b_slice)
	args.C = gvek.New_f32(ctx, uint(len(a_slice)))
	gvek.Add_f32(args)
	c_expect := []float32{5, 11, 15, 11, 12, 12, 12, 11, 5}
	c_actual := args.C
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
