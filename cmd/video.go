package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/login"
	"github.com/marcus-crane/wwenet/storage"
)

func Video(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	var token string
	accessToken, err := db.GetToken(ctx, "accessToken")
	if err != nil {
		// We don't have a token so we need to fetch one
		credentials, loginErr := login.Login(ctx, cfg)
		if loginErr != nil {
			return loginErr
		}
		if accessUpsertErr := login.UpsertToken(ctx, login.AccessTokenID, credentials.AuthorisationToken, db); accessUpsertErr != nil {
			return accessUpsertErr
		}
		if refreshTokenErr := login.UpsertToken(ctx, login.RefreshTokenID, credentials.RefreshToken, db); refreshTokenErr != nil {
			return refreshTokenErr
		}
		token = credentials.AuthorisationToken
	} else {
		token = accessToken.Value
	}
	fmt.Printf("Fetched access token %s\n", token)
	return nil
}
