// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .
//go:build !plus

package services

import "github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"

func (this *HTTPAccessLogService) canWriteAccessLogsToDB() bool {
	return true
}

func (this *HTTPAccessLogService) writeAccessLogsToPolicy(pbAccessLogs []*pb.HTTPAccessLog) error {
	return nil
}
