package modules

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sort"

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
	//score을 기반으로 sort 시작
	sort.Sort(sort.Reverse(will_send_ARR))
	ai_list_num := 0
	for numi, _ := range will_send_ARR {
		will_send_ARR[numi].AI_List_num = ai_list_num
		ai_list_num++
	}

	var Outputsize_var int               //결과 슬라이싱시 인덱스 바깥으로 튀는것 방지하기 위함
	if len(will_send_ARR) < OUTPUTSIZE { //결과 슬라이싱시 인덱스 바깥으로 튀는것 방지하기 위함
		Outputsize_var = len(will_send_ARR)
	} else {
		Outputsize_var = OUTPUTSIZE
	}

	//필요한 만큼 outputsize로 자르고 메인에서 필요한 데이터만 남김
	var will_send_refined []SO_jobs_refined
	for ir, vr := range will_send_ARR {
		if ir > Outputsize_var {
			break
		}
		tmp_address1 := vr.Loc_1 + " " + vr.Loc_2
		tmp := SO_jobs_refined{
			AI_List_num:    vr.AI_List_num,
			ID:             vr.ID,
			NameofFacility: vr.NameofFacility,
			TypeofFacility: vr.TypeofFacility,
			Address:        tmp_address1,
			Operator:       vr.Operator,
			ViewCount:      vr.ViewCount,
		}
		will_send_refined = append(will_send_refined, tmp)
	}

	will_send_json, _ := json.MarshalIndent(will_send_refined, " ", "	")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(will_send_json)
}
