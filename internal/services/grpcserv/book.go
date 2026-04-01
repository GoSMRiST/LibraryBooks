package grpcserv

import (
	"2/internal/core"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DbInterface interface {
	CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

type GrpcBookService struct {
	bookRepository DbInterface
}

func NewGrpcBookService(bookRepository DbInterface) *GrpcBookService {
	return &GrpcBookService{bookRepository: bookRepository}
}

func (bs *GrpcBookService) CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error) {
	response, err := bs.bookRepository.CheckAvailabilityByAuthorTitle(ctx, request)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return response, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}
