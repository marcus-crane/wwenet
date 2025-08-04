package login

import (
	"context"
	"database/sql"
	"time"

	"github.com/marcus-crane/wwenet/storage"
)

const (
	AccessTokenID        = "accessToken"
	RefreshTokenID       = "refreshToken"
	TokenLifetimeMinutes = 10
)

func UpsertToken(ctx context.Context, id string, token string, db *storage.Queries) error {
	expiresAt := time.Now().Add(TokenLifetimeMinutes * time.Minute).Unix()

	existingToken, err := db.GetToken(ctx, id)
	if err != nil {
		// we don't have an existing token so let's insert it
		params := storage.CreateTokenParams{
			ID:        id,
			Value:     token,
			ExpiresAt: sql.NullInt64{Int64: expiresAt, Valid: true},
		}
		_, err := db.CreateToken(ctx, params)
		return err
	}
	if existingToken.Value != token {
		params := storage.UpdateTokenParams{
			ID:        id,
			Value:     token,
			ExpiresAt: sql.NullInt64{Int64: expiresAt, Valid: true},
		}
		return db.UpdateToken(ctx, params)
	}
	// no-op as there is nothing to update
	return nil
}

func IsTokenExpired(token storage.Token) bool {
	if !token.ExpiresAt.Valid {
		return true
	}

	// We'll refresh 1 minute in advance for safety
	bufferSeconds := int64(60)
	return time.Now().Unix() >= (token.ExpiresAt.Int64 - bufferSeconds)
}
