package servers

import "net/http"

// ResponseWriter は、通常のhttp.ResponseWriterの機能に加え、ステータスコードを保持できる構造体です。
type ResponseWriter struct {
	writer http.ResponseWriter
	status int
}

func (w *ResponseWriter) Write(bytes []byte) (int, error) {
	return w.writer.Write(bytes)
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
	w.writer.WriteHeader(status)
}

func (w *ResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *ResponseWriter) Status() int {
	return w.status
}
