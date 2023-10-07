package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type logout_struct struct {
	Logout string `json:"logout"`
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // 모든 도메인에 대해 허용
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "POST" {
		//json을 수신하여 settings_struct에 저장
		var logout_json logout_struct
		body, err := io.ReadAll(r.Body)
		fmt.Println("json 수신 원본: ", string(body))
		ErrOK(err)
		err = json.Unmarshal(body, &logout_json)
		ErrOK(err)
		defer r.Body.Close()
		if logout_json.Logout == "" {
			return
		}

		err = godotenv.Load()
		URI := os.Getenv("MONGODB_URI")
		if URI == "" {
			Critical(err)
		}

		user_oid := SessionTO_oid(w, r)

		db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
		Critical(err)
		defer func() {
			err := db.Disconnect(context.TODO())
			Critical(err)
		}()

		cookie := &http.Cookie{Name: "dj_session", Value: "", Path: "/", MaxAge: -1}
		http.SetCookie(w, cookie)

		//db에 세션 클리어
		coll_dj_session := db.Database("gd_users").Collection("sessions")
		filter := bson.D{{"dj_user_id", user_oid}}
		result, err := coll_dj_session.DeleteMany(context.TODO(), filter)
		ErrOK(err)
		log.Println("로그아웃: session 삭제:", result.DeletedCount)
	}
}
