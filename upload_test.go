package assemblyai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestUpload(t *testing.T) {
	client, handler, teardown := setup()
	defer teardown()

	handler.HandleFunc("/v2/upload", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		want := "data"

		if got := string(b); got != want {
			t.Errorf("Request body = %+v, want = %+v", got, want)
		}

		fmt.Fprintf(w, "{\"upload_url\": %q}", fakeAudioURL)
	})

	ctx := context.Background()

	buf := bytes.NewBufferString("data")

	got, err := client.Upload(ctx, buf)
	if err != nil {
		t.Errorf("Upload returned error: %v", err)
	}

	want := fakeAudioURL

	if got != want {
		t.Errorf("Upload URL = %v, want %v", got, want)
	}
}
