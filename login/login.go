package login

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/storage"
	_ "modernc.org/sqlite"
)

const (
	LOGIN_API_URL   = "https://dce-frontoffice.imggaming.com/api/v2/login"
	REFRESH_API_URL = "https://dce-frontoffice.imggaming.com/api/v2/token/refresh"
)

type loginPayload struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

type tokenResponse struct {
	AuthorisationToken       string `json:"authorisationToken"`
	RefreshToken             string `json:"refreshToken"`
	MissingStatusInformation string `json:"missingInformationStatus"`
}

type loginError struct {
	Status    int      `json:"status"`
	Code      string   `json:"code"`
	Messages  []string `json:"messages"`
	RequestID string   `json:"requestId"`
}

func (le loginError) Error() string {
	return fmt.Sprintf("login failed (status %d, code %s): %v", le.Status, le.Code, le.Messages)
}

type refreshPayload struct {
	RefreshToken string `json:"refreshToken"`
}

func Login(ctx context.Context, cfg config.Config) (tokenResponse, error) {
	lBody := loginPayload{
		Id:     cfg.Credentials.Username,
		Secret: cfg.Credentials.Password,
	}
	payload, err := json.Marshal(lBody)
	if err != nil {
		return tokenResponse{}, err
	}
	req, err := http.NewRequest(http.MethodPost, LOGIN_API_URL, bytes.NewBuffer(payload))
	if err != nil {
		return tokenResponse{}, err
	}
	req.Header.Add("Realm", "dce.wwe")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-app-var", cfg.Network.XAppVar)
	req.Header.Add("x-api-key", cfg.Network.XApiKey)
	req.Header.Add("User-Agent", cfg.Network.UserAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return tokenResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return tokenResponse{}, err
	}

	if res.StatusCode >= 400 {
		var lErr loginError
		if err := json.Unmarshal(body, &lErr); err != nil {
			return tokenResponse{}, fmt.Errorf("login failed with status %d: %s", res.StatusCode, string(body))
		}
		return tokenResponse{}, lErr
	}

	var lResp tokenResponse
	if err = json.Unmarshal(body, &lResp); err != nil {
		return tokenResponse{}, err
	}
	return lResp, nil
}

func RefreshToken(ctx context.Context, cfg config.Config, refreshToken string) (tokenResponse, error) {
	rBody := refreshPayload{
		RefreshToken: refreshToken,
	}
	payload, err := json.Marshal(rBody)
	if err != nil {
		return tokenResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, REFRESH_API_URL, bytes.NewBuffer(payload))
	if err != nil {
		return tokenResponse{}, err
	}
	req.Header.Add("Realm", "dce.wwe")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-app-var", cfg.Network.XAppVar)
	req.Header.Add("x-api-key", cfg.Network.XApiKey)
	req.Header.Add("User-Agent", cfg.Network.UserAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return tokenResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return tokenResponse{}, err
	}

	if res.StatusCode >= 400 {
		var lErr loginError
		if err := json.Unmarshal(body, &lErr); err != nil {
			return tokenResponse{}, fmt.Errorf("token refresh failed with status %d: %s", res.StatusCode, string(body))
		}
		return tokenResponse{}, lErr
	}

	var tResp tokenResponse
	if err = json.Unmarshal(body, &tResp); err != nil {
		return tokenResponse{}, err
	}
	return tResp, nil
}

func GetAuthToken(ctx context.Context, cfg config.Config, db *storage.Queries) (string, error) {
	accessToken, err := db.GetToken(ctx, "accessToken")
	if err != nil || IsTokenExpired(accessToken) {
		fmt.Println("Access token is missing or expired")

		if err == nil {
			if refreshed, refreshErr := tryRefreshToken(ctx, cfg, db); refreshErr != nil {
				fmt.Println("Successfully refreshed token")
				return refreshed, nil
			}
			fmt.Println("Refresh failed, performing full login...")
		}

		return performFullLogin(ctx, cfg, db)
	}

	return accessToken.Value, nil
}

func tryRefreshToken(ctx context.Context, cfg config.Config, db *storage.Queries) (string, error) {
	refreshToken, err := db.GetToken(ctx, "refreshToken")
	if err != nil {
		return "", err
	}

	credentials, err := RefreshToken(ctx, cfg, refreshToken.Value)
	if err != nil {
		return "", err
	}

	if err := UpsertToken(ctx, AccessTokenID, credentials.AuthorisationToken, db); err != nil {
		return "", err
	}

	if err := UpsertToken(ctx, RefreshTokenID, credentials.RefreshToken, db); err != nil {
		return "", err
	}

	return credentials.AuthorisationToken, nil
}

func performFullLogin(ctx context.Context, cfg config.Config, db *storage.Queries) (string, error) {
	credentials, err := Login(ctx, cfg)
	if err != nil {
		return "", fmt.Errorf("login failed: %w", err)
	}

	if err := UpsertToken(ctx, AccessTokenID, credentials.AuthorisationToken, db); err != nil {
		return "", err
	}
	if err := UpsertToken(ctx, RefreshTokenID, credentials.RefreshToken, db); err != nil {
		return "", err
	}

	return credentials.AuthorisationToken, nil
}
