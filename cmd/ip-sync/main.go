package main

import (
	"cf-npm-ip-sync/internal"
	"context"
	"log"
)

func main() {
	c, err := internal.NewConfig()
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}

	ctx := context.Background()
	ips, err := internal.GetCloudflareIPs(ctx)
	if err != nil {
		log.Fatalf("failed to get cloudflare IPs: %v", err)
	}

	nc, err := internal.NewNPMClient(c.NPMHost, c.NPMEmail, c.NPMPassword)
	if err != nil {
		log.Fatalf("failed to create npm client: %v", err)
	}
	diff, err := nc.UpdateAccessList(c.NPMAccessListID, ips)
	if err != nil {
		log.Fatalf("failed to update access list: %v", err)
	}

	if diff > 0 {
		log.Printf("added %d IPs to access list", diff)
	} else if diff < 0 {
		log.Printf("removed %d IPs from access list", -diff)
	} else {
		log.Println("no changes to access list")
	}
}
