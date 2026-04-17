package metrics

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Reporter prints metric snapshots in a human-readable table.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter writing to w. If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Print writes a formatted metrics snapshot to the reporter's writer.
func (r *Reporter) Print(s Snapshot) {
	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "METRIC\tVALUE")
	fmt.Fprintln(tw, "------\t-----")
	fmt.Fprintf(tw, "Uptime\t%s\n", s.Uptime)
	fmt.Fprintf(tw, "Scans Total\t%d\n", s.ScansTotal)
	fmt.Fprintf(tw, "Changes Total\t%d\n", s.ChangesTotal)
	fmt.Fprintf(tw, "Alerts Total\t%d\n", s.AlertsTotal)
	fmt.Fprintf(tw, "Open Ports (last scan)\t%d\n", s.OpenPorts)
	if !s.LastScanAt.IsZero() {
		fmt.Fprintf(tw, "Last Scan At\t%s\n", s.LastScanAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Fprintf(tw, "Last Scan At\t%s\n", "never")
	}
	tw.Flush()
}
