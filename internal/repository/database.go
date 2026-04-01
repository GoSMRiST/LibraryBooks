package repository

import (
	"2/internal/core"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type DataBase struct {
	conn *pgx.Conn
}

func InitDataBase(ctx context.Context, connString string) (*DataBase, error) {
	dataBase := &DataBase{}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return dataBase, err
	}

	dataBase.conn = conn
	return dataBase, nil
}

func (dataBase *DataBase) CloseDataBase(ctx context.Context) error {
	return dataBase.conn.Close(ctx)
}

func (dataBase *DataBase) CreateTable(ctx context.Context) error {
	strQuery := `
		CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		author TEXT,
		title TEXT
		);
	`

	_, err := dataBase.conn.Exec(ctx, strQuery)
	if err != nil {
		return err
	}

	return nil
}

func (dataBase *DataBase) AddBook(ctx context.Context, book *core.Book) error {
	var exists bool

	err := dataBase.conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM books WHERE title=$1 AND author=$2)",
		book.Title, book.Author,
	).Scan(&exists)
	if err != nil {
		return err // ошибка БД
	} else if exists {
		return core.ErrDuplicateBook
	}

	strQuery := "INSERT INTO books (author, title) VALUES ($1, $2)"

	_, err = dataBase.conn.Exec(ctx, strQuery, book.Author, book.Title)
	if err != nil {
		return err
	}

	return nil
}

func (dataBase *DataBase) GetAllBooks(ctx context.Context) ([]core.Book, error) {
	var books []core.Book
	rows, err := dataBase.conn.Query(ctx, "SELECT id, author, title FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var book core.Book
		if err := rows.Scan(&book.ID, &book.Author, &book.Title); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (dataBase *DataBase) GetBookByID(ctx context.Context, bookID int) (*core.Book, error) {
	book := &core.Book{}
	strQuery := "SELECT id, author, title FROM books WHERE id = $1"

	row := dataBase.conn.QueryRow(ctx, strQuery, bookID)
	if err := row.Scan(&book.ID, &book.Author, &book.Title); err != nil {
		return nil, core.ErrNotFound
	}

	return book, nil
}

func (dataBase *DataBase) CheckAvailabilityByAuthorTitle(ctx context.Context, request *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error) {
	result := &core.CheckAvailabilityResponse{Result: false}
	strQuery := "SELECT 1 FROM books WHERE author = $1 AND title = $2"

	// просто проверяем наличие записи
	err := dataBase.conn.QueryRow(ctx, strQuery, request.Author, request.Title).Scan(new(int))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, core.ErrNotFound // книги нет
		}

		return result, err // другая ошибка
	}

	result.Result = true // запись найдена
	return result, nil
}
