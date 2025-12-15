package main

import (
	"math"
	"testing"
)

func TestBasics_GetIdentifier(t *testing.T) {
	tests := []struct {
		name string
		b    *Basics
		want string
	}{
		{
			name: "Regular Identifier",
			b:    &Basics{Identifier: "A"},
			want: "A",
		},
		{
			name: "Empty Identifier",
			b:    &Basics{Identifier: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetIdentifier(); got != tt.want {
				t.Errorf("Basics.GetIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasics_GetColor(t *testing.T) {
	tests := []struct {
		name  string
		b     *Basics
		want  byte
		want1 byte
		want2 byte
	}{
		{
			name:  "Red Color",
			b:     &Basics{Color: [3]byte{255, 0, 0}},
			want:  255,
			want1: 0,
			want2: 0,
		},
		{
			name:  "Black Color",
			b:     &Basics{Color: [3]byte{0, 0, 0}},
			want:  0,
			want1: 0,
			want2: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := tt.b.GetColor()
			if got != tt.want {
				t.Errorf("Basics.GetColor() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Basics.GetColor() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("Basics.GetColor() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestReservoir_GetStructurals(t *testing.T) {
	s := Structurals{MaxVolume: 100}
	r := &Reservoir{
		Structurals: s,
	}
	t.Run("Get Structural Data", func(t *testing.T) {
		got := r.GetStructurals()
		if got == nil {
			t.Fatal("GetStructurals returned nil")
		}
		if got.MaxVolume != 100 {
			t.Errorf("GetStructurals().MaxVolume = %v, want 100", got.MaxVolume)
		}
	})
}

func TestNewPipe(t *testing.T) {
	type args struct {
		from   Component
		to     Component
		len    float64
		radius float64
	}
	tests := []struct {
		name string
		args args
		want *Pipe
	}{
		{
			name: "Create Pipe",
			args: args{
				from:   nil,
				to:     nil,
				len:    10.0,
				radius: 2.0,
			},
			want: nil, // We'll assert properties manually
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPipe(tt.args.from, tt.args.to, tt.args.len, tt.args.radius)
			if got == nil {
				t.Fatal("NewPipe() returned nil")
			}
			expectedArea := math.Pi * 4.0 // pi * r^2 = pi * 2^2
			if math.Abs(got.Structurals.Area-expectedArea) > 0.0001 {
				t.Errorf("NewPipe() Area = %v, want %v", got.Structurals.Area, expectedArea)
			}
			expectedVol := expectedArea * 10.0
			if math.Abs(got.Structurals.MaxVolume-expectedVol) > 0.0001 {
				t.Errorf("NewPipe() MaxVolume = %v, want %v", got.Structurals.MaxVolume, expectedVol)
			}
			if got.Length != 10.0 {
				t.Errorf("NewPipe() Length = %v, want 10.0", got.Length)
			}
		})
	}
}

func TestPipe_GetStructurals(t *testing.T) {
	s := Structurals{MaxVolume: 50}
	p := &Pipe{
		Structurals: s,
	}
	t.Run("Get Structural Data", func(t *testing.T) {
		got := p.GetStructurals()
		if got == nil {
			t.Fatal("GetStructurals returned nil")
		}
		if got.MaxVolume != 50 {
			t.Errorf("GetStructurals().MaxVolume = %v, want 50", got.MaxVolume)
		}
	})
}
