// Package notify provides notification backends for portwatch.
//
// Supported notifiers:
//
//   - LogNotifier  – writes formatted lines to an io.Writer (default: stdout)
//   - WebhookNotifier – HTTP POST JSON payload to a remote URL
//   - EmailNotifier – sends SMTP email alerts via PlainAuth
//
// Multiple notifiers can be composed with NewMulti, which fans out a single
// Message to every registered backend. Any backend error is collected and
// returned as a combined error.
//
// # Message
//
// Each Send call receives a Message that wraps a state.Change (the port event)
// together with a Timestamp so backends can render consistent time strings
// regardless of when the message is ultimately delivered.
package notify
