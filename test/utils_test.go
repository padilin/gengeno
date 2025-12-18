package test

import (
	"testing"

	"github.com/padilin/gengeno/game"
)

func TestSamePtr(t *testing.T) {
	var i1 = 1
	var i2 = 1
	a := &i1
	b := &i2
	c := a

	if same, _ := game.SamePtr(a, c); !same {
		t.Error("SamePtr(a, a) returned false")
	}
	if same, _ := game.SamePtr(a, b); same {
		t.Error("SamePtr(a, b) returned true")
	}
	if same, _ := game.SamePtr(a, nil); same {
		t.Error("SamePtr(a, nil) returned true")
	}
}

func Test_identifier(t *testing.T) {
	// t.Skip("Skipping Test_identifier (unexported function)")
	if got := game.Identifier(nil); got != "nil" {
		t.Errorf("identifier(nil) = %q, want %q", got, "nil")
	}

	c := &game.Reservoir{Basics: game.Basics{Identifier: "test-id"}}
	if got := game.Identifier(c); got != "test-id" {
		t.Errorf("identifier(c) = %q, want %q", got, "test-id")
	}
}

func Test_cap(t *testing.T) {
	// t.Skip("Skipping Test_cap (unexported function)")
	if got := game.Cap(nil); got != 0 {
		t.Errorf("cap(nil) = %f, want 0", got)
	}
	s := &game.Structurals{MaxVolume: 100}
	if got := game.Cap(s); got != 100 {
		t.Errorf("cap(s) = %f, want 100", got)
	}
}

func Test_pres(t *testing.T) {
	// t.Skip("Skipping Test_pres (unexported function)")
	if got := game.Pres(nil); got != 0 {
		t.Errorf("pres(nil) = %f, want 0", got)
	}
	s := &game.Structurals{MaxVolume: 100, Quantity: 50}
	if got := game.Pres(s); got != 0.5 {
		t.Errorf("pres(s) = %f, want 0.5", got)
	}
	if got := game.Pres(&game.Structurals{MaxVolume: 0}); got != 0 {
		t.Errorf("pres(s) = %f, want 0", got)
	}
}
