package pkg

import "testing"

func TestQueuePushPopMaintainsFIFOAndLen(t *testing.T) {
	queue := NewQueue()

	if queue.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", queue.Len())
	}

	queue.Push("first")
	queue.Push("second")
	queue.Push("third")

	if queue.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", queue.Len())
	}

	if got := queue.Pop(); got != "first" {
		t.Fatalf("first Pop() = %v, want %q", got, "first")
	}
	if got := queue.Pop(); got != "second" {
		t.Fatalf("second Pop() = %v, want %q", got, "second")
	}
	if got := queue.Pop(); got != "third" {
		t.Fatalf("third Pop() = %v, want %q", got, "third")
	}
	if got := queue.Pop(); got != nil {
		t.Fatalf("empty Pop() = %v, want nil", got)
	}
}
