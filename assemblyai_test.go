package assemblyai

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func setup() (*Client, *http.ServeMux, func()) {
	handler := http.NewServeMux()

	server := httptest.NewServer(handler)

	client := NewClientWithOptions(WithBaseURL(server.URL))

	return client, handler, server.Close
}

func writeFileResponse(t *testing.T, w http.ResponseWriter, filename string) {
	t.Helper()

	if filepath.Ext(filename) == ".json" {
		w.Header().Set("Content-Type", "application/json")
	}

	b, err := os.ReadFile(filename)
	require.NoError(t, err)

	_, err = w.Write(b)
	require.NoError(t, err)
}
