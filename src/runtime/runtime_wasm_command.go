//go:build wasip1 && !tinygo_wasm_reactor

package runtime

import "unsafe"

//export _start
func _start() {
	// These need to be initialized early so that the heap can be initialized.
	heapStart = uintptr(unsafe.Pointer(&heapStartSymbol))
	heapEnd = uintptr(wasm_memory_size(0) * wasmPageSize)
	run()
	__stdio_exit()
}

func sleepTicks(d timeUnit) {
	wasiSleepTicks(d)
}

// TODO: we should define wasmExportRun and wasmExportExit here so that
// //go:wasmexport also works without -buildmode=c-shared (for when a
// //go:wasmexport call). Calling //go:wasmexport functions is only allowed
// before main.main returns.
