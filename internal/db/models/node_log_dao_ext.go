//go:build !plus

package models

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/iwind/TeaGo/dbs"
)

func (this *NodeLogDAO) deleteNodeLogsWithCluster(tx *dbs.Tx, role nodeconfigs.NodeRole, clusterId int64) error {
	return nil
}
