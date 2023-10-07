package modules

import (
	"log"
	"net/http"
)

func Http_TO_https(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:] //sv_urlpath에 유저가 어떤 url을 요청했는지 저장됨
	log.Printf("%v/%v : http 요청임으로 https로 리다이렉션", GetIP(r), url)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msg := "<meta http-equiv=\"refresh\" content=\"0;url=https://pi.imdhson.com/"
	msg += url
	msg += "\"></meta>"
	w.Write([]byte(msg))
}
