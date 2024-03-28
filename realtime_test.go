package assemblyai

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const testTimeout = 5 * time.Second

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

	handler := &mockRealtimeHandler{}

	client := NewRealTimeClientWithOptions(
		WithRealTimeAPIKey("api-key"),
		WithRealTimeBaseURL(ts.URL),
		WithHandler(handler),
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

		var err error

		err = beginSession(ctx, conn)
		require.NoError(t, err)

		_, got, _ := conn.Read(ctx)

		require.Equal(t, []byte("foo"), got)

		err = terminateSession(ctx, conn)
		require.NoError(t, err)
	}))
	defer ts.Close()

	handler := &mockRealtimeHandler{}

	client := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithHandler(handler),
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

	done := make(chan bool)

	handler := &mockRealtimeHandler{
		FinalTranscriptFunc: func(transcript FinalTranscript) {
			require.Equal(t, "foo", transcript.Text)

			done <- true
		},
	}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(ts.URL), WithHandler(handler))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	<-done

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	require.True(t, handler.FinalTranscriptCalled)
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

	handler := &mockRealtimeHandler{
		SessionTerminatedFunc: func(event SessionTerminated) {},
	}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(ts.URL), WithHandler(handler))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	err = client.Connect(ctx)
	require.NoError(t, err)

	err = client.Disconnect(ctx, true)
	require.NoError(t, err)

	require.True(t, handler.SessionTerminatedCalled)
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

	handler := &mockRealtimeHandler{}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(ts.URL), WithHandler(handler))

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

	handler := &mockRealtimeHandler{}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(ts.URL), WithHandler(handler))

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
		if r.URL.Path == "/v2/realtime/token" {
			fmt.Fprintln(w, `{"token":"temp-token"}`)
			return
		}

		token := r.URL.Query().Get("token")
		require.Equal(t, "temp-token", token)

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		ctx := context.Background()

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

	var handler mockRealtimeHandler

	// Use the temporary token.
	tokenClient := NewRealTimeClientWithOptions(
		WithRealTimeBaseURL(ts.URL),
		WithRealTimeAuthToken(ToString(resp.Token)),
		WithHandler(&handler),
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

type mockRealtimeHandler struct {
	SessionBeginsCalled bool
	SessionBeginsFunc   func(event SessionBegins)

	SessionTerminatedCalled bool
	SessionTerminatedFunc   func(event SessionTerminated)

	FinalTranscriptCalled bool
	FinalTranscriptFunc   func(transcript FinalTranscript)

	PartialTranscriptCalled bool
	PartialTranscriptFunc   func(transcript PartialTranscript)

	ErrorCalled bool
	ErrorFunc   func(err error)
}

func (h *mockRealtimeHandler) SessionBegins(event SessionBegins) {
	h.SessionBeginsCalled = true
	if h.SessionBeginsFunc != nil {
		h.SessionBeginsFunc(event)
	}
}

func (h *mockRealtimeHandler) SessionTerminated(event SessionTerminated) {
	h.SessionTerminatedCalled = true
	if h.SessionTerminatedFunc != nil {
		h.SessionTerminatedFunc(event)
	}
}

func (h *mockRealtimeHandler) FinalTranscript(transcript FinalTranscript) {
	h.FinalTranscriptCalled = true
	h.FinalTranscriptFunc(transcript)
}

func (h *mockRealtimeHandler) PartialTranscript(transcript PartialTranscript) {
	h.PartialTranscriptCalled = true
	if h.PartialTranscriptFunc != nil {
		h.PartialTranscriptFunc(transcript)
	}
}

func (h *mockRealtimeHandler) Error(err error) {
	h.ErrorCalled = true
	if h.ErrorFunc != nil {
		h.ErrorFunc(err)
	}
}
