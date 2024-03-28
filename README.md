# SSL Certificate Expiration Checker

## Setup

1. Clone the repository:

```bash
git clone https://github.com/oscar3x39/ssl-expiration-checker.git
```

2. Install dependencies:

```bash
go mod tidy
```

3. Compile the Go code:
```bash
go build -o ssl-checker main.go
```

4. Create a config.yaml file with the following structure:
```bash
slack_webhook_url: "https://hooks.slack.com/services/your-webhook-url"
domains:
  - name: "Example Domain"
    url: "example.com"
    contact: "John Doe"
  # Add more domains as needed
```

5. Replace your-webhook-url with your actual Slack webhook URL.

# Usag
To manually run the program, execute the following command:
```
./ssl-checker
```

# Scheduler

Use cron jobs to schedule the execution of the ssl-checker binary. Edit the crontab file by running:
```bash
crontab -e
```

Then, add the following line to execute the program daily at midnight:
```bash
0 0 * * * /path/to/ssl-checker
```


# Functionality
- The program reads the configuration from config.yaml, which includes the Slack webhook URL and a list of domains to check.
- For each domain, it checks the validity of SSL certificates.
- If a certificate is expiring soon (within 30 days) or if there's a TLS connection issue, a notification is sent to Slack using the provided webhook URL.
- The notification includes details such as domain name, URL, contact person, and error message.

# Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or create a pull request.

# License
This project is licensed under the MIT License.