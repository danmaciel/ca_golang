package events

import "time"

// PayloadEvent is a per-dispatch event value. Creating a separate instance
// prevents unrelated requests from overwriting each other's payload.
type PayloadEvent struct {
	name    string
	payload interface{}
	created time.Time
}

func NewPayloadEvent(name string, payload interface{}) *PayloadEvent {
	return &PayloadEvent{name: name, payload: payload, created: time.Now()}
}

func (e *PayloadEvent) GetName() string { return e.name }
func (e *PayloadEvent) GetPayload() interface{} { return e.payload }
func (e *PayloadEvent) GetDateTime() time.Time { return e.created }

// Retained for compatibility with EventInterface. Dispatch producers should
// construct an event per request and not modify it after dispatch.
func (e *PayloadEvent) SetPayload(payload interface{}) { e.payload = payload }
