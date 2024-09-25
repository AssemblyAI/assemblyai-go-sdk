package assemblyai

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const testTimeout = 5 * time.Second

func TestRealTime_Handler(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		_, got, err := conn.Read(ctx)
		require.NoError(t, err)

		require.Equal(t, []byte("foo"), got)

		err = wsjson.Write(ctx, conn, PartialTranscript{
			MessageType: MessageTypePartialTranscript,
			RealTimeBaseTranscript: RealTimeBaseTranscript{
				Text: "foo",
			},
		})
		require.NoError(t, err)

		err = wsjson.Write(ctx, conn, FinalTranscript{
			MessageType: MessageTypeFinalTranscript,
			RealTimeBaseTranscript: RealTimeBaseTranscript{
				Text: "Foo.",
			},
		})
		require.NoError(t, err)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	var partialTranscriptReceived, finalTranscriptReceived bool

	wg := sync.WaitGroup{}
	wg.Add(2)

	client := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{
			OnPartialTranscript: func(event PartialTranscript) {
				partialTranscriptReceived = true
				require.Equal(t, "foo", event.Text)
				wg.Done()
			},
			OnFinalTranscript: func(event FinalTranscript) {
				finalTranscriptReceived = true
				require.Equal(t, "Foo.", event.Text)
				wg.Done()
			},
			OnError: func(err error) {
				require.NoError(t, err)
			},
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Send(ctx, []byte("foo"))
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	wg.Wait()

	require.True(t, partialTranscriptReceived)
	require.True(t, finalTranscriptReceived)
}

func TestRealTime_Connect(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		apiKey := r.Header.Get("Authorization")

		require.Equal(t, "api-key", apiKey)

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		_, got, _ := conn.Read(ctx)

		require.Equal(t, []byte("foo"), got)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client := NewRealTimeClientWithOptions(
		WithRealTimeAPIKey("api-key"),
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Send(ctx, []byte("foo"))
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)
}

func TestRealTime_ConnectFailsDisconnect(t *testing.T) {
	t.Parallel()

	// setup webhook server that fails to make a connection
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		require.Equal(t, "api-key", apiKey)
	}))
	defer ts.Close()

	client := NewRealTimeClientWithOptions(
		WithRealTimeAPIKey("api-key"),
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error
	// try to connect, note error is returned, but ignore it and proceed to Disconnect
	err = client.Connect(ctx)
	err = client.Disconnect(ctx, true)
	require.Errorf(t, err, "client connection does not exist")
}

func TestRealTime_Send(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		wordBoost := r.URL.Query().Get("word_boost")
		require.Equal(t, `["foo","bar"]`, wordBoost)

		encoding := r.URL.Query().Get("encoding")
		require.Equal(t, "pcm_mulaw", encoding)

		sampleRate := r.URL.Query().Get("sample_rate")
		require.Equal(t, "8000", sampleRate)

		disablePartialTranscripts := r.URL.Query().Get("disable_partial_transcripts")
		require.Equal(t, "true", disablePartialTranscripts)

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		_, got, _ := conn.Read(ctx)

		require.Equal(t, []byte("foo"), got)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{}),
		WithRealTimeWordBoost([]string{"foo", "bar"}),
		WithRealTimeEncoding(RealTimeEncodingPCMMulaw),
		WithRealTimeSampleRate(8_000),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Send(ctx, []byte("foo"))
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)
}

func TestRealTime_Receive(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		transcript := FinalTranscript{
			MessageType: MessageTypeFinalTranscript,
			RealTimeBaseTranscript: RealTimeBaseTranscript{
				Text: "foo",
			},
		}

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		err = wsjson.Write(ctx, conn, transcript)
		require.NoError(t, err)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	var finalTranscriptInvoked bool
	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(ts.URL), WithRealTimeTranscriber(&RealTimeTranscriber{
		OnFinalTranscript: func(transcript FinalTranscript) {
			finalTranscriptInvoked = true
			require.Equal(t, "foo", transcript.Text)

			wg.Done()
		},
	}))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	wg.Wait()

	require.True(t, finalTranscriptInvoked)
}

func TestRealTime_TerminateSessionOnDisconnect(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	var sessionTerminatedInvoked bool

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(ts.URL), WithRealTimeTranscriber(&RealTimeTranscriber{
		OnSessionTerminated: func(_ SessionTerminated) {
			sessionTerminatedInvoked = true
		},
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	require.True(t, sessionTerminatedInvoked)
}

func TestRealTime_ForceEndUtterance(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		_, b, _ := conn.Read(ctx)
		require.Equal(t, `{"force_end_utterance":true}`, strings.TrimSpace(string(b)))

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.ForceEndUtterance(ctx)
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)
}

func TestRealTime_SetEndUtteranceSilenceThreshold(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		_, b, _ := conn.Read(ctx)
		require.Equal(t, `{"end_utterance_silence_threshold":350}`, strings.TrimSpace(string(b)))

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.SetEndUtteranceSilenceThreshold(ctx, 350)
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)
}

func TestRealTime_TemporaryToken(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.URL.Path == "/v2/realtime/token" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `{"token":"temp-token"}`)
			return
		}

		token := r.URL.Query().Get("token")
		require.Equal(t, "temp-token", token)

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Acquire a temporary token.
	client := NewClientWithOptions(
		WithBaseURL(ts.URL),
		WithAPIKey("api-key"),
	)

	resp, err := client.RealTime.CreateTemporaryToken(ctx, 480)
	require.NoError(t, err)

	// Use the temporary token.
	tokenClient := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeAuthToken(ToString(resp.Token)),
		WithRealTimeTranscriber(&RealTimeTranscriber{}),
	)

	err = tokenClient.Connect(ctx)
	require.NoError(t, err)

	err = tokenClient.Disconnect(ctx, true)
	require.NoError(t, err)
}

func beginSession(ctx context.Context, conn *websocket.Conn) error {
	return wsjson.Write(ctx, conn, map[string]interface{}{"message_type": string(MessageTypeSessionBegins)})
}

func terminateSession(ctx context.Context, conn *websocket.Conn) error {
	var message map[string]bool

	if err := wsjson.Read(ctx, conn, &message); err != nil {
		return err
	}

	if _, ok := message["terminate_session"]; ok {
		sessionTerminatedMessage := map[string]interface{}{"message_type": MessageTypeSessionTerminated}

		if err := wsjson.Write(ctx, conn, sessionTerminatedMessage); err != nil {
			return err
		}

		_, _, err := conn.Read(ctx)
		if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
			return err
		}

	}

	return nil
}

func upgradeRequest(w http.ResponseWriter, r *http.Request) (*websocket.Conn, func() error) {
	conn, _ := websocket.Accept(w, r, nil)

	return conn, func() error {
		return conn.Close(websocket.StatusInternalError, "websocket closed unexpectedly")
	}
}

func TestRealTime_DisablePartialTranscriptsIfNoCallback(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		disablePartialTranscripts := r.URL.Query().Get("disable_partial_transcripts")
		require.Equal(t, "true", disablePartialTranscripts)

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)
}

func TestRealTime_EnablePartialTranscriptsIfCallback(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		require.False(t, r.URL.Query().Has("disable_partial_transcripts"))

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeTranscriber(&RealTimeTranscriber{
			OnPartialTranscript: func(_ PartialTranscript) {},
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)
}
