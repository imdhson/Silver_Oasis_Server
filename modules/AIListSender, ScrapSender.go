package modules

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BATCHSIZE        = 3000
	OUTPUTSIZE       = 50
	LOC1_MATCH_SCORE = 500
	LOC2_MATCH_SCORE = 150
)

func serviceCheck(Dj_users_users_settings, *GD_jobs_detail) {
	//타입1,2,3 선택한 것에서 중복 제거하고
	// gd_services_score에 맞게 더해야함
	//더해야할 서비스 score만 남기고 0으로 만들기

	//레퍼런스로 찾아가서 점수를 더해줌.

}

func AIListSender(w http.ResponseWriter, r *http.Request) { //메인화면 시설 리스트
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
	splited_loc := strings.Split(user_struct.Settings.Loc, " ")
	if !IsHeLogin(w, r) { //인덱스 런타임 에러 방지
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		err_msg := map[string]string{"error": "Not LOGIN"}
		err_msg_json, _ := json.MarshalIndent(err_msg, " ", "	")
		w.Write(err_msg_json)
		return
	}
	//collection_GD_services := db.Database("gd_facilities").Collection("gd_services_score")
	// !!!!! filter_for_TypeofFac := bson.D{{"시설종류", "1"}} // 변수 수정 필요 ! <- 리스트에서 가져온 시설종류를 기입하면 가중치 반환하려고 만든 필터

	// 시설목록가져와서 시설종류-weightDB에서 선택된 서비스만 +=score
	// 그리고 splited_loc[0] 과 [1]에 각각 포함되어있는지 확인해서 +=score

	//splited[0]과 loc1을 상대로 검색 시작
	collection_GD_list := db.Database("gd_facilities").Collection("gd_fac_list")
	filter_for_GD_list := bson.D{{"지역구분1", splited_loc[0]}}
	cursor_for_GD_list, err := collection_GD_list.Find(context.TODO(), filter_for_GD_list)

	//next 이전에 willsendappend 수행
	var dbres_GD_Detail_t1 GD_jobs_detail
	cursor_for_GD_list.Decode(&dbres_GD_Detail_t1)
	var will_send_ARR GD_jobs_detail_s                                   //willsend 객체 선언
	will_send_ARR.will_send_append(dbres_GD_Detail_t1, LOC1_MATCH_SCORE) //지역1 매치스코어 만큼 더해지게됨.
	cnt := 0                                                             //지역 1을 대상으로 순회하며 willsendappend 수행
	for cursor_for_GD_list.Next(context.TODO()) {                        //커서next가 성공하면 참
		if cnt > BATCHSIZE {
			break
		}
		var dbres_GD_Detail_t2 GD_jobs_detail
		cursor_for_GD_list.Decode(&dbres_GD_Detail_t2)
		will_send_ARR.will_send_append(dbres_GD_Detail_t2, LOC1_MATCH_SCORE)
		cnt++
	}

}
