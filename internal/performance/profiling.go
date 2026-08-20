package performance

import (
	"fmt"
	"os"
	"runtime"
)

func ReportHeapSize(label string) {
	runtime.GC() // Force gc so that we don't report garbage that hasn't been collected yet

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Fprintf(os.Stderr, "[heap:%s] Alloc=%.2fMB Sys=%.2fMB NumGC=%d\n",
		label,
		float64(m.HeapAlloc)/1024/1024,
		float64(m.Sys)/1024/1024,
		m.NumGC,
	)
}
