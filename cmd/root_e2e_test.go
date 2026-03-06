package cmd

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/spaolacci/murmur3"
	"github.com/yyhuni/xingfinger/pkg"
)

type eHoleFingerprint struct {
	CMS      string   `json:"cms"`
	Method   string   `json:"method"`
	Location string   `json:"location"`
	Keyword  []string `json:"keyword"`
}

type eHoleFingerprintFile struct {
	Fingerprint []eHoleFingerprint `json:"fingerprint"`
}

func resetRootCommandForTest() {
	targetURL = ""
	urlFile = ""
	thread = 50
	timeout = 10
	output = ""
	proxy = ""
	silent = false
	jsonOutput = false
	noDefault = false
	eholeFile = ""
	gobyFile = ""
	wappalyzerFile = ""
	fingersFile = ""
	fingerprintFile = ""
	arlFile = ""
	rootCmd.SetArgs(nil)
}

func TestExecuteHelperProcess(t *testing.T) {
	if os.Getenv("TEST_EXECUTE_HELPER") != "1" {
		return
	}

	resetRootCommandForTest()

	var args []string
	if err := json.Unmarshal([]byte(os.Getenv("TEST_EXECUTE_ARGS")), &args); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	rootCmd.SetArgs(args)
	Execute()
}

func newLocalTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	server := httptest.NewUnstartedServer(handler)
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	server.Listener = listener
	server.Start()
	t.Cleanup(server.Close)
	return server
}

func newLocalTLSTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	server := httptest.NewUnstartedServer(handler)
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error = %v", err)
	}
	server.Listener = listener
	server.StartTLS()
	t.Cleanup(server.Close)
	return server
}

func runCLICommand(t *testing.T, args []string) string {
	t.Helper()

	argsJSON, err := json.Marshal(args)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestExecuteHelperProcess$")
	cmd.Env = append(os.Environ(),
		"TEST_EXECUTE_HELPER=1",
		"TEST_EXECUTE_ARGS="+string(argsJSON),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CombinedOutput() error = %v, output = %s", err, output)
	}
	return string(output)
}

func writeEHoleFingerprintsFile(t *testing.T, fingerprints ...eHoleFingerprint) string {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "ehole.json")
	data, err := json.Marshal(eHoleFingerprintFile{Fingerprint: fingerprints})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if err := os.WriteFile(filename, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return filename
}

func writeARLFile(t *testing.T, content string) string {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "arl.yaml")
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return filename
}

func writeURLListFile(t *testing.T, urls ...string) string {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "urls.txt")
	content := strings.Join(urls, "\n") + "\n"
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return filename
}

func parseFirstJSONResult(t *testing.T, output string) pkg.Result {
	t.Helper()

	for _, candidate := range strings.Split(strings.TrimSpace(output), "\n") {
		candidate = strings.TrimSpace(candidate)
		if !strings.HasPrefix(candidate, "{") {
			continue
		}
		var result pkg.Result
		if err := json.Unmarshal([]byte(candidate), &result); err != nil {
			t.Fatalf("Unmarshal() error = %v, output = %s", err, output)
		}
		return result
	}

	t.Fatalf("expected JSON output from CLI scan, output = %s", output)
	return pkg.Result{}
}

func readResultsFile(t *testing.T, filename string) []pkg.Result {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var results []pkg.Result
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatalf("Unmarshal() error = %v, data = %s", err, data)
	}
	return results
}

func calcFaviconHashForTest(data []byte) string {
	b64 := base64.StdEncoding.EncodeToString(data)
	var buf bytes.Buffer
	for i := 0; i < len(b64); i += 76 {
		end := i + 76
		if end > len(b64) {
			end = len(b64)
		}
		buf.WriteString(b64[i:end])
		buf.WriteString("\n")
	}
	hash := murmur3.Sum32(buf.Bytes())
	return strconv.FormatInt(int64(int32(hash)), 10)
}

