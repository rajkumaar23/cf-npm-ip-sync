package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NPMClient struct {
	Host  string
	Token string
}

func NewNPMClient(host string, email string, password string) (*NPMClient, error) {
	c := &NPMClient{
		Host: host,
	}

	err := c.acquireAccessToken(email, password)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *NPMClient) acquireAccessToken(email string, password string) error {
	payload := map[string]string{
		"scope":    "user",
		"identity": email,
		"secret":   password,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/api/tokens", c.Host), "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("unexpected status code: %d, response body: %s", res.StatusCode, string(body))
	}

	var response struct {
		Token   string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	c.Token = response.Token
	return nil
}
