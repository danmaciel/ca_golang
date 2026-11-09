package event

import "testing"

func TestListOrderHasDistinctEventName(t *testing.T) {
	if got := NewListOrder().GetName(); got != "ListOrder" {
		t.Fatalf("expected ListOrder event name, got %q", got)
	}
	if NewListOrder().GetName() == NewOrderCreated().GetName() {
		t.Fatal("list and create operations cannot share the same event name")
	}
}