func findResultByURL(results []pkg.Result, url string) (pkg.Result, bool) {
	for _, result := range results {
		if result.URL == url {
			return result, true
		}
	}
	return pkg.Result{}, false
}

func TestExecuteScansLocalServerE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><head><title>Local Test Site</title></head><body>hello-keyword</body></html>`)
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "KeywordCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"hello-keyword"},
	})

	output := runCLICommand(t, []string{
		"-u", server.URL,
		"--ehole", eholePath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if result.URL != server.URL {
		t.Fatalf("result.URL = %q, want %q", result.URL, server.URL)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("result.StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
	}
	if result.Title != "Local Test Site" {
		t.Fatalf("result.Title = %q, want %q", result.Title, "Local Test Site")
	}
	if result.Length == 0 {
		t.Fatalf("result.Length = %d, want > 0", result.Length)
	}
	if !strings.Contains(strings.ToLower(result.CMS), "keywordcms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "KeywordCMS")
	}
}

func TestExecuteReadsURLListAndWritesJSONFileE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch r.URL.Path {
		case "/header":
			w.Header().Set("X-Powered-By", "HeaderValue")
			fmt.Fprint(w, `<html><head><title>Ignore Me</title></head><body>header-body</body></html>`)
		case "/title":
			fmt.Fprint(w, `<html><head><title>TitleValue</title></head><body>title-body</body></html>`)
		default:
			http.NotFound(w, r)
		}
	})

	eholePath := writeEHoleFingerprintsFile(t,
		eHoleFingerprint{CMS: "HeaderCMS", Method: "keyword", Location: "header", Keyword: []string{"HeaderValue"}},
		eHoleFingerprint{CMS: "TitleCMS", Method: "keyword", Location: "title", Keyword: []string{"TitleValue"}},
	)
	listPath := writeURLListFile(t, server.URL+"/header", server.URL+"/title")
	outputFile := filepath.Join(t.TempDir(), "result.json")

	_ = runCLICommand(t, []string{
		"-l", listPath,
		"-o", outputFile,
		"--ehole", eholePath,
		"--no-default",
		"-t", "1",
		"--timeout", "2",
	})

	results := readResultsFile(t, outputFile)
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}

	headerResult, ok := findResultByURL(results, server.URL+"/header")
	if !ok {
		t.Fatalf("missing result for %s/header", server.URL)
	}
	titleResult, ok := findResultByURL(results, server.URL+"/title")
	if !ok {
		t.Fatalf("missing result for %s/title", server.URL)
	}

	if !strings.Contains(strings.ToLower(headerResult.CMS), "headercms") {
		t.Fatalf("headerResult.CMS = %q, want contains %q", headerResult.CMS, "HeaderCMS")
	}
	if !strings.Contains(strings.ToLower(titleResult.CMS), "titlecms") {
		t.Fatalf("titleResult.CMS = %q, want contains %q", titleResult.CMS, "TitleCMS")
	}
	if titleResult.Title != "TitleValue" {
		t.Fatalf("titleResult.Title = %q, want %q", titleResult.Title, "TitleValue")
	}
}

func TestExecuteSilentModeOnlyPrintsHitsE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch r.URL.Path {
		case "/hit":
			fmt.Fprint(w, `<html><head><title>Hit</title></head><body>silent-hit-keyword</body></html>`)
		case "/miss":
			fmt.Fprint(w, `<html><head><title>Miss</title></head><body>no-match</body></html>`)
		default:
			http.NotFound(w, r)
		}
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "SilentCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"silent-hit-keyword"},
	})
	listPath := writeURLListFile(t, server.URL+"/hit", server.URL+"/miss")

	output := runCLICommand(t, []string{
		"-l", listPath,
		"--ehole", eholePath,
		"--no-default",
		"-s",
		"-t", "1",
		"--timeout", "2",
	})

	lowerOutput := strings.ToLower(output)
	if !strings.Contains(lowerOutput, strings.ToLower(server.URL+"/hit [silentcms]")) {
		t.Fatalf("output = %q, want hit line", output)
	}
	if strings.Contains(lowerOutput, strings.ToLower(server.URL+"/miss")) {
		t.Fatalf("output = %q, want miss URL to be silent", output)
	}
}

func TestExecuteDetectsFaviconE2E(t *testing.T) {
	iconData := []byte("e2e-favicon")
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, `<html><head><title>Favicon Site</title><link rel="icon" href="/favicon.ico"></head><body>favicon</body></html>`)
		case "/favicon.ico":
			w.Write(iconData)
		default:
			http.NotFound(w, r)
		}
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "FaviconCMS",
		Method:   "faviconhash",
		Location: "body",
		Keyword:  []string{calcFaviconHashForTest(iconData)},
	})

	output := runCLICommand(t, []string{
		"-u", server.URL,
		"--ehole", eholePath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if !strings.Contains(strings.ToLower(result.CMS), "faviconcms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "FaviconCMS")
	}
	if result.Title != "Favicon Site" {
		t.Fatalf("result.Title = %q, want %q", result.Title, "Favicon Site")
	}
}

func TestExecuteDetectsARLE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Test-ARL", "arl-header")
		fmt.Fprint(w, `<html><head><title>ARL Site</title></head><body>arl-body</body></html>`)
	})

	arlPath := writeARLFile(t, "- name: ARLCMS_header\n  rule: header=\"arl-header\" && body=\"arl-body\"\n")

	output := runCLICommand(t, []string{
		"-u", server.URL,
		"--arl", arlPath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if !strings.Contains(strings.ToLower(result.CMS), "arlcms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "ARLCMS")
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("result.StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
	}
}

func TestExecuteScansLocalTLSServerE2E(t *testing.T) {
	server := newLocalTLSTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><head><title>TLS Site</title></head><body>tls-keyword</body></html>`)
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "TLSCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"tls-keyword"},
	})

	output := runCLICommand(t, []string{
		"-u", server.URL,
		"--ehole", eholePath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if result.URL != server.URL {
		t.Fatalf("result.URL = %q, want %q", result.URL, server.URL)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("result.StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
	}
	if result.Title != "TLS Site" {
		t.Fatalf("result.Title = %q, want %q", result.Title, "TLS Site")
	}
	if !strings.Contains(strings.ToLower(result.CMS), "tlscms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "TLSCMS")
	}
}

func TestExecuteFollowsHTTPRedirectE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/landing", http.StatusFound)
		case "/landing":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, `<html><head><title>Redirect Site</title></head><body>redirect-keyword</body></html>`)
		default:
			http.NotFound(w, r)
		}
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "RedirectCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"redirect-keyword"},
	})

	output := runCLICommand(t, []string{
		"-u", server.URL + "/start",
		"--ehole", eholePath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if result.URL != server.URL+"/landing" {
		t.Fatalf("result.URL = %q, want %q", result.URL, server.URL+"/landing")
	}
	if result.Title != "Redirect Site" {
		t.Fatalf("result.Title = %q, want %q", result.Title, "Redirect Site")
	}
	if !strings.Contains(strings.ToLower(result.CMS), "redirectcms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "RedirectCMS")
	}
}

func TestExecuteFollowsJSRedirectE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch r.URL.Path {
		case "/start":
			fmt.Fprint(w, `<html><head><title>JS Start</title><script>window.location.href='/landing'</script></head><body>start</body></html>`)
		case "/landing":
			fmt.Fprint(w, `<html><head><title>JS Landing</title></head><body>js-keyword</body></html>`)
		default:
			http.NotFound(w, r)
		}
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "JSCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"js-keyword"},
	})
	outputFile := filepath.Join(t.TempDir(), "js-results.json")

	_ = runCLICommand(t, []string{
		"-u", server.URL + "/start",
		"-o", outputFile,
		"--ehole", eholePath,
		"--no-default",
		"-t", "1",
		"--timeout", "2",
	})

	results := readResultsFile(t, outputFile)
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}

	startResult, ok := findResultByURL(results, server.URL+"/start")
	if !ok {
		t.Fatalf("missing result for %s/start", server.URL)
	}
	landingResult, ok := findResultByURL(results, server.URL+"/landing")
	if !ok {
		t.Fatalf("missing result for %s/landing", server.URL)
	}

	if startResult.Title != "JS Start" {
		t.Fatalf("startResult.Title = %q, want %q", startResult.Title, "JS Start")
	}
	if landingResult.Title != "JS Landing" {
		t.Fatalf("landingResult.Title = %q, want %q", landingResult.Title, "JS Landing")
	}
	if !strings.Contains(strings.ToLower(landingResult.CMS), "jscms") {
		t.Fatalf("landingResult.CMS = %q, want contains %q", landingResult.CMS, "JSCMS")
	}
}

