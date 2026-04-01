package middleware

import (
	"2/internal/core"
	"context"
	"github.com/GoSMRiST/protosLibary/gen/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(authClient auth.AuthClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		token = strings.TrimSpace(token)

		validateData := &auth.ValidateTokenRequest{
			Token: token,
		}

		resp, err := authClient.ValidateToken(ctx, validateData)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unauthenticated {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": st.Message()})
				return
			}

			slog.Info("error", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		tokenData := core.TokenData{
			UserID: int(resp.UserId),
			Role:   resp.Role,
		}

		ctx.Set("token_data", tokenData)

		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), core.TokenDataKey, tokenData),
		)

		ctx.Next()
	}
}
