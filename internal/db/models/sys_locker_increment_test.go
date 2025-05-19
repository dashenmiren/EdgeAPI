package models_test

import (
	"testing"

	"github.com/dashenmiren/EdgeAPI/internal/db/models"
)

func TestNewSysLockerIncrement(t *testing.T) {
	var increment = models.NewSysLockerIncrement(10)
	increment.Push("key", 1, 10)
	t.Log(increment.MaxValue("key"))
	for i := 0; i < 11; i++ {
		result, value := increment.Pop("key")
		t.Log(i, "=>", result, value)
	}

	for i := 0; i < 11; i++ {
		result, value := increment.Pop("key1")
		t.Log(i, "=>", result, value)
	}
}
