package services

import (
	"context"
	"github.com/dashenmiren/EdgeAPI/internal/db/models"
	"github.com/dashenmiren/EdgeAPI/internal/errors"
	rpcutils "github.com/dashenmiren/EdgeAPI/internal/rpc/utils"
	"github.com/dashenmiren/EdgeAPI/internal/utils"
	"github.com/dashenmiren/EdgeCommon/pkg/iputils"
	"github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"
	"github.com/dashenmiren/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/dashenmiren/EdgeCommon/pkg/serverconfigs/ipconfigs"
	"net"
	"time"
)

// IPItemService IP条目相关服务
type IPItemService struct {
	BaseService
}

// CreateIPItem 创建IP
func (this *IPItemService) CreateIPItem(ctx context.Context, req *pb.CreateIPItemRequest) (*pb.CreateIPItemResponse, error) {
	// 校验请求
	userType, _, userId, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin, rpcutils.UserTypeUser, rpcutils.UserTypeNode, rpcutils.UserTypeDNS)
	if err != nil {
		return nil, err
	}

	if len(req.Value) > 0 {
		newValue, ipFrom, ipTo, ok := models.SharedIPItemDAO.ParseIPValue(req.Value)
		if !ok {
			return nil, errors.New("invalid 'value' format")
		}

		req.Value = newValue
		req.IpFrom = ipFrom
		req.IpTo = ipTo
	} else if req.Type != models.IPItemTypeAll {
		if !iputils.IsValid(req.IpFrom) {
			return nil, errors.New("invalid 'ipFrom'")
		}
		if len(req.IpTo) > 0 {
			if !iputils.IsValid(req.IpTo) {
				return nil, errors.New("invalid 'ipTo'")
			}

			if !iputils.IsSameVersion(req.IpFrom, req.IpTo) {
				return nil, errors.New("'ipFrom' and 'ipTo' should be in same version")
			}

			if iputils.CompareIP(req.IpFrom, req.IpTo) > 0 {
				req.IpFrom, req.IpTo = req.IpTo, req.IpFrom
			}
		}
	}

	var tx = this.NullTx()

	if userType == rpcutils.UserTypeUser {
		if userId <= 0 {
			return nil, errors.New("invalid userId")
		} else {
			err = models.SharedIPListDAO.CheckUserIPList(tx, userId, req.IpListId)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(req.Type) == 0 {
		req.Type = models.IPItemTypeIPv4
	}

	// 删除以前的
	err = models.SharedIPItemDAO.DeleteOldItem(tx, req.IpListId, req.IpFrom, req.IpTo)
	if err != nil {
		return nil, err
	}

	itemId, err := models.SharedIPItemDAO.CreateIPItem(tx, req.IpListId, req.Value, req.IpFrom, req.IpTo, req.ExpiredAt, req.Reason, req.Type, req.EventLevel, req.NodeId, req.ServerId, req.SourceNodeId, req.SourceServerId, req.SourceHTTPFirewallPolicyId, req.SourceHTTPFirewallRuleGroupId, req.SourceHTTPFirewallRuleSetId, true)
	if err != nil {
		return nil, err
	}

	return &pb.CreateIPItemResponse{IpItemId: itemId}, nil
}

// CreateIPItems 创建一组IP
func (this *IPItemService) CreateIPItems(ctx context.Context, req *pb.CreateIPItemsRequest) (*pb.CreateIPItemsResponse, error) {
	// 校验请求
	userType, _, userId, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin, rpcutils.UserTypeUser, rpcutils.UserTypeNode, rpcutils.UserTypeDNS)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	// 校验
	for _, item := range req.IpItems {
		if len(item.Value) > 0 {
			newValue, ipFrom, ipTo, ok := models.SharedIPItemDAO.ParseIPValue(item.Value)
			if !ok {
				return nil, errors.New("invalid 'value': " + item.Value)
			}
			item.Value = newValue
			item.IpFrom = ipFrom
			item.IpTo = ipTo
		} else if item.Type != models.IPItemTypeAll {
			if !iputils.IsValid(item.IpFrom) {
				return nil, errors.New("invalid 'ipFrom': " + item.IpFrom)
			}
			if len(item.IpTo) > 0 {
				if !iputils.IsValid(item.IpTo) {
					return nil, errors.New("invalid 'ipTo': " + item.IpTo)
				}

				if !iputils.IsSameVersion(item.IpFrom, item.IpTo) {
					return nil, errors.New("'ipFrom' (" + item.IpFrom + ") and 'ipTo' (" + item.IpTo + ") should be in same version")
				}

				if iputils.CompareIP(item.IpFrom, item.IpTo) > 0 {
					item.IpFrom, item.IpTo = item.IpTo, item.IpFrom
				}
			}
		}

		if userType == rpcutils.UserTypeUser {
			if userId <= 0 {
				return nil, errors.New("invalid userId")
			} else {
				err = models.SharedIPListDAO.CheckUserIPList(tx, userId, item.IpListId)
				if err != nil {
					return nil, err
				}
			}
		}

		if len(item.Type) == 0 {
			item.Type = models.IPItemTypeIPv4
		}
	}

	// 创建
	var ipItemIds = []int64{}
	for index, item := range req.IpItems {
		var shouldNotify = index == len(req.IpItems)-1

		// 删除以前的
		if len(item.Value) > 0 {
			err = models.SharedIPItemDAO.DeleteOldItemWithValue(tx, item.IpListId, item.Value)
		} else {
			err = models.SharedIPItemDAO.DeleteOldItem(tx, item.IpListId, item.IpFrom, item.IpTo)
		}
		if err != nil {
			return nil, err
		}

		itemId, err := models.SharedIPItemDAO.CreateIPItem(tx, item.IpListId, item.Value, item.IpFrom, item.IpTo, item.ExpiredAt, item.Reason, item.Type, item.EventLevel, item.NodeId, item.ServerId, item.SourceNodeId, item.SourceServerId, item.SourceHTTPFirewallPolicyId, item.SourceHTTPFirewallRuleGroupId, item.SourceHTTPFirewallRuleSetId, shouldNotify)
		if err != nil {
			return nil, err
		}
		ipItemIds = append(ipItemIds, itemId)
	}

	return &pb.CreateIPItemsResponse{
		IpItemIds: ipItemIds,
	}, nil
}

