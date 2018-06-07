package appclient

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

func dial(urlStr string, token Token) (*websocket.Conn, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("constructing request %s: %v", urlStr, err)
	}

	token.Set(req)

	conn, _, err := websocket.DefaultDialer.Dial(urlStr, req.Header)
	if err != nil {
		return nil, fmt.Errorf("dial error %s: %v", urlStr, err)
	}

	return conn, nil
}

func dialWS(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err != nil {
		return nil, resp, fmt.Errorf("dial error %s: %v", urlStr, err)
	}
	return conn, resp, nil
}

func IsExpectedWSCloseError(err error) bool {
	return err == io.EOF || err == io.ErrClosedPipe || websocket.IsCloseError(err,
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseNoStatusReceived,
		websocket.CloseAbnormalClosure,
	)
}
