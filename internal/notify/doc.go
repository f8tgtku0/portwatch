// Package notify provides pluggable notification backends for portwatch.
//
// The core abstraction is the Notifier interface, which exposes a single
// Send(Message) method. Two implementations are included:
//
//   - LogNotifier — writes human-readable lines to any io.Writer (default
//     os.Stdout). Suitable for terminal output and log aggregation.
//
//   - Multi — fan-out wrapper that forwards every message to a list of
//     Notifiers. Useful when multiple sinks (e.g. log file + webhook) must
//     receive the same event.
//
// # Usage
//
//	n := notify.NewLogNotifier(os.Stderr)
//	_ = n.Send(notify.Message{
//		Level: notify.LevelAlert,
//		Title: "Port opened",
//		Body:  "port 22 is now open",
//	})
package notify