// UpdateIPItem 修改IP
func (this *IPItemService) UpdateIPItem(ctx context.Context, req *pb.UpdateIPItemRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	// validate ip
	if len(req.Value) > 0 {
		newValue, ipFrom, ipTo, ok := models.SharedIPItemDAO.ParseIPValue(req.Value)
		if !ok {
			return nil, errors.New("invalid 'value' format")
		}
		req.Value = newValue
		req.IpFrom = ipFrom
		req.IpTo = ipTo
	} else if req.Type != models.IPItemTypeAll {
		if !iputils.IsValid(req.IpFrom) {
			return nil, errors.New("invalid 'ipFrom'")
		}
		if len(req.IpTo) > 0 {
			if !iputils.IsValid(req.IpTo) {
				return nil, errors.New("invalid 'ipTo'")
			}

			if !iputils.IsSameVersion(req.IpFrom, req.IpTo) {
				return nil, errors.New("'ipFrom' and 'ipTo' should be in same version")
			}

			if iputils.CompareIP(req.IpFrom, req.IpTo) > 0 {
				req.IpFrom, req.IpTo = req.IpTo, req.IpFrom
			}
		}
	}

	if userId > 0 {
		listId, err := models.SharedIPItemDAO.FindItemListId(tx, req.IpItemId)
		if err != nil {
			return nil, err
		}

		err = models.SharedIPListDAO.CheckUserIPList(tx, userId, listId)
		if err != nil {
			return nil, err
		}
	}

	if len(req.Type) == 0 {
		req.Type = models.IPItemTypeIPv4
	}

	err = models.SharedIPItemDAO.UpdateIPItem(tx, req.IpItemId, req.Value, req.IpFrom, req.IpTo, req.ExpiredAt, req.Reason, req.Type, req.EventLevel)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// DeleteIPItem 删除IP
func (this *IPItemService) DeleteIPItem(ctx context.Context, req *pb.DeleteIPItemRequest) (*pb.RPCSuccess, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if req.IpItemId <= 0 && len(req.Value) == 0 && len(req.IpFrom) == 0 {
		return nil, errors.New("one of 'ipItemId', 'value' or 'ipFrom' params required")
	}

	// 如果是使用IPItemId删除
	if req.IpItemId > 0 {
		err = models.SharedIPItemDAO.DisableIPItem(tx, req.IpItemId, userId)
		if err != nil {
			return nil, err
		}
		return this.Success()
	}

	// 使用value删除
	if len(req.Value) > 0 {
		// 检查IP列表
		if req.IpListId > 0 && userId > 0 && !firewallconfigs.IsGlobalListId(req.IpListId) {
			err = models.SharedIPListDAO.CheckUserIPList(tx, userId, req.IpListId)
			if err != nil {
				return nil, err
			}
		}

		err = models.SharedIPItemDAO.DisableIPItemsWithIPValue(tx, req.Value, userId, req.IpListId)
		if err != nil {
			return nil, err
		}
		return this.Success()
	}

	// 如果是使用ipFrom+ipTo删除
	if len(req.IpFrom) > 0 {
		// 检查IP列表
		if req.IpListId > 0 && userId > 0 && !firewallconfigs.IsGlobalListId(req.IpListId) {
			err = models.SharedIPListDAO.CheckUserIPList(tx, userId, req.IpListId)
			if err != nil {
				return nil, err
			}
		}

		err = models.SharedIPItemDAO.DisableIPItemsWithIP(tx, req.IpFrom, req.IpTo, userId, req.IpListId)
		if err != nil {
			return nil, err
		}
		return this.Success()
	}

	return this.Success()
}

// DeleteIPItems 批量删除IP
func (this *IPItemService) DeleteIPItems(ctx context.Context, req *pb.DeleteIPItemsRequest) (*pb.RPCSuccess, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()
	for _, itemId := range req.IpItemIds {
		err = models.SharedIPItemDAO.DisableIPItem(tx, itemId, userId)
		if err != nil {
			return nil, err
		}
	}
	return this.Success()
}

// CountIPItemsWithListId 计算IP数量
func (this *IPItemService) CountIPItemsWithListId(ctx context.Context, req *pb.CountIPItemsWithListIdRequest) (*pb.RPCCountResponse, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		// 检查用户所属名单
		if !firewallconfigs.IsGlobalListId(req.IpListId) {
			err = models.SharedIPListDAO.CheckUserIPList(tx, userId, req.IpListId)
			if err != nil {
				return nil, err
			}
		}
	}

	count, err := models.SharedIPItemDAO.CountIPItemsWithListId(tx, req.IpListId, userId, req.Keyword, req.IpFrom, req.IpTo, req.EventLevel)
	if err != nil {
		return nil, err
	}
	return this.SuccessCount(count)
}

