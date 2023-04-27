package servers

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Logger はロギングを行います。
type Logger struct {
	accessLogsDir string
	errorLogsDir  string
}

// Println は、現在時刻付きで標準出力に、与えられた文字列を出力します。
func (l *Logger) Println(content string) {
	content = fmt.Sprintf("%s%s", l.getTimeString(time.Now().Local()), content)
	fmt.Println(content)
}

// FPrintErrorLog は、エラーログファイルにエラーの内容を現在時刻付きで書き込みます。
func (l *Logger) FPrintErrorLog(err error, additional string) {
	now := l.getLocalNow()
	log := fmt.Sprintf("%s, %s", l.getTimeString(now), err)
	if additional != "" {
		log += fmt.Sprintf(", %s", additional)
	}

	fileName := fmt.Sprintf("%s/ErrorLog-%s.txt", l.errorLogsDir, now.Format("2006-01-02"))
	fp, err := l.openLogFile(fileName, "time, error, additional")
	if err != nil {
		fmt.Println(NewLoggingError("could not print the error log to file.", err.Error()).Error())
	}
	fmt.Fprintln(fp, log)
	fmt.Println(log)
	fp.Close()
}

// FPrintaccessLog は、アクセスログファイルにアクセス内容、レスポンス内容を現在時刻付きで書き込みます。
func (l *Logger) FPrintAccessLog(req *http.Request, res *ResponseLog) {
	now := l.getLocalNow()
	log := fmt.Sprintf("%s, %s, %s, %s, %s, %d, %d, %d, %s", l.getTimeString(now), req.RemoteAddr, req.RequestURI, req.Method, req.Proto, res.Status, res.Since, res.ContentLength, req.UserAgent())
	fileName := fmt.Sprintf("%s/AcccesLog-%s.txt", l.accessLogsDir, now.Format("2006-01-02"))
	fp, err := l.openLogFile(fileName, "time, remoteAddr, reqUri, method, proto, status, elapsed, length, userAgent")
	if err != nil {
		l.FPrintErrorLog(NewLoggingError("could not print the access log to file.", err.Error()), "")
	}
	fmt.Fprintln(fp, log)
	fp.Close()
}

// openLogFile はログファイルを開き、ファイルポインタを返却します。
func (l *Logger) openLogFile(fileName string, firstContent string) (fp *os.File, err error) {
	fp, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)
	if os.IsNotExist(err) {
		fp, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}
		fmt.Fprintln(fp, firstContent)
	} else if err != nil {
		return nil, err
	}
	return fp, nil
}

// getTimeString は、与えた時刻をフォーマットした文字列で返却します。
func (l Logger) getTimeString(t time.Time) string {
	return fmt.Sprintf("[%s]", t.Format("2006-01-02 15:04:05"))
}

// getLocalNow は、現在のローカルタイムを返却します。
func (l Logger) getLocalNow() time.Time {
	return time.Now().Local()
}

// tryCreateLogsDirs は、ログを格納するディレクトリを作成します。
func (l Logger) tryCreateLogsDirs() {
	paths := []string{
		l.accessLogsDir,
		l.errorLogsDir,
	}
	for _, path := range paths {
		err := os.Mkdir(path, 0666)
		if os.IsExist(err) {
			continue
		}
		if err != nil {
			l.Println(err.Error())
		}
	}
}

func NewLogger() *Logger {
	logger := Logger{accessLogsDir: "AccessLogs", errorLogsDir: "ErrorLogs"}
	logger.tryCreateLogsDirs()
	return &logger
}
