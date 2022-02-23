package debug

import (
	"os"
	"runtime"
	"runtime/pprof"
)

func Profile(f func()) {
	c, _ := os.Create("cpuprofile.prof")
	defer c.Close() // error handling omitted for example
	pprof.StartCPUProfile(c)
	defer pprof.StopCPUProfile()

	f()

	m, _ := os.Create("memprofile.prof")
	defer m.Close()
	runtime.GC()
	pprof.WriteHeapProfile(m)
}