// ListIPItemsWithListId 列出单页的IP
func (this *IPItemService) ListIPItemsWithListId(ctx context.Context, req *pb.ListIPItemsWithListIdRequest) (*pb.ListIPItemsWithListIdResponse, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	if userId > 0 {
		// 检查用户所属名单
		if !firewallconfigs.IsGlobalListId(req.IpListId) {
			err = models.SharedIPListDAO.CheckUserIPList(tx, userId, req.IpListId)
			if err != nil {
				return nil, err
			}
		}
	}

	items, err := models.SharedIPItemDAO.ListIPItemsWithListId(tx, req.IpListId, userId, req.Keyword, req.IpFrom, req.IpTo, req.EventLevel, req.Offset, req.Size)
	if err != nil {
		return nil, err
	}
	var result = []*pb.IPItem{}
	for _, item := range items {
		if len(item.Type) == 0 {
			item.Type = models.IPItemTypeIPv4
		}

		// server
		var pbSourceServer *pb.Server
		if item.SourceServerId > 0 {
			serverName, err := models.SharedServerDAO.FindEnabledServerName(tx, int64(item.SourceServerId))
			if err != nil {
				return nil, err
			}
			pbSourceServer = &pb.Server{
				Id:   int64(item.SourceServerId),
				Name: serverName,
			}
		}

		// WAF策略
		var pbSourcePolicy *pb.HTTPFirewallPolicy
		if item.SourceHTTPFirewallPolicyId > 0 {
			policy, err := models.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyBasic(tx, int64(item.SourceHTTPFirewallPolicyId))
			if err != nil {
				return nil, err
			}
			if policy != nil {
				pbSourcePolicy = &pb.HTTPFirewallPolicy{
					Id:       int64(item.SourceHTTPFirewallPolicyId),
					Name:     policy.Name,
					ServerId: int64(policy.ServerId),
				}
			}
		}

		// WAF分组
		var pbSourceGroup *pb.HTTPFirewallRuleGroup
		if item.SourceHTTPFirewallRuleGroupId > 0 {
			groupName, err := models.SharedHTTPFirewallRuleGroupDAO.FindHTTPFirewallRuleGroupName(tx, int64(item.SourceHTTPFirewallRuleGroupId))
			if err != nil {
				return nil, err
			}
			pbSourceGroup = &pb.HTTPFirewallRuleGroup{
				Id:   int64(item.SourceHTTPFirewallRuleGroupId),
				Name: groupName,
			}
		}

		// WAF规则集
		var pbSourceSet *pb.HTTPFirewallRuleSet
		if item.SourceHTTPFirewallRuleSetId > 0 {
			setName, err := models.SharedHTTPFirewallRuleSetDAO.FindHTTPFirewallRuleSetName(tx, int64(item.SourceHTTPFirewallRuleSetId))
			if err != nil {
				return nil, err
			}
			pbSourceSet = &pb.HTTPFirewallRuleSet{
				Id:   int64(item.SourceHTTPFirewallRuleSetId),
				Name: setName,
			}
		}

		result = append(result, &pb.IPItem{
			Id:                            int64(item.Id),
			Value:                         item.ComposeValue(),
			IpFrom:                        item.IpFrom,
			IpTo:                          item.IpTo,
			Version:                       int64(item.Version),
			CreatedAt:                     int64(item.CreatedAt),
			ExpiredAt:                     int64(item.ExpiredAt),
			Reason:                        item.Reason,
			Type:                          item.Type,
			EventLevel:                    item.EventLevel,
			NodeId:                        int64(item.NodeId),
			ServerId:                      int64(item.ServerId),
			SourceNodeId:                  int64(item.SourceNodeId),
			SourceServerId:                int64(item.SourceServerId),
			SourceHTTPFirewallPolicyId:    int64(item.SourceHTTPFirewallPolicyId),
			SourceHTTPFirewallRuleGroupId: int64(item.SourceHTTPFirewallRuleGroupId),
			SourceHTTPFirewallRuleSetId:   int64(item.SourceHTTPFirewallRuleSetId),
			SourceServer:                  pbSourceServer,
			SourceHTTPFirewallPolicy:      pbSourcePolicy,
			SourceHTTPFirewallRuleGroup:   pbSourceGroup,
			SourceHTTPFirewallRuleSet:     pbSourceSet,
			IsRead:                        item.IsRead,
		})
	}

	return &pb.ListIPItemsWithListIdResponse{IpItems: result}, nil
}

