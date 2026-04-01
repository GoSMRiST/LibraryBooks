package grpc

import (
	"2/internal/core"
	"context"
	book "github.com/GoSMRiST/protosLibary/gen/go/book"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcService interface {
	CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

type Server struct {
	book.UnimplementedBookServer
	grpcService GrpcService
}

func NewServer(gRPC *grpc.Server, bookService GrpcService) {
	book.RegisterBookServer(gRPC, &Server{grpcService: bookService})
}

func (serv *Server) CheckAvailability(ctx context.Context, request *book.CheckRequest) (*book.CheckResponse, error) {
	coreRequest := &core.CheckAvailabilityRequest{
		Author: request.Author,
		Title:  request.Title,
	}

	resp, err := serv.grpcService.CheckAvailabilityByAuthorTitle(ctx, coreRequest)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return &book.CheckResponse{
				Result: false,
			}, st.Err()
		}

		return &book.CheckResponse{
			Result: false,
		}, status.Error(codes.Internal, err.Error())
	}

	return &book.CheckResponse{
		Result: resp.Result,
	}, nil
}
