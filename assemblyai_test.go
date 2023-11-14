package assemblyai

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setup() (*Client, *http.ServeMux, func()) {
	handler := http.NewServeMux()

	server := httptest.NewServer(handler)

	client := NewClientWithOptions(WithBaseURL(server.URL))

	return client, handler, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()

	if r.Method != want {
		t.Errorf("Request method = %v, want = %v", r.Method, want)
	}
}

func testQuery(t *testing.T, r *http.Request, want string) {
	t.Helper()

	if r.URL.RawQuery != want {
		t.Errorf("Request query = %v, want = %v", r.URL.RawQuery, want)
	}
}

func writeFileResponse(t *testing.T, w http.ResponseWriter, filename string) {
	t.Helper()

	b, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("ReadFile returned error: %v", err)
	}

	w.Write(b)
}