// FindEnabledIPItem 查找单个IP
func (this *IPItemService) FindEnabledIPItem(ctx context.Context, req *pb.FindEnabledIPItemRequest) (*pb.FindEnabledIPItemResponse, error) {
	// 校验请求
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	item, err := models.SharedIPItemDAO.FindEnabledIPItem(tx, req.IpItemId)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return &pb.FindEnabledIPItemResponse{IpItem: nil}, nil
	}

	if userId > 0 {
		err = models.SharedIPListDAO.CheckUserIPList(tx, userId, int64(item.ListId))
		if err != nil {
			return nil, err
		}
	}

	if len(item.Type) == 0 {
		item.Type = models.IPItemTypeIPv4
	}

	return &pb.FindEnabledIPItemResponse{IpItem: &pb.IPItem{
		Id:         int64(item.Id),
		Value:      item.ComposeValue(),
		IpFrom:     item.IpFrom,
		IpTo:       item.IpTo,
		Version:    int64(item.Version),
		CreatedAt:  int64(item.CreatedAt),
		ExpiredAt:  int64(item.ExpiredAt),
		Reason:     item.Reason,
		Type:       item.Type,
		EventLevel: item.EventLevel,
		NodeId:     int64(item.NodeId),
		ServerId:   int64(item.ServerId),
	}}, nil
}

