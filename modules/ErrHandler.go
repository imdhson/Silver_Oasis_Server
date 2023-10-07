package modules

import (
	"net/http"
	"os"
)

func ErrHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	wwwfile, err := os.ReadFile("./www/error.html")
	Critical(err)
	w.Write(wwwfile)
}
