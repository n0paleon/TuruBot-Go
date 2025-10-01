package logstream

import (
	"context"
	"fmt"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

var (
	BaseUrl = "https://logs.cettalabs.com"
)

type actionType string

var (
	actionPushLog      actionType = "push_log"
	actionSetNote      actionType = "set_note"
	actionPushFinished actionType = "push_finished"
)

type wartaLogStream_ActionRequest struct {
	Action actionType `json:"action"`
	Data   string     `json:"data"`
}

type WartaLogStream struct {
	id         string
	token      string
	context    context.Context
	httpClient *resty.Client
	wsConn     *websocket.Conn
}

type wartaLogStream_CreateResponse struct {
	SessionId    string `json:"session_id"`
	SessionToken string `json:"token"`
}

func NewWartaLogStream(ctx context.Context) (*WartaLogStream, error) {
	l := &WartaLogStream{
		httpClient: resty.New(),
	}

	var response wartaLogStream_CreateResponse
	rctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := l.httpClient.R().
		SetContext(rctx).
		SetResult(&response).
		Post(BaseUrl + "/session")

	if err != nil {
		return nil, fmt.Errorf("failed to create new log stream: %v", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to create new log stream, status code: %v", resp.StatusCode())
	}

	l.id = response.SessionId
	l.token = response.SessionToken
	l.context = ctx

	if err := l.connect(ctx); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *WartaLogStream) GetStreamUrl() string {
	return fmt.Sprintf("%s/?session=%s", BaseUrl, l.id)
}

func (l *WartaLogStream) PushLog(data string) {
	go l.write(actionPushLog, data)
}

func (l *WartaLogStream) SetNote(data string) {
	go l.write(actionSetNote, data)
}

func (l *WartaLogStream) write(action actionType, data string) error {
	return wsjson.Write(l.context, l.wsConn, map[string]interface{}{
		"action": action,
		"data":   data,
	})
}

func (l *WartaLogStream) connect(ctx context.Context) error {
	headers := http.Header{}
	headers.Add("x-writer-token", l.token)

	opts := &websocket.DialOptions{
		HTTPHeader: headers,
	}

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("%s/session/%s/ws/writer", BaseUrl, l.id), opts)
	if err != nil {
		return fmt.Errorf("failed to open websocket connection: %v", err)
	}
	l.wsConn = c

	go func() {
		<-ctx.Done()
		_ = l.wsConn.Close(websocket.StatusNormalClosure, "context canceled")
	}()

	return nil
}

func (l *WartaLogStream) Close() {
	_ = l.write(actionPushFinished, "")
	_ = l.wsConn.Close(websocket.StatusNormalClosure, "context canceled")
}
