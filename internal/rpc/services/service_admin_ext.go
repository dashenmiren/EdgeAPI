// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .
//go:build !plus

package services

import (
	"context"
	"github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/dbs"
)

// ComposeAdminDashboard方法扩展
func (this *AdminService) composeAdminDashboardExt(tx *dbs.Tx, ctx context.Context, result *pb.ComposeAdminDashboardResponse) error {
	return nil
}
