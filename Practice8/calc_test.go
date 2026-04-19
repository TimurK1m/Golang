package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTableDriven(t *testing.T) {
 	tests := []struct {
	name string
	a, b int
	want int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		got := Add(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got,
			tt.want) 
			}
		})
	}
}

func TestSubtractTableDriven(t *testing.T) {
	
	tests := []struct {
		name     string
		a,b      int
		expected int
	}{
		{"Both positive", 10, 5, 5},
		{"Positive minus zero", 10, 0, 10},
		{"Negative minus positive", -5, 10, -15},
		{"Both negative", -5, -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDivide(t *testing.T) {
	
	res, err := Divide(10, 2)
	assert.NoError(t, err)
	assert.Equal(t, 5, res)

	
	res, err = Divide(10, 0)
	assert.Error(t, err)
	assert.Equal(t, "division by zero", err.Error())
}