package shadow

import (
	"fmt"
	"io"
	"os"
	"time"
)

func defaultWriter() io.Writer {
	return os.Stderr
}

func logShadowError(w io.Writer, err error) {
	fmt.Fprintf(w, "[shadow] %s error: %v\n", time.Now().Format(time.RFC3339), err)
}
