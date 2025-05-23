// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package services

import (
	"context"
	"errors"
	"github.com/dashenmiren/EdgeAPI/internal/db/models"
	"github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"
)

// LoginSessionService 登录SESSION服务
type LoginSessionService struct {
	BaseService
}

// WriteLoginSessionValue 写入SESSION数据
func (this *LoginSessionService) WriteLoginSessionValue(ctx context.Context, req *pb.WriteLoginSessionValueRequest) (*pb.RPCSuccess, error) {
	_, _, err := this.ValidateAdminAndUser(ctx, false)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()
	err = models.SharedLoginSessionDAO.WriteSessionValue(tx, req.Sid, req.Key, req.Value)
	if err != nil {
		return nil, err
	}
	return this.Success()
}

// DeleteLoginSession 删除SESSION
func (this *LoginSessionService) DeleteLoginSession(ctx context.Context, req *pb.DeleteLoginSessionRequest) (*pb.RPCSuccess, error) {
	_, _, err := this.ValidateAdminAndUser(ctx, false)
	if err != nil {
		return nil, err
	}

	if len(req.Sid) == 0 {
		return nil, errors.New("'sid' should not be empty")
	}

	var tx = this.NullTx()
	err = models.SharedLoginSessionDAO.DeleteSession(tx, req.Sid)
	if err != nil {
		return nil, err
	}

	return this.Success()
}

// FindLoginSession 查找SESSION
func (this *LoginSessionService) FindLoginSession(ctx context.Context, req *pb.FindLoginSessionRequest) (*pb.FindLoginSessionResponse, error) {
	_, _, err := this.ValidateAdminAndUser(ctx, false)
	if err != nil {
		return nil, err
	}

	if len(req.Sid) == 0 {
		return nil, errors.New("'token' should not be empty")
	}

	var tx = this.NullTx()
	session, err := models.SharedLoginSessionDAO.FindSession(tx, req.Sid)
	if err != nil {
		return nil, err
	}
	if session == nil || !session.IsAvailable() {
		return &pb.FindLoginSessionResponse{
			LoginSession: nil,
		}, nil
	}

	return &pb.FindLoginSessionResponse{
		LoginSession: &pb.LoginSession{
			Id:         int64(session.Id),
			Sid:        session.Sid,
			AdminId:    int64(session.AdminId),
			UserId:     int64(session.UserId),
			Ip:         session.Ip,
			CreatedAt:  int64(session.CreatedAt),
			ExpiresAt:  int64(session.ExpiresAt),
			ValuesJSON: session.Values,
		},
	}, nil
}

// ClearOldLoginSessions 清理老的SESSION
func (this *LoginSessionService) ClearOldLoginSessions(ctx context.Context, req *pb.ClearOldLoginSessionsRequest) (*pb.RPCSuccess, error) {
	_, _, err := this.ValidateAdminAndUser(ctx, false)
	if err != nil {
		return nil, err
	}

	if len(req.Sid) == 0 {
		return nil, errors.New("'token' should not be empty")
	}

	var tx = this.NullTx()
	session, err := models.SharedLoginSessionDAO.FindSession(tx, req.Sid)
	if err != nil {
		return nil, err
	}
	if session == nil || !session.IsAvailable() {
		return nil, errors.New("invalid sid")
	}

	err = models.SharedLoginSessionDAO.ClearOldSessions(tx, int64(session.AdminId), int64(session.UserId), req.Sid, req.Ip)
	if err != nil {
		return nil, err
	}

	return this.Success()
}
