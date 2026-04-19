package triage_test

import (
	"testing"

	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/triage"
)

func ptr(l triage.Level) *triage.Level { return &l }

func defaultRules() []triage.Rule {
	return []triage.Rule{
		{MinPort: 1, MaxPort: 1023, Opened: ptr(triage.LevelCritical), Closed: ptr(triage.LevelWarning)},
		{MinPort: 1024, MaxPort: 49151, Opened: ptr(triage.LevelWarning), Closed: ptr(triage.LevelInfo)},
	}
}

func TestClassify_PrivilegedOpenedIsCritical(t *testing.T) {
	tr := triage.New(defaultRules())
	c := state.Change{Port: 22, Action: state.Opened}
	if got := tr.Classify(c); got != triage.LevelCritical {
		t.Fatalf("expected critical, got %s", got)
	}
}

func TestClassify_PrivilegedClosedIsWarning(t *testing.T) {
	tr := triage.New(defaultRules())
	c := state.Change{Port: 80, Action: state.Closed}
	if got := tr.Classify(c); got != triage.LevelWarning {
		t.Fatalf("expected warning, got %s", got)
	}
}

func TestClassify_RegisteredOpenedIsWarning(t *testing.T) {
	tr := triage.New(defaultRules())
	c := state.Change{Port: 8080, Action: state.Opened}
	if got := tr.Classify(c); got != triage.LevelWarning {
		t.Fatalf("expected warning, got %s", got)
	}
}

func TestClassify_NoMatchingRuleIsInfo(t *testing.T) {
	tr := triage.New(defaultRules())
	c := state.Change{Port: 60000, Action: state.Opened}
	if got := tr.Classify(c); got != triage.LevelInfo {
		t.Fatalf("expected info, got %s", got)
	}
}

func TestAnnotate_MapsAllChanges(t *testing.T) {
	tr := triage.New(defaultRules())
	changes := []state.Change{
		{Port: 22, Action: state.Opened},
		{Port: 8080, Action: state.Opened},
		{Port: 60000, Action: state.Opened},
	}
	annotated := tr.Annotate(changes)
	if len(annotated) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(annotated))
	}
	if annotated[changes[0]] != triage.LevelCritical {
		t.Errorf("port 22 should be critical")
	}
	if annotated[changes[1]] != triage.LevelWarning {
		t.Errorf("port 8080 should be warning")
	}
	if annotated[changes[2]] != triage.LevelInfo {
		t.Errorf("port 60000 should be info")
	}
}

func TestLevel_String(t *testing.T) {
	cases := map[triage.Level]string{
		triage.LevelInfo:     "info",
		triage.LevelWarning:  "warning",
		triage.LevelCritical: "critical",
	}
	for level, want := range cases {
		if got := level.String(); got != want {
			t.Errorf("Level(%d).String() = %q, want %q", level, got, want)
		}
	}
}
