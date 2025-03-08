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
		Token string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	c.Token = response.Token
	return nil
}

type AccessListClient struct {
	Address   string `json:"address"`
	Directive string `json:"directive"`
}

type AccessList struct {
	Name       string             `json:"name"`
	SatisfyAny bool               `json:"satisfy_any"`
	PassAuth   bool               `json:"pass_auth"`
	Items      []interface{}      `json:"items"`
	Clients    []AccessListClient `json:"clients"`
}

func (c *NPMClient) getAccessList(id int) (*AccessList, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/nginx/access-lists/%d?expand=items,clients", c.Host, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access list: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("unexpected status code: %d, response body: %s", res.StatusCode, string(body))
	}

	var response AccessList
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (c *NPMClient) UpdateAccessList(id int, ips []string) error {
	accessList, err := c.getAccessList(id)
	if err != nil {
		return err
	}

	var clients []AccessListClient
	for _, ip := range ips {
		clients = append(clients, AccessListClient{
			Address:   ip,
			Directive: "allow",
		})
	}

	accessList.Clients = clients
	jsonPayload, err := json.Marshal(accessList)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/nginx/access-lists/%d", c.Host, id), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update access list: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return fmt.Errorf("unexpected status code: %d, response body: %s", res.StatusCode, string(body))
	}
	return nil
}
