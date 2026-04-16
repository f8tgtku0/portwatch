// Package notify provides pluggable notifiers for portwatch.
//
// Each notifier implements a Send method that accepts a history.Entry
// and delivers an alert via its configured channel.
//
// Available notifiers:
//
//   - LogNotifier  — writes formatted lines to an io.Writer (default: stdout)
//   - WebhookNotifier — HTTP POST JSON payload to a URL
//   - EmailNotifier — sends SMTP email alerts
//   - SlackNotifier — posts messages to a Slack incoming webhook
//
// Multiple notifiers can be composed with NewMulti, which fans out a
// single Send call to all registered notifiers and collects errors.
package notify
