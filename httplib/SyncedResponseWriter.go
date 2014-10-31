package httplib

import "net/http"

// SyncedResponseWriter is a wrapper for the http.ResponseWriter
// A call to writeHeader will not actually send the header, but instead store the status-code until Write or
// SendHeader is called
type SyncedResponseWriter struct{
	writer        http.ResponseWriter
	headerWritten bool
	status        int
}

func NewSyncedResponseWriter(rw http.ResponseWriter) *SyncedResponseWriter{
	return &SyncedResponseWriter{writer: rw, headerWritten: false, status: 200}
}

// Header returns the header map that will be sent by WriteHeader.
func (w *SyncedResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply,
func (w *SyncedResponseWriter) Write(p []byte) (int, error) {
	if !w.headerWritten {
		w.SendHeader()
	}
	return w.writer.Write(p)
}

// WriteHeader stores an http status code for future write
func (w *SyncedResponseWriter) WriteHeader(code int) {
	w.status = code
}

// SendHeader sends an HTTP response header with status code,
// and sets `started` to true.
func (w *SyncedResponseWriter) SendHeader() {
	if w.status == 0 {
		w.WriteHeader(200)
	}
	w.headerWritten = true
	w.writer.WriteHeader(w.status)
}

// IsHeaderWritten returns weather the header has been sent or not
func (w *SyncedResponseWriter) IsHeaderWritten() bool {
	return w.headerWritten
}
