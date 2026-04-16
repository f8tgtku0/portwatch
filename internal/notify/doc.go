// Package notify provides multiple notification backends for portwatch.
//
// Supported backends:
//   - Log (stdout/writer)
//   - Webhook (generic HTTP)
//   - Slack
//   - PagerDuty
//   - OpsGenie
//   - Email (SMTP)
//
// All backends implement the Notifier interface and can be composed
// using NewMulti to fan out alerts to multiple destinations.
package notify
