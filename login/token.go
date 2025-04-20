package login

import (
	"context"

	"github.com/marcus-crane/wwenet/storage"
)

const (
	AccessTokenID  = "accessToken"
	RefreshTokenID = "refreshToken"
)

func UpsertToken(ctx context.Context, id string, token string, db *storage.Queries) error {
	existingToken, err := db.GetToken(ctx, id)
	if err != nil {
		// we don't have an existing token so let's insert it
		params := storage.CreateTokenParams{
			ID:    id,
			Value: token,
		}
		_, err := db.CreateToken(ctx, params)
		return err
	}
	if existingToken.Value != token {
		params := storage.UpdateTokenParams{
			ID:    id,
			Value: token,
		}
		return db.UpdateToken(ctx, params)
	}
	// no-op as there is nothing to update
	return nil
}
