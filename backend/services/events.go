package services

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	progressEventSuffix = ":progress"
	completeEventSuffix = ":complete"
)

// EventSink delivers a named event with an optional payload. It mirrors the
// shape of runtime.EventsEmit with the Wails context already bound, which keeps
// OperationEmitter decoupled from the Wails runtime and testable in isolation.
type EventSink func(eventName string, optionalData ...interface{})

// OperationEmitter enforces the shell-out event convention (Req 12.8): a single
// shell-out operation emits zero or more "<area>:progress" events followed by
// exactly one terminal "<area>:complete" event. It guarantees the terminal
// guarantee by deduplicating completion and suppressing any progress emitted
// after completion, so the complete event is always the final event.
type OperationEmitter struct {
	area      string
	sink      EventSink
	completed bool
}

// NewOperationEmitter builds an emitter for the given area backed by an
// arbitrary sink. A nil sink yields a no-op emitter that still tracks state.
func NewOperationEmitter(area string, sink EventSink) *OperationEmitter {
	return &OperationEmitter{area: area, sink: sink}
}

// NewRuntimeOperationEmitter builds an emitter that delivers events through the
// Wails runtime bound to ctx. A nil ctx yields a no-op sink so callers without
// an active frontend context emit nothing while preserving the convention.
func NewRuntimeOperationEmitter(ctx context.Context, area string) *OperationEmitter {
	if ctx == nil {
		return NewOperationEmitter(area, nil)
	}
	return NewOperationEmitter(area, func(eventName string, optionalData ...interface{}) {
		runtime.EventsEmit(ctx, eventName, optionalData...)
	})
}

// Progress emits an "<area>:progress" event. It is a no-op once Complete has
// been called, ensuring no progress event can follow the terminal event.
func (e *OperationEmitter) Progress(optionalData ...interface{}) {
	if e == nil || e.completed || e.sink == nil {
		return
	}
	e.sink(e.area+progressEventSuffix, optionalData...)
}

// Complete emits the single terminal "<area>:complete" event. Repeated calls
// are ignored so exactly one completion event is emitted per operation.
func (e *OperationEmitter) Complete(optionalData ...interface{}) {
	if e == nil || e.completed {
		return
	}
	e.completed = true
	if e.sink == nil {
		return
	}
	e.sink(e.area+completeEventSuffix, optionalData...)
}

// Completed reports whether the terminal event has already been emitted.
func (e *OperationEmitter) Completed() bool {
	return e != nil && e.completed
}
