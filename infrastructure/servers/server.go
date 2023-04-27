package servers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Server 構造体は、サーバの基本操作を提供します。
type Server struct {
	port   int
	logger *Logger
}

// Start はサーバをスタートし、HTTPリクエストを受け付けられる状態にします。
func (s *Server) Start() {
	loggingHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// ステータスコードを知りたいので、http.ResponseWriterをラップしたものに差し替える。
			newWriter := ResponseWriter{writer: w}
			h.ServeHTTP(&newWriter, r)
			length, err := strconv.Atoi(w.Header().Get("Content-Length"))
			if err != nil {
				length = 0
			}

			responseLog := ResponseLog{newWriter.Status(), time.Since(start), length}
			s.logger.FPrintAccessLog(r, &responseLog)
		})
	}
	s.logger.Println("start server.")
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), loggingHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

// AddHandlers は、サーバにハンドラを追加します。
func (s *Server) AddHandlers(handlers []Handler) {
	for _, handler := range handlers {
		handleFunc := handler.HandlerFunc
		http.HandleFunc(handler.Pattern, func(w http.ResponseWriter, r *http.Request) {
			// CORS用設定
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,DELETE,OPTIONS,PATCH")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			status, body := handleFunc(w, r, s.logger)
			w.WriteHeader(status)
			w.Header().Set("Content-Length", fmt.Sprint(len(body)))
			w.Write(body)
		})
	}
}

func NewServer(port int, logger *Logger) (*Server, error) {
	if port > 65535 || port < 0 {
		return nil, fmt.Errorf("port must be (0 ~ 65535)")
	}
	return &Server{logger: logger, port: port}, nil
}
