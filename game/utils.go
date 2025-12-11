package main

import (
	"fmt"
	"reflect"
)

// SamePtr reports whether a and b are the exact same pointer value.
func SamePtr(a, b interface{}) (bool, string) {
	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)
	if ra.IsValid() && rb.IsValid() && ra.Kind() == reflect.Ptr && rb.Kind() == reflect.Ptr {
		eq := ra.Pointer() == rb.Pointer()
		desc := fmt.Sprintf("%T(0x%x) == %T(0x%x) => %v", a, ra.Pointer(), b, rb.Pointer(), eq)
		return eq, desc
	}
	// Fallback: try direct interface compare (works for comparable underlying types)
	eq := false
	defer func() {
		// silence panic if not comparable
		_ = recover()
	}()
	eq = (a == b)
	desc := fmt.Sprintf("fallback compare %T == %T => %v", a, b, eq)
	return eq, desc
}

func identifier(c Component) string {
	if c == nil {
		return "nil"
	}
	return c.GetIdentifier()
}

func cap(s *Structurals) float64 {
	if s == nil {
		return 0
	}
	return s.Volume
}

func pres(s *Structurals) float64 {
	if s == nil {
		return 0
	}
	return s.Pressure
}