// ListIPItemsAfterVersion 根据版本列出一组IP
func (this *IPItemService) ListIPItemsAfterVersion(ctx context.Context, req *pb.ListIPItemsAfterVersionRequest) (*pb.ListIPItemsAfterVersionResponse, error) {
	// 校验请求
	_, _, _, err := rpcutils.ValidateRequest(ctx, rpcutils.UserTypeAdmin, rpcutils.UserTypeNode)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()

	var result = []*pb.IPItem{}
	items, err := models.SharedIPItemDAO.ListIPItemsAfterVersion(tx, req.Version, req.Size)
	if err != nil {
		return nil, err
	}

	var latestVersion = req.Version

	for _, item := range items {
		latestVersion = int64(item.Version)

		// 是否已过期
		if item.ExpiredAt > 0 && int64(item.ExpiredAt) <= time.Now().Unix() {
			item.State = models.IPItemStateDisabled
		}

		if len(item.Type) == 0 {
			item.Type = models.IPItemTypeIPv4
		}

		// List类型
		list, err := models.SharedIPListDAO.FindIPListCacheable(tx, int64(item.ListId))
		if err != nil {
			return nil, err
		}
		if list == nil {
			continue
		}

		// 跳过灰名单
		if list.Type == ipconfigs.IPListTypeGrey {
			continue
		}

		// 如果已经删除
		if list.State != models.IPListStateEnabled {
			item.State = models.IPItemStateDisabled
		}

		result = append(result, &pb.IPItem{
			Id:         int64(item.Id),
			Value:      item.ComposeValue(),
			IpFrom:     item.IpFrom,
			IpTo:       item.IpTo,
			Version:    int64(item.Version),
			CreatedAt:  int64(item.CreatedAt),
			ExpiredAt:  int64(item.ExpiredAt),
			Reason:     "", // 这里我们不需要这个数据
			ListId:     int64(item.ListId),
			IsDeleted:  item.State == 0,
			Type:       item.Type,
			EventLevel: item.EventLevel,
			ListType:   list.Type,
			IsGlobal:   list.IsPublic && list.IsGlobal,
			NodeId:     int64(item.NodeId),
			ServerId:   int64(item.ServerId),
		})
	}

	return &pb.ListIPItemsAfterVersionResponse{
		IpItems: result,
		Version: latestVersion,
	}, nil
}

