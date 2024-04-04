package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	grpcclient "github.com/shamank/ai-marketplace-api-gateway/internal/clients/stats-service/grpc"
	"github.com/shamank/ai-marketplace-api-gateway/internal/domain/models"
	"log/slog"
	"net/url"
	"strconv"
)

type Handler struct {
	grpcClient *grpcclient.StatsServiceClient
	log        *slog.Logger
}

func NewHandler(grpcClient *grpcclient.StatsServiceClient, log *slog.Logger) *Handler {
	return &Handler{
		grpcClient: grpcClient,
		log:        log,
	}
}

func (h *Handler) InitAPIRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(CORS)

	api := r.Group("/api")

	{
		api.POST("/service", h.createService)
		api.POST("/call", h.call)
		api.GET("/calls", h.getCalls)
	}

	return r
}

type createServiceRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required"`
}

type createServiceResponse struct {
	UID string `json:"uid"`
}

func (h *Handler) createService(c *gin.Context) {
	var req createServiceRequest

	if err := bindAndValidate(c, &req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.grpcClient.CreateService(c.Request.Context(), models.AIServiceCreate{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, createServiceResponse{
		UID: result,
	})
}

type callRequest struct {
	UserUID      string `json:"user_uid" validate:"required,uuid"`
	AIServiceUID string `json:"service_uid" validate:"required,uuid"`
}

type callResponse struct {
	Message string `json:"message"`
}

func (h *Handler) call(c *gin.Context) {

	var req callRequest
	if err := bindAndValidate(c, &req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.grpcClient.Call(c.Request.Context(), req.UserUID, req.AIServiceUID); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, callResponse{
		Message: "OK",
	})
}

type getCallsRequest struct {
	UserUID      *string `form:"user_uid"`
	AIServiceUID *string `form:"service_uid"`
	Order        *string `form:"order"`
	PageSize     *uint32 `form:"page_size"`
	PageNumber   *uint32 `form:"page_number"`
}

type getCallResponse struct {
	UserUID      string  `json:"user_uid"`
	AIServiceUID string  `json:"service_uid"`
	Count        uint32  `json:"count"`
	FullAmount   float64 `json:"full_amount"`
}

type getCallsResponse struct {
	Calls []getCallResponse `json:"calls"`
}

func (h *Handler) getCalls(c *gin.Context) {

	var params getCallsRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.grpcClient.GetCalls(c.Request.Context(), models.StatisticFilter{
		UserUID:      params.UserUID,
		AIServiceUID: params.AIServiceUID,
		Order:        params.Order,
		PageSize:     params.PageSize,
		PageNumber:   params.PageNumber,
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	calls := make([]getCallResponse, 0, len(result))
	for _, call := range result {
		calls = append(calls, getCallResponse{
			UserUID:      call.UserUID,
			AIServiceUID: call.AIServiceUID,
			Count:        call.Count,
			FullAmount:   call.FullAmount,
		})
	}

	c.JSON(200, getCallsResponse{
		Calls: calls,
	})
}

func bindAndValidate(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(&req); err != nil {
		return err
	}
	if err := validateRequest(req); err != nil {
		return err
	}
	return nil
}

func validateRequest(req interface{}) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func parseFilterAttributes(params url.Values) (models.StatisticFilter, error) {
	userUID := params.Get("user_uid")
	serviceUID := params.Get("service_uid")
	pageNumber, err := parseUInt32Param(params.Get("page_number"))
	if err != nil {
		return models.StatisticFilter{}, err
	}
	pageSize, err := parseUInt32Param(params.Get("page_size"))
	if err != nil {
		return models.StatisticFilter{}, err
	}
	order := params.Get("order")

	filter := models.StatisticFilter{}

	if userUID != "" {
		filter.UserUID = &userUID
	}
	if serviceUID != "" {
		filter.AIServiceUID = &serviceUID
	}
	if pageNumber != 0 {
		uPageNumber := uint32(pageNumber)
		filter.PageNumber = &uPageNumber
	}
	if pageSize != 0 {
		uPageSize := uint32(pageSize)
		filter.PageSize = &uPageSize
	}
	if order != "" {
		filter.Order = &order
	}

	return filter, nil
}

func parseUInt32Param(param string) (uint32, error) {
	if param == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}
