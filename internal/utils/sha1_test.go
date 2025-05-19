// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package utils_test

import (
	"github.com/dashenmiren/EdgeAPI/internal/utils"
	"github.com/iwind/TeaGo/types"
	"testing"
)

func TestSha1Random(t *testing.T) {
	for i := 0; i < 10; i++ {
		var s = utils.Sha1RandomString()
		t.Log("["+types.String(len(s))+"]", s)
	}
}
