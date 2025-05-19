package utils_test

import (
	"testing"

	"github.com/dashenmiren/EdgeAPI/internal/utils"
	"github.com/iwind/TeaGo/types"
)

func TestSha1Random(t *testing.T) {
	for i := 0; i < 10; i++ {
		var s = utils.Sha1RandomString()
		t.Log("["+types.String(len(s))+"]", s)
	}
}
