// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package services

import (
	"context"
	"github.com/dashenmiren/EdgeAPI/internal/db/models"
	"github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"
)

// LatestItemService 最近使用的条目服务
type LatestItemService struct {
	BaseService
}

// IncreaseLatestItem 记录最近使用的条目
func (this *LatestItemService) IncreaseLatestItem(ctx context.Context, req *pb.IncreaseLatestItemRequest) (*pb.RPCSuccess, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}
	var tx = this.NullTx()
	err = models.SharedLatestItemDAO.IncreaseItemCount(tx, req.ItemType, req.ItemId)
	if err != nil {
		return nil, err
	}
	return this.Success()
}
