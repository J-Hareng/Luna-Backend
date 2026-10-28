# Luna Backend

Luna Backend is a Go-based HTTP API for user accounts, teams, tasks, collections, and file management. It uses Gin for HTTP routing, MongoDB for persistence, AWS S3 for file storage, and email-based account/session functionality.

> **Current status:** This project is **not currently running**.

## Requirements

- Go 1.21 or newer
- MongoDB
- An AWS account with access to the configured S3 bucket
- SMTP/email credentials

## Configuration

Create a `.env` file in the project root with the values required by the application:

```env
PORT=8080
MONGOURI=mongodb://localhost:27017
AWS_ACCESS_KEY=your-access-key
AWS_SECRET_KEY=your-secret-key
```

Email configuration is also required by the email service. Do not commit credentials or other secrets to the repository.

## Running locally

Install dependencies and start the backend with:

```bash
go mod download
go run src/main.go
```

The server listens on `0.0.0.0:$PORT`. Protected endpoints are served under `/api` and require the session-token middleware.

## Main API areas

- Authentication and session tokens
- Users and group management
- Teams and team membership
- Tasks and task applications
- Collections and tags
- File upload, download, and deletion through S3

## Project layout

```text
src/
├── api/          Database, email, and S3 integrations
├── helper/       Shared helpers and environment access
└── httpd/        HTTP server, handlers, and security middleware
```
