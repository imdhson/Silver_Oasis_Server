package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	BATCHSIZE_FOR_AILIST = 100
	OUTPUTSIZE           = 30
	LOC1_MATCH_SCORE     = 500
	LOC2_MATCH_SCORE     = 150
	VIEWCOUNT_RATIOBYONE = 2 // 지역구분1,2가 겹치면  score에 2번 더해짐을 고려해야함.
)

func (a *SO_jobs_detail_s) will_send_append(i_detail *SO_jobs_detail, score int) {
	var now_index uint8
	for i, v := range *a {
		if v.ID == (*i_detail).ID {
			fmt.Println("겹침!!!!!------")
			(*a)[i].AI_List_score += score
			return
		} else if i+1 >= BATCHSIZE_FOR_AILIST {
			//isFull 추가해서 성능이슈 해결
			return
		} else if v.ID == primitive.NilObjectID {
			now_index = uint8(i)
			break
		}
	}
	fmt.Println(now_index)
	i_detail.AI_List_score += score
	(*a)[now_index] = *i_detail //버그 있는 곳
}

func (a *SO_jobs_detail_s) serviceScoreAdd(i_settings Dj_users_users_settings, i_detail SO_jobs_detail) {
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
	// gd_services_score에 맞게 더해야함
	//더해야할 서비스 score만 남기고 0으로 만들기

	//타입1,2,3 선택한 것에서 중복 제거, map의 키가 겹칠수 없음을 집합처럼 응용
	setofServices := map[string]bool{}
	setofServices[i_settings.Service1] = true
	setofServices[i_settings.Service2] = true
	setofServices[i_settings.Service3] = true
	//TypeofFacility 가져와서 gd_services_score를 구함
	collection_SO_services := db.Database("gd_facilities").Collection("gd_services_score")
	filter_for_SO_services := bson.D{{"시설종류", i_detail.TypeofFacility}}
	var services_score_board SO_service_type
	err = collection_SO_services.FindOne(context.TODO(), filter_for_SO_services).Decode(&services_score_board)
	ErrOK(err)
	for service := range setofServices {
		//service에 겹치지않은 서비스가 저장되어있음.
		switch service {
		case "방문서비스":
			//will_send_append를 이용하여 각 서비스 만큼 스코어를 더하여줌.
			a.will_send_append(&i_detail, services_score_board.VisitService)
		case "물품지원":
			a.will_send_append(&i_detail, services_score_board.ObjectSupport)
		case "복지서비스":
			a.will_send_append(&i_detail, services_score_board.WelfareService)
		case "생활지원":
			a.will_send_append(&i_detail, services_score_board.LifeSupport)
		}
	}
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
	var filter_loc_0 string
	var filter_loc_1 string
	if len(splited_loc) <= 0 { //빈칸일경우 모든 지역 포함간주
		filter_loc_0 = ""
		filter_loc_1 = ""
	} else if len(splited_loc) == 1 {
		filter_loc_0 = splited_loc[0]
		filter_loc_1 = ""
	} else {
		filter_loc_0 = splited_loc[0]
		filter_loc_1 = splited_loc[1]
	}

	//splited[0] and splited[1] = loc 1 and loc2 인 경우를 cursor.next형 탐색
	collection_SO_list := db.Database("gd_facilities").Collection("gd_fac_list")
	filter_for_SO_list := bson.D{
		{"$and", bson.A{
			bson.D{{"지역구분1", bson.D{{"$regex", filter_loc_0}}}},
			bson.D{{"지역구분2", bson.D{{"$regex", filter_loc_1}}}},
		}}}
	cursor_for_SO_list, err := collection_SO_list.Find(context.TODO(), filter_for_SO_list)
	ErrOK(err)
	defer cursor_for_SO_list.Close(context.TODO())

	now_batch := 0
	var will_send_ARR SO_jobs_detail_s //willsend 객체 선언

	//순회하며 배열에 담기 시작
	for cursor_for_SO_list.Next(context.TODO()) { //커서next가 성공하면 참
		if now_batch > BATCHSIZE_FOR_AILIST {
			break
		}
		var dbres_GD_Detail_t SO_jobs_detail
		cursor_for_SO_list.Decode(&dbres_GD_Detail_t)
		will_send_ARR.will_send_append(&dbres_GD_Detail_t, LOC2_MATCH_SCORE)
		will_send_ARR[len(will_send_ARR)-1].AI_List_score += dbres_GD_Detail_t.ViewCount * VIEWCOUNT_RATIOBYONE //조회수만큼 점수 더해줌.
		now_batch++
	}

	//filterloc[0] = loc1 인 경우를 cursor.next형 탐색
	filter_for_SO_list = bson.D{
		{"$and", bson.A{
			bson.D{{"지역구분1", bson.D{{"$regex", filter_loc_0}}}},
		}}}
	cursor_for_SO_list, err = collection_SO_list.Find(context.TODO(), filter_for_SO_list)
	//ErrOK(err)
	defer cursor_for_SO_list.Close(context.TODO())

	//순회하며 배열에 담기 시작
	for cursor_for_SO_list.Next(context.TODO()) { //커서next가 성공하면 참
		if now_batch > BATCHSIZE_FOR_AILIST {
			break
		}
		var dbres_GD_Detail_t SO_jobs_detail
		cursor_for_SO_list.Decode(&dbres_GD_Detail_t)
		will_send_ARR.will_send_append(&dbres_GD_Detail_t, LOC1_MATCH_SCORE)
		will_send_ARR[len(will_send_ARR)-1].AI_List_score += dbres_GD_Detail_t.ViewCount * VIEWCOUNT_RATIOBYONE //조회수만큼 점수 더해줌.
		now_batch++
	}

	//will_send_ARR 순회하며 scoreADD 호출.
	//ScoreAdd는 선택한 서비스에 대한 것만 점수에 추가하여줌.
	for _, v := range will_send_ARR {
		will_send_ARR.serviceScoreAdd(user_struct.Settings, v)
	}

	//score을 기반으로 sort 시작
	sort.Sort(sort.Reverse(will_send_ARR))

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
	ai_list_num := 0
	for numi := range will_send_refined {
		will_send_refined[numi].AI_List_num = ai_list_num
		ai_list_num++
	}

	will_send_json, _ := json.MarshalIndent(will_send_refined, " ", "	")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(will_send_json)
}
