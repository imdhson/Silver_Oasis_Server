package modules

import (
	"net/http"
)

func WebViewExit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msg := "<style> *{  background-color: #ffe8c4; }</style>"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
