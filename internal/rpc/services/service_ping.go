// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package services

import (
	"context"
	"github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"
)

// PingService Ping服务
// 用来测试连接是否可用
type PingService struct {
	BaseService
}

// Ping 发起Ping
func (this *PingService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	_, _, err := this.ValidateNodeId(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.PingResponse{}, nil
}
