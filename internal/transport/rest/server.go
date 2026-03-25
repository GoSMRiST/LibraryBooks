package rest

import (
	"2/internal/core"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ServiceInterface interface {
	AddBook(ctx context.Context, book core.Book) error
	GetAllBooks(ctx context.Context) ([]core.Book, error)
	GetBookByID(ctx context.Context, bookID int) (core.Book, error)
}

type bookHandler struct {
	service ServiceInterface
}

func NewBookHandler(service ServiceInterface) *bookHandler {
	return &bookHandler{service: service}
}

func (h *bookHandler) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/books", h.GetAllBooks)
	engine.GET("/books/:id", h.GetBookByID)
	engine.POST("/books", h.AddBook)
}

func (h *bookHandler) GetAllBooks(ctx *gin.Context) {
	books, err := h.service.GetAllBooks(ctx.Request.Context())
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	booksDTO := make([]core.BookResponseDTO, len(books))
	for i, book := range books {
		booksDTO[i] = core.BookResponseDTO{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": booksDTO})
}

func (h *bookHandler) AddBook(ctx *gin.Context) {
	var bookDTO core.BookCreateDTO

	if err := ctx.ShouldBindJSON(&bookDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	book := core.Book{
		Title:  bookDTO.Title,
		Author: bookDTO.Author,
	}

	err := h.service.AddBook(ctx.Request.Context(), book)
	if err != nil {
		if errors.Is(err, core.ErrInvalidInput) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
			return
		} else if errors.Is(err, core.ErrDuplicateBook) {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": bookDTO})
}

func (h *bookHandler) GetBookByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	bookID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "invalid id",
		})
		return
	}

	book, err := h.service.GetBookByID(ctx.Request.Context(), bookID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "error": err.Error()})
			return
		} else if errors.Is(err, core.ErrInvalidInput) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	bookDTO := core.BookResponseDTO{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": bookDTO})
}