// CheckIPItemStatus 检查IP状态
func (this *IPItemService) CheckIPItemStatus(ctx context.Context, req *pb.CheckIPItemStatusRequest) (*pb.CheckIPItemStatusResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// 校验IP
	var ip = net.ParseIP(req.Ip)
	if len(ip) == 0 {
		return &pb.CheckIPItemStatusResponse{
			IsOk:  false,
			Error: "请输入正确的IP",
		}, nil
	}

	var tx = this.NullTx()

	// 名单类型
	list, err := models.SharedIPListDAO.FindEnabledIPList(tx, req.IpListId, nil)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return &pb.CheckIPItemStatusResponse{
			IsOk:  false,
			Error: "IP名单不存在",
		}, nil
	}
	var isAllowed = list.Type == ipconfigs.IPListTypeWhite || list.Type == ipconfigs.IPListTypeGrey

	// 检查IP名单
	item, err := models.SharedIPItemDAO.FindEnabledItemContainsIP(tx, req.IpListId, req.Ip)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return &pb.CheckIPItemStatusResponse{
			IsOk:      true,
			Error:     "",
			IsFound:   true,
			IsAllowed: isAllowed,
			IpItem: &pb.IPItem{
				Id:         int64(item.Id),
				Value:      item.ComposeValue(),
				IpFrom:     item.IpFrom,
				IpTo:       item.IpTo,
				CreatedAt:  int64(item.CreatedAt),
				ExpiredAt:  int64(item.ExpiredAt),
				Reason:     item.Reason,
				Type:       item.Type,
				EventLevel: item.EventLevel,
				ListType:   list.Type,
			},
		}, nil
	}

	return &pb.CheckIPItemStatusResponse{
		IsOk:      true,
		Error:     "",
		IsFound:   false,
		IsAllowed: false,
		IpItem:    nil,
	}, nil
}

// ExistsEnabledIPItem 检查IP是否存在
func (this *IPItemService) ExistsEnabledIPItem(ctx context.Context, req *pb.ExistsEnabledIPItemRequest) (*pb.ExistsEnabledIPItemResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()
	b, err := models.SharedIPItemDAO.ExistsEnabledItem(tx, req.IpItemId)
	if err != nil {
		return nil, err
	}
	return &pb.ExistsEnabledIPItemResponse{Exists: b}, nil
}

// CountAllEnabledIPItems 计算所有IP数量
func (this *IPItemService) CountAllEnabledIPItems(ctx context.Context, req *pb.CountAllEnabledIPItemsRequest) (*pb.RPCCountResponse, error) {
	adminId, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if adminId > 0 {
		userId = req.UserId
	}

	var tx = this.NullTx()
	count, err := models.SharedIPItemDAO.CountAllEnabledIPItems(tx, userId, req.Keyword, req.Ip, 0, req.Unread, req.EventLevel, req.ListType, req.GlobalOnly)
	if err != nil {
		return nil, err
	}
	return this.SuccessCount(count)
}

