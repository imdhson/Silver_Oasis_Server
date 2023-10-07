package modules

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type receive_scrap struct {
	ID string `json:"id"`
}

func ScrapAddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // 모든 도메인에 대해 허용
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "POST" {
		//json을 수신하여 Recive Scrap에 저장
		var scrap_struct receive_scrap
		body, err := io.ReadAll(r.Body)
		log.Println("json 수신 원본: ", string(body))
		ErrOK(err)
		err = json.Unmarshal(body, &scrap_struct)
		ErrOK(err)
		defer r.Body.Close()
		if !IsHeLogin(w, r) { //로그인 안했으면
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			msg := map[string]string{"email": "Not LOGIN"}
			msg_json, _ := json.Marshal(msg)
			w.Write(msg_json)
			return
		}
		user_oid := SessionTO_oid(w, r)
		godotenv.Load()
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
		coll := db.Database("dj_users").Collection("users")

		filter := bson.D{{"_id", user_oid}}
		/*update := bson.D{
			{"$set", bson.D{{"settings.loc", settings_struct.Loc},
				{"settings.type1", settings_struct.Type1},
				{"settings.type2", settings_struct.Type2},
				{"settings.type3", settings_struct.Type3}}},
		}*/
		sid, err := primitive.ObjectIDFromHex(scrap_struct.ID)
		ErrOK(err)
		update := bson.D{
			{"$push", bson.D{{"scrapList", sid}}},
		}
		_, err = coll.UpdateOne(context.TODO(), filter, update)
		ErrOK(err)

		//직장리스트 목록에 scarpCount 추가
		func(sid primitive.ObjectID) {
			coll_for_scrapCount := db.Database("dj_jobs").Collection("job_list")
			filter_for_scrapCount := bson.D{{"_id", sid}}
			update_for_scrapCount := bson.D{
				{"$inc", bson.D{{"scrapCount", 1}}},
			}
			_, err := coll_for_scrapCount.UpdateOne(context.TODO(), filter_for_scrapCount, update_for_scrapCount)
			ErrOK(err)

		}(sid)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("good"))
	}

}
