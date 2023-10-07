package modules

import (
	"net/http"
	"os"
)

func ArticlesInsertPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	wwwfile, err := os.ReadFile("./www/articles_insert.html")
	Critical(err)
	w.Write(wwwfile)
}
