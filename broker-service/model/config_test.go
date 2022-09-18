package model

import (
	"os"
	"testing"
)

func TestGetConmfig(t *testing.T) {
	buf, _ := os.ReadFile("../configExample.json")
	if len(buf) == 0 {
		t.Fatal("cant load json")
	}
}
