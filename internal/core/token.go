package core

type contextKey string

const TokenDataKey contextKey = "token_data"

type TokenData struct {
	UserID int
	Role   string
}
