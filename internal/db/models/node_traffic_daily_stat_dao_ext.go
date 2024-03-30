//go:build !plus

package models

import (
	"github.com/dashenmiren/EdgeCommon/pkg/nodeconfigs"
	"github.com/iwind/TeaGo/dbs"
)

// 增加日统计Hook
func (this *NodeTrafficDailyStatDAO) increaseDailyStatHook(tx *dbs.Tx, role nodeconfigs.NodeRole, nodeId int64) error {
	return nil
}
