package modules

import (
	"net/http"
)

func WebViewExit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msg := "<style> *{  background-color: #F2DCFF; }</style>"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
