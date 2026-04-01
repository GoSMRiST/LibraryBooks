package rest

import (
	"2/internal/core"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type ServiceInterface interface {
	AddBook(ctx context.Context, book *core.Book) error
	GetAllBooks(ctx context.Context) ([]core.Book, error)
	GetBookByID(ctx context.Context, bookID int) (*core.Book, error)
}

type BookHandler struct {
	log     *slog.Logger
	service ServiceInterface
}

func NewBookHandler(log *slog.Logger, service ServiceInterface) *BookHandler {
	return &BookHandler{
		log:     log,
		service: service,
	}
}

func (h *BookHandler) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/books", h.GetAllBooks)
	engine.GET("/books/:id", h.GetBookByID)
	engine.POST("/books", h.AddBook)
}

func (h *BookHandler) GetAllBooks(ctx *gin.Context) {
	books, err := h.service.GetAllBooks(ctx.Request.Context())
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error":  "books not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "internal server error",
		})
		return
	}

	booksDTO := make([]core.BookResponseDTO, len(books))
	for i, book := range books {
		booksDTO[i] = core.BookResponseDTO{
			ID:     book.ID,
			Author: book.Author,
			Title:  book.Title,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   booksDTO,
	})
}

func (h *BookHandler) AddBook(ctx *gin.Context) {
	var bookDTO core.BookCreateDTO

	if err := ctx.ShouldBindJSON(&bookDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  core.ErrInvalidInput,
		})
		return
	}

	book := &core.Book{
		Title:  bookDTO.Title,
		Author: bookDTO.Author,
	}

	err := h.service.AddBook(ctx.Request.Context(), book)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrUnauthorized):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"error":  "unauthorized",
			})
			return

		case errors.Is(err, core.ErrNoRights):
			ctx.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "no rights",
			})
			return

		case errors.Is(err, core.ErrDuplicateBook):
			ctx.JSON(http.StatusConflict, gin.H{
				"status": "error",
				"error":  "book already exists",
			})
			return

		case errors.Is(err, core.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid input",
			})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "internal server error",
			})
		}

		h.log.Error("book not added", "error", err)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": bookDTO})
}

func (h *BookHandler) GetBookByID(ctx *gin.Context) {
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
		switch {
		case errors.Is(err, core.ErrNoRights):
			ctx.JSON(http.StatusForbidden, gin.H{
				"status": "error",
				"error":  "no rights",
			})
			return

		case errors.Is(err, core.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error":  "books not found",
			})
			return

		case errors.Is(err, core.ErrInvalidInput):
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid input",
			})
			return

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  "internal server error",
			})
		}

		h.log.Error("error in book output", "error", err)

		return
	}

	bookDTO := core.BookResponseDTO{
		ID:     book.ID,
		Author: book.Author,
		Title:  book.Title,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": bookDTO})
}
