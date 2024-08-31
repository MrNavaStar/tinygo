//go:build wasip1 && tinygo_wasm_reactor

package runtime

import (
	"internal/task"
	"unsafe"
)

var packagesInitialized = false

//export _initialize
func _initialize() {
	// This function is called before any //go:wasmexport functions are called
	// to initialize everything. It must not block.

	// Initialize the heap.
	heapStart = uintptr(unsafe.Pointer(&heapStartSymbol))
	heapEnd = uintptr(wasm_memory_size(0) * wasmPageSize)
	initHeap()

	if hasScheduler {
		// A package initializer might do funky stuff like start a goroutine and
		// wait until it completes, so we have to run package initializers in a
		// goroutine.
		go func() {
			initAll()
			packagesInitialized = true
		}()
		scheduler(true)
		if !packagesInitialized {
			// Unlikely, but if package initializers do something blocking (like
			// time.Sleep()), that's a bug.
			runtimePanic("package initializer blocks")
		}
	} else {
		// There are no goroutines (except for the main one, if you can call it
		// that), so we can just run all the package initializers.
		initAll()
	}
}

func sleepTicks(d timeUnit) {
	// See the proposal:
	// > When the goroutine running the exported function blocks for any reason,
	// > the function will yield to the Go runtime. The Go runtime will schedule
	// > other goroutines as necessary. If there are no other goroutines, the
	// > application will crash with a deadlock, as there is no way to proceed,
	// > and Wasm code cannot block.
	// We can only get here when there are no runnable goroutines. In other
	// words, when the exported function blocks and there are no other
	// goroutines to run. So we crash with a deadlock.
	runtimePanic("all goroutines are asleep - deadlock!")
}

// Called from within a //go:wasmexport wrapper (the one that's exported from
// the wasm module) after the goroutine has been queued. Just run the scheduler,
// and check that the goroutine finished when the scheduler is idle (as required
// by the //go:wasmexport proposal).
//
// This function is not called when the scheduler is disabled.
func wasmExportRun(done *bool) {
	scheduler(true)
	if !*done {
		runtimePanic("//go:wasmexport function did not finish")
	}
}

// Called from the goroutine wrapper for the //go:wasmexport function. It just
// signals to the runtime that the //go:wasmexport call has finished, and can
// switch back to the wasmExportRun function.
//
// This function is not called when the scheduler is disabled.
func wasmExportExit() {
	task.Pause()

	// TODO: we could cache the allocated stack so we don't have to keep
	// allocating a new stack on every //go:wasmexport call.
}
