package maputils_test

import (
	"testing"

	maputils "github.com/TeaOSLab/EdgeAPI/internal/utils/maps"
	"github.com/iwind/TeaGo/assert"
)

func TestNewFixedMap(t *testing.T) {
	var a = assert.NewAssertion(t)

	{
		var m = maputils.NewFixedMap(5)
		m.Set("a", 1)
		m.Set("b", 2)
		a.IsTrue(m.Has("a"))
		a.IsTrue(m.Has("b"))
		a.IsFalse(m.Has("c"))
	}

	{
		var m = maputils.NewFixedMap(5)
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		m.Set("d", 4)
		m.Set("e", 5)
		a.IsTrue(m.Size() == 5)
		m.Set("f", 6)
		a.IsTrue(m.Size() == 5)
		a.IsFalse(m.Has("a"))
	}

	{
		var m = maputils.NewFixedMap(5)
		m.Set("a", 1)
		t.Log(m.Get("a"))
		t.Log(m.Get("b"))
	}
}