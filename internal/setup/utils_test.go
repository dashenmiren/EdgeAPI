

package setup_test

import (
	"github.com/TeaOSLab/EdgeAPI/internal/setup"
	"testing"
)

func TestComposeSQLVersion(t *testing.T) {
	t.Log(setup.ComposeSQLVersion())
}
