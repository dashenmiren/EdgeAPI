package clients

import (
	"context"

	"github.com/dashenmiren/EdgeAPI/internal/db/models/clients"
	"github.com/dashenmiren/EdgeAPI/internal/rpc/services"
	"github.com/dashenmiren/EdgeCommon/pkg/rpc/pb"
)

// ClientAgentService Agent服务
type ClientAgentService struct {
	services.BaseService
}

// FindAllClientAgents 查找所有Agent
func (this *ClientAgentService) FindAllClientAgents(ctx context.Context, req *pb.FindAllClientAgentsRequest) (*pb.FindAllClientAgentsResponse, error) {
	_, err := this.ValidateAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var tx = this.NullTx()
	agents, err := clients.SharedClientAgentDAO.FindAllAgents(tx)
	if err != nil {
		return nil, err
	}
	var pbAgents = []*pb.ClientAgent{}
	for _, agent := range agents {
		pbAgents = append(pbAgents, &pb.ClientAgent{
			Id:          int64(agent.Id),
			Name:        agent.Name,
			Code:        agent.Code,
			Description: agent.Description,
			CountIPs:    int64(agent.CountIPs),
		})
	}
	return &pb.FindAllClientAgentsResponse{ClientAgents: pbAgents}, nil
}
