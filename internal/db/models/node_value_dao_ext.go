//go:build !plus

package models

import (
	"github.com/dashenmiren/EdgeCommon/pkg/nodeconfigs"
	"github.com/iwind/TeaGo/dbs"
)

// 节点值变更Hook
func (this *NodeValueDAO) nodeValueHook(tx *dbs.Tx, role nodeconfigs.NodeRole, nodeId int64, item nodeconfigs.NodeValueItem, valueJSON []byte) error {
	return nil
}
