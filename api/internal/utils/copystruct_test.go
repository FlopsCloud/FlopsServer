package utils

import (
	"testing"
)

func TestCopyStruct(t *testing.T) {
	// Test case 1: Basic struct with simple fields
	type Source struct {
		Name    string
		Age     int
		Active  bool
		Balance float64
	}

	type Destination struct {
		Name    string
		Age     int
		Active  bool
		Balance float64
	}

	src := &Source{
		Name:    "John Doe",
		Age:     30,
		Active:  true,
		Balance: 100.50,
	}

	dst := &Destination{}

	err := CopyStruct(src, dst)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if dst.Name != src.Name {
		t.Errorf("Expected Name to be %s, got %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("Expected Age to be %d, got %d", src.Age, dst.Age)
	}
	if dst.Active != src.Active {
		t.Errorf("Expected Active to be %v, got %v", src.Active, dst.Active)
	}
	if dst.Balance != src.Balance {
		t.Errorf("Expected Balance to be %f, got %f", src.Balance, dst.Balance)
	}

	// Test case 2: Struct with different field sets
	type PartialDestination struct {
		Name string
		Age  int
		// Omitting Active and Balance fields
	}

	partialDst := &PartialDestination{}
	err = CopyStruct(src, partialDst)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if partialDst.Name != src.Name {
		t.Errorf("Expected Name to be %s, got %s", src.Name, partialDst.Name)
	}
	if partialDst.Age != src.Age {
		t.Errorf("Expected Age to be %d, got %d", src.Age, partialDst.Age)
	}

	// Test case 3: Invalid input (nil pointers)
	err = CopyStruct(nil, nil)
	if err == nil {
		t.Error("Expected error for nil inputs, got nil")
	}
}