// ListAllEnabledIPItems 搜索IP
func (this *IPItemService) ListAllEnabledIPItems(ctx context.Context, req *pb.ListAllEnabledIPItemsRequest) (*pb.ListAllEnabledIPItemsResponse, error) {
	adminId, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if adminId > 0 {
		userId = req.UserId
	}

	var results = []*pb.ListAllEnabledIPItemsResponse_Result{}
	var tx = this.NullTx()
	items, err := models.SharedIPItemDAO.ListAllEnabledIPItems(tx, userId, req.Keyword, req.Ip, 0, req.Unread, req.EventLevel, req.ListType, req.GlobalOnly, req.Offset, req.Size)
	if err != nil {
		return nil, err
	}

	var cacheMap = utils.NewCacheMap()
	for _, item := range items {
		// server
		var pbSourceServer *pb.Server
		if item.SourceServerId > 0 {
			serverName, err := models.SharedServerDAO.FindEnabledServerName(tx, int64(item.SourceServerId))
			if err != nil {
				return nil, err
			}
			pbSourceServer = &pb.Server{
				Id:   int64(item.SourceServerId),
				Name: serverName,
			}
		}

		// WAF策略
		var pbSourcePolicy *pb.HTTPFirewallPolicy
		if item.SourceHTTPFirewallPolicyId > 0 {
			policy, err := models.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyBasic(tx, int64(item.SourceHTTPFirewallPolicyId))
			if err != nil {
				return nil, err
			}
			if policy != nil {
				pbSourcePolicy = &pb.HTTPFirewallPolicy{
					Id:       int64(item.SourceHTTPFirewallPolicyId),
					Name:     policy.Name,
					ServerId: int64(policy.ServerId),
				}
			}
		}

		// WAF分组
		var pbSourceGroup *pb.HTTPFirewallRuleGroup
		if item.SourceHTTPFirewallRuleGroupId > 0 {
			groupName, err := models.SharedHTTPFirewallRuleGroupDAO.FindHTTPFirewallRuleGroupName(tx, int64(item.SourceHTTPFirewallRuleGroupId))
			if err != nil {
				return nil, err
			}
			pbSourceGroup = &pb.HTTPFirewallRuleGroup{
				Id:   int64(item.SourceHTTPFirewallRuleGroupId),
				Name: groupName,
			}
		}

		// WAF规则集
		var pbSourceSet *pb.HTTPFirewallRuleSet
		if item.SourceHTTPFirewallRuleSetId > 0 {
			setName, err := models.SharedHTTPFirewallRuleSetDAO.FindHTTPFirewallRuleSetName(tx, int64(item.SourceHTTPFirewallRuleSetId))
			if err != nil {
				return nil, err
			}
			pbSourceSet = &pb.HTTPFirewallRuleSet{
				Id:   int64(item.SourceHTTPFirewallRuleSetId),
				Name: setName,
			}
		}

		// 节点
		var pbSourceNode *pb.Node
		if item.SourceNodeId > 0 {
			node, err := models.SharedNodeDAO.FindEnabledBasicNode(tx, int64(item.SourceNodeId))
			if err != nil {
				return nil, err
			}
			if node != nil {
				pbSourceNode = &pb.Node{
					Id:          int64(node.Id),
					Name:        node.Name,
					NodeCluster: &pb.NodeCluster{Id: int64(node.ClusterId)},
				}
			}
		}

		var pbItem = &pb.IPItem{
			Id:                            int64(item.Id),
			Value:                         item.ComposeValue(),
			IpFrom:                        item.IpFrom,
			IpTo:                          item.IpTo,
			Version:                       int64(item.Version),
			CreatedAt:                     int64(item.CreatedAt),
			ExpiredAt:                     int64(item.ExpiredAt),
			Reason:                        item.Reason,
			Type:                          item.Type,
			EventLevel:                    item.EventLevel,
			NodeId:                        int64(item.NodeId),
			ServerId:                      int64(item.ServerId),
			SourceNodeId:                  int64(item.SourceNodeId),
			SourceServerId:                int64(item.SourceServerId),
			SourceHTTPFirewallPolicyId:    int64(item.SourceHTTPFirewallPolicyId),
			SourceHTTPFirewallRuleGroupId: int64(item.SourceHTTPFirewallRuleGroupId),
			SourceHTTPFirewallRuleSetId:   int64(item.SourceHTTPFirewallRuleSetId),
			SourceServer:                  pbSourceServer,
			SourceHTTPFirewallPolicy:      pbSourcePolicy,
			SourceHTTPFirewallRuleGroup:   pbSourceGroup,
			SourceHTTPFirewallRuleSet:     pbSourceSet,
			SourceNode:                    pbSourceNode,
			IsRead:                        item.IsRead,
		}

		// 所属名单
		list, err := models.SharedIPListDAO.FindEnabledIPList(tx, int64(item.ListId), cacheMap)
		if err != nil {
			return nil, err
		}
		if list == nil {
			err = models.SharedIPItemDAO.DisableIPItem(tx, int64(item.Id), 0)
			if err != nil {
				return nil, err
			}
			continue
		}
		var pbList = &pb.IPList{
			Id:       int64(list.Id),
			Name:     list.Name,
			Type:     list.Type,
			IsPublic: list.IsPublic,
			IsGlobal: list.IsGlobal,
		}

		// 所属服务（注意与SourceServer不同）
		var pbFirewallServer *pb.Server

		// 所属策略（注意与SourceHTTPFirewallPolicy不同）
		var pbFirewallPolicy *pb.HTTPFirewallPolicy
		if !list.IsPublic {
			policy, err := models.SharedHTTPFirewallPolicyDAO.FindEnabledFirewallPolicyWithIPListId(tx, int64(list.Id))
			if err != nil {
				return nil, err
			}
			if policy == nil {
				err = models.SharedIPItemDAO.DisableIPItem(tx, int64(item.Id), 0)
				if err != nil {
					return nil, err
				}
				continue
			}

			pbFirewallPolicy = &pb.HTTPFirewallPolicy{
				Id:   int64(policy.Id),
				Name: policy.Name,
			}

			if policy.ServerId > 0 {
				serverName, err := models.SharedServerDAO.FindEnabledServerName(tx, int64(policy.ServerId))
				if err != nil {
					return nil, err
				}
				if len(serverName) == 0 {
					serverName = "[已删除]"
				}
				pbFirewallServer = &pb.Server{
					Id:   int64(policy.ServerId),
					Name: serverName,
				}
			}
		}

		results = append(results, &pb.ListAllEnabledIPItemsResponse_Result{
			IpList:             pbList,
			IpItem:             pbItem,
			Server:             pbFirewallServer,
			HttpFirewallPolicy: pbFirewallPolicy,
		})
	}

	return &pb.ListAllEnabledIPItemsResponse{Results: results}, nil
}

