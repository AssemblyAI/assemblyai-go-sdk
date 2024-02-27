package assemblyai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const testTimeout = 5 * time.Second

func TestRealTime_Send(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		if err := beginSession(ctx, conn); err != nil {
			t.Error(err)
		}

		_, b, _ := conn.Read(ctx)

		got := strings.TrimSpace(string(b))
		want := `{"audio_data":"Zm9v"}`

		if got != want {
			t.Errorf("message = %v, want %v", got, want)
		}

		if err := terminateSession(t, ctx, conn); err != nil {
			t.Errorf("terminateSession returned error: %v", err)
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	handler := &mockRealtimeHandler{}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(wsURL), WithHandler(handler))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect returned error: %v", err)
	}

	if err := client.Send(ctx, []byte("foo")); err != nil {
		t.Errorf("Send returned error: %v", err)
	}

	if err := client.Disconnect(ctx, true); err != nil {
		t.Errorf("Disconnect returned error: %v", err)
	}
}

func TestRealTime_Receive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		if err := beginSession(ctx, conn); err != nil {
			t.Error(err)
		}

		transcript := FinalTranscript{
			MessageType: MessageTypeFinalTranscript,
			RealTimeBaseTranscript: RealTimeBaseTranscript{
				Text: "foo",
			},
		}

		if err := wsjson.Write(ctx, conn, transcript); err != nil {
			t.Error(err)
		}

		if err := terminateSession(t, ctx, conn); err != nil {
			t.Errorf("terminateSession returned error: %v", err)
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	done := make(chan bool)

	handler := &mockRealtimeHandler{
		FinalTranscriptFunc: func(transcript FinalTranscript) {
			want := "foo"

			if transcript.Text != want {
				t.Errorf("transcript.Text = %v, want %v", transcript.Text, want)
			}

			done <- true
		},
	}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(wsURL), WithHandler(handler))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect returned error: %v", err)
	}

	<-done

	if !handler.FinalTranscriptCalled {
		t.Error("missing final transcript")
	}

	if err := client.Disconnect(ctx, true); err != nil {
		t.Errorf("Disconnect returned error: %v", err)
	}

}

func TestRealTime_TerminateSessionOnDisconnect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		if err := beginSession(ctx, conn); err != nil {
			t.Error(err)
		}

		if err := terminateSession(t, ctx, conn); err != nil {
			t.Errorf("terminateSession returned error: %v", err)
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	handler := &mockRealtimeHandler{
		SessionTerminatedFunc: func(event SessionTerminated) {},
	}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(wsURL), WithHandler(handler))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	if err := client.Disconnect(ctx, true); err != nil {
		t.Errorf("Disconnect returned error: %v", err)
	}

	if !handler.SessionTerminatedCalled {
		t.Error("missing SessionTerminated event")
	}
}

func TestRealTime_ForceEndUtterance(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		if err := beginSession(ctx, conn); err != nil {
			t.Error(err)
		}

		_, b, _ := conn.Read(ctx)

		got := strings.TrimSpace(string(b))
		want := `{"force_end_utterance":true}`

		if got != want {
			t.Errorf("message = %v, want %v", got, want)
		}

		if err := terminateSession(t, ctx, conn); err != nil {
			t.Errorf("terminateSession returned error: %v", err)
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	handler := &mockRealtimeHandler{}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(wsURL), WithHandler(handler))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect returned error: %v", err)
	}

	if err := client.ForceEndUtterance(ctx); err != nil {
		t.Errorf("ForceEndUtterance returned error: %v", err)
	}

	if err := client.Disconnect(ctx, true); err != nil {
		t.Errorf("Disconnect returned error: %v", err)
	}
}

func TestRealTime_SetEndUtteranceSilenceThreshold(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, teardown := upgradeRequest(w, r)
		defer teardown()

		if err := beginSession(ctx, conn); err != nil {
			t.Error(err)
		}

		_, b, _ := conn.Read(ctx)

		got := strings.TrimSpace(string(b))
		want := `{"end_utterance_silence_threshold":350}`

		if got != want {
			t.Errorf("message = %v, want %v", got, want)
		}

		if err := terminateSession(t, ctx, conn); err != nil {
			t.Errorf("terminateSession returned error: %v", err)
		}
	}))
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	handler := &mockRealtimeHandler{}

	client := NewRealTimeClientWithOptions(WithRealTimeBaseURL(wsURL), WithHandler(handler))

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatalf("Connect returned error: %v", err)
	}

	if err := client.SetEndUtteranceSilenceThreshold(ctx, 350); err != nil {
		t.Errorf("SetEndUtteranceSilenceThreshold returned error: %v", err)
	}

	if err := client.Disconnect(ctx, true); err != nil {
		t.Errorf("Disconnect returned error: %v", err)
	}
}

func beginSession(ctx context.Context, conn *websocket.Conn) error {
	return wsjson.Write(ctx, conn, map[string]interface{}{"message_type": string(MessageTypeSessionBegins)})
}

func terminateSession(t *testing.T, ctx context.Context, conn *websocket.Conn) error {
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
