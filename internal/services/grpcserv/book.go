package grpcserv

import (
	"2/internal/core"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DBRepository interface {
	CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

type GrpcBookService struct {
	bookRepository DBRepository
}

func NewGrpcBookService(bookRepository DBRepository) *GrpcBookService {
	return &GrpcBookService{bookRepository: bookRepository}
}

func (bs *GrpcBookService) CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error) {
	response, err := bs.bookRepository.CheckAvailabilityByAuthorTitle(ctx, request)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}
