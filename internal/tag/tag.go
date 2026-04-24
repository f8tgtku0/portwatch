// Package tag provides lightweight key-value tagging for port change events,
// allowing operators to attach metadata (e.g. environment, owner) to alerts.
package tag

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/portwatch/internal/state"
)

// Tag is a key-value pair attached to a port change.
type Tag struct {
	Key   string
	Value string
}

// String returns the tag in "key=value" form.
func (t Tag) String() string {
	return fmt.Sprintf("%s=%s", t.Key, t.Value)
}

// Tagger holds a fixed set of tags and annotates change sets with them.
type Tagger struct {
	tags []Tag
}

// New creates a Tagger from a map of key-value pairs.
// Keys and values are trimmed of surrounding whitespace.
func New(pairs map[string]string) *Tagger {
	tags := make([]Tag, 0, len(pairs))
	for k, v := range pairs {
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if k == "" {
			continue
		}
		tags = append(tags, Tag{Key: k, Value: v})
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].Key < tags[j].Key })
	return &Tagger{tags: tags}
}

// Tags returns a copy of the configured tags.
func (t *Tagger) Tags() []Tag {
	out := make([]Tag, len(t.tags))
	copy(out, t.tags)
	return out
}

// Annotate returns a new slice of changes where each entry's Labels map
// is augmented with the tagger's key-value pairs.
// Original change values are not mutated.
func (t *Tagger) Annotate(changes []state.Change) []state.Change {
	if len(changes) == 0 || len(t.tags) == 0 {
		return changes
	}
	out := make([]state.Change, len(changes))
	for i, c := range changes {
		merged := make(map[string]string, len(c.Labels)+len(t.tags))
		for k, v := range c.Labels {
			merged[k] = v
		}
		for _, tag := range t.tags {
			merged[tag.Key] = tag.Value
		}
		c.Labels = merged
		out[i] = c
	}
	return out
}
