package core

type BookCreateDTO struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type BookResponseDTO struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
}
