package restserv

import (
	"2/internal/core"
	"context"
)

type DbInterface interface {
	AddBook(ctx context.Context, books *core.Book) error
	GetAllBooks(ctx context.Context) ([]core.Book, error)
	GetBookByID(ctx context.Context, bookID int) (*core.Book, error)
}

type RestBookService struct {
	repo DbInterface
}

func NewRestBookService(repo DbInterface) *RestBookService {
	return &RestBookService{repo: repo}
}

func (s *RestBookService) AddBook(ctx context.Context, book *core.Book) error {
	tokenData, ok := ctx.Value(core.TokenDataKey).(core.TokenData)
	if !ok {
		return core.ErrUnauthorized
	}

	if tokenData.Role == "user" {
		return core.ErrNoRights
	}

	if book.Author == "" || book.Title == "" {
		return core.ErrInvalidInput
	}

	return s.repo.AddBook(ctx, book)
}

func (s *RestBookService) GetAllBooks(ctx context.Context) ([]core.Book, error) {
	return s.repo.GetAllBooks(ctx)
}

func (s *RestBookService) GetBookByID(ctx context.Context, bookID int) (*core.Book, error) {
	if ctx.Value("role") == "user" {
		return nil, core.ErrNoRights
	}

	if bookID <= 0 {
		return nil, core.ErrInvalidInput
	}

	return s.repo.GetBookByID(ctx, bookID)
}
