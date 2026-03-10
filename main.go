package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: check-mail-domain <domain> <aws-region>")
		os.Exit(1)
	}

	domain := os.Args[1]
	region := os.Args[2]

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Printf("Error loading AWS config: %v\n", err)
		os.Exit(1)
	}

	client := sesv2.NewFromConfig(cfg)

	// Get DKIM attributes from SES
	resp, err := client.GetEmailIdentity(ctx, &sesv2.GetEmailIdentityInput{
		EmailIdentity: &domain,
	})
	if err != nil {
		fmt.Printf("Error getting SES identity: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== SES DKIM Configuration ===")
	if resp.DkimAttributes != nil {
		fmt.Printf("Status: %s\n", resp.DkimAttributes.Status)
		if resp.DkimAttributes.Tokens != nil {
			fmt.Println("\nDKIM Tokens (selectors):")
			for _, token := range resp.DkimAttributes.Tokens {
				fmt.Printf("  %s._domainkey.%s\n", token, domain)
			}
		}
	}

	// Custom MAIL FROM domain
	var mailFromDomain string
	fmt.Println("\n=== Custom MAIL FROM Domain ===")
	if resp.MailFromAttributes != nil {
		mailFromDomain = *resp.MailFromAttributes.MailFromDomain
		fmt.Printf("MAIL FROM Domain: %s\n", mailFromDomain)
		fmt.Printf("Status: %s\n", resp.MailFromAttributes.MailFromDomainStatus)
		fmt.Printf("Behavior on MX Failure: %s\n", resp.MailFromAttributes.BehaviorOnMxFailure)
	} else {
		fmt.Println("No custom MAIL FROM domain configured")
	}

	// Check DNS records
	records := []string{
		"_amazonses." + domain,
		"_dmarc." + domain,
	}

	// Add MAIL FROM domain MX record check
	if mailFromDomain != "" {
		records = append(records, "MX:"+mailFromDomain)
	}

	// Add DKIM tokens from SES
	if resp.DkimAttributes != nil && resp.DkimAttributes.Tokens != nil {
		for _, token := range resp.DkimAttributes.Tokens {
			records = append(records, token+"._domainkey."+domain)
		}
	}

	fmt.Println("\n=== DNS Records ===")
	for _, record := range records {
		// Check if this is an MX record check
		if len(record) > 3 && record[:3] == "MX:" {
			mxDomain := record[3:]
			mxRecords, err := net.LookupMX(mxDomain)
			if err != nil || len(mxRecords) == 0 {
				fmt.Printf("%-60s \033[31mFail\033[0m\n", mxDomain+" (MX)")
				if err != nil {
					fmt.Printf("  Error: %v\n", err)
				}
			} else {
				fmt.Printf("%-60s \033[32mPass\033[0m\n", mxDomain+" (MX)")
				for _, mx := range mxRecords {
					fmt.Printf("  Priority: %d, Host: %s\n", mx.Pref, mx.Host)
				}
			}
		} else {
			txtRecords, err := net.LookupTXT(record)
			if err != nil || len(txtRecords) == 0 {
				fmt.Printf("%-60s \033[31mFail\033[0m\n", record)
			} else {
				fmt.Printf("%-60s \033[32mPass\033[0m\n", record)
				for _, txt := range txtRecords {
					fmt.Printf("  %s\n", txt)
				}
			}
		}
	}

	// Check MX records
	fmt.Println("\n=== MX Records ===")
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		fmt.Printf("%-60s \033[31mFail\033[0m\n", domain)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		}
	} else {
		fmt.Printf("%-60s \033[32mPass\033[0m\n", domain)
		for _, mx := range mxRecords {
			fmt.Printf("  Priority: %d, Host: %s\n", mx.Pref, mx.Host)
		}
	}
}
