package event

import "time"

type ListOrder struct {
	Name    string
	Payload interface{}
}

func NewListOrder() *ListOrder {
	return &ListOrder{
		Name: "OrderCreated",
	}
}

func (e *ListOrder) GetName() string {
	return e.Name
}

func (e *ListOrder) GetPayload() interface{} {
	return e.Payload
}

func (e *ListOrder) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *ListOrder) GetDateTime() time.Time {
	return time.Now()
}
