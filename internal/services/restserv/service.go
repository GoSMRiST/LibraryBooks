package restserv

import (
	"2/internal/core"
	"context"
)

type DBRepository interface {
	CreateTable(ctx context.Context) error
	AddBook(ctx context.Context, books core.Book) error
	GetAllBooks(ctx context.Context) ([]core.Book, error)
	GetBookByID(ctx context.Context, bookID int) (core.Book, error)
}

type RestBookService struct {
	repo DBRepository
}

func NewBookService(repo DBRepository) *RestBookService {
	return &RestBookService{repo: repo}
}

func (s *RestBookService) CreateTable(ctx context.Context) error {
	return s.repo.CreateTable(ctx)
}

func (s *RestBookService) AddBook(ctx context.Context, book core.Book) error {
	if book.Author == "" || book.Title == "" {
		return core.ErrInvalidInput
	}

	return s.repo.AddBook(ctx, book)
}

func (s *RestBookService) GetAllBooks(ctx context.Context) ([]core.Book, error) {
	return s.repo.GetAllBooks(ctx)
}

func (s *RestBookService) GetBookByID(ctx context.Context, bookID int) (core.Book, error) {
	if bookID <= 0 {
		return core.Book{}, core.ErrInvalidInput
	}

	return s.repo.GetBookByID(ctx, bookID)
}
