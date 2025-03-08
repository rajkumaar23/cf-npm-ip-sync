package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetCloudflareIPs(ctx context.Context) ([]string, error) {
	ipv4Ranges, err := getIPRanges("v4")
	if err != nil {
		return nil, err
	}

	ipv6Ranges, err := getIPRanges("v6")
	if err != nil {
		return nil, err
	}

	return append(ipv4Ranges, ipv6Ranges...), nil
}

func getIPRanges(ipType string) ([]string, error) {
	url := fmt.Sprintf("https://www.cloudflare.com/ips-%s", ipType)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s ranges: %w", ipType, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
		if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	ips := strings.Split(string(body), "\n")
	return ips, nil
}
