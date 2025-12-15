package main

import "testing"

func TestSamePtr(t *testing.T) {
	var i1 int = 1
	var i2 int = 1
	a := &i1
	b := &i2
	c := a

	if same, _ := SamePtr(a, c); !same {
		t.Error("SamePtr(a, a) returned false")
	}
	if same, _ := SamePtr(a, b); same {
		t.Error("SamePtr(a, b) returned true")
	}
	if same, _ := SamePtr(a, nil); same {
		t.Error("SamePtr(a, nil) returned true")
	}
}

func Test_identifier(t *testing.T) {
	// Basics does not likely implement Component fully (missing GetStructurals potentially if not embedded)
	// Using Reservoir which definitely implements Component
	c := &Reservoir{Basics: Basics{Identifier: "ID"}}
	if got := identifier(c); got != "ID" {
		t.Errorf("identifier() = %v, want ID", got)
	}
	if got := identifier(nil); got != "nil" {
		t.Errorf("identifier(nil) = %v, want nil", got)
	}
}

func Test_cap(t *testing.T) {
	s := &Structurals{MaxVolume: 100}
	if got := cap(s); got != 100 {
		t.Errorf("cap() = %v, want 100", got)
	}
	if got := cap(nil); got != 0 {
		t.Errorf("cap(nil) = %v, want 0", got)
	}
}

func Test_pres(t *testing.T) {
	// pres uses Quantity / MaxVolume, ignoring explicit Pressure field
	s := &Structurals{Quantity: 50, MaxVolume: 100}
	if got := pres(s); got != 0.5 {
		t.Errorf("pres() = %v, want 0.5", got)
	}
	if got := pres(nil); got != 0 {
		t.Errorf("pres(nil) = %v, want 0", got)
	}
}
