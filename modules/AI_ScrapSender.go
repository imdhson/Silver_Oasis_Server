package modules

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ScrapSender(w http.ResponseWriter, r *http.Request) { //메인화면 시설 리스트
	err := godotenv.Load()
	Critical(err)
	URI := os.Getenv("MONGODB_URI")
	if URI == "" {
		Critical(err)
	}
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	Critical(err)
	defer func() {
		err := db.Disconnect(context.TODO())
		Critical(err)
	}()
	user_struct := OidTOuser_struct(SessionTO_oid(w, r))
	if !IsHeLogin(w, r) { //인덱스 런타임 에러 방지
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err_msg := map[string]string{"error": "Not LOGIN"}
		err_msg_json, _ := json.MarshalIndent(err_msg, " ", "	")
		w.Write(err_msg_json)
		return
	}
	// 지역관련 쿼리는 수행안함. 서비스유형으로만 sort 예정

	var will_send_ARR SO_jobs_detail_s
	//user.scrapList 쿼리 시작
	for _, v := range user_struct.ScrapList {
		SO_fac_detail_inRange, err := OidTOjobDetail(v)
		will_send_ARR.will_send_append(&SO_fac_detail_inRange, 0)
		ErrOK(err)
	}
	//scoreADD 수행
	for _, v := range will_send_ARR {
		will_send_ARR.serviceScoreAdd(user_struct.Settings, v)
	}
}
