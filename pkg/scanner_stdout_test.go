package pkg

import (
	"errors"
	"os"
	"testing"
)

func TestWithSilentStdoutRestoresStdoutAndClosesWriter(t *testing.T) {
	oldOpenSilentWriter := openSilentWriter
	defer func() { openSilentWriter = oldOpenSilentWriter }()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	defer reader.Close()

	openSilentWriter = func() (*os.File, error) {
		return writer, nil
	}

	oldStdout := os.Stdout
	err = withSilentStdout(func() error {
		if os.Stdout != writer {
			t.Fatalf("os.Stdout = %v, want temporary writer", os.Stdout)
		}
		_, writeErr := os.Stdout.Write([]byte("hidden output"))
		return writeErr
	})
	if err != nil {
		t.Fatalf("withSilentStdout() error = %v", err)
	}

	if os.Stdout != oldStdout {
		t.Fatalf("os.Stdout not restored")
	}

	if _, err := writer.Write([]byte("x")); err == nil {
		t.Fatalf("temporary writer should be closed after withSilentStdout")
	}
}

func TestWithSilentStdoutReturnsOpenErrorWithoutChangingStdout(t *testing.T) {
	oldOpenSilentWriter := openSilentWriter
	defer func() { openSilentWriter = oldOpenSilentWriter }()

	wantErr := errors.New("open devnull failed")
	openSilentWriter = func() (*os.File, error) {
		return nil, wantErr
	}

	oldStdout := os.Stdout
	err := withSilentStdout(func() error {
		t.Fatal("callback should not run when openSilentWriter fails")
		return nil
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("withSilentStdout() error = %v, want %v", err, wantErr)
	}
	if os.Stdout != oldStdout {
		t.Fatalf("os.Stdout changed on open error")
	}
}
