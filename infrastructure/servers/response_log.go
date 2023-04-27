package servers

import "time"

// ResponseLog は、レスポンス内容のログを書き込むためのデータを表現しています。
type ResponseLog struct {
	Status        int
	Since         time.Duration
	ContentLength int
}
