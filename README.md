# A dead simple, no frills, but useful, Go watcher for jobs updates in the EMSI community's job feed

## Quick Start

To quickly get started, you could use docker to run the watcher.

```bash
docker run -d --name watcher -e EC_EMAIL=<your-email> -e EC_PASSWORD=<your-password> -e EC_DISCORD_WEBHOOK=<your-discord-webhook-url> xenedium/ec-job-watcher
```

## Description

This watcher will poll the EMSI Community job feed every 5 minutes and send a notification to a Discord webhook or a function URL when a new job is posted, keeping you up to date with the latest job postings as soon as they are posted giving you a competitive edge.

## Configuration

The watcher is configured via environment variables.

| Variable | Description |
| --- | --- |
| `EC_EMAIL` | Your EMSI Community email address. |
| `EC_PASSWORD` | Your EMSI Community password. |
| `EC_DISCORD_WEBHOOK` | The Discord webhook URL to send notifications to. |
| `EC_FUNC_URL` | Your function URL to send notifications to, e.g AWS Lambda. |

Note: If you specify both `EC_DISCORD_WEBHOOK` and `EC_FUNC_URL`, the watcher will send notifications to both.
