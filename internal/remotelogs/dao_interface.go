// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package remotelogs

import (
	"github.com/dashenmiren/EdgeCommon/pkg/nodeconfigs"
	"github.com/iwind/TeaGo/dbs"
)

type DAOInterface interface {
	CreateLog(tx *dbs.Tx, nodeRole nodeconfigs.NodeRole, nodeId int64, serverId int64, originId int64, level string, tag string, description string, createdAt int64, logType string, paramsJSON []byte) error
}
