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
	var args gvek.Apply_Args
	gvek.Add_f32(&args)
}
