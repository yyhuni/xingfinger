package pkg

import (
	"net/http"
	"strings"
	"testing"
)

func TestExtractTitleStripsWhitespace(t *testing.T) {
	body := "<html><head><title>\n\t Example Title \r\n</title></head></html>"
	got := extractTitle(body)
	want := "Example Title"

	if got != want {
		t.Fatalf("extractTitle() = %q, want %q", got, want)
	}
}

func TestBuildRawResponseIncludesStatusHeadersAndBody(t *testing.T) {
	resp := &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
		StatusCode: http.StatusCreated,
		Header: http.Header{
			"Content-Type": []string{"text/plain"},
			"X-Test":       []string{"yes"},
		},
	}
	body := []byte("hello")

	got := string(buildRawResponse(resp, body))

	if !strings.Contains(got, "HTTP/1.1 201 Created\r\n") {
		t.Fatalf("raw response missing status line: %q", got)
	}
	if !strings.Contains(got, "Content-Type: text/plain\r\n") {
		t.Fatalf("raw response missing Content-Type header: %q", got)
	}
	if !strings.Contains(got, "X-Test: yes\r\n") {
		t.Fatalf("raw response missing X-Test header: %q", got)
	}
	if !strings.HasSuffix(got, "\r\nhello") {
		t.Fatalf("raw response missing body suffix: %q", got)
	}
}
