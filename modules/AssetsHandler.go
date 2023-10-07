package modules

import (
	"net/http"
	"os"
)

func AssetsHanlder(w http.ResponseWriter, r *http.Request, url *string) {
	docsfile := DotFileType(*url)
	switch docsfile {
	case "jpg":
		w.Header().Set("Content-Type", "image/jpg; charset=utf-8")
	case "png":
		w.Header().Set("Content-Type", "image/png; charset=utf-8")
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg; charset=utf-8")
	case "js":
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	case "css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case "scss":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case "html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case "mp4":
		w.Header().Set("Content-Type", "video/mp4; charset=utf-8")
	case "avi":
		w.Header().Set("Content-Type", "video/avi; charset=utf-8")
	case "mov":
		w.Header().Set("Content-Type", "video/mov; charset=utf-8")
	case "webm":
		w.Header().Set("Content-Type", "video/webm; charset=utf-8")
	}
	wwwfile, err := os.ReadFile("www/" + *url) // www/assets/main.css 와 같이 작동하게 됨
	ErrOK(err)
	w.Write(wwwfile)
}
