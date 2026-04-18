package ratelimit

import (
	"fmt"

	"github.com/user/portwatch/internal/state"
)

// ChangeKey returns a stable string key for a port change event,
// suitable for use as a rate limit key.
func ChangeKey(c state.Change) string {
	action := "opened"
	if !c.Opened {
		action = "closed"
	}
	return fmt.Sprintf("port:%d:%s", c.Port, action)
}
