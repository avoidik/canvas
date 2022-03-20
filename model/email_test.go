package model_test

import (
	"canvas/model"
	"testing"

	"github.com/matryer/is"
)

func TestEmail_IsValid(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"me@example.com", true},
		{"@example.com", false},
		{"me@", false},
		{"@", false},
		{"", false},
		{"`@example.com", true},
		{"!@example.com", true},
		{"{}@example.com", true},
	}
	t.Run("email pairs", func(t *testing.T) {
		is := is.New(t)
		for _, test := range tests {
			e := model.Email(test.address)
			is.Equal(test.valid, e.IsValid())
		}
	})
}
