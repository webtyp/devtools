//go:build dev && wasm

package devtools

import (
	"runtime"
	"syscall/js"
)

var (
	memStatsRegistered bool
	memStatsBuf        runtime.MemStats
)

// registerMemStats exposes window.__webtypMemStats() in dev WASM builds so
// webtyp's MCP (browser_get_performance) can measure Go allocations per action.
// Both memStatsBuf and the returned JS object are reused to avoid allocating.
func registerMemStats() {
	if memStatsRegistered {
		return
	}
	memStatsRegistered = true
	out := js.Global().Get("Object").New()
	fn := js.FuncOf(func(js.Value, []js.Value) any {
		runtime.ReadMemStats(&memStatsBuf)
		out.Set("totalAlloc", float64(memStatsBuf.TotalAlloc))
		out.Set("mallocs", float64(memStatsBuf.Mallocs))
		out.Set("heapAlloc", float64(memStatsBuf.HeapAlloc))
		out.Set("sys", float64(memStatsBuf.Sys))
		out.Set("numGC", float64(memStatsBuf.NumGC))
		return out
	})
	js.Global().Set("__webtypMemStats", fn)
}
