package main

import (
	"cf-npm-ip-sync/internal"
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	ips, err := internal.GetCloudflareIPs(ctx)
	if err != nil {
		log.Fatalf("failed to get cloudflare IPs: %v", err)
	}

	log.Printf("cloudflare IPs: %v", ips)
}