// ListAllIPItemIds 列出所有名单中的IP ID
func (this *IPItemService) ListAllIPItemIds(ctx context.Context, req *pb.ListAllIPItemIdsRequest) (*pb.ListAllIPItemIdsResponse, error) {
	adminId, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	if adminId > 0 {
		userId = req.UserId
	}

	var tx = this.NullTx()
	itemIds, err := models.SharedIPItemDAO.ListAllIPItemIds(tx, userId, req.Keyword, req.Ip, 0, req.Unread, req.EventLevel, req.ListType, req.Offset, req.Size)
	if err != nil {
		return nil, err
	}
	return &pb.ListAllIPItemIdsResponse{IpItemIds: itemIds}, nil
}

// UpdateIPItemsRead 设置所有为已读
func (this *IPItemService) UpdateIPItemsRead(ctx context.Context, req *pb.UpdateIPItemsReadRequest) (*pb.RPCSuccess, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()
	err = models.SharedIPItemDAO.UpdateItemsRead(tx, userId)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// FindServerIdWithIPItemId 查找IP对应的名单所属网站ID
func (this *IPItemService) FindServerIdWithIPItemId(ctx context.Context, req *pb.FindServerIdWithIPItemIdRequest) (*pb.FindServerIdWithIPItemIdResponse, error) {
	_, userId, err := this.ValidateAdminAndUser(ctx, true)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()
	listId, err := models.SharedIPItemDAO.FindItemListId(tx, req.IpItemId)
	if err != nil {
		return nil, err
	}

	if listId > 0 {
		var serverId int64
		serverId, err = models.SharedIPListDAO.FindServerIdWithListId(tx, listId)
		if err != nil {
			return nil, err
		}

		if serverId > 0 {
			// check user
			if userId > 0 {
				err = models.SharedServerDAO.CheckUserServer(tx, userId, serverId)
				if err != nil {
					return nil, err
				}
			}
			return &pb.FindServerIdWithIPItemIdResponse{ServerId: serverId}, nil
		}
	}

	return &pb.FindServerIdWithIPItemIdResponse{ServerId: 0}, nil
}