func TestExecuteHandlesDeflateResponseE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		writer := zlib.NewWriter(&buf)
		_, _ = writer.Write([]byte(`<html><head><title>Deflate Site</title></head><body>deflate-keyword</body></html>`))
		_ = writer.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Encoding", "deflate")
		_, _ = w.Write(buf.Bytes())
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "DeflateCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"deflate-keyword"},
	})

	output := runCLICommand(t, []string{
		"-u", server.URL,
		"--ehole", eholePath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if result.Title != "Deflate Site" {
		t.Fatalf("result.Title = %q, want %q", result.Title, "Deflate Site")
	}
	if !strings.Contains(strings.ToLower(result.CMS), "deflatecms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "DeflateCMS")
	}
}

func TestExecuteHandlesGzipResponseE2E(t *testing.T) {
	server := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, _ = writer.Write([]byte(`<html><head><title>Gzip Site</title></head><body>gzip-keyword</body></html>`))
		_ = writer.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(buf.Bytes())
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "GzipCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"gzip-keyword"},
	})

	output := runCLICommand(t, []string{
		"-u", server.URL,
		"--ehole", eholePath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if result.Title != "Gzip Site" {
		t.Fatalf("result.Title = %q, want %q", result.Title, "Gzip Site")
	}
	if !strings.Contains(strings.ToLower(result.CMS), "gzipcms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "GzipCMS")
	}
}

func TestExecuteUsesProxyE2E(t *testing.T) {
	backend := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html><head><title>Proxy Site</title></head><body>proxy-keyword</body></html>`)
	})

	var proxyHits int32
	proxy := newLocalTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&proxyHits, 1)

		targetURL := r.RequestURI
		if parsed, err := url.Parse(targetURL); err != nil || parsed.Scheme == "" {
			targetURL = backend.URL + r.URL.RequestURI()
		}

		req, err := http.NewRequest(r.Method, targetURL, nil)
		if err != nil {
			t.Fatalf("NewRequest() error = %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("proxy upstream error = %v", err)
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	})

	eholePath := writeEHoleFingerprintsFile(t, eHoleFingerprint{
		CMS:      "ProxyCMS",
		Method:   "keyword",
		Location: "body",
		Keyword:  []string{"proxy-keyword"},
	})

	output := runCLICommand(t, []string{
		"-u", backend.URL,
		"-p", proxy.URL,
		"--ehole", eholePath,
		"--no-default",
		"-j",
		"-t", "1",
		"--timeout", "2",
	})

	result := parseFirstJSONResult(t, output)
	if !strings.Contains(strings.ToLower(result.CMS), "proxycms") {
		t.Fatalf("result.CMS = %q, want contains %q", result.CMS, "ProxyCMS")
	}
	if result.Title != "Proxy Site" {
		t.Fatalf("result.Title = %q, want %q", result.Title, "Proxy Site")
	}
	if atomic.LoadInt32(&proxyHits) == 0 {
		t.Fatal("expected request to pass through proxy")
	}
}
