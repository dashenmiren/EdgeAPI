package utils_test

import (
	"testing"

	"github.com/TeaOSLab/EdgeAPI/internal/utils"
)

func TestSystemMemoryGB(t *testing.T) {
	t.Log(utils.SystemMemoryGB())
	t.Log(utils.SystemMemoryGB())
	t.Log(utils.SystemMemoryGB())
}