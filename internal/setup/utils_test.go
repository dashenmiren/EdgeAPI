package setup_test

import (
	"testing"

	"github.com/dashenmiren/EdgeAPI/internal/setup"
)

func TestComposeSQLVersion(t *testing.T) {
	t.Log(setup.ComposeSQLVersion())
}
