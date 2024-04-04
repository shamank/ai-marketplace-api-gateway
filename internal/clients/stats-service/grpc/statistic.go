package grpcclient

import (
	"context"
	"github.com/shamank/ai-marketplace-api-gateway/internal/domain/models"
	statsv1 "github.com/shamank/ai-marketplace-protos/gen/go/stats-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type StatsServiceClient struct {
	api statsv1.StatisticServiceClient
	log *slog.Logger
}

func NewStatsServiceClient(ctx context.Context, log *slog.Logger, addr string) (*StatsServiceClient, error) {

	client, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &StatsServiceClient{
		api: statsv1.NewStatisticServiceClient(client),
		log: log,
	}, nil
}

func (c *StatsServiceClient) CreateService(ctx context.Context, service models.AIServiceCreate) (string, error) {

	req := generateCreateServiceRequest(service)

	resp, err := c.api.Create(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetServiceUid(), nil
}

func (c *StatsServiceClient) GetCalls(ctx context.Context, filter models.StatisticFilter) ([]models.StatisticRead, error) {

	req := generateGetCallsRequest(filter)

	resp, err := c.api.GetCalls(ctx, req)
	if err != nil {
		return nil, err
	}

	result := make([]models.StatisticRead, 0)

	for _, stat := range resp.Calls {
		result = append(result, models.StatisticRead{
			UserUID:      stat.UserUid,
			AIServiceUID: stat.ServiceUid,
			Count:        stat.Count,
			FullAmount:   stat.Amount,
		})
	}

	return result, nil
}

func (c *StatsServiceClient) Call(ctx context.Context, serviceUID string, userUID string) error {

	req := &statsv1.CallRequest{
		UserUid:    userUID,
		ServiceUid: serviceUID,
	}

	_, err := c.api.Call(ctx, req)

	return err
}

func generateCreateServiceRequest(service models.AIServiceCreate) *statsv1.CreateAIServiceRequest {
	req := statsv1.CreateAIServiceRequest{}

	req.Title = service.Title
	if service.Description != nil {
		req.Description = *service.Description
	}
	req.Price = service.Price

	return &req
}

func generateGetCallsRequest(filter models.StatisticFilter) *statsv1.GetCallsRequest {
	req := statsv1.GetCallsRequest{}

	if filter.UserUID != nil {
		req.UserUid = *filter.UserUID
	}
	if filter.AIServiceUID != nil {
		req.ServiceUid = *filter.AIServiceUID
	}

	if filter.Order != nil {
		if *filter.Order == "desc" {
			req.Order = statsv1.OrderEnum_DESC
		} else {
			req.Order = statsv1.OrderEnum_ASC
		}
	}

	if filter.PageSize != nil {
		req.PageSize = *filter.PageSize
	}
	if filter.PageNumber != nil {
		req.PageNumber = *filter.PageNumber
	}

	return &req
}
