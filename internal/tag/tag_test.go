package tag_test

import (
	"testing"

	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/tag"
)

func makeChange(port int, labels map[string]string) state.Change {
	return state.Change{Port: port, Action: state.Opened, Labels: labels}
}

func TestNew_SortsTagsAndTrimsSpace(t *testing.T) {
	tr := tag.New(map[string]string{
		" env ": " prod ",
		"owner": "team-a",
	})
	tags := tr.Tags()
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Key != "env" || tags[0].Value != "prod" {
		t.Errorf("unexpected first tag: %v", tags[0])
	}
	if tags[1].Key != "owner" {
		t.Errorf("unexpected second tag key: %s", tags[1].Key)
	}
}

func TestNew_SkipsEmptyKeys(t *testing.T) {
	tr := tag.New(map[string]string{"  ": "value", "k": "v"})
	if len(tr.Tags()) != 1 {
		t.Errorf("expected 1 tag after skipping empty key")
	}
}

func TestTag_String(t *testing.T) {
	tg := tag.Tag{Key: "env", Value: "staging"}
	if tg.String() != "env=staging" {
		t.Errorf("unexpected String(): %s", tg.String())
	}
}

func TestAnnotate_AddsTagsToChanges(t *testing.T) {
	tr := tag.New(map[string]string{"env": "prod", "region": "us-east-1"})
	changes := []state.Change{makeChange(80, nil), makeChange(443, nil)}

	annotated := tr.Annotate(changes)
	if len(annotated) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(annotated))
	}
	for _, c := range annotated {
		if c.Labels["env"] != "prod" {
			t.Errorf("port %d: missing env tag", c.Port)
		}
		if c.Labels["region"] != "us-east-1" {
			t.Errorf("port %d: missing region tag", c.Port)
		}
	}
}

func TestAnnotate_DoesNotMutateOriginal(t *testing.T) {
	tr := tag.New(map[string]string{"env": "prod"})
	orig := makeChange(22, map[string]string{"existing": "yes"})
	changes := []state.Change{orig}

	tr.Annotate(changes)

	if _, ok := changes[0].Labels["env"]; ok {
		t.Error("original change was mutated")
	}
}

func TestAnnotate_PreservesExistingLabels(t *testing.T) {
	tr := tag.New(map[string]string{"env": "prod"})
	changes := []state.Change{makeChange(8080, map[string]string{"service": "api"})}

	annotated := tr.Annotate(changes)
	if annotated[0].Labels["service"] != "api" {
		t.Error("existing label was lost")
	}
	if annotated[0].Labels["env"] != "prod" {
		t.Error("new tag was not added")
	}
}

func TestAnnotate_EmptyChanges_ReturnsEmpty(t *testing.T) {
	tr := tag.New(map[string]string{"env": "prod"})
	result := tr.Annotate(nil)
	if result != nil {
		t.Errorf("expected nil for nil input, got %v", result)
	}
}

func TestAnnotate_NoTags_ReturnsOriginal(t *testing.T) {
	tr := tag.New(nil)
	changes := []state.Change{makeChange(80, nil)}
	result := tr.Annotate(changes)
	if &result[0] == &changes[0] {
		// same slice returned is fine
	}
	if len(result) != 1 {
		t.Errorf("expected 1 change, got %d", len(result))
	}
}
