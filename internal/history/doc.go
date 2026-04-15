// Package history provides persistent event logging for portwatch.
//
// It records port open/close events to a JSON file on disk with a
// configurable maximum size (oldest entries are dropped when the limit
// is exceeded). A Reporter is provided to format and display the log
// in a human-readable tabular form.
//
// Typical usage:
//
//	h, err := history.New("/var/lib/portwatch/history.json", 500)
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = h.Record(8080, "opened")
//
//	r := history.NewReporter(os.Stdout)
//	_ = r.Print(h.Entries())
package history
