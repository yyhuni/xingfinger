package pkg

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("reader.Close() error = %v", err)
	}
	return buf.String()
}

func TestPrintResultJSONOutput(t *testing.T) {
	scanner := &Scanner{jsonOutput: true}
	result := Result{URL: "https://example.com", CMS: "nginx", StatusCode: 200}

	output := captureStdout(t, func() {
		scanner.printResult(result)
	})

	if !strings.Contains(output, `"url":"https://example.com"`) {
		t.Fatalf("output = %q, want JSON url field", output)
	}
	if !strings.Contains(output, `"cms":"nginx"`) {
		t.Fatalf("output = %q, want JSON cms field", output)
	}
}

func TestPrintResultSilentOutputsOnlyHits(t *testing.T) {
	scanner := &Scanner{silent: true}
	hit := Result{URL: "https://example.com", CMS: "nginx"}
	miss := Result{URL: "https://example.com", CMS: ""}

	hitOutput := captureStdout(t, func() {
		scanner.printResult(hit)
	})
	missOutput := captureStdout(t, func() {
		scanner.printResult(miss)
	})

	if !strings.Contains(hitOutput, "https://example.com [nginx]") {
		t.Fatalf("hit output = %q, want silent hit line", hitOutput)
	}
	if missOutput != "" {
		t.Fatalf("miss output = %q, want empty", missOutput)
	}
}

func TestPrintResultNormalOutputWithoutCMSUsesPlainFormat(t *testing.T) {
	scanner := &Scanner{}
	result := Result{URL: "https://example.com", StatusCode: 204, Length: 0, Server: "nginx", Title: "No Content"}

	output := captureStdout(t, func() {
		scanner.printResult(result)
	})

	if !strings.Contains(output, "https://example.com [204] [0] [nginx] [No Content]") {
		t.Fatalf("output = %q, want plain formatted line", output)
	}
}
