package login

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/marcus-crane/wwenet/config"
	_ "modernc.org/sqlite"
)

const (
	LOGIN_API_URL = "https://dce-frontoffice.imggaming.com/api/v2/login"
)

type loginPayload struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
}

type loginResponse struct {
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

type refreshPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type refreshResponse struct{}

func Login(ctx context.Context, cfg config.Config) (loginResponse, error) {
	lBody := loginPayload{
		Id:     cfg.Credentials.Username,
		Secret: cfg.Credentials.Password,
	}
	payload, err := json.Marshal(lBody)
	if err != nil {
		return loginResponse{}, err
	}
	req, err := http.NewRequest(http.MethodPost, LOGIN_API_URL, bytes.NewBuffer(payload))
	if err != nil {
		return loginResponse{}, err
	}
	req.Header.Add("Realm", "dce.wwe")
	req.Header.Add("Content-Type", "application/json") // Required
	req.Header.Add("x-app-var", "6.0.1.f8add0e")
	req.Header.Add("x-api-key", cfg.Credentials.APIKey)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:138.0) Gecko/20100101 Firefox/138.0")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return loginResponse{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return loginResponse{}, err
	}
	var lResp loginResponse
	if err = json.Unmarshal(body, &lResp); err != nil {
		return loginResponse{}, err
	}
	return lResp, nil
}
