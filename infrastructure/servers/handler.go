package servers

import "net/http"

// Handler は、ハンドラの構造を表現した構造体です。
type Handler struct {
	Pattern     string
	HandlerFunc func(http.ResponseWriter, *http.Request, *Logger) (status int, body []byte)
}
