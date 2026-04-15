package history

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// Reporter formats and writes history entries to an output sink.
type Reporter struct {
	out io.Writer
}

// NewReporter creates a Reporter writing to the given writer.
func NewReporter(out io.Writer) *Reporter {
	return &Reporter{out: out}
}

// Print writes a formatted table of history entries to the reporter's output.
func (r *Reporter) Print(entries []Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(r.out, "No history recorded.")
		return err
	}

	w := tabwriter.NewWriter(r.out, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tPORT\tEVENT")
	fmt.Fprintln(w, "---------\t----\t-----")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%d\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Port,
			e.Event,
		)
	}
	return w.Flush()
}

// Summary returns a brief string summarising the entries.
func Summary(entries []Entry) string {
	opened, closed := 0, 0
	for _, e := range entries {
		switch e.Event {
		case "opened":
			opened++
		case "closed":
			closed++
		}
	}
	return fmt.Sprintf("%d event(s): %d opened, %d closed", len(entries), opened, closed)
}
