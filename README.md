# SMTP Server & Frontend

A lightweight SMTP server written in Go that captures emails in memory and displays them in a modern web frontend.

## 1. Installation

### Setup
Clone the repository: `git clone https://github.com/gnpaone/basic-smtp.git`

```bash
cd basic-smtp
go mod download
```

## 2. Running the Server

Start the application:

```bash
go run .
```

You should see output indicating both servers are running:
```text
Starting SMTP Server at :1025
Starting HTTP Frontend at :8080
```

- **SMTP Server**: `localhost:1025`
- **Web Interface**: `http://localhost:8080`

## 3. Testing (How to send emails)

You can send test emails using various tools.

### Option A: Using Python
```bash
python3 -c 'import smtplib; from email.message import EmailMessage; msg = EmailMessage(); msg.set_content("This is a test body.\n\nCheers!"); msg["Subject"] = "Hello World!"; msg["From"] = "sender@example.com"; msg["To"] = "receiver@example.com"; s=smtplib.SMTP("localhost", 1025); s.send_message(msg); s.quit(); print("Email sent successfully!")'
```

### Option B: Using `swaks` (Swiss Army Knife for SMTP)
If you have `swaks` installed:
```bash
swaks --to user@example.com --from me@example.com --server localhost --port 1025 --header "Subject: Test email" --body "Testing the SMTP server."
```

### Option C: Using `curl`
```bash
curl --url "smtp://localhost:1025" \
  --mail-from "sender@example.com" \
  --mail-rcpt "receiver@example.com" \
  --upload-file - <<EOF
From: sender@example.com
To: receiver@example.com
Subject: Test email

This is a test email.
EOF
```

## 4. Viewing Emails
Open your browser to [http://localhost:8080](http://localhost:8080). The interface updates automatically every 3 seconds.
