package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

// Config struct is used for parsing the configuration file
type Config struct {
	SlackWebhookURL string   `yaml:"slack_webhook_url"`
	Domains         []Domain `yaml:"domains"`
}

// Domain struct represents each domain entry in the configuration file
type Domain struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Contact string `yaml:"contact"`
}

func main() {
	// Read the configuration file
	config, err := readConfig("config.yaml")
	if err != nil {
		log.Fatalf("Unable to read the configuration file: %v", err)
	}

	// Check the validity of certificates for each domain
	for _, domain := range config.Domains {
		err := checkCertValidity(domain.URL)
		if err != nil {
			// Generate Slack message based on the error type
			var errorMessage string
			if err == ErrTLSConnection {
				errorMessage = fmt.Sprintf("Unable to establish TLS connection: %v", err)
			} else if err == ErrCertExpiring {
				errorMessage = fmt.Sprintf("Certificate is expiring soon: %v", err)
			}
			errSlack := sendSlackMessage(domain, errorMessage) // Pass domain as the first argument
			if errSlack != nil {
				log.Printf("Unable to send message to Slack: %v", errSlack)
			}
		} else {
			log.Printf("Certificate check passed for: %s", domain.URL)
		}
	}
}

// readConfig function is used to read the configuration file
func readConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// ErrTLSConnection represents the error for TLS connection issues
var ErrTLSConnection = fmt.Errorf("Unable to establish TLS connection")

// ErrCertExpiring represents the error for expiring certificates
var ErrCertExpiring = fmt.Errorf("Certificate is expiring soon")

// checkCertValidity function is used to check the certificate validity for a specified domain
func checkCertValidity(domain string) error {
	// Connect to the HTTPS port of the domain
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		if _, ok := err.(tls.RecordHeaderError); ok {
			return ErrTLSConnection
		}
		return err
	}
	defer conn.Close()

	// Get the expiry time of the certificate
	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter

	// Calculate the days until certificate expiration
	daysUntilExpiration := int(expiry.Sub(time.Now()).Hours() / 24)

	// Return ErrCertExpiring error if the certificate expires within a month
	if daysUntilExpiration <= 30 {
		return ErrCertExpiring
	}

	return nil
}

func sendSlackMessage(domain Domain, errorMessage string) error {
	// Slack Webhook URL (hardcoded)
	webhookURL := "https://hooks.slack.com/services/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"

	// Construct the message content, using Markdown formatting, and adding emoji
	message := fmt.Sprintf(":warning: *Certificate check failed:*\n*Name:* %s\n*URL:* <%s|%s>\n*Contact:* %s\n*Error:* %s", domain.Name, domain.URL, domain.URL, domain.Contact, errorMessage)

	// Construct the payload for the message
	payload := map[string]string{"text": message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Unable to serialize JSON: %v", err)
	}

	// Send the POST request to Slack Webhook
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("Unable to send POST request to Slack Webhook: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Slack Webhook response status code: %d", resp.StatusCode)
	}

	return nil
}