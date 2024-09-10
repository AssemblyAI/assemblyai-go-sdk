package assemblyai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpload(t *testing.T) {
	t.Parallel()

	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/upload", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)

		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		require.Equal(t, "data", string(b))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"upload_url\": %q}", fakeAudioURL)
	})

	ctx := context.Background()

	buf := bytes.NewBufferString("data")

	got, err := client.Upload(ctx, buf)
	require.NoError(t, err)

	require.Equal(t, fakeAudioURL, got)
}
