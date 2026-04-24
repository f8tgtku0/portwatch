// Package shadow implements a shadow-mode notifier that mirrors traffic
// to a secondary notification channel without affecting the primary one.
//
// Use shadow mode to validate a new notifier (e.g. a new webhook endpoint
// or a different alert channel) against real production events before
// cutting over fully.
//
// Errors from the shadow notifier are logged but never propagated to the
// caller, ensuring that a misconfigured shadow target cannot disrupt
// normal alerting.
package shadow
